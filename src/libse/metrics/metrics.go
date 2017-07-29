package Metrics

import (
	"strings"
	"golang.org/x/net/html"
)

func isMixedNode(node *html.Node, activeOnly bool) bool {
	// activeOnly
	if node.Type == html.ElementNode && activeOnly == true {
		switch node.Data {
		case "script", "iframe":
			for _, attr := range node.Attr {
				if key, val := attr.Key, attr.Val; key == "src" {
					return strings.HasPrefix(val, "http:")
				}
			}

		case "link":
			for _, attr := range node.Attr {
				if key, val := attr.Key, attr.Val; key == "href" {
					return strings.HasPrefix(val, "http:")
				}
			}

		case "object":
			for _, attr := range node.Attr {
				if key, val := attr.Key, attr.Val; key == "data" {
					return strings.HasPrefix(val, "http:")
				}
			}
		}
	}

	// !activeOnly
	if node.Type == html.ElementNode && activeOnly == false {
		switch node.Data {
		case "img", "audio", "video":
			for _, attr := range node.Attr {
				if key, val := attr.Key, attr.Val; key == "src" {
					return strings.HasPrefix(val, "http:")
				}
			}

		// TODO
		// case "object":
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if isMixedNode(c, activeOnly) {
			return true
		}
	}
	return false
}

func HasMixedContent(ctx string, activeOnly bool) (bool, error) {
	// Ref: https://developer.mozilla.org/en-US/docs/Web/Security/Mixed_content
	doc, err := html.Parse(strings.NewReader(ctx))
	if err != nil {
		return false, err
	}

	return isMixedNode(doc, activeOnly), nil
}
