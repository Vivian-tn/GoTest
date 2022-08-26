package pier

import (
	"sort"

	"github.com/hashicorp/golang-lru"
)

type CDNDomainParameter struct {
	Weight int16
	Regions []string
}

type PierStrategy struct {
	VersionID int64
	DefaultSpec string
	DefaultQuality int16
	CDNDomains map[string]CDNDomainParameter
	SpecMapping map[string]string
	QualityMapping map[int16]int16
}

const (
	maxStrategyCacheSize   = (1 << 10) * 4
	strategyVersionIDKey = "strategyVersionID"

	specOf1280x640 = "1280x640"
)

var (
	strategyLruCache, _ = lru.New(maxStrategyCacheSize)

	pierVersion = "v1_3_13"

	defaultStrategy = PierStrategy{
		VersionID: 1248670622311821312,
		DefaultSpec: "720w",
		DefaultQuality: 75,
		SpecMapping: map[string]string{
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
			"1168x300": specOf1280x640,
			"1200x500": specOf1280x640,
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
			"1200x300": specOf1280x640,
			"fhd": "1440w",
		},
		QualityMapping: map[int16]int16{
			10: 50,
			20: 50,
			30: 50,
			40: 50,
			50: 50,
			60: 80,
			70: 80,
			80: 80,
			90: 80,
		},
		CDNDomains: map[string]CDNDomainParameter{
			"pic1.zhimg.com": CDNDomainParameter{
				Weight: 10,
				Regions: []string{},
			},
			"pic2.zhimg.com": CDNDomainParameter{
				Weight: 10,
				Regions: []string{},
			},
			"pic3.zhimg.com": CDNDomainParameter{
				Weight: 10,
				Regions: []string{},
			},
			"pic4.zhimg.com": CDNDomainParameter{
				Weight: 10,
				Regions: []string{},
			},
		},
	}
)

func NewPierStrategy() *PierStrategy {
	return &defaultStrategy
}

func (s *PierStrategy) getVersionID() int64 {
	versionID, ok := strategyLruCache.Get(strategyVersionIDKey)
	if !ok {
		return s.VersionID
	}
	return versionID.(int64)
}

func (s *PierStrategy) getStrategy() *PierStrategy {
	strategy, ok := strategyLruCache.Get(s.getVersionID())
	if !ok {
		return s
	}
	return strategy.(*PierStrategy)
}

func (s *PierStrategy) getCDNDomains() []string {
	domains, strategy := []string{}, s.getStrategy()
	for domain, parameter := range strategy.CDNDomains {
		for i:=0; i<int(parameter.Weight); i++ {
			domains = append(domains, domain)
		}
	}
	return domains
}

func GetCDNDomainsByLRU() []string {
	versionID, ok := strategyLruCache.Get(strategyVersionIDKey)
	if !ok {
		versionID = defaultStrategy.VersionID
	}

	cdnDomains := defaultStrategy.CDNDomains
	strategy, ok := strategyLruCache.Get(versionID)
	if ok {
		cdnDomains = strategy.(*PierStrategy).CDNDomains
	}

	domains := []string{}
	for domain, _ := range cdnDomains {
		domains = append(domains, domain)
	}
	sort.Strings(domains)
	return domains
}
