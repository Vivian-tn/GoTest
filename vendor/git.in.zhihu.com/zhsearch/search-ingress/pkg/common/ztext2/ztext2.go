package ztext2

import (
	"fmt"

	"git.in.zhihu.com/production_backend/pier"
	"git.in.zhihu.com/go/ztext/node"
	"git.in.zhihu.com/go/ztext/ztext"
	"git.in.zhihu.com/zhsearch/search-ingress/pkg/log"
)

type ApiZOutputText struct {
	ztext.ZOutputText
}

func (z *ApiZOutputText) OnATag(n *node.Node) {
	eHref := n.GetAttrOrDefault("href", "")
	if eHref != "" {
		if eHref[:2] == "//" {
			n.SetAttr("href", "https:"+eHref)
		} else if eHref[:1] == "/" {
			n.SetAttr("href", "https://www.zhihu.com"+eHref)
		}
	}

	eClass := n.GetAttrOrDefault("class", "")
	switch eClass {
	case "member_mention":
		z.ProcessMemberMentionOutput(n)
	default:
		z.ProcessLinkOutput(n)
	}
}

func NewApiZOutputText(content string) (*ApiZOutputText, error) {
	sz := ztext.ZOutputText{
		RootDomain: "zhihu.com",
		Protocol:   "https",
	}
	ztext.StrictBMP(true)(&sz)
	ztext.ImageURLResolver(pier.GetFullURL)(&sz)

	z := &ApiZOutputText{
		sz,
	}

	content = ztext.ZSafeText(content, true)

	fragment, err := node.NewFragmentFromString(content)
	if err != nil {
		return nil, fmt.Errorf("initializing ZOutputText failed while parsing content: %s", err)
	}
	z.ZDom = ztext.NewZDom(fragment, ztext.DomStrictBMP(true))

	return z, nil
}

type ZQRCodeReplaceText struct {
	*ztext.ZDom
}

func (z *ZQRCodeReplaceText) OnImgTag(n *node.Node) {
	tag := n.GetAttrOrDefault("data-tag", "")
	var replacement *node.Node
	qrvalue := n.GetAttrOrDefault("data-qrvalue", "")
	if tag == "qrcode" && qrvalue != "" {
		replacement = node.NewP()
		replacement.SetText(qrvalue + " (二维码自动识别)")
		n.SetText("")
		n.ParentNode().InsertBefore(replacement, n)
		n.DropTag()
	}
}

func NewZQRCodeReplaceText(content string) (*ZQRCodeReplaceText, error) {
	dom := ztext.ZText2(content)
	z := &ZQRCodeReplaceText{
		dom,
	}

	return z, nil
}

func ZTextPlaintext(content string) string {
	return GetZTextPlaintext(content)
}

func GetZTextPlaintext(content string) string {
	if content == "" {
		return ""
	}

	zot, err := ztext.NewZOutputText(content)
	if err != nil {
		log.Warnf("GetZTextPlaintext err:%v", err)
		return content
	}

	s, err := ztext.RenderPlainText(zot)
	if err != nil {
		log.Warnf("GetZTextPlaintext err:%v", err)
		return content
	}

	return s
}

func APIZOutputText(content string) string {
	return GetAPIZOutputText(content)
}

func GetAPIZOutputText(content string) string {
	if content == "" {
		return ""
	}

	zot, err := NewApiZOutputText(content)
	if err != nil {
		log.Warnf("GetAPIZOutputText err:%v", err)
		return content
	}

	s, err := ztext.Render(zot)
	if err != nil {
		log.Warnf("GetAPIZOutputText err:%v", err)
		return content
	}

	return s
}

func WebZOutputText(content string) string {
	return GetWebZOutputContent(content)
}

func GetWebZOutputContent(content string) string {
	if content == "" {
		return ""
	}

	zot, err := NewZQRCodeReplaceText(content)
	if err != nil {
		log.Warnf("NewZQRCodeReplaceText err:%v", err)
		return content
	}

	s, err := ztext.Render(zot)
	if err != nil {
		log.Warnf("RenderZQRCodeReplaceText err:%v", err)
		return content
	}

	lazyImage, err := ztext.NewZLazyImgOutputText(s, ztext.AllowGif(true))
	if err != nil {
		log.Warnf("NewZLazyImgOutputText err:%v", err)
		return content
	}

	lazyImage.SvgPlaceholder = true

	res, err:=ztext.Render(lazyImage)
	if err != nil {
		log.Warnf("RenderZLazyImgOutputText err:%v", err)
		return content
	}

	return res
}

func ZExcerptPlaintext(content string) string {
	return GetZExcerptPlaintext(content)
}

func GetZExcerptPlaintext(content string) string {
	if content == "" {
		return ""
	}

	zot, err := ztext.NewZExcerptText(content, ztext.ExcerptLength(200))
	if err != nil {
		log.Warnf("GetZExcerptPlaintext err:%v", err)
		return content
	}

	s, err := ztext.RenderPlainText(zot)
	if err != nil {
		log.Warnf("GetZExcerptPlaintext err:%v", err)
		return content
	}

	return s
}

func WebZOutputPlaintext(content string) string {
	return GetWebZOutputPlaintext(content)
}

func GetWebZOutputPlaintext(content string) string {
	if content == "" {
		return ""
	}

	zot, err := ztext.NewZOutputText(content, ztext.StrictBMP(false))
	if err != nil {
		log.Warnf("GetWebZOutputPlaintext err:%v", err)
		return content
	}

	s, err := ztext.RenderPlainText(zot)
	if err != nil {
		log.Warnf("GetWebZOutputPlaintext err:%v", err)
		return content
	}

	return s
}
