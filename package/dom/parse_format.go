package dom

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
)

func Format(node DOMNode) string {
	return doFormat(node, 1)
}

func doFormat(node DOMNode, indent int) string {
	var attrs []string
	for key, val := range node.Attrs() {
		attrs = append(attrs, fmt.Sprintf("%s=%#v", key, val))
	}
	sort.Strings(attrs)
	attrsStr := strings.Join(attrs, " ")
	if len(attrs) > 0 {
		attrsStr = " " + attrsStr
	}

	children := node.Children()
	if len(children) > 0 {
		indentStr := strings.Repeat("  ", indent)
		var childrenLines []string
		for _, child := range children {
			childrenLines = append(childrenLines, indentStr+doFormat(child, indent+1))
		}
		childrenStr := strings.Join(childrenLines, "\n")
		return fmt.Sprintf("<%s%s>\n%s\n</%s>", node.Name(), attrsStr, childrenStr, node.Name())
	}
	return fmt.Sprintf("<%s%s />", node.Name(), attrsStr)
}

func Parse(data []byte) (DOMNode, error) {
	g := TextNode{}
	err := xml.Unmarshal(data, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
