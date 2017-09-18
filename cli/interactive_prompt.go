package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var InteractiveOutput io.Writer = os.Stdout
var InteractiveInput io.Reader = os.Stdin

type InteractivePrompt struct {
	Input  io.Reader
	Prompt string
}

func (ip InteractivePrompt) Get() (string, error) {
	fmt.Fprint(InteractiveOutput, ip.Prompt+": ")

	reader := bufio.NewReader(InteractiveInput)
	val, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(val), nil
}
