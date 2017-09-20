package utils

import "strings"

const (
	commaSpace = ", "
	commaOnly  = ","
	spaceOnly  = " "
)

func Arrayify(input string) []string {
	input = strings.TrimSpace(input)

	if input == "" {
		return []string{}
	} else if strings.Contains(input, commaSpace) {
		return removeEmpty(strings.Split(input, commaSpace))
	} else if strings.Contains(input, commaOnly) {
		return removeEmpty(strings.Split(input, commaOnly))
	} else if strings.Contains(input, spaceOnly) {
		return removeEmpty(strings.Split(input, spaceOnly))
	}
	return []string{input}
}

func removeEmpty(input []string) []string {
	output := []string{}
	for _, item := range input {
		if item != "" {
			output = append(output, strings.TrimSpace(item))
		}
	}
	return output
}
