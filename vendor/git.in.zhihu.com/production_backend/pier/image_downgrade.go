package pier

type ImageQualityType struct {
	LowQuality int
	NormalQuality int
	HighQuality int
	LowLimit int
	Interval int
	LowQualityLiteralValue string
	NormalQualityLiteralValue string
	HighQualityLiteralValue string
}

var (
	ImageQuality = ImageQualityType{
		HighQuality: 100,
		NormalQuality: 80,
		LowQuality: 50,
		LowLimit: 1,
		Interval: 10,
		LowQualityLiteralValue: "50",
		NormalQualityLiteralValue: "80",
		HighQualityLiteralValue: "100",
	}

	// 更新后严格对照规格，防止写错 http://origin1.zhimg.com/api/v1/specs
	DowngradeSpecs = map[string]string{
		"112x74": "bh",
		"540x300": "bh",
		"64x64": "l",
		"600x250": "bh",
		"280x280": "xl",
		"b": "720w",
		"224x148": "bh",
		"96x96": "l",
		"1080x1920": "xxdpi",
		"370x100": "bh",
		"485x275": "bh",
		"380x214": "bh",
		"150x50": "bh",
		"324x182": "bh",
		"220x110": "bh",
		"200x112": "bh",
		"270x150": "bh",
		"mh": "bh",
		"t": "xl",
		"600x500": "720w",
		"294x245": "200x0",
		"740x200": "720w",
		"300x250": "200x0",
		"450x300": "qhd",
		"xdpi": "720w",
		"720x4096": "hd",
		"588x490": "720w",
		"547x308": "bh",
		"300x100": "200x0",
		"750w": "720w",
		"540x450": "720w",
		"600x150": "720w",
		"hdpi": "720w",
		"b2m": "720w",
		"120x160": "200x0",
		"xll": "qhd",
		"mdpi": "qhd",
		"270x225": "200x0",
		"584x150": "720w",
		"bs": "720w",
		"mt": "200x0",
		"ms": "qhd",
		"970x550": "720w",
		"xs": "l",
		"im": "l",
		"1168x300": "1280x640",
		"1200x500": "1280x640",
		"xld": "200x0",
		"ipico": "l",
		"s": "l",
		"250x0": "200x0",
		"250x250": "xl",
		"400x224": "bh",
		"60w": "200x0",
		"100w": "200x0",
		"is": "l",
		"m": "l",
		"640w": "720w",
		"180x120": "bh",
		"bl": "bh",
		"1200x300": "1280x640",
		"fhd": "1440w",
	}
)
