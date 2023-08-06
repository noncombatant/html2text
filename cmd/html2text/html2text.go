/*
html2text converts HTML files to plain text.

Without an explicit pathname, it converts the standard input. Given the name of
1 or more files, it converts those files. It always writes to the standard
output.

Usage:

    html2text [pathname ...]

html2text uses ANSI colors to distinguish some HTML tags. If you don’t want
color, set the `NO_COLOR` environment variable to some non-empty value. For
example, if you’re using a POSIX shell:

    NO_COLOR=true html2text [pathname ...]

When not using color, html2text displays some HTML tags with markup that looks
not entirely unlike Markdown.
*/
package main

import (
  "flag"
  "fmt"
  "os"

	"golang.org/x/net/html"

  "html2text"
)

func main() {
	flag.Parse()
	arguments := flag.Args()
	if len(arguments) == 0 {
		if document, e := html.Parse(os.Stdin); e == nil {
			html2text.Render(os.Stdout, document, document)
		} else {
			fmt.Fprintln(os.Stderr, e)
		}
	}
	for _, pathname := range flag.Args() {
		file, e := os.Open(pathname)
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			continue
		}
		if document, e := html.Parse(file); e == nil {
			html2text.Render(os.Stdout, document, document)
		} else {
			fmt.Fprintln(os.Stderr, e)
		}
		file.Close()
	}
}
