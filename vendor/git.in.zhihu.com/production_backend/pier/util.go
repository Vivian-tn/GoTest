package pier

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"git.in.zhihu.com/go/base/telemetry/statsd"
)

var (
	PierStatsd, _   = statsd.New("pico3")
	SafePierVersion = getSafeVersion()
	ZaeAppName      = getZaeAppName()
	SafeZaeAppName  = SafeColumn(ZaeAppName)
)

const (
	DefaultFmt  = "jpg"
	DefaultSpec = "r"

	ImageContentLimit = 32

	productionEnvMark = "production"
)

func getSafeVersion() string {
	fmtVersion := strings.Replace(PierVersion, ".", "_", -1)
	return fmtVersion
}

func isURL(token string) bool {
	return strings.HasPrefix(token, "http://") || strings.HasPrefix(token, "https://") || strings.HasPrefix(token, "//")
}

func parseToken(token string) (string, string) {
	if strings.Contains(token, "?") {
		parts := strings.SplitN(token, "?", 2)
		token = parts[0]
	}

	rawToken, fmt := token, "jpg"

	if strings.Contains(token, ".") {
		parts := strings.SplitN(token, ".", 2)
		rawToken, fmt = parts[0], parts[1]
	}

	if strings.Contains(rawToken, "_") {
		parts := strings.SplitN(token, "_", 2)
		rawToken = parts[0]
	}

	return rawToken, fmt
}

func md5SumImage(plaintext []byte) string {
	m := md5.New()
	m.Write(plaintext)
	return hex.EncodeToString(m.Sum(nil))
}

func formatToken(md5Sum string) string {
	return "v2-" + md5Sum
}

func getValidSpec(specMapping map[string]string, spec string) string {
	if _, ok := ImageSpecs[spec]; !ok {
		if _, ok := WaterMarks[spec]; !ok {
			return ""
		}
	}

	downgradeSpec, ok := specMapping[spec]
	if !ok {
		return spec
	}
	return downgradeSpec
}

func getValidQuality(qualityMapping map[int16]int16, qualityLiteralValue string) string {
	quality, err := strconv.Atoi(qualityLiteralValue)
	if err != nil {
		return ""
	}

	downgradeQuality, ok := qualityMapping[int16(quality)]
	if !ok {
		return ""
	}
	return strconv.Itoa(int(downgradeQuality))
}

func sniffMIMEAndFormat(content []byte) (string, string) {
	var (
		mime        = ""
		imageFormat = ""
	)

	if len(content) > ImageContentLimit {
		content = content[:ImageContentLimit]
	}

	for fileHeader, fileFormat := range ImageHeaderFormatMapping {
		header := hex.EncodeToString(content)
		if len(header) > len(fileHeader) {
			header = header[:len(fileHeader)]
		}

		if fileHeader == strings.ToUpper(header) {
			return ImageFormatMIMEMapping[fileFormat], fileFormat
		}
	}
	return mime, imageFormat
}

func getZaeAppName() string {
	zaeAppName := os.Getenv("ZAE_APP_NAME")
	if zaeAppName == "" {
		zaeAppName = os.Args[0]
	}
	return zaeAppName
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

func MaskAppName(appName string) string {
	m := md5.New()
	m.Write([]byte(appName))
	return hex.EncodeToString(m.Sum(nil))[:8]
}

func SafeColumn(arg string) string {
	formatArg := strings.Replace(arg, ".", "_", -1)
	formatArg = strings.Replace(formatArg, ":", "_", -1)
	formatArg = strings.Replace(formatArg, "|", "_", -1)
	return formatArg
}

func IsProductionEnv() bool {
	envVal := os.Getenv("ZAE_ENV")
	return envVal == productionEnvMark
}

func StatsdGetFullUrl(domain string, quality int, spec string, format string) {
	formatDomain := SafeColumn(domain)

	var strQuality string
	if quality == 0 {
		strQuality = "default"
	} else {
		strQuality = strconv.Itoa(quality)
	}
	if spec == "" {
		spec = "default"
	}
	if format == "" {
		format = "default"
	}

	PierStatsd.Increment(fmt.Sprintf("pier-go.%s.%s.%s.%s.%s.count", SafeZaeAppName, formatDomain, strQuality, spec, format))

}
