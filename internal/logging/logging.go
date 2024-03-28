package logging

import "github.com/fatih/color"

var Error = color.New(color.FgRed).PrintfFunc()
var Success = color.New(color.FgGreen).PrintfFunc()
var Info = color.New(color.FgCyan).PrintfFunc()
var Warning = color.New(color.FgYellow).PrintfFunc()
var Prompt = color.New(color.FgWhite).Add(color.Bold).PrintfFunc()
