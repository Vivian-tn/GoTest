package node

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// NewFragmentsFromString parses several HTML elements, returning a list of elements.
func NewFragmentsFromString(content string) ([]*Node, error) {
	var err error
	var fragments []*Node

	if len(content) == 0 {
		return fragments, err
	}

	var nodes []*html.Node
	if nodes, err = html.ParseFragment(strings.NewReader(content), (*html.Node)(NewNode(atom.Div))); err != nil {
		return fragments, err
	}

	fragments = make([]*Node, len(nodes))
	for idx, node := range nodes {
		fragments[idx] = (*Node)(node)
	}
	return fragments, nil
}

// NewFragmentFromString parses a single HTML element
func NewFragmentFromString(content string) (*Node, error) {
	fragment := NewNode(atom.Div)

	nodes, err := NewFragmentsFromString(content)
	if err != nil {
		return fragment, err
	}

	for _, node := range nodes {
		fragment.AppendChild(node)
	}
	return fragment, nil
}
