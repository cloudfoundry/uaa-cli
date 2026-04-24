package cmd

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"github.com/spf13/cobra"
)

var signingKey string
var decodeTimes bool

// knownTimestampFields maps JWT claim names to human-readable labels.
var knownTimestampFields = []struct {
	field string
	label string
}{
	{"iat", "Issued At"},
	{"exp", "Expires At"},
	{"nbf", "Not Before"},
	{"auth_time", "Auth Time"},
	{"updated_at", "Updated At"},
}

func decodeJWTPayload(tokenStr string) (map[string]interface{}, []byte, []byte, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, nil, nil, errors.New("invalid JWT: expected 3 parts separated by '.'")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid JWT payload encoding: %v", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid JWT payload JSON: %v", err)
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid JWT header encoding: %v", err)
	}

	sigBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid JWT signature encoding: %v", err)
	}
	_ = sigBytes

	return claims, headerBytes, sigBytes, nil
}

func verifyJWTSignature(tokenStr string, keyPEM string) error {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return errors.New("invalid JWT format")
	}

	var header map[string]interface{}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return fmt.Errorf("invalid JWT header: %v", err)
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return fmt.Errorf("invalid JWT header JSON: %v", err)
	}

	alg, _ := header["alg"].(string)
	signingInput := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("invalid JWT signature: %v", err)
	}

	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return errors.New("failed to decode PEM block from key")
	}

	switch alg {
	case "RS256", "RS384", "RS512":
		pub, err := parseRSAPublicKey(block)
		if err != nil {
			return err
		}
		return verifyRSA(alg, signingInput, sig, pub)
	case "ES256", "ES384", "ES512":
		pub, err := parseECPublicKey(block)
		if err != nil {
			return err
		}
		return verifyECDSA(alg, signingInput, sig, pub)
	default:
		return fmt.Errorf("unsupported algorithm: %s", alg)
	}
}

func parseRSAPublicKey(block *pem.Block) (*rsa.PublicKey, error) {
	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %v", err)
		}
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("key is not an RSA public key")
		}
		return rsaKey, nil
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %v", err)
		}
		rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("certificate does not contain an RSA public key")
		}
		return rsaKey, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type for RSA: %s", block.Type)
	}
}

func parseECPublicKey(block *pem.Block) (*ecdsa.PublicKey, error) {
	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %v", err)
		}
		ecKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("key is not an EC public key")
		}
		return ecKey, nil
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %v", err)
		}
		ecKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("certificate does not contain an EC public key")
		}
		return ecKey, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type for EC: %s", block.Type)
	}
}

func verifyRSA(alg, signingInput string, sig []byte, pub *rsa.PublicKey) error {
	var hash crypto.Hash
	var h interface {
		Write([]byte) (int, error)
		Sum([]byte) []byte
	}
	switch alg {
	case "RS256":
		hash = crypto.SHA256
		s := sha256.New()
		h = s
	case "RS384":
		hash = crypto.SHA384
		s := sha512.New384()
		h = s
	case "RS512":
		hash = crypto.SHA512
		s := sha512.New()
		h = s
	}
	h.Write([]byte(signingInput))
	digest := h.Sum(nil)
	if err := rsa.VerifyPKCS1v15(pub, hash, digest, sig); err != nil {
		return errors.New("invalid token signature")
	}
	return nil
}

func verifyECDSA(alg, signingInput string, sig []byte, pub *ecdsa.PublicKey) error {
	var h interface {
		Write([]byte) (int, error)
		Sum([]byte) []byte
	}
	var keySize int
	switch alg {
	case "ES256":
		s := sha256.New()
		h = s
		keySize = 32
	case "ES384":
		s := sha512.New384()
		h = s
		keySize = 48
	case "ES512":
		s := sha512.New()
		h = s
		keySize = 66
	default:
		return fmt.Errorf("unsupported EC algorithm: %s", alg)
	}
	h.Write([]byte(signingInput))
	digest := h.Sum(nil)

	if len(sig) != 2*keySize {
		return errors.New("invalid token signature: unexpected signature length")
	}
	r := new(big.Int).SetBytes(sig[:keySize])
	s := new(big.Int).SetBytes(sig[keySize:])

	// Validate that the public key curve matches the expected key size.
	switch alg {
	case "ES256":
		if pub.Curve != elliptic.P256() {
			return errors.New("key curve does not match ES256")
		}
	case "ES384":
		if pub.Curve != elliptic.P384() {
			return errors.New("key curve does not match ES384")
		}
	case "ES512":
		if pub.Curve != elliptic.P521() {
			return errors.New("key curve does not match ES512")
		}
	}

	if !ecdsa.Verify(pub, digest, r, s) {
		return errors.New("invalid token signature")
	}
	return nil
}

func printDecodedTimestamps(claims map[string]interface{}) {
	now := time.Now()
	printed := false
	for _, tf := range knownTimestampFields {
		v, ok := claims[tf.field]
		if !ok {
			continue
		}
		var epoch int64
		switch n := v.(type) {
		case float64:
			epoch = int64(n)
		case json.Number:
			epoch, _ = n.Int64()
		default:
			continue
		}
		t := time.Unix(epoch, 0).UTC()
		rel := relativeTime(t, now)
		if !printed {
			log.Info("--- Decoded timestamps ---")
			printed = true
		}
		log.Info(fmt.Sprintf("%-12s %-16s %s  (%s)", tf.field, "("+tf.label+"):", t.Format("2006-01-02 15:04:05 UTC"), rel))
	}
}

func relativeTime(t, now time.Time) string {
	diff := t.Sub(now)
	abs := diff
	if abs < 0 {
		abs = -abs
	}

	var unit string
	var n int64
	switch {
	case abs < time.Minute:
		n = int64(abs.Seconds())
		unit = "second"
	case abs < time.Hour:
		n = int64(abs.Minutes())
		unit = "minute"
	case abs < 24*time.Hour:
		n = int64(abs.Hours())
		unit = "hour"
	default:
		n = int64(abs.Hours() / 24)
		unit = "day"
	}
	if n != 1 {
		unit += "s"
	}
	if diff < 0 {
		return fmt.Sprintf("%d %s ago", n, unit)
	}
	return fmt.Sprintf("in %d %s", n, unit)
}

func DecodeTokenCmd(cfg config.Config, args []string) error {
	var tokenStr string

	if len(args) > 0 {
		tokenStr = args[0]
	} else {
		ctx := cfg.GetActiveContext()
		if ctx.Token.AccessToken == "" {
			return errors.New("no token provided and no token found in active context")
		}
		tokenStr = ctx.Token.AccessToken
	}

	claims, _, _, err := decodeJWTPayload(tokenStr)
	if err != nil {
		return err
	}

	if signingKey != "" {
		if err := verifyJWTSignature(tokenStr, signingKey); err != nil {
			return err
		}
		log.Info("Valid token signature.")
	}

	if err := cli.NewJsonPrinter(log).Print(claims); err != nil {
		return err
	}

	if decodeTimes {
		printDecodedTimestamps(claims)
	}

	return nil
}

var decodeTokenCmd = &cobra.Command{
	Use:   "decode-token [TOKEN]",
	Short: "Decode a JWT token and display its claims",
	Long: `Decode a JWT token and display its claims as JSON.

If TOKEN is not provided, the access token from the active context is used.
Use --key to verify the token signature against a PEM-encoded public key.
Use --decode-times to print human-readable timestamps for iat, exp, nbf, and other date fields.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyErrorsWithRetry(DecodeTokenCmd(cfg, args), log, cfg)
	},
}

func init() {
	RootCmd.AddCommand(decodeTokenCmd)
	decodeTokenCmd.Flags().StringVarP(&signingKey, "key", "", "", "PEM-encoded public key or certificate for signature verification")
	decodeTokenCmd.Flags().BoolVarP(&decodeTimes, "decode-times", "d", false, "Print human-readable timestamps for date fields (iat, exp, nbf, etc.)")
	decodeTokenCmd.Annotations = make(map[string]string)
	decodeTokenCmd.Annotations[TOKEN_CATEGORY] = "true"
}
