package zae

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type ResourceType string

func (d ResourceType) Lower() string {
	return strings.ToLower(d.String())
}

func (d ResourceType) Upper() string {
	return strings.ToUpper(d.String())
}

func (d ResourceType) String() string {
	return string(d)
}

const (
	ResourceMySQL   ResourceType = "MySQL"
	ResourceRWRedis ResourceType = "RWRedis"
	ResourceStatsd  ResourceType = "Statsd"
)

func Register(resource ResourceType, name, kind string, dsns []string) (err error) {
	name, kind = upper(name), upper(kind)
	for i, dsn := range dsns {
		var key string
		if kind == "" {
			key = fmt.Sprintf("ZAE_RES_%s_%s_%d_ADDR", resource.Upper(), name, i)
		} else {
			key = fmt.Sprintf("ZAE_RES_%s_%s_%s_%d_ADDR", resource.Upper(), name, kind, i)
		}
		if !strings.Contains(dsn, "://") {
			dsn = fmt.Sprintf("%s://%s", resource.Lower(), dsn)
		}
		if _, err = url.Parse(dsn); err != nil {
			return err
		}
		if edsn := os.Getenv(key); edsn != "" {
			return fmt.Errorf("%s.%s.%s address already exist: %v", resource, name, kind, edsn)
		}
		err = os.Setenv(key, dsn)
		if err != nil {
			return err
		}
	}
	return nil
}

func DiscoveryOne(resource ResourceType, name, kind string) (addr *url.URL, err error) {
	addrs, err := DiscoveryMany(resource, name, kind)
	if err != nil {
		return nil, err
	}
	if len(addrs) != 1 {
		return nil, fmt.Errorf("%s.%s.%s address not found: %v", resource, name, kind, addrs)
	}
	return addrs[0], nil
}

func DiscoveryMany(resource ResourceType, name, kind string) (addrs []*url.URL, err error) {
	name, kind = upper(name), upper(kind)
	for i := 0; ; i++ {
		var key string
		if kind == "" {
			key = fmt.Sprintf("ZAE_RES_%s_%s_%d_ADDR", resource.Upper(), name, i)
		} else {
			key = fmt.Sprintf("ZAE_RES_%s_%s_%s_%d_ADDR", resource.Upper(), name, kind, i)
		}
		dsn := os.Getenv(key)
		if dsn == "" {
			break
		}
		if !strings.Contains(dsn, "://") {
			dsn = fmt.Sprintf("%s://%s", resource.Lower(), dsn)
		}
		addr, err := url.Parse(dsn)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func upper(str string) string {
	return strings.ToUpper(regexp.MustCompile("[^0-9a-zA-Z]+").ReplaceAllString(str, "_"))
}
