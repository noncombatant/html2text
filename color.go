package html2text

// I adapted this code from Fatih Arlsan’s color package
// (https://github.com/fatih/color), which he released under the MIT license.
// Thank you, Fatih!

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	NoColor = os.Getenv("NO_COLOR") != ""
)

func SetColor(a ...Attribute) string {
	if NoColor {
		return ""
	}
	format := make([]string, len(a))
	for i, v := range a {
		format[i] = strconv.Itoa(int(v))
	}
	sequence := strings.Join(format, ";")
	return fmt.Sprintf("%s[%sm", escape, sequence)
}

func UnsetColor() string {
	if NoColor {
		return ""
	}
	return fmt.Sprintf("%s[%dm", escape, Reset)
}

const escape = "\x1b"

type Attribute int

const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)
