package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var InteractiveOutput io.Writer = os.Stdout
var InteractiveInput io.Reader = os.Stdin
var ReadPassword func(fd int) ([]byte, error) = terminal.ReadPassword

type InteractiveSecret struct {
	Prompt string
}

func (is *InteractiveSecret) Get() (string, error) {
	fmt.Fprint(InteractiveOutput, is.Prompt+": ")

	bytePassword, err := ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(bytePassword), nil
}

type InteractivePrompt struct {
	Prompt string
}

func (ip InteractivePrompt) Get() (string, error) {
	prompt := color.CyanString(ip.Prompt + ": ")
	fmt.Fprint(InteractiveOutput, prompt)

	reader := bufio.NewReader(InteractiveInput)
	val, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(val), nil
}
