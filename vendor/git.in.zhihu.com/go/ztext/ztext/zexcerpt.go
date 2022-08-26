package ztext

import (
	"regexp"
	"strings"
	"unicode"

	"git.in.zhihu.com/go/ztext/node"
	"git.in.zhihu.com/go/ztext/util"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// ExcerptAllowedTags contains tags allowed in ZExcerptText by default
	ExcerptAllowedTags = map[atom.Atom]struct{}{
		atom.B:      {},
		atom.I:      {},
		atom.U:      {},
		atom.Em:     {},
		atom.Strong: {},
		atom.A:      {},
	}
	spacesPattern = regexp.MustCompile(`(?m)([\n \x{a0}\x{2000}-\x{200b}\x{202f}\x{205f}\x{3000}]|\\n)+`)
)

// ExcerptOption for ZExcerptText config
type ExcerptOption func(*ZExcerptText)

// AllowedTags sets allowed html tags
func AllowedTags(tags map[atom.Atom]struct{}) ExcerptOption {
	return ExcerptOption(func(z *ZExcerptText) {
		z.allowedTags = tags
	})
}

// Trail sets trail char
func Trail(trail string) ExcerptOption {
	return ExcerptOption(func(z *ZExcerptText) {
		z.trail = trail
	})
}

// ExcerptLength sets excerpt length
func ExcerptLength(n int) ExcerptOption {
	return ExcerptOption(func(z *ZExcerptText) {
		z.length = n
	})
}

// ZOutputTextOptions set ZOutputText options in ZExcerptText
func ZOutputTextOptions(opts []Option) ExcerptOption {
	return ExcerptOption(func(z *ZExcerptText) {
		z.zOutputOptions = opts
	})
}

// ZExcerptText ...
type ZExcerptText struct {
	*ZOutputText

	// allowed tags, default is [b, i ,u, em, strong, a]
	allowedTags map[atom.Atom]struct{}

	// trailing char, default is "…"
	trail string

	// excerpt length, default is 40
	length int

	// if has been shortened
	Shortened bool

	// set ZOutputText options
	zOutputOptions []Option
}

// NewZExcerptText returns a instanse with given ExcerptOption, if no option given,
// return with default option, which is listed in the comments of ZExcerptText definition.
func NewZExcerptText(content string, opts ...ExcerptOption) (*ZExcerptText, error) {
	z := &ZExcerptText{nil, ExcerptAllowedTags, "…", 40, false, nil}
	for _, opt := range opts {
		opt(z)
	}
	o, err := NewZOutputText(content, z.zOutputOptions...)
	if err != nil {
		return nil, err
	}
	z.ZOutputText = o
	return z, nil
}

// PrepareWalkDOM 遍历 DOM 树，按照给定长度截取摘要
func (z *ZExcerptText) PrepareWalkDOM() error {
	err := z.ZOutputText.PrepareWalkDOM()
	if err != nil {
		return err
	}

	var hasMultimedia bool
	z.Walk(func(n *node.Node) bool {
		if n == z.Root() {
			return true
		}

		tag := n.DataAtom
		if tag == atom.Img || tag == atom.Video {
			hasMultimedia = true
		}

		//drop all tags not in allowedTags
		if n.Type != html.ElementNode {
			return true
		}
		if tag == atom.A && n.GetAttrOrDefault("data-draft-type", "") == "mcn-link-card" {
			n.DropTag()
		}
		if tag == atom.A && n.GetAttrOrDefault("data-draft-type", "") == "km-sku-card" {
			n.DropTag()
		}
		if _, ok := z.allowedTags[tag]; ok {
			return true
		}
		// tag not in allowedTags
		if tag == atom.Br {
			sib := (*node.Node)(n.PrevSibling)
			if sib != nil && sib.Type == html.TextNode {
				// if prev node is textnode, we just combine text
				sib.SetText(sib.Text() + " ")
			} else {
				// or else, create a new text node with " "
				n.Parent.InsertBefore((*html.Node)(node.NewTextNode(" ")), (*html.Node)(n))
			}
		}
		if tag == atom.P {
			if n.Text() != "" {
				n.SetText(n.Text() + " ")
			} else {
				n.SetText(" ")
			}
		}
		n.DropTag()

		return true
	})

	if t := z.Text(); t != "" {
		// strip left most spaces

		// strings.TrimLeftFunc(_, unicode.IsSpace) is equivalent
		// of .lstrip() in python
		t = strings.TrimLeftFunc(t, unicode.IsSpace)
		z.SetText(t)

		if t != "" {
			// convert container's text to '__text__'
			div := node.NewDiv()
			div.SetData("__text__")
			div.SetText(t)

			z.InsertBefore(div, (*node.Node)(z.FirstChild))
			// remove original text node
			z.RemoveChild((*node.Node)(z.FirstChild.NextSibling))
		}

	}

	z.Walk(func(n *node.Node) bool {
		// convert element tail to '__text__' elem
		// and sqeeze multiple spaces into one
		// for convenience
		if n.Type != html.TextNode {
			return true
		}
		t := n.Text()
		if t != "" {
			n.SetText(spacesPattern.ReplaceAllLiteralString(t, " "))
		}

		// TODO (ybb) Handle tail?

		return true
	})

	preserve := false
	shortened := false
	currentLength := 0
	excerptLength := z.length * 2
	var lastNode *node.Node
	z.Walk(func(n *node.Node) bool {

		if n.Type == html.TextNode {
			p := n.Parent
			if p != nil && p.FirstChild == (*html.Node)(n) {
				return true
			}
		}
		// count link as 10 chars
		lastNode = n
		text := n.Text()
		// plain hyperlink, will be abbreviated by ZOutputText
		if n.Atom() == atom.A && n.GetAttrOrDefault("href", "") == text {
			if excerptLength-currentLength < 20 {
				// exceeds limit
				preserve = false
				_ = n.SetTail(z.trail)
				n.SetText("")
				shortened = true
				lastNode = (*node.Node)(n.NextSibling) // point lastNode to tail node "…"
				return false
			}

			currentLength += 20
		} else if text != "" {
			space := excerptLength - currentLength

			_len := util.GB18030Len(text)
			if _len > space {
				// exceeds limit
				preserve = true
				n.SetText(util.GB18030Cut(text, space-2) + z.trail)
				shortened = true
				return false
			}

			currentLength += _len
		}

		return true
	})

	z.Shortened = shortened || hasMultimedia

	// remove all elems behind pinned elem
	pinned := lastNode
	if shortened {
		for nn := lastNode; nn != z.Root(); {
			p := nn.Parent
			if p == nil {
				break
			}

			for sib := (*html.Node)(nn.NextSibling); sib != nil; {
				nextSib := sib.NextSibling
				p.RemoveChild(sib)
				sib = nextSib
			}

			nn = (*node.Node)(p)
		}

		if !preserve && pinned.Parent != nil {
			// when "<a></a>…", lastNode is "…", but we should dropTag on <a>
			(*node.Node)(pinned.PrevSibling).DropTag()
		}
	}

	// post process
	// replace <br> tag with space
	// some spacing adjustment
	// drop '__text__'
	z.Walk(func(n *node.Node) bool {
		if n.Data == "__text__" {
			n.DropTag()
		}
		return true
	})
	return nil
}
