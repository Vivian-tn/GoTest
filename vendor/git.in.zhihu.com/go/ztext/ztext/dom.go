package ztext

import (
	"strconv"
	"strings"

	"git.in.zhihu.com/go/ztext/escape"
	"git.in.zhihu.com/go/ztext/node"
	"git.in.zhihu.com/go/ztext/util"
	"golang.org/x/net/html/atom"
)

// DomOption is option type for ZDom
type DomOption func(*ZDom)

// DomStrictBMP sets strictBMP
func DomStrictBMP(strict bool) DomOption {
	return DomOption(func(c *ZDom) {
		c.StrictBMP = strict
	})
}

// ZDom implement Python ZDom class, and implement golang's ZText interface
type ZDom struct {
	*node.Node
	walked    bool
	StrictBMP bool
}

var _ ZText = (*ZDom)(nil)

// NewZDom returns ZDom instance
func NewZDom(node *node.Node, opts ...DomOption) *ZDom {
	z := &ZDom{node, false, true}
	for _, opt := range opts {
		opt(z)
	}
	return z
}

// Render output z.Root dom tree to rendered html
func (z *ZDom) Render() (string, error) {
	return z.Node.Render()
}

// Walked returns whether dom tree has been walked
func (z *ZDom) Walked() bool {
	return z.walked
}

// PlainText returns plain text
func (z *ZDom) PlainText() string {
	plain := z.Node.PlainText()
	if z.StrictBMP {
		plain = escape.CollapseSpacesSafeTextBMP(plain)
	} else {
		plain = escape.CollapseSpacesSafeText(plain)
	}
	return plain
}

// OnEquationTag is the callback function when visiting <equation> tag
func (z *ZDom) OnEquationTag(n *node.Node) {}

// OnATag is the callback function when visiting <a> tag
func (z *ZDom) OnATag(n *node.Node) {}

// OnPTag is the callback function when visiting <p> tag
func (z *ZDom) OnPTag(n *node.Node) {}

// OnCodeTag is the callback function when visiting <code> tag
func (z *ZDom) OnCodeTag(n *node.Node) {}

// OnVideoTag is the callback function when visiting <video> tag
func (z *ZDom) OnVideoTag(n *node.Node) {}

// OnImgTag is the callback function when visiting <img> tag
func (z *ZDom) OnImgTag(n *node.Node) {}

// OnUlTag is the callback function when visiting <ul> tag
func (z *ZDom) OnUlTag(n *node.Node) {}

// WalkDOM traverse dom tree with node.Vistor
func (z *ZDom) WalkDOM(n node.Visitor) error {
	z.Walk(n)
	z.walked = true
	return nil
}

// PrepareWalkDOM traverse z.Root dom tree, and do some preparation work.
func (z *ZDom) PrepareWalkDOM() error {
	return nil
}

// Thumbnail returns thumbnail without filter
func (z *ZDom) Thumbnail(minWidth, minHeight int) *node.Node {
	return z.ThumbnailWithPredicate(minWidth, minHeight, func(w, h int) bool {
		return true
	})
}

// ThumbnailWithPredicate returns thumbnail of the content
func (z *ZDom) ThumbnailWithPredicate(minWidth, minHeight int, predicate func(int, int) bool) (result *node.Node) {
	for c := z.FirstChild; c != nil; c = c.NextSibling {
		n := (*node.Node)(c)
		if n.Atom() != atom.Img {
			continue
		}
		width := -1
		height := -1
		for _, attribute := range n.Attr {
			if width == -1 && attribute.Key == "data-rawwidth" {
				width, _ = strconv.Atoi(strings.TrimSpace(attribute.Val))
				if height != -1 {
					break
				}
			} else if height == -1 && attribute.Key == "data-rawheight" {
				height, _ = strconv.Atoi(strings.TrimSpace(attribute.Val))
				if width != -1 {
					break
				}
			}
		}
		if width >= minWidth && height >= minHeight && predicate(width, height) {
			imageHash, ok := n.GetAttr("src")
			if !ok {
				continue
			}
			if strings.HasSuffix(imageHash, ".gif") {
				imageHash = util.URL2Token(imageHash, false)
				n.SetAttr("src", imageHash)
			}

			n.SetText("")
			if n.FirstChild != nil {
				n.RemoveChild((*node.Node)(n.FirstChild))
			}

			return n
		}
		continue
	}
	return nil
}

// ZText2 对应于 Python 的 Ztext2 类, 没有保留类的形式。如果你需要继承 Python
// 的 ZText2, 请组合 Golang 版本的 ZDom，代码范例请参考 ZOutputText。
// strictBMP 默认值为 true
func ZText2(content string, opts ...DomOption) *ZDom {
	zdom := NewZDom(nil, opts...)

	content = ZSafeText(content, zdom.StrictBMP)
	fragment, err := node.NewFragmentFromString(content)
	if err != nil {
		return nil
	}
	zdom.Node = fragment
	return zdom
}

// ZSafeText 对应于 Python 的 ZSafeText, 过滤掉非法字符，例如一些零宽字符，
// `strictBMP` 决定是否使用严格的字符面，影响是否能使用 emoji
func ZSafeText(content string, strictBMP bool) string {
	if strictBMP {
		return escape.SafeTextBMP(content)
	}
	return escape.SafeText(content)
}
