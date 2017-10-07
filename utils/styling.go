package utils

import "github.com/fatih/color"

var Emphasize = color.New(color.FgCyan, color.Bold).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
