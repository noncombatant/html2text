// Package html2text converts HTML to plain text.
package html2text

import (
	"fmt"
	"io"
	"regexp"

	"golang.org/x/net/html"
)

type set map[string]bool

func (s set) contains(x string) bool {
	_, ok := s[x]
	return ok
}

var (
	spaces               = regexp.MustCompile(`\s+`)
	blockDisplayElements = set{
		"h1":         true,
		"h2":         true,
		"h3":         true,
		"h4":         true,
		"h5":         true,
		"header":     true,
		"hr":         true,
		"p":          true,
		"pre":        true,
		"figcaption": true,
		"footer":     true,
		"nav":        true,
		"title":      true,
	}
	ignoreElements = set{
		"meta":   true,
		"script": true,
		"style":  true,
		"title":  true,
	}
)

func isElement(n *html.Node, tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func isElementInSet(n *html.Node, s set) bool {
	return n.Type == html.ElementNode && s.contains(n.Data)
}

func isBlockDisplayElement(n *html.Node) bool {
	return isElementInSet(n, blockDisplayElements)
}

func shouldIgnoreElement(n *html.Node) bool {
	return isElementInSet(n, ignoreElements)
}

func hasParent(node *html.Node, tag string) bool {
	for p := node.Parent; p != nil; p = p.Parent {
		if p.Type == html.ElementNode && p.Data == tag {
			return true
		}
	}
	return false
}

func getAttribute(as []html.Attribute, key string) string {
	for _, a := range as {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// Render walks down the HTML tree starting at n and writes a plain text
// rendition of it to w.
//
// When passing a root node, pass that node as the parent, as well.
//
// See also the Parse function in the [golang.org/x/net/html] package.
func Render(w io.Writer, n, parent *html.Node) {
	if shouldIgnoreElement(n) {
		return
	}
	if isBlockDisplayElement(n) {
		fmt.Fprint(w, "\n")
	}

	if isElement(n, "img") {
		alt := getAttribute(n.Attr, "alt")
		if alt != "" {
			fmt.Fprintf(w, "(image: %s) ", alt)
		} else {
			fmt.Fprint(w, "(image) ")
		}
	} else if isElement(n, "h1") {
		fmt.Fprint(w, "# ")
	} else if isElement(n, "h2") {
		fmt.Fprint(w, "## ")
	} else if isElement(n, "h3") {
		fmt.Fprint(w, "### ")
	} else if isElement(n, "h4") {
		fmt.Fprint(w, "#### ")
	} else if isElement(n, "h5") {
		fmt.Fprint(w, "##### ")
	} else if isElement(n, "br") {
		fmt.Fprintln(w)
	} else if isElement(n, "hr") {
		fmt.Fprintln(w, "------")
	} else if isElement(n, "figcaption") {
		fmt.Fprint(w, "[")
	} else if isElement(n, "i") || isElement(n, "cite") {
		fmt.Fprint(w, "_")
	} else if isElement(n, "b") || isElement(n, "em") {
		fmt.Fprint(w, "*")
	} else if isElement(n, "code") || isElement(n, "tt") {
		fmt.Fprint(w, "`")
	} else if isElement(n, "pre") {
		fmt.Fprint(w, "```\n")
	}

	if n.Type == html.TextNode && parent.Type == html.ElementNode && parent.Data != "html" && parent.Data != "body" {
		data := n.Data
		if !hasParent(n, "pre") {
			data = spaces.ReplaceAllString(data, " ")
		}
		fmt.Fprint(w, data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Render(w, c, n)
	}

	if isElement(n, "a") {
		if href := getAttribute(n.Attr, "href"); href != "" {
			if rel := getAttribute(n.Attr, "rel"); rel == "" {
				fmt.Fprintf(w, " (%s)", href)
			}
		}
	} else if isElement(n, "figcaption") {
		fmt.Fprint(w, "]")
	} else if isElement(n, "i") || isElement(n, "cite") {
		fmt.Fprint(w, "_")
	} else if isElement(n, "b") || isElement(n, "em") {
		fmt.Fprint(w, "*")
	} else if isElement(n, "code") || isElement(n, "tt") {
		fmt.Fprint(w, "`")
	} else if isElement(n, "pre") {
		fmt.Fprint(w, "\n```")
	}
	if isBlockDisplayElement(n) {
		fmt.Fprint(w, "\n")
	}
}
