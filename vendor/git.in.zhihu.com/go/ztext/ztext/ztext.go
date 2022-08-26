package ztext

import (
	"regexp"

	"git.in.zhihu.com/go/ztext/node"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// OutsideDomains means same as its name
	OutsideDomains = []string{"tmail.zhihu.com"}
	// ProtocolMap is map between Protocol and its prefix value in url
	ProtocolMap = map[protocolType]string{
		HTTP:     "http:",
		HTTPS:    "https:",
		Relative: "",
	}
	// LinkdDomain is link domain of zhihu
	LinkdDomain     = "link.zhihu.com"
	linkdURLPattern = regexp.MustCompile(`((http:)|(https:))?//link.zhihu.com/\?target=`)

	linkLabelDomainMap = map[string]string{
		"zhuanlan.zhihu.com": "zhuanlan",
		"www.zhihu.com":      "zhihu",
		"zhihu.com":          "zhihu",
	}
	linkLabelPatternMap = map[string]map[*regexp.Regexp]string{
		"zhuanlan": {
			regexp.MustCompile(`^/p/\d+`):   "文章",
			regexp.MustCompile(`^/[\w-]*$`): "专栏",
		},
		"zhihu": {
			regexp.MustCompile(`^/question/\d+/?$`):                      `问题`,
			regexp.MustCompile(`^/question/\d+/answer/\d+/?$`):           `回答`,
			regexp.MustCompile(`^/collection/\d+/?$`):                    `收藏夹`,
			regexp.MustCompile(`^/publications/(nacl|weekly|hour)/\d+$`): `电子书`,
			regexp.MustCompile(`^/roundtable/\w+$`):                      `圆桌`,
			regexp.MustCompile(`^/lives/\d+$`):                           `Live`,
			regexp.MustCompile(`^/topic/(\d+|\d+/hot)?$`):                `话题`,
			regexp.MustCompile(`^/people/[^/]+$`):                        `个人主页`,
			regexp.MustCompile(`^/org/[^/]+$`):                           `机构帐号`,
		},
	}
)

// ZText is the underlying interface, ZDom(including his children) must
// implement it.
type ZText interface {
	PrepareWalkDOM() error
	WalkDOM(node.Visitor) error
	Render() (string, error)

	OnATag(n *node.Node)
	OnImgTag(n *node.Node)
	OnPTag(n *node.Node)
	OnVideoTag(n *node.Node)
	OnCodeTag(n *node.Node)
	OnEquationTag(n *node.Node)
	OnUlTag(n *node.Node)

	Thumbnail(minWidth, minHeight int) *node.Node
	PlainText() string

	Walked() bool
}

func walk(z ZText) error {
	if !z.Walked() {
		err := z.PrepareWalkDOM()
		if err != nil {
			return err
		}

		err = z.WalkDOM(func(n *node.Node) bool {
			switch n.DataAtom {
			case atom.A:
				z.OnATag(n)
			case atom.P:
				z.OnPTag(n)
			case atom.Video:
				z.OnVideoTag(n)
			case atom.Img:
				z.OnImgTag(n)
			case atom.Code:
				z.OnCodeTag(n)
			case atom.Ul:
				z.OnUlTag(n)
			case 0:
				// equation tag is not a standard tag
				// go html parser will set DataAtom to 0 for unrecognized tag.
				// 线上有内容为 equation 的 textnode, 因此这里必须加上 Type 判断
				if n.Type == html.ElementNode && n.Data == "equation" {
					z.OnEquationTag(n)
				}
			}
			return true
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Render a ZText instance to its corresponding output in string format
func Render(z ZText) (string, error) {
	if err := walk(z); err != nil {
		return "", err
	}

	return z.Render()
}

// RenderPlainText a ZText instance to its corresponding output in string format
func RenderPlainText(z ZText) (string, error) {
	if err := walk(z); err != nil {
		return "", err
	}

	return z.PlainText(), nil
}
