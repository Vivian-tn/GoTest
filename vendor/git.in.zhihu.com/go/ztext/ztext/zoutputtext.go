package ztext

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"git.in.zhihu.com/production_backend/pier"
	"git.in.zhihu.com/go/ztext/escape"
	"git.in.zhihu.com/go/ztext/link"
	"git.in.zhihu.com/go/ztext/node"
	"git.in.zhihu.com/go/ztext/util"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ZOutputText struct takes string content to document tree, processed it, and
// renders to a processed string.
type ZOutputText struct {
	// root domain, default is zhihu.com
	RootDomain string

	// https/http/relative, default is https
	Protocol string

	// watermark style
	watermark string

	// root node of document tree
	*ZDom

	// IMG_URL_RESOLVER
	imageURLResolver ImageURLResolverFunc

	// if allow gif
	allowGif bool

	// if enable Strict BMP
	strictBMP bool
}

// ImageURLResolverFunc defines IMG_URL_RESOLVER
type ImageURLResolverFunc func(token string, size string, fmt string, quality string, secure bool) string

var _ ZText = (*ZOutputText)(nil)

var (
	forbiddenTag = map[atom.Atom]struct{}{
		atom.A:    {},
		atom.Pre:  {},
		atom.Code: {},
		0:         {}, // for equation tag
	}

	imgCommentAllowedAttrs = map[string]struct{}{
		"href":            {},
		"class":           {},
		"data-width":      {},
		"data-height":     {},
		"data-sticker-id": {},
	}

	dataAttrs = map[string]struct{}{
		"data-draft-node":     {},
		"data-draft-type":     {},
		"data-image":          {},
		"data-image-width":    {},
		"data-image-height":   {},
		"data-entity-type":    {},
		"data-entity-data":    {},
		"data-mcn-id":         {},
		"data-metalink-id":    {},
		"data-ad-id":          {},
		"data-sku-id":         {},
		"data-file-source":    {},
		"data-file-size":      {},
		"data-file-extension": {},
		"data-file-type":      {},
		"data-size":           {},
	}

	autoLinkPattern = regexp.MustCompile(`(?i)(?P<body>https?://(?P<host>[a-z0-9._\-]+(:\d{2,6})?)(?:/[/\-_.,a-z0-9%&?!;=~#:@\+]*)?(?:\([/\-_.,a-z0-9%&?;=~#:@\+]*\))?)`)

	boxHTMLTmpl = `<a class="video-box" href="%s" target="_blank" data-video-id="%s" data-video-playable="%s" data-name="%s" data-poster="%s" data-lens-id="%s">` +
		`<img class="thumbnail" src="%s">` +
		`<span class="content">` +
		`<span class="title">%s<span class="z-ico-extern-gray"></span><span class="z-ico-extern-blue"></span></span>` +
		`<span class="url"><span class="z-ico-video"></span>%s</span>` +
		`</span>` +
		`</a>`
	linkTemplate       = `<a href="%s" class="link-box" target="_blank">%s<span class="content"><span class="title">%s</span><span class="url">%s</span>%s</span></a>`
	zVideoLinkCardHTML = `<a href="https://www.zhihu.com/zvideo/%s" data-draft-node="block" data-draft-type="link-card"></a>`
	watermarkMap       = map[string]string{
		"original":          "data-original-src",
		"watermark":         "data-watermark-src",
		"private_watermark": "data-private-watermark-src",
	}
)

// NewZOutputText returns a ZOutputText instance with default params
func NewZOutputText(content string, opts ...Option) (*ZOutputText, error) {
	z := &ZOutputText{
		RootDomain:       "zhihu.com",
		Protocol:         ProtocolMap[HTTPS],
		watermark:        "",
		imageURLResolver: pier.GetFullURL,
		allowGif:         false,
		strictBMP:        true,
	}

	for _, opt := range opts {
		opt(z)
	}

	content = ZSafeText(content, z.strictBMP)

	fragment, err := node.NewFragmentFromString(content)
	if err != nil {
		return nil, fmt.Errorf("initializing ZOutputText failed while parsing content: %s", err)
	}
	z.ZDom = NewZDom(fragment, DomStrictBMP(z.strictBMP))

	return z, nil
}

// PrepareWalkDOM traverse z.Root dom tree, and do some preparation work.
func (z *ZOutputText) PrepareWalkDOM() error {
	z.Walk(func(n *node.Node) bool {
		a := n.DataAtom
		if _, ok := forbiddenTag[a]; ok {
			if a != 0 {
				return true
			} else if a == 0 && n.Data == "equation" {
				return true
			}
		}

		// Not handling TextNode
		if n.Type == html.TextNode {
			return true
		}

		if t := n.Text(); t != "" {
			n.SetText(link.UnifyLink(t))
		}
		if t := n.Tail(); t != "" {
			_ = n.SetTail(link.UnifyLink(t))
		}

		return true
	})

	link.AutoLink(z.Node, []*regexp.Regexp{autoLinkPattern})
	return nil
}

func (z *ZOutputText) detectLinkLabel(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return ""
	}

	if patternMapKey, ok := linkLabelDomainMap[parsed.Host]; ok {
		for pattern, category := range linkLabelPatternMap[patternMapKey] {
			if len(pattern.FindString(parsed.Path)) > 0 {
				return category
			}
		}
	}

	return ""
}

// ProcessMemberMentionOutput process mention output
func (z *ZOutputText) ProcessMemberMentionOutput(n *node.Node) {
	memberHash, _ := n.GetAttr("data-hash")
	if memberHash != "" {
		n.SetAttr("href",
			fmt.Sprintf("%s//www.%s/people/%s",
				z.Protocol,
				z.RootDomain,
				memberHash))
		n.SetAttr("data-hovercard", "p$b$"+memberHash)
		n.DelAttr("data_id")
	} else {
		n.DropTag()
	}
}

func (z *ZOutputText) prefixedHref(lnk string) string {
	linkInfo, err := url.Parse(lnk)
	if err != nil || len(linkInfo.Scheme) == 0 {
		lnk = "http://" + lnk
	}
	return lnk
}

// ProcessShareLinkOutput process share link output
func (z *ZOutputText) ProcessShareLinkOutput(n *node.Node) {
	href := n.GetAttrOrDefault("href", "")
	lnk := href
	imageSrc := n.GetAttrOrDefault("data-img-src", "")
	title := n.GetAttrOrDefault("data-title", "")
	prefixedHref := z.prefixedHref(href)
	label := z.detectLinkLabel(prefixedHref)

	// http://ph.in.zhihu.com/T61547#962932
	localShortenLink := func(link string) string {
		linkInfo, err := url.Parse(link)
		if err != nil || len(linkInfo.Scheme) == 0 {
			return ""
		}

		return strings.TrimPrefix(linkInfo.Host, "www.")
	}

	if len(title) == 0 {
		title = href
		lnk = ""
		label = ""
	}

	lnk = localShortenLink(lnk)
	imageDom := ""
	if len(imageSrc) > 0 {
		imageDom = fmt.Sprintf(`<img src="%s" class="thumbnail">`, escape.XhtmlEscape(imageSrc))
	}

	labelDom := ""
	if len(label) > 0 {
		labelDom = fmt.Sprintf(`<span class="label">%s</span>`, escape.XhtmlEscape(label))
	}

	linkBox := fmt.Sprintf(linkTemplate,
		escape.XhtmlEscape(href),
		imageDom,
		escape.XhtmlEscape(title),
		escape.XhtmlEscape(lnk),
		labelDom,
	)

	boxElement, err := node.NewFragmentsFromString(linkBox)
	if err != nil {
		log.Printf("Creating box node failed: %s", linkBox)
		return
	}

	n.SetText("")
	n.InsertBefore(boxElement[0], (*node.Node)(n.FirstChild))
	n.DropTag()
}

// ProcessImgCommentOutput process image comment output
func (z *ZOutputText) ProcessImgCommentOutput(n *node.Node) {
	delAttrs := make([]string, 0)
	for _, attr := range n.Attr {
		key := attr.Key
		if _, ok := imgCommentAllowedAttrs[key]; !ok {
			delAttrs = append(delAttrs, key)
		}
	}
	n.DelAttrs(delAttrs...)

	var href string
	hash, _ := n.GetAttr("href")
	if eClass, _ := n.GetAttr("class"); eClass == "comment_sticker" {
		href = pier.GetFullURL(hash, "", "", "", true)
	} else {
		href = pier.GetFullURL(hash, "qhd", "", "", true)
	}
	n.SetAttr("href", href)
}

func (z *ZOutputText) shortenLink(n *node.Node) {
	showEllipsis := false
	longText := n.GetAttrOrDefault("href", "")
	aText := n.TextContent()
	n.SetTextContent(aText)

	parent := n.ParentNode()
	// TODO(yangbo01) 需要重构下，分离数据结构操作与业务代码
	for child := n.FirstChild; child != nil; {
		next := child.NextSibling
		if child.DataAtom == atom.Img {
			// 将 A 的子 tag 中 Img 移到 A 的后面，即变为 A.NextSibling
			if parent != nil {
				// 链表删除元素的常规操作
				// e.prev.next = e.next
				// e.next.prev = e.prev
				if prev := child.PrevSibling; prev != nil {
					prev.NextSibling = child.NextSibling
				}
				if next != nil {
					next.PrevSibling = child.PrevSibling
				}

				// InsertAfter 要求 child 的这三个属性都要是 nil
				child.Parent, child.NextSibling, child.PrevSibling = nil, nil, nil
				_ = parent.InsertAfter((*node.Node)(child), n)
			}
		} else if child.Type != html.TextNode {
			// 移除所有不是 TextNode 或者 Img 的标签
			n.RemoveChild((*node.Node)(child))
		}
		child = next
	}

	if len(aText) > 0 && aText != longText {
		// Customize a tag.
		n.SetAttr("class", n.GetAttrOrDefault("class", "")+" wrap")
	} else {
		shortText, err := link.ShortenLink(longText)
		if err != nil {
			shortText = ""
		}
		if strings.HasSuffix(shortText, "...") {
			shortText = shortText[:len(shortText)-3]
			showEllipsis = true
		}
		startIndex := strings.Index(longText, shortText)
		if startIndex == -1 {
			log.Printf("Impossible shortenLink: %s->%s\n", longText, shortText)
			return
		}

		stopIndex := startIndex + len(shortText)
		n.SetText("")

		spanElement1 := node.NewNode(atom.Span, "class", "invisible")
		spanElement1.SetText(longText[:startIndex])
		n.AppendChild(spanElement1)

		spanElement2 := node.NewNode(atom.Span, "class", "visible")
		spanElement2.SetText(longText[startIndex:stopIndex])
		n.AppendChild(spanElement2)

		spanElement3 := node.NewNode(atom.Span, "class", "invisible")
		spanElement3.SetText(longText[stopIndex:])
		n.AppendChild(spanElement3)
	}

	if showEllipsis {
		spanElement := node.NewNode(atom.Span, "class", "ellipsis")
		n.AppendChild(spanElement)
	}
}

// ProcessLinkOutput process link output
func (z *ZOutputText) ProcessLinkOutput(n *node.Node) {

	dataImg := n.GetAttrOrDefault("data-image", "")

	if len(dataImg) > 0 {
		size := n.GetAttrOrDefault("data-image-size", "")
		n.SetAttr("data-image", pier.GetFullURL(dataImg, size, "jpg", "", true))
	}
	delAttrs := make([]string, 0)
	for _, attr := range n.Attr {
		key := attr.Key
		if strings.HasPrefix(key, "data-") {
			if _, ok := dataAttrs[key]; !ok {
				delAttrs = append(delAttrs, key)
			}
		}
	}
	n.DelAttrs(delAttrs...)

	href := n.GetAttrOrDefault("href", "")

	// 相对路径不做外链处理
	if parsedHref, err := url.Parse(href); err != nil {
		processedHref, err := url.PathUnescape(href)
		if err == nil {
			if parsedHref, err = url.Parse(processedHref); err == nil && parsedHref.Host == "" {
				return
			}

		}
	} else if parsedHref.Host == "" { // err == nil
		return
	}

	n.SetAttr("class", "")

	aText := n.TextContent()
	if len(linkdURLPattern.FindString(aText)) > 0 {
		unquotedText, err := url.QueryUnescape(linkdURLPattern.ReplaceAllLiteralString(aText, ""))
		if err != nil {
			log.Printf("Failed to remove linkd domain: %s\n", aText)
		}
		n.SetTextContent(unquotedText)
	}

	quotedColon := "%3A"
	if link.IsLinkOfDomain(href, LinkdDomain) {
		href = linkdURLPattern.ReplaceAllLiteralString(href, "")
		if strings.Index(href, quotedColon) > 0 {
			if newHref, err := url.QueryUnescape(href); err == nil {
				href = newHref
			}
		}
		n.SetAttr("href", href)
	}

	z.shortenLink(n)

	if (href != "") &&
		(!link.IsLinkOfDomain(href, z.RootDomain)) &&
		(!link.IsLinkUseSpecialScheme(href, "zhihu")) ||
		link.IsOutLinkDomain(href, OutsideDomains) {
		quotedHref := strings.Replace(url.QueryEscape(href), "%2F", "/", -1)
		finalHref := fmt.Sprintf("%s//%s/?target=%s",
			z.Protocol,
			LinkdDomain,
			quotedHref)

		n.SetAttr("href", finalHref)
		n.SetAttr("target", "_blank")
		n.SetAttr("class", n.GetAttrOrDefault("class", "")+" external")
		n.SetAttr("rel", "nofollow noreferrer")
	} else {
		n.SetAttr("class", "internal")
	}
}

// OnATag is the callback function when visiting <a> tag
func (z *ZOutputText) OnATag(n *node.Node) {
	eClass := n.GetAttrOrDefault("class", "")
	switch eClass {
	case "member_mention":
		z.ProcessMemberMentionOutput(n)
	case "share-link":
		// Seems that no longer used.
		z.ProcessShareLinkOutput(n)
	case "comment_sticker", "comment_gif", "comment_img":
		z.ProcessImgCommentOutput(n)
	default:
		switch n.GetAttrOrDefault("data-draft-type", "") {
		case "mcn-link-card":
			z.ProcessMcnCardOutput(n)
		case "metalink":
			z.ProcessMetalinkOutput(n)
		case "km-sku-card":
			z.ProcessSkuCardOutput(n)
		case "file-link-card":
			z.ProcessFileLinkCardOutput(n)
		default:
			z.ProcessLinkOutput(n)
		}
	}
}

func (z *ZOutputText) ProcessMcnCardOutput(n *node.Node) {
	if n.GetAttrOrDefault("data-draft-node", "") != "block" ||
		n.GetAttrOrDefault("data-draft-type", "") != "mcn-link-card" ||
		n.GetAttrOrDefault("data-mcn-id", "") == "" {
		n.DropTag()
	}
}

func (z *ZOutputText) ProcessMetalinkOutput(n *node.Node) {
	if n.GetAttrOrDefault("data-draft-node", "") != "block" ||
		n.GetAttrOrDefault("data-draft-type", "") != "metalink" ||
		n.GetAttrOrDefault("data-metalink-id", "") == "" {
		n.DropTag()
	}
}

func (z *ZOutputText) ProcessSkuCardOutput(n *node.Node) {
	if n.GetAttrOrDefault("data-draft-node", "") != "block" ||
		n.GetAttrOrDefault("data-draft-type", "") != "km-sku-card" ||
		n.GetAttrOrDefault("data-sku-id", "") == "" {
		n.DropTag()
	}
	for _, attr := range n.Attr {
		if attr.Key != "data-draft-node" &&
			attr.Key != "data-draft-type" &&
			attr.Key != "data-sku-id" {
			n.DelAttr(attr.Key)
		}
	}
}

func (z *ZOutputText) ProcessFileLinkCardOutput(n *node.Node) {
	if n.GetAttrOrDefault("data-draft-node", "") != "block" ||
		n.GetAttrOrDefault("data-draft-type", "") != "file-link-card" {
		n.DropTag()
	}
	for _, attr := range n.Attr {
		if attr.Key != "data-draft-node" &&
			attr.Key != "data-draft-type" &&
			attr.Key != "href" &&
			attr.Key != "data-file-source" &&
			attr.Key != "data-file-size" &&
			attr.Key != "data-file-extension" &&
			attr.Key != "data-file-type" {
			n.DelAttr(attr.Key)
		}
	}
}

// OnImgTag is the callback function when visiting <img> tag
func (z *ZOutputText) OnImgTag(n *node.Node) {
	gifMode := false
	imgSrc := n.GetAttrOrDefault("src", "")
	if imgSrc == "" {
		n.DropTag()
		return
	}

	if strings.HasSuffix(imgSrc, ".gif") {
		if z.allowGif {
			gifMode = true
			imgHash := util.URL2Token(imgSrc, false)
			n.SetAttr(
				"data-thumbnail",
				z.imageURLResolver(imgHash, "b", "", "", true),
			)

		} else {
			imgSrc = imgSrc[:len(imgSrc)-4]
		}
	}

	if strings.HasSuffix(imgSrc, ".webp") {
		imgSrc = util.URL2Token(imgSrc, false)
	}

	// JIRA: http://jira.in.zhihu.com/browse/AL-671
	// 过滤图像类型属性
	n.DelAttr("data-tags")
	n.DelAttr("data-qrcode-value")

	n.SetAttr("src", z.imageURLResolver(imgSrc, "b", "", "", true))

	// 选择水印类型，优先选择传进来的参数，其次才使用 dom 里的参数，dom 里为 original 时不可被覆盖
	watermarkType := n.GetAttrOrDefault("data-watermark", "original")
	if watermarkType != "original" {
		if z.watermark != "" {
			watermarkType = z.watermark
		} else if _, ok := watermarkMap[watermarkType]; !ok {
			watermarkType = "original"
		}
	}
	// i'm sure watermarkType in watermarkMap
	src, _ := n.GetAttr(watermarkMap[watermarkType])
	if src != "" {
		format := ""
		if gifMode {
			format = "gif"
		}
		n.SetAttr("src", z.imageURLResolver(src, "b", format, "", true))
	}

	privateSrc, _ := n.GetAttr(watermarkMap["private_watermark"])
	if privateSrc != "" {
		n.SetAttr(
			"data-default-watermark-src",
			z.imageURLResolver(privateSrc, "b", "", "", true),
		)
	}

	n.DelAttrs("data-watermark",
		"data-original-src",
		"data-watermark-src",
		"data-private-watermark-src",
	)

	n.SetAttr("class", "content_image")

	dataRawWidth, ok1 := n.GetAttr("data-rawwidth")
	_, ok2 := n.GetAttr("data-rawheight")

	if ok1 && ok2 {
		n.SetAttr("width", dataRawWidth)

		width, err := strconv.Atoi(dataRawWidth)
		if err != nil {
			width = 0
		}

		if width > 420 {
			src, _ := n.GetAttr("src")
			n.SetAttr("data-original", z.imageURLResolver(util.URL2Token(src, false), "r", "", "", true))
			n.SetAttr("class", "origin_image zh-lightbox-thumb")
		}
	}

	dataCaption := n.GetAttrOrDefault("data-caption", "")
	figureNode := node.NewNode(atom.Figure)

	if dataSize, ok := n.GetAttr("data-size"); ok {
		figureNode.SetAttr("data-size", dataSize)
	}

	if len(dataCaption) > 0 {
		n.DelAttr("data-caption")
		figureNode.InsertBefore(
			node.NewNodeWithText(atom.Figcaption, dataCaption),
			(*node.Node)(figureNode.FirstChild),
		)
	}

	parent := n.Parent
	tail := n.Tail()
	tailNode := n.NextSibling
	parent.InsertBefore((*html.Node)(figureNode), (*html.Node)(n))
	if tail != "" {
		parent.RemoveChild(tailNode)
		_ = (*node.Node)(parent).InsertAfter(
			node.NewTextNode(tailNode.Data),
			figureNode,
		)
	}

	parent.RemoveChild((*html.Node)(n))

	figureNode.InsertBefore(n, (*node.Node)(figureNode.FirstChild))
}

// OnPTag is the callback function when visiting <p> tag
func (z *ZOutputText) OnPTag(n *node.Node) {
	// add class when matching `<p><br></p>`
	if _, ok := n.GetAttr("class"); ok {
		return
	}
	if n.FirstChild != nil {
		child := (*node.Node)(n.FirstChild)
		if child.NextSibling == nil && child.DataAtom == atom.Br {
			n.SetAttr("class", "ztext-empty-paragraph")
		}
	}
}

// OnVideoTag is the callback function when visiting <video> tag
func (z *ZOutputText) OnVideoTag(n *node.Node) {
	if value, _ := n.GetAttr("data-zvideo-id"); value != "" {
		linkCardElement, _ := node.NewFragmentsFromString(fmt.Sprintf(zVideoLinkCardHTML,
			escape.XhtmlEscape(value)),
		)
		n.InsertBefore(linkCardElement[0], (*node.Node)(n.FirstChild))
	} else if value, _ := n.GetAttr("data-sourceurl"); value != "" {
		escapedValue := strings.Replace(url.QueryEscape(value), "%2F", "/", -1)
		sourceURL := fmt.Sprintf("%s//%s/?target=%s",
			z.Protocol,
			LinkdDomain,
			escapedValue)

		// very c style TODO: can we find some optimized way of doing this?
		boxElement, _ := node.NewFragmentsFromString(fmt.Sprintf(boxHTMLTmpl,
			escape.XhtmlEscape(sourceURL),
			escape.XhtmlEscape(n.GetAttrOrDefault("data-video-id", "")),
			escape.XhtmlEscape(n.GetAttrOrDefault("data-video-playable", "")),
			escape.XhtmlEscape(n.GetAttrOrDefault("data-name", "")),
			escape.XhtmlEscape(n.GetAttrOrDefault("poster", "")),
			escape.XhtmlEscape(n.GetAttrOrDefault("data-lens-id", "")),
			escape.XhtmlEscape(n.GetAttrOrDefault("poster", "")),
			escape.XhtmlEscape(n.GetAttrOrDefault("data-name", "")),
			escape.XhtmlEscape(n.GetAttrOrDefault("data-sourceurl", ""))),
		)

		lensStatus := n.GetAttrOrDefault("data-lens-status", "")
		if lensStatus == "deleted" {
			boxElement[0].SetAttr("data-lens-status", lensStatus)
			boxElement[0].SetAttr("data-description", "[资源不存在]")
		}

		n.InsertBefore(boxElement[0], (*node.Node)(n.FirstChild))
	}
	n.DropTag()
}

// OnCodeTag is the callback function when visiting <code> tag
func (z *ZOutputText) OnCodeTag(n *node.Node) {
	inline := n.GetAttrOrDefault("class", "")
	if inline == "inline" {
		n.DelAttr("class")
		return
	}

	// 处理代码块
	if !n.HasText() {
		n.DropTag()
		return
	}

	lang, _ := n.GetAttr("lang")
	if lang == "" {
		lang = "text"
	}

	highlighted, lang, err := util.Highlight(n.Text(), lang)
	if err != nil {
		panic(fmt.Errorf("highlight code failed. code: %s, error: %s", n.Text(), err))
	}

	var tree []*node.Node
	if tree, err = node.NewFragmentsFromString(highlighted); err != nil {
		panic(fmt.Errorf("generate dom tree from highlighted code failed. code: %s, error: %s", highlighted, err))
	}

	// code-block>pre => code-block>pre>code
	codeBlock := tree[0]
	code := (*node.Node)(codeBlock.FirstChild)
	code.SetAtom(atom.Code)
	code.SetData("code")
	code.SetAttr("class", "language-"+lang)
	codeBlock.RemoveChild(code)
	pre := node.NewNode(atom.Pre)
	pre.AppendChild(code)
	codeBlock.AppendChild(pre)

	p := n.ParentNode()
	p.InsertBefore(codeBlock, n)
	n.DropTree()
}

// OnEquationTag is the callback function when visiting <equation> tag
func (z *ZOutputText) OnEquationTag(n *node.Node) {
	equation := n.Text()
	if len(equation) == 0 {
		n.DropTree()
		return
	}

	n.SetAtom(atom.Img)
	n.SetData("img")

	var urlFormat string
	var eeimg string
	var display string
	if display, _ = n.GetAttr("display"); display != "" {
		urlFormat = "%s//%s/equation?%s&display=2"
		eeimg = "2"
	} else {
		urlFormat = "%s//%s/equation?%s"
		eeimg = "1"
	}

	query := url.Values{"tex": []string{equation}}
	n.SetAttr("src", fmt.Sprintf(urlFormat, z.Protocol, "www."+z.RootDomain, query.Encode()))
	n.SetAttr("alt", equation)
	n.SetAttr("eeimg", eeimg)

	if display != "" {
		n.DelAttr("display")
		container := node.NewDiv()
		container.SetAttr("class", "ee-displaymath")
		p := n.ParentNode()
		p.InsertBefore(container, n)
		p.RemoveChild(n)
		container.InsertBefore(n, (*node.Node)(container.FirstChild))
	}
	// Img node has no child.
	n.RemoveAllChildren()
}

// AllowGif returns allowGif field
func (z *ZOutputText) AllowGif() bool {
	return z.allowGif
}

// ImageURLResolver returns imageURLResolver field
func (z *ZOutputText) ImageURLResolver() ImageURLResolverFunc {
	return z.imageURLResolver
}
