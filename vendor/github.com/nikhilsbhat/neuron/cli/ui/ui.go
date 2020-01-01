// Package ui helps one to make their output look colorful, import this package and start using the methods this package implements to make your output colorful
package ui

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// UiWriter holds the writer for the cli ui.
type UiWriter struct {
	Writer io.Writer
}

// Color will just a wrap of the struct Attribute of package color.
type Color struct {
	colour color.Attribute
}

// NeuronUi is just a wrap of the struct UiWriter of the same package.
// For the people who doubts why this wraps the struct of the same package. We are plannig to add few more components under NeuronUi.
type NeuronUi struct {
	*UiWriter
}

var (
	// Red Blue Green Cyan Yellow, Magenta are the variables which holds the respective colors.
	red     = Color{color.FgRed}
	blue    = Color{color.FgBlue}
	green   = Color{color.FgGreen}
	cyan    = Color{color.FgCyan}
	yellow  = Color{color.FgYellow}
	magenta = Color{color.FgMagenta}
)

func (c *Color) makeitcolorful(msg string) string {
	return color.New(c.colour).SprintFunc()(msg)
}

// NeuronSaysItsInfo will be defining the color while printing informations.
func (n NeuronUi) NeuronSaysItsInfo(msg string) {
	n.neuronPrints(green.makeitcolorful(msg))
}

// NeuronSaysItsError will be defining the color while printing errors.
func (n NeuronUi) NeuronSaysItsError(msg string) {
	n.neuronPrints(red.makeitcolorful("Error: " + msg))
}

// NeuronSaysItsWarn will be defining the color while printing warnings.
func (n NeuronUi) NeuronSaysItsWarn(msg string) {
	n.neuronPrints(yellow.makeitcolorful(msg))
}

// NeuronSaysItsDebug will be defining the color while printing warnings.
func (n NeuronUi) NeuronSaysItsDebug(msg string) {
	n.neuronPrints(magenta.makeitcolorful(msg))
}

// Info will be just defining the color for informations.
func Info(msg string) string {
	return green.makeitcolorful(msg)
}

// Error will be just defining the color for errors.
func Error(msg string) string {
	return red.makeitcolorful(msg)
}

// Warn will be just defining the color for warnings.
func Warn(msg string) string {
	return yellow.makeitcolorful(msg)
}

// Debug will be just defining the color for debugging.
func Debug(msg string) string {
	return magenta.makeitcolorful(msg)
}

func (n NeuronUi) neuronPrints(msg string) {
	fmt.Fprint(n.Writer, msg)
	fmt.Fprintf(n.Writer, "\n")
}
