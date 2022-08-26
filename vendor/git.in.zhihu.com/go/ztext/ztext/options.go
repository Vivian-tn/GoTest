package ztext

type protocolType int64

const (
	// HTTPS Type for Protocol
	HTTPS protocolType = iota
	// HTTP Type for Protocol
	HTTP
	// Relative Type for Protocol
	Relative
)

// Option for zoutputtext config
type Option func(*ZOutputText)

// Protocol sets Protocol used, including "http", "https", and "relative".
// Default is https, and illegal protocol is also https.
func Protocol(protocol protocolType) Option {
	return Option(func(c *ZOutputText) {
		if p, ok := ProtocolMap[protocol]; ok {
			c.Protocol = p
		} else {
			c.Protocol = ProtocolMap[HTTPS]
		}
	})
}

// Watermark sets watermark style, optoins are "original", "watermark",
// and "private_watermark", default is empty string.
func Watermark(watermark string) Option {
	return Option(func(c *ZOutputText) {
		c.watermark = watermark
	})
}

// AllowGif determines if allow gif, default is false
func AllowGif(allow bool) Option {
	return Option(func(c *ZOutputText) {
		c.allowGif = allow
	})
}

// ImageURLResolver sets IMG_URL_RESOLVER function, default is pier.GetFullURL
func ImageURLResolver(resolver ImageURLResolverFunc) Option {
	return Option(func(c *ZOutputText) {
		c.imageURLResolver = resolver
	})
}

// StrictBMP determines whether enable strict BMP, default is true
func StrictBMP(strict bool) Option {
	return Option(func(c *ZOutputText) {
		c.strictBMP = strict
	})
}

// RootDomain sets root domain of site, default is zhihu.com
func RootDomain(domain string) Option {
	return Option(func(c *ZOutputText) {
		c.RootDomain = domain
	})
}
