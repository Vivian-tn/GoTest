package ztext

import (
	"fmt"

	"git.in.zhihu.com/go/ztext/escape"
	"git.in.zhihu.com/go/ztext/node"
)

type ZThumbnailText struct {
	*ZOutputText
}

var (
	boxTmpl = `<img class="thumbnail" src="%s"><span class="z-ico-play-video"></span>`
)

const (
	//	thumbnailImgWidth   = 200
	//	thumbnailImgHeight  = 112
	thumbnailImgSizeStr = "200x112"
)

// NewZThumbnail returns a instance of ZEditText
func NewZThumbnail(content string, opts ...Option) (*ZThumbnailText, error) {
	z, err := NewZOutputText(content, opts...)
	if err != nil {
		return nil, err
	}

	return &ZThumbnailText{
		ZOutputText: z,
	}, nil
}

// PrepareWalkDOM overwrites ZOutputText.PrepareWalkDOM
func (z *ZThumbnailText) PrepareWalkDOM() error {
	// do nothing
	return nil
}

// OnImgTag overwrites ZOutputText.OnImgTag
func (z *ZThumbnailText) OnImgTag(n *node.Node) {
	/*
		统一缩略图尺寸，增加相关属性。
		一般是配合 ZDom.thumbnail 一起使用
	*/
	imgSrc, _ := n.GetAttr("src")
	if imgSrc == "" {
		n.DropTag()
		return
	}
	n.SetAttrs(
		map[string]string{
			"src":           z.imageURLResolver(imgSrc, thumbnailImgSizeStr, "", "", true),
			"class":         "origin_image inline-img zh-lightbox-thumb",
			"data-original": z.imageURLResolver(imgSrc, "r", "", "", true),
		},
	)

}

// OnVideoTag is the callback function when visiting <video> tag
func (z *ZThumbnailText) OnVideoTag(n *node.Node) {
	if n.GetAttrOrDefault("data-lens-status", "") == "deleted" {
		n.DropTag()
		return
	}
	box, _ := node.NewFragmentFromString(
		fmt.Sprintf(boxTmpl, escape.XhtmlEscape(n.GetAttrOrDefault("poster", ""))),
	)
	// NewFragmentFromString will surround boxTmpl with <div> tag
	box.SetAttr("class", "video-box-thumbnail")
	n.InsertBefore((*node.Node)(box), (*node.Node)(n.FirstChild))
	n.DropTag()
}
