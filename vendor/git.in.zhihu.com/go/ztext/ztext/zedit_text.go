package ztext

import (
	"fmt"
	"strings"

	"git.in.zhihu.com/production_backend/pier"
	"git.in.zhihu.com/go/ztext/node"
	"git.in.zhihu.com/go/ztext/util"
	"golang.org/x/net/html/atom"
)

// ZEditText impl ZEditText
type ZEditText struct {
	*ZOutputText
}

// NewZEditText returns a instance of ZEditText
func NewZEditText(content string, opts ...Option) (*ZEditText, error) {
	z, err := NewZOutputText(content, opts...)
	if err != nil {
		return nil, err
	}

	return &ZEditText{
		ZOutputText: z,
	}, nil
}

// PrepareWalkDOM overwrites PrepareWalkDOM
func (z *ZEditText) PrepareWalkDOM() error {
	// do nothing
	return nil
}

// OnATag is the callback function when visiting <a> tag
func (z *ZEditText) OnATag(n *node.Node) {
	eClass := n.GetAttrOrDefault("class", "")
	switch eClass {
	case "member_mention":
		z.ProcessMemberMentionOutput(n)
	case "share-link":
		// Seems that no longer used.
		z.processShareLinkOutput(n)
	case "comment_sticker", "comment_gif", "comment_img":
		z.ProcessImgCommentOutput(n)
	default:
		z.processLinkOutput(n)
	}
}

// OnImgTag overwrites OnImgTag
func (z *ZEditText) OnImgTag(n *node.Node) {
	imgSrc, _ := n.GetAttr("src")
	if imgSrc == "" {
		n.DropTag()
		return
	}

	if strings.HasSuffix(imgSrc, ".gif") {
		imgHash := util.URL2Token(imgSrc, false)
		n.SetAttr("data-thumbnail", z.imageURLResolver(imgHash, "b", "", "", true))
	} else if strings.HasSuffix(imgSrc, ".webp") {
		imgSrc = util.URL2Token(imgSrc, false)
	}

	n.SetAttr("src", z.imageURLResolver(imgSrc, "b", "", "", true))
	n.SetAttr("class", "content_image")

	watermarks := []string{"data-original-src", "data-watermark-src", "data-private-watermark-src"}
	for _, w := range watermarks {
		if src, _ := n.GetAttr(w); src != "" {
			n.SetAttr(w, z.imageURLResolver(src, "b", "", "", true))
		}
	}
}

// OnVideoTag overwrites OnVideoTag
func (z *ZEditText) OnVideoTag(n *node.Node) {
	if zVideoID := n.GetAttrOrDefault("data-zvideo-id", ""); zVideoID != "" {
		htmlText := `<a href="https://www.zhihu.com/zvideo/%s" data-draft-node="block" data-draft-type="link-card"></a>`
		vElement, _ := node.NewFragmentsFromString(fmt.Sprintf(htmlText, zVideoID))
		n.InsertBefore(vElement[0], (*node.Node)(n.FirstChild))
		n.DropTag()
		return
	}

	lensStatus := n.GetAttrOrDefault("data-lens-status", "")
	lensID, _ := n.GetAttr("data-lens-id")
	censorItems := n.GetAttrOrDefault("data-lens-censor-fail-items", "")
	if lensID != "" && (lensStatus == "uploading_fail" || lensStatus == "reviewing_fail" || lensStatus == "deleted") {
		// <a class="unprocessable_video" href="javascript:;"></a>
		vElement := node.NewA()
		vElement.SetAttr("class", "unprocessable_video")
		vElement.SetAttr("href", "javascript:;")

		vElement.SetAttr("data-lens-status", lensStatus)
		vElement.SetAttr("data-lens-id", lensID)
		switch lensStatus {
		case "uploading_fail":
			vElement.SetAttr("data-description", "[视频上传失败]")
		case "reviewing_fail":
			vElement.SetAttr("data-description", "[视频审核未通过]")
			vElement.SetAttr("data-lens-censor-fail-items", censorItems)
			if attr, ok := n.GetAttr("data-name"); ok {
				vElement.SetAttr("data-name", attr)
			}
		case "deleted":
			vElement.SetAttr("data-description", "[资源不存在]")
		}

		n.InsertBefore(vElement, (*node.Node)(n.FirstChild))
		n.DropTag()
		return
	}

	vElement := node.NewA()
	vElement.AddAttr("class", "video-link")
	vElement.AddAttr("data-src", n.GetAttrOrDefault("data-swfurl", ""))
	vElement.AddAttr("href", n.GetAttrOrDefault("data-sourceurl", ""))
	vElement.AddAttr("data-videoid", n.GetAttrOrDefault("id", ""))
	vElement.AddAttr("data-poster", n.GetAttrOrDefault("poster", ""))
	vElement.AddAttr("data-name", n.GetAttrOrDefault("data-name", ""))
	vElement.AddAttr("data-video-id", n.GetAttrOrDefault("data-video-id", ""))
	vElement.AddAttr("data-lens-id", n.GetAttrOrDefault("data-lens-id", ""))
	vElement.AddAttr("data-video-playable", n.GetAttrOrDefault("data-video-playable", ""))
	if censorItems != "" {
		vElement.SetAttr("data-lens-censor-fail-items", censorItems)
	}

	vElement.SetText(n.GetAttrOrDefault("data-name", ""))
	n.InsertBefore(vElement, (*node.Node)(n.FirstChild))
	n.DropTag()
}

// OnCodeTag overwrites OnCodeTag
func (z *ZEditText) OnCodeTag(n *node.Node) {
	n.SetAtom(atom.Pre)
}

// processLinkOutput overwrites processLinkOutput
func (z *ZEditText) processLinkOutput(n *node.Node) {
	if dataImage, _ := n.GetAttr("data-image"); dataImage != "" {
		dataImageSize := n.GetAttrOrDefault("data-image-size", "")
		n.SetAttr("data-image", pier.GetFullURL(dataImage, dataImageSize, "jpg", "", true))
	}

	delAttrs := make([]string, 0)
	for _, k := range n.Attr {
		key := k.Key
		_, ok := dataAttrs[key]
		if strings.HasPrefix(key, "data-") && !ok {
			delAttrs = append(delAttrs, key)
		}
	}
	n.DelAttrs(delAttrs...)
}

func (z *ZEditText) processShareLinkOutput(n *node.Node) {
	// do nothing
}
