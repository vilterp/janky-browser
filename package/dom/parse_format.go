package dom

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
)

func Format(node Node) string {
	return doFormat(node, 0)
}

func doFormat(node Node, indent int) string {
	attrsStr := formatAttrs(node)
	indentStr := strings.Repeat("  ", indent)
	children := node.Children()
	if len(children) > 0 {
		var childrenLines []string
		for _, child := range children {
			childrenLines = append(childrenLines, indentStr+doFormat(child, indent+1))
		}
		childrenStr := strings.Join(childrenLines, "\n")
		return fmt.Sprintf(
			"%s<%s%s>\n%s\n%s</%s>",
			indentStr, node.Name(), attrsStr, childrenStr, indentStr, node.Name(),
		)
	}
	return fmt.Sprintf("%s<%s%s />", indentStr, node.Name(), attrsStr)
}

func FormatWithoutChildren(node Node) string {
	// Format node
	if len(node.Children()) == 0 {
		return fmt.Sprintf("<%s%s />", node.Name(), formatAttrs(node))
	}
	return fmt.Sprintf("<%s%s>", node.Name(), formatAttrs(node))
}

func formatAttrs(node Node) string {
	var attrs []string
	for key, val := range node.Attrs() {
		attrs = append(attrs, fmt.Sprintf("%s=%#v", key, val))
	}
	sort.Strings(attrs)
	attrsStr := strings.Join(attrs, " ")
	if len(attrs) > 0 {
		attrsStr = " " + attrsStr
	}
	return attrsStr
}

func Parse(data []byte) (Node, error) {
	g := GroupNode{}
	err := xml.Unmarshal(data, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
