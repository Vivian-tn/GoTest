package pier

var (
	FormatSet = map[string]struct{}{
		"jpg": struct{}{},
		"png": struct{}{},
		"gif": struct{}{},
		"webp": struct{}{},
		"jpeg": struct{}{},
	}

	ImageHeaderFormatMapping = map[string]string{
		"FFD8FF": "jpeg",
		"89504E470D0A1A0A": "png",
		"47494638": "gif",
		"52494646": "webp",
	}

	ImageFormatMIMEMapping = map[string]string{
		"jpeg": "image/jpeg",
		"png": "image/png",
		"gif": "image/gif",
		"webp": "image/webp",
	}
)

type WaterMarkObject struct {
	Image string
	Object string
	Watermark int64
}

var WaterMarks = map[string]WaterMarkObject{
	"adx1": WaterMarkObject{
		Image: "v2-0455aaa2cec0016cc2d2f036b910b0c5?x-oss-process=image/resize,P_15",
		Object: "v2-0455aaa2cec0016cc2d2f036b910b0c5@15P",
		Watermark: 1,
	},
	"adx4": WaterMarkObject{
		Image: "v2-088d3bff72fc2fbc63cc00e97a76f84e?x-oss-process=image/resize,P_15",
		Object: "v2-088d3bff72fc2fbc63cc00e97a76f84e@15P",
		Watermark: 1,
	},
	"adx8": WaterMarkObject{
		Image: "v2-a6770a2acdca232c38b742e492fa5007?x-oss-process=image/resize,P_15",
		Object: "v2-a6770a2acdca232c38b742e492fa5007@15P",
		Watermark: 1,
	},
}

type ImageSpecObject struct {
	Width int64
	Height int64
	Quality int64
	Mode string
}

var ImageSpecs = map[string]ImageSpecObject{
	"112x74": ImageSpecObject{
		Width: 112,
		Quality: 90,
		Mode: "crop",
		Height: 74,
	},
	"600x500": ImageSpecObject{
		Width: 600,
		Quality: 98,
		Mode: "crop",
		Height: 500,
	},
	"540x300": ImageSpecObject{
		Width: 540,
		Quality: 98,
		Mode: "crop",
		Height: 300,
	},
	"xl": ImageSpecObject{
		Width: 200,
		Quality: 98,
		Mode: "crop",
		Height: 200,
	},
	"720w": ImageSpecObject{
		Width: 720,
		Quality: 90,
		Mode: "crop",
		Height: 0,
	},
	"is": ImageSpecObject{
		Width: 34,
		Quality: 98,
		Mode: "crop",
		Height: 34,
	},
	"400x224": ImageSpecObject{
		Width: 400,
		Quality: 98,
		Mode: "crop",
		Height: 224,
	},
	"370x100": ImageSpecObject{
		Width: 370,
		Quality: 90,
		Mode: "crop",
		Height: 100,
	},
	"294x245": ImageSpecObject{
		Width: 294,
		Quality: 90,
		Mode: "crop",
		Height: 245,
	},
	"740x200": ImageSpecObject{
		Width: 740,
		Quality: 90,
		Mode: "crop",
		Height: 200,
	},
	"300x250": ImageSpecObject{
		Width: 300,
		Quality: 98,
		Mode: "crop",
		Height: 250,
	},
	"485x275": ImageSpecObject{
		Width: 485,
		Quality: 98,
		Mode: "crop",
		Height: 275,
	},
	"180x120": ImageSpecObject{
		Width: 180,
		Quality: 90,
		Mode: "crop",
		Height: 120,
	},
	"64x64": ImageSpecObject{
		Width: 64,
		Quality: 98,
		Mode: "crop",
		Height: 64,
	},
	"im": ImageSpecObject{
		Width: 68,
		Quality: 98,
		Mode: "crop",
		Height: 68,
	},
	"380x214": ImageSpecObject{
		Width: 380,
		Quality: 98,
		Mode: "crop",
		Height: 214,
	},
	"450x300": ImageSpecObject{
		Width: 450,
		Quality: 98,
		Mode: "crop",
		Height: 300,
	},
	"xdpi": ImageSpecObject{
		Width: 720,
		Quality: 98,
		Mode: "crop",
		Height: 1080,
	},
	"xs": ImageSpecObject{
		Width: 50,
		Quality: 98,
		Mode: "crop",
		Height: 50,
	},
	"150x50": ImageSpecObject{
		Width: 150,
		Quality: 90,
		Mode: "crop",
		Height: 50,
	},
	"720x4096": ImageSpecObject{
		Width: 720,
		Quality: 90,
		Mode: "lfit",
		Height: 4096,
	},
	"1168x300": ImageSpecObject{
		Width: 1168,
		Quality: 98,
		Mode: "crop",
		Height: 300,
	},
	"640w": ImageSpecObject{
		Width: 640,
		Quality: 90,
		Mode: "crop",
		Height: 0,
	},
	"588x490": ImageSpecObject{
		Width: 588,
		Quality: 90,
		Mode: "crop",
		Height: 490,
	},
	"324x182": ImageSpecObject{
		Width: 324,
		Quality: 98,
		Mode: "crop",
		Height: 182,
	},
	"547x308": ImageSpecObject{
		Width: 547,
		Quality: 98,
		Mode: "crop",
		Height: 308,
	},
	"250x0": ImageSpecObject{
		Width: 250,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"300x100": ImageSpecObject{
		Width: 300,
		Quality: 90,
		Mode: "crop",
		Height: 100,
	},
	"bl": ImageSpecObject{
		Width: 640,
		Quality: 98,
		Mode: "crop",
		Height: 260,
	},
	"750w": ImageSpecObject{
		Width: 750,
		Quality: 90,
		Mode: "crop",
		Height: 0,
	},
	"hd": ImageSpecObject{
		Width: 720,
		Quality: 98,
		Mode: "lfit",
		Height: 4096,
	},
	"540x450": ImageSpecObject{
		Width: 540,
		Quality: 98,
		Mode: "crop",
		Height: 450,
	},
	"100w": ImageSpecObject{
		Width: 100,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"250x250": ImageSpecObject{
		Width: 250,
		Quality: 98,
		Mode: "crop",
		Height: 250,
	},
	"220x110": ImageSpecObject{
		Width: 220,
		Quality: 98,
		Mode: "crop",
		Height: 110,
	},
	"600x150": ImageSpecObject{
		Width: 600,
		Quality: 98,
		Mode: "crop",
		Height: 150,
	},
	"60w": ImageSpecObject{
		Width: 60,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"hdpi": ImageSpecObject{
		Width: 540,
		Quality: 98,
		Mode: "crop",
		Height: 960,
	},
	"1200x300": ImageSpecObject{
		Width: 1200,
		Quality: 98,
		Mode: "crop",
		Height: 300,
	},
	"600x250": ImageSpecObject{
		Width: 600,
		Quality: 98,
		Mode: "crop",
		Height: 250,
	},
	"200x112": ImageSpecObject{
		Width: 200,
		Quality: 98,
		Mode: "crop",
		Height: 112,
	},
	"xxdpi": ImageSpecObject{
		Width: 1080,
		Quality: 98,
		Mode: "crop",
		Height: 1920,
	},
	"bm": ImageSpecObject{
		Width: 640,
		Quality: 98,
		Mode: "crop",
		Height: 960,
	},
	"bh": ImageSpecObject{
		Width: 640,
		Quality: 98,
		Mode: "crop",
		Height: 320,
	},
	"280x280": ImageSpecObject{
		Width: 280,
		Quality: 100,
		Mode: "crop",
		Height: 280,
	},
	"b2m": ImageSpecObject{
		Width: 640,
		Quality: 98,
		Mode: "crop",
		Height: 1136,
	},
	"b": ImageSpecObject{
		Width: 600,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"qhd": ImageSpecObject{
		Width: 480,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"120x160": ImageSpecObject{
		Width: 120,
		Quality: 90,
		Mode: "crop",
		Height: 160,
	},
	"ipico": ImageSpecObject{
		Width: 120,
		Quality: 100,
		Mode: "crop",
		Height: 120,
	},
	"200x0": ImageSpecObject{
		Width: 200,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"224x148": ImageSpecObject{
		Width: 224,
		Quality: 90,
		Mode: "crop",
		Height: 148,
	},
	"270x150": ImageSpecObject{
		Width: 270,
		Quality: 98,
		Mode: "crop",
		Height: 150,
	},
	"l": ImageSpecObject{
		Width: 100,
		Quality: 98,
		Mode: "crop",
		Height: 100,
	},
	"xll": ImageSpecObject{
		Width: 400,
		Quality: 98,
		Mode: "crop",
		Height: 400,
	},
	"mdpi": ImageSpecObject{
		Width: 360,
		Quality: 98,
		Mode: "crop",
		Height: 640,
	},
	"270x225": ImageSpecObject{
		Width: 270,
		Quality: 98,
		Mode: "crop",
		Height: 225,
	},
	"584x150": ImageSpecObject{
		Width: 584,
		Quality: 98,
		Mode: "crop",
		Height: 150,
	},
	"fhd": ImageSpecObject{
		Width: 1080,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"1440w": ImageSpecObject{
		Width: 1440,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"xld": ImageSpecObject{
		Width: 200,
		Quality: 98,
		Mode: "crop",
		Height: 0,
	},
	"m": ImageSpecObject{
		Width: 75,
		Quality: 98,
		Mode: "crop",
		Height: 75,
	},
	"mh": ImageSpecObject{
		Width: 320,
		Quality: 98,
		Mode: "crop",
		Height: 160,
	},
	"bs": ImageSpecObject{
		Width: 640,
		Quality: 98,
		Mode: "crop",
		Height: 640,
	},
	"96x96": ImageSpecObject{
		Width: 96,
		Quality: 98,
		Mode: "crop",
		Height: 96,
	},
	"mt": ImageSpecObject{
		Width: 216,
		Quality: 98,
		Mode: "crop",
		Height: 160,
	},
	"r": ImageSpecObject{
		Width: 0,
		Quality: 100,
		Mode: "crop",
		Height: 0,
	},
	"1080x1920": ImageSpecObject{
		Width: 1080,
		Quality: 98,
		Mode: "crop",
		Height: 1920,
	},
	"t": ImageSpecObject{
		Width: 150,
		Quality: 70,
		Mode: "crop",
		Height: 150,
	},
	"ms": ImageSpecObject{
		Width: 320,
		Quality: 98,
		Mode: "crop",
		Height: 320,
	},
	"1200x500": ImageSpecObject{
		Width: 1200,
		Quality: 98,
		Mode: "crop",
		Height: 500,
	},
	"s": ImageSpecObject{
		Width: 25,
		Quality: 98,
		Mode: "crop",
		Height: 25,
	},
	"970x550": ImageSpecObject{
		Width: 970,
		Quality: 98,
		Mode: "crop",
		Height: 550,
	},
	"1280x640": ImageSpecObject{
		Width: 1280,
		Quality: 80,
		Mode: "crop",
		Height: 640,
	},
}
