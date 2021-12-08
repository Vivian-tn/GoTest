package zae

import (
	"log"
	"net/url"
	"os"
	"regexp"
)

var (
	isDebug                  bool
	hostname                 string
	discoveryURI             string
	host                     = os.Getenv("HOST_IP")
	env                      = os.Getenv("ZAE_ENV")
	region                   = os.Getenv("ZAE_REGION_NAME")
	app                      = os.Getenv("ZAE_APP_NAME")
	appToken                 = os.Getenv("ZAE_SEC_APP_TOKEN")
	service                  = os.Getenv("ZAE_UNIT_NAME")
	sentryDSN, sentryRelease = os.Getenv("SENTRY_DSN"), os.Getenv("SENTRY_RELEASE")
)

// http://wiki.in.zhihu.com/pages/viewpage.action?pageId=11832684

const (
	EnvProd    = "production" // 生产环境
	EnvTesting = "testing"    // 测试环境
	EnvCI      = "test"       // CI 环境
	EnvDevelop = "develop"    // 开发环境
)

func init() {
	var err error

	{
		hostname, err = os.Hostname()
		if err != nil {
			panic(err)
		}
		hostname = regexp.MustCompile("[^0-9a-zA-Z]+").ReplaceAllString(hostname, "_")
	}

	{
		discoveryURI = os.Getenv("SERVICE_DISCOVERY_URI")
		if discoveryURI == "" {
			discoveryURI = "127.0.0.1:8500"
		} else {
			parsedURL, err := url.Parse(discoveryURI)
			if err != nil || parsedURL.Scheme != "consul" {
				log.Fatalf("malformed discovery url: %s", discoveryURI)
			}
			discoveryURI = parsedURL.Host
		}
	}
}

func Hostname() string {
	return hostname
}

func App() string {
	if IsDevelopEnv() && app == "" {
		return "develop-app"
	}
	return app
}

func AppToken() string {
	return appToken
}

func Service() string {
	if (IsDevelopEnv() || IsCIEnv()) && service == "" {
		return "develop-service"
	}
	return service
}

func Region() string {
	return region
}

func SetRegion(s string) {
	region = s
}

func DiscoveryURI() string {
	return discoveryURI
}

func Sentry() (string, string) {
	return sentryDSN, sentryRelease
}

func Environment() string {
	if env == "" {
		return EnvDevelop
	}
	return env
}

func IsProdEnv() bool {
	return Environment() == EnvProd
}

func IsTestingEnv() bool {
	return Environment() == EnvTesting
}

func IsCIEnv() bool {
	return Environment() == EnvCI
}

func IsDevelopEnv() bool {
	return Environment() == EnvDevelop
}

func Host() string {
	return host
}

func IsDebug() bool {
	return isDebug
}

func EnableDebug() {
	isDebug = true
}
