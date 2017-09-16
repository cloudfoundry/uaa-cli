package utils

import "strings"

func StringSliceStringifier(stringsList []string) string {
	return "[" + strings.Join(stringsList, ", ") + "]"
}
