package ztext

import (
	"fmt"
	"strconv"

	"git.in.zhihu.com/go/ztext/node"
)

// ZLazyImgOutputText impl ZLazyImgOutputText
type ZLazyImgOutputText struct {
	*ZOutputText

	SvgPlaceholder bool
}

// NewZLazyImgOutputText returns instance of ZLazyImgOutputText
func NewZLazyImgOutputText(content string, opts ...Option) (*ZLazyImgOutputText, error) {
	z, err := NewZOutputText(content, opts...)
	if err != nil {
		return nil, err
	}

	return &ZLazyImgOutputText{
		ZOutputText:    z,
		SvgPlaceholder: false,
	}, nil
}

// OnImgTag overwrites ZOutputText.OnImgTag
func (z *ZLazyImgOutputText) OnImgTag(n *node.Node) {
	z.ZOutputText.OnImgTag(n)

	origin := n.CloneTree()
	if n.ParentNode() != nil {
		n.SetAttr("class", n.GetAttrOrDefault("class", "")+" lazy")
		n.SetAttr("data-actualsrc", n.GetAttrOrDefault("src", ""))

		if z.SvgPlaceholder {
			// http://jira.in.zhihu.com/browse/COMMUNITY-46
			var width, height int
			w := n.GetAttrOrDefault("data-rawwidth", "0")
			width, err := strconv.Atoi(w)
			if err != nil {
				width = 0
			}
			h := n.GetAttrOrDefault("data-rawheight", "0")
			height, err = strconv.Atoi(h)
			if err != nil {
				height = 0
			}
			s := fmt.Sprintf("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='%d' height='%d'></svg>", width, height)
			n.SetAttr("src", s)
		} else {
			// whitedot placeholder img
			n.SetAttr("src", "//zhstatic.zhihu.com/assets/zhihu/ztext/whitedot.jpg")
		}

		origin.Parent = nil // hack for AppendChild
		noscriptNode := node.NewNoscript()
		noscriptNode.AppendChild(origin)
		// e.addprevious(noscript)
		// Adds the element as a preceding sibling directly before this element.
		n.ParentNode().InsertBefore(noscriptNode, n)
	}
}
