package node

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// ErrAssociated is returned when manipulating a node which should *not*
	// have a parent associated.
	ErrAssociated = errors.New("node shouldnot have a parent associated")
	// ErrNotAssociated is returned when manipulating a node which should
	// have a parent associated.
	ErrNotAssociated = errors.New("node should have a parent associated")
	// ErrWrongNodeType is returned when manipulating a node which doesn't have
	// desired type
	ErrWrongNodeType = errors.New("called on wrong node type")
)

// PlainText returns the plaintext version of the dom tree starting from here
func (n *Node) PlainText() string {
	cpy := n.CloneTree()
	cpy.Walk(func(nn *Node) bool {
		switch nn.DataAtom {
		case atom.A, atom.Br:
			_ = nn.SetTail(" " + nn.Tail())
		}
		return true
	})
	buffer := bytes.Buffer{}
	cpy.textContent(&buffer)
	return buffer.String()
}

// Node is golang.org/x/net/html.Node with some other methods
type Node html.Node

// Visitor is function to visit and manipulate nodes
type Visitor func(*Node) bool

// NewNode returns new *Node with given parameter.
func NewNode(dataAtom atom.Atom, attributes ...string) *Node {
	result := &Node{
		DataAtom: dataAtom,
		Type:     html.ElementNode,
		Data:     dataAtom.String(),
	}

	if size := len(attributes); size > 0 && size&1 == 0 {
		result.Attr = make([]html.Attribute, size/2)
		for idx := 0; idx < size/2; idx++ {
			aidx := 2 * idx
			result.Attr[idx] = html.Attribute{
				Key: attributes[aidx],
				Val: attributes[aidx+1],
			}
		}
	}
	return result
}

// SetAtom sets a new DataAtom
func (n *Node) SetAtom(dataAtom atom.Atom) {
	n.DataAtom = dataAtom
}

// SetData sets node's data as given text
func (n *Node) SetData(text string) {
	n.Data = text
}

// DropTree TODO(zheng)
func (n *Node) DropTree() {
	if n.Parent != nil {
		n.Parent.RemoveChild((*html.Node)(n))
	}
}

// ParentNode returns it's parent node
func (n *Node) ParentNode() *Node {
	return (*Node)(n.Parent)
}

// RemoveChild removes given element node
func (n *Node) RemoveChild(child *Node) {
	(*html.Node)(n).RemoveChild((*html.Node)(child))
}

// RemoveAllChildren removes all of n's child nodes
func (n *Node) RemoveAllChildren() {
	for n.FirstChild != nil {
		n.RemoveChild((*Node)(n.FirstChild))
	}
}

// AppendChild appends given element node as its last child.
func (n *Node) AppendChild(element *Node) {
	(*html.Node)(n).AppendChild((*html.Node)(element))
}

// InsertBefore inserts newChild before oldChild
func (n *Node) InsertBefore(newChild *Node, oldChild *Node) {
	(*html.Node)(n).InsertBefore((*html.Node)(newChild), (*html.Node)(oldChild))
}

// InsertAfter inserts newChild as a child of n, immediately after oldChild
// in the sequence of n's children. oldChild may be nil, in which case newChild
// is appended to the end of n's children.
//
// return ErrAssociated if newChild already has a parent or siblings.
func (n *Node) InsertAfter(newChild, oldChild *Node) error {
	newC, oldC := (*html.Node)(newChild), (*html.Node)(oldChild)

	if newC.Parent != nil || newC.PrevSibling != nil || newC.NextSibling != nil {
		return ErrAssociated
	}

	var prev, next *html.Node
	if oldC != nil {
		prev, next = oldC, oldC.NextSibling
	} else {
		prev = n.LastChild
	}
	if prev != nil {
		prev.NextSibling = newC
	} else {
		n.FirstChild = newC
	}
	if next != nil {
		next.PrevSibling = newC
	} else {
		n.LastChild = newC
	}

	newC.Parent = (*html.Node)(n)
	newC.PrevSibling = prev
	newC.NextSibling = next
	return nil
}

// Text get the text content (Data field of Node) of a node, depending on the type of node, there're 3 cases:
// 1, TextNode: text content
// 2, ElementType: return text content of its 1st child, if its type is not TextNode, return empty string.
// 3, other node: empty string.
func (n *Node) Text() string {
	switch t := n.Type; t {
	case html.TextNode:
		return n.Data
	case html.ElementNode:
		if ch := n.FirstChild; ch != nil && ch.Type == html.TextNode {
			return ch.Data
		}
	default:
		return ""
	}

	return ""
}

// HasText check text depends on node's type.
func (n *Node) HasText() bool {
	if n.Type == html.TextNode {
		if len(n.Data) > 0 {
			return true
		}
	} else if ch := n.FirstChild; ch != nil && ch.Type == html.TextNode {
		if len(ch.Data) > 0 {
			return true
		}
	}
	return false
}

// SetText set text depends on node's type.
func (n *Node) SetText(text string) {
	if n.Type == html.TextNode {
		n.Data = text
	} else {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			n.FirstChild.Data = text
		} else {
			//if n.FirstChild == nil {
			//n.AppendChild(NewTextNode(text))
			//} else {
			//n.InsertBefore(NewTextNode(text), (*Node)(n.FirstChild))
			//}
			// see docuement of html.insertBefore, there's no need to use if/else
			n.InsertBefore(NewTextNode(text), (*Node)(n.FirstChild))
		}
	}
}

// SetTextContent set text depends on node's type.
func (n *Node) SetTextContent(text string) {
	// TODO(yangbo)
	if n.Type == html.TextNode {
		n.Data = text
	} else {
		remove := make([]*html.Node, 0)
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				remove = append(remove, child)
			}
		}
		for _, r := range remove {
			n.RemoveChild((*Node)(r))
		}
		n.InsertBefore(NewTextNode(text), (*Node)(n.FirstChild))
	}
}

// TextContent returns the text content of the node, including the text content
// of its children, with no markup.
func (n *Node) TextContent() string {
	var contentBuffer bytes.Buffer
	n.textContent(&contentBuffer)

	return contentBuffer.String()
}

func (n *Node) textContent(buffer *bytes.Buffer) {
	if n.Type == html.TextNode {
		buffer.WriteString(n.Data)
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		(*Node)(child).textContent(buffer)
	}
}

// Walk traverse dom tree use n as its root node, and invoke visitor function
// on every child node.
func (n *Node) Walk(visitor Visitor) bool {

	children := []*Node{}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		children = append(children, (*Node)(child))
	}

	if !visitor(n) {
		return false
	}

	for _, current := range children {
		// skip deleted node
		if current.Parent == nil {
			continue
		}

		if !current.Walk(visitor) {
			return false
		}
	}
	return true
}

// Tail get the tail text of a node. Concepts stolen from lxml
func (n *Node) Tail() string {
	if n.NextSibling == nil || (*Node)(n.NextSibling).Type != html.TextNode {
		return ""
	}
	return (*Node)(n.NextSibling).Text()
}

// SetTail require the `n` has a parent associated,
// node calling SetTail without parent will retrun ErrNotAssociated
// textnode calling SetTail will return ErrWrongNodeType
func (n *Node) SetTail(text string) error {
	if sibling := n.NextSibling; sibling != nil && sibling.Type == html.TextNode {
		(*Node)(sibling).SetText(text)
	} else if n.Type != html.TextNode {
		parentNode := (*Node)(n.Parent)
		if parentNode != nil {
			textNode := NewTextNode(text)
			_ = parentNode.InsertAfter(textNode, n)
		} else {
			return ErrNotAssociated
		}
	} else {
		return ErrWrongNodeType
	}
	return nil
}

// DelAttr deletes attribute with given key
func (n *Node) DelAttr(key string) {
	if len(n.Attr) == 0 {
		return
	}
	attrs := make([]html.Attribute, 0, len(n.Attr))
	deleted := false
	for _, attr := range n.Attr {
		if attr.Key != key {
			attrs = append(attrs, attr)
		} else {
			deleted = true
		}
	}
	if deleted {
		n.Attr = attrs
	}
}

// DelAttrs deletes attributes with given keys, in one traverse.
func (n *Node) DelAttrs(keys ...string) {
	if len(n.Attr) == 0 {
		return
	}
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}

	attrs := make([]html.Attribute, 0, len(n.Attr))
	deleted := false
	for _, attr := range n.Attr {
		if _, ok := m[attr.Key]; !ok {
			attrs = append(attrs, attr)
		} else {
			deleted = true
		}
	}
	if deleted {
		n.Attr = attrs
	}
}

// GetAttr returns atrribute value and a flag by given key. If not found,
// return empty value and false flag.
func (n *Node) GetAttr(key string) (string, bool) {
	for _, attribute := range n.Attr {
		if attribute.Key == key {
			return attribute.Val, true
		}
	}
	return "", false
}

// GetAttrOrDefault returns atrribute value by given key. If not found,
// return defaultValue
func (n *Node) GetAttrOrDefault(name, defaultValue string) string {
	value, ok := n.GetAttr(name)
	if ok {
		return value
	}
	return defaultValue
}

// SetAttr sets attribute with key and value. If attribute with key name
// already exists, update it with new value, otherwise, append new attribute.
func (n *Node) SetAttr(key string, value string) {
	for index, attr := range n.Attr {
		if attr.Key == key {
			n.Attr[index].Val = value
			return
		}
	}
	// no attr found, so we add a new one
	n.Attr = append(n.Attr, html.Attribute{
		Key: key,
		Val: value,
	})
}

// SetAttrs sets some attributes with key and value. If attribute with key name
// already exists, update it with new value, otherwise, append new attribute.
func (n *Node) SetAttrs(attrs map[string]string) {
	for index, attr := range n.Attr {
		if attrVal, ok := attrs[attr.Key]; ok {
			n.Attr[index].Val = attrVal
			delete(attrs, attr.Key)
		}
	}

	for attrKey, attrVal := range attrs {
		n.Attr = append(n.Attr, html.Attribute{
			Key: attrKey,
			Val: attrVal,
		})
	}
}

// AddAttr appends new attribute.
func (n *Node) AddAttr(key string, value string) {
	n.Attr = append(n.Attr, html.Attribute{
		Key: key,
		Val: value,
	})
}

// DropTag remove node itself, and add all its children to its parent. If it
// doesn't parent, just return with no action(not remove itself).
// Reimplementation of drop_tag in lxml.
func (n *Node) DropTag() {
	parent := n.Parent
	if parent == nil {
		return
	}
	node := (*html.Node)(n)
	var sibling *html.Node
	for current := node.FirstChild; current != nil; current = sibling {
		sibling = current.NextSibling
		node.RemoveChild(current)

		parentPrevSibling := (*Node)(node.PrevSibling)
		if parentPrevSibling != nil && parentPrevSibling.Type == html.TextNode && current.Type == html.TextNode {
			//merge text node
			parentPrevSibling.Data += current.Data
		} else {
			parent.InsertBefore(current, node)
		}
	}
	prev := (*Node)(node.PrevSibling)
	next := (*Node)(node.NextSibling)
	parent.RemoveChild(node)

	// also merge text node with tail, if tail exists
	if prev != nil && next != nil && prev.Type == html.TextNode && next.Type == html.TextNode {
		prev.Data += next.Data
		parent.RemoveChild((*html.Node)(next))
	}
}

// FindTag return the node that match given tag
func (n *Node) FindTag(tag string) *Node {
	if n == nil || n.DataAtom.String() == strings.ToLower(tag) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom.String() == strings.ToLower(tag) {
			return (*Node)(c)
		}
	}
	return nil
}

// Atom returns node's DataAtom
func (n *Node) Atom() atom.Atom {
	return n.DataAtom
}

// Render dom tree to html string.
func (n *Node) Render() (string, error) {
	var result string
	var err error
	var buffer bytes.Buffer

	if err = html.Render(&buffer, (*html.Node)(n)); err != nil {
		return result, err
	}
	result = strings.TrimSpace(buffer.String())

	// remove the outer <div></div> added by NewFragmentFromString
	result = result[5 : len(result)-6]
	return result, err
}

func (n *Node) cloneTree() *html.Node {
	nn := &html.Node{
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]html.Attribute, len(n.Attr)),
	}

	copy(nn.Attr, n.Attr)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nn.AppendChild((*Node)(c).cloneTree())
	}

	return nn
}

// CloneTree deep copy a node. The new node has clones of all the original node's
// children but none of its parents or siblings.
func (n *Node) CloneTree() *Node {
	return (*Node)(n.cloneTree())
}

func printTree(element *Node, depth int) {
	if element.Type != html.TextNode {
		fmt.Println(strings.Repeat("  ", depth), element.DataAtom, "=>", element.Text(), "+", element.Tail())

		for child := element.FirstChild; child != nil; child = child.NextSibling {
			printTree((*Node)(child), depth+1)
		}
	} else {
		fmt.Println(strings.Repeat("  ", depth), fmt.Sprintf("%+v", element))
	}
}

// Print dom tree, for debug purpose
func (n *Node) Print() {
	printTree(n, 0)
}

// Root returns self
func (n *Node) Root() *Node {
	return n
}

// NewTextNode return a new TextNode with given text
func NewTextNode(text string) *Node {
	return &Node{
		DataAtom: atom.Atom(0),
		Type:     html.TextNode,
		Data:     text,
	}
}

// NewNodeWithText returns a ElementNode with given text, which means with a
// TextNode as its child.
func NewNodeWithText(dataAtom atom.Atom, text string, attributes ...string) *Node {
	n := NewNode(dataAtom, attributes...)
	textNode := NewTextNode(text)
	n.InsertBefore(textNode, (*Node)(n.FirstChild))
	return n
}

// NewDiv returns a Div ElementNode
func NewDiv(attributes ...string) *Node {
	return NewNode(atom.Div, attributes...)
}

// NewP returns a P ElementNode
func NewP(attributes ...string) *Node {
	return NewNode(atom.P, attributes...)
}

// NewA returns a P ElementNode
func NewA(attributes ...string) *Node {
	return NewNode(atom.A, attributes...)
}

// NewNoscript returns a noscript ElementNode
func NewNoscript(attributes ...string) *Node {
	return NewNode(atom.Noscript, attributes...)
}
