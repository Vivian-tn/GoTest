package diplomat

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.in.zhihu.com/go/base/internal/request"
	"git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/base/zae"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var ErrServiceEntryNotFound = errors.New("service entry not found")

type Strategy int
type Consistency string

const (
	Random Strategy = iota
	Roundrobin
)

// ServiceEntry represent an available service info registered on consul.
type ServiceEntry struct {
	Name string
	Host string
	Port int
}

func (s *ServiceEntry) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Diplomat can get target info on consul.
type Diplomat struct {
	cacheTTL   time.Duration
	backoffTTL time.Duration
	cache      *Cache
	req        *request.Request
	counts     sync.Map
	group      singleflight.Group
}

// New return a new Diplomat instance pointer.
//
// Param: address should be in format: host:port
func New(addr string, cacheTTL time.Duration) *Diplomat {
	return &Diplomat{
		cacheTTL:   cacheTTL,
		backoffTTL: 2 * time.Second,
		cache:      NewCache(),
		req:        request.New().SetBaseURL(fmt.Sprintf("http://%s", addr)),
	}
}

// Discover automatically find a proper consul address and make a Diplomat with it.
// This method is recommended to use instead of call New manually.
// Panic on malformed url found.
func Discover() *Diplomat {
	return New(zae.DiscoveryURI(), 30*time.Second)
}

func (d *Diplomat) RegisterLocal(service, address string) {
	_ = os.Setenv(fmt.Sprintf("%s_SERVICE_ADDRESS", upper(service)), address)
}

func (d *Diplomat) FindLocal(service string) ([]*ServiceEntry, error) {
	localAddr := os.Getenv(fmt.Sprintf("%s_SERVICE_ADDRESS", upper(service)))
	if localAddr != "" {
		splitedAddr := strings.Split(localAddr, ":")
		port, err := strconv.Atoi(splitedAddr[1])
		if err != nil {
			return nil, err
		}
		return []*ServiceEntry{
			{
				Name: service,
				Host: splitedAddr[0],
				Port: port,
			},
		}, nil
	}
	return nil, nil
}

// Find find a list of ServiceEntry with specified name.
// Use this method unless you know what you want, otherwise you
// should always call Select to find a proper ServiceEntry info.
func (d *Diplomat) Find(ctx context.Context, service string, refresh bool) ([]*ServiceEntry, error) {
	entries, err := d.FindLocal(service)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return d.healthWithCache(ctx, service, refresh)
	}
	return entries, nil
}

// Select is a wrapper for Find with a select strategy to help you get a proper ServiceEntry.
//
// Return nil if no address available.
func (d *Diplomat) Select(ctx context.Context, name string, strategy Strategy, refresh bool) (*ServiceEntry, error) {
	entries, err := d.Find(ctx, name, refresh)
	if len(entries) == 0 && err != nil {
		return nil, err
	}
	if err != nil {
		log.Errorf(ctx, "consul request failed %s", err)
	}
	if len(entries) == 0 {
		return nil, ErrServiceEntryNotFound
	}

	switch strategy {
	case Roundrobin:
		count := 1
		if v, ok := d.counts.LoadOrStore(name, 1); ok {
			count = v.(int) + 1
			d.counts.Store(name, count)
		}
		idx := count % len(entries)
		return entries[idx], nil
	default:
		return entries[rand.Intn(len(entries))], nil
	}
}

func (d *Diplomat) Discard(discarded *ServiceEntry) {
	d.cache.Discard(discarded.Name, discarded)
}

func (d *Diplomat) healthWithCache(ctx context.Context, name string, refresh bool) (entries []*ServiceEntry, err error) {
	var expired bool
	if refresh {
		expired = true
	} else {
		entries, expired = d.cache.Get(name)
	}
	if expired || len(entries) == 0 {
		newEntries, err := d.health(ctx, name)
		if err != nil {
			d.cache.Expire(name, d.backoffTTL)
			return entries, err
		}
		if d.cacheTTL > 0 {
			d.cache.Store(name, newEntries, d.cacheTTL)
		}
		entries = newEntries
	}
	return entries, err
}

type entryResponse struct {
	Node struct {
		ID              string `json:"ID"`
		Node            string `json:"Node"`
		Address         string `json:"Address"`
		Datacenter      string `json:"Datacenter"`
		TaggedAddresses struct {
			Lan string `json:"lan"`
			Wan string `json:"wan"`
		} `json:"TaggedAddresses"`
		Meta struct {
			ConsulNetworkSegment string `json:"consul-network-segment"`
		} `json:"Meta"`
		CreateIndex int `json:"CreateIndex"`
		ModifyIndex int `json:"ModifyIndex"`
	} `json:"Node"`
	Service struct {
		ID      string      `json:"ID"`
		Service string      `json:"Service"`
		Tags    []string    `json:"Tags"`
		Address string      `json:"Address"`
		Meta    interface{} `json:"Meta"`
		Port    int         `json:"Port"`
		Weights struct {
			Passing int `json:"Passing"`
			Warning int `json:"Warning"`
		} `json:"Weights"`
		EnableTagOverride bool `json:"EnableTagOverride"`
		Proxy             struct {
			MeshGateway struct {
			} `json:"MeshGateway"`
			Expose struct {
			} `json:"Expose"`
		} `json:"Proxy"`
		Connect struct {
		} `json:"Connect"`
		CreateIndex int `json:"CreateIndex"`
		ModifyIndex int `json:"ModifyIndex"`
	} `json:"Service"`
	Checks []struct {
		Node        string        `json:"Node"`
		CheckID     string        `json:"CheckID"`
		Name        string        `json:"Name"`
		Status      string        `json:"Status"`
		Notes       string        `json:"Notes"`
		Output      string        `json:"Output"`
		ServiceID   string        `json:"ServiceID"`
		ServiceName string        `json:"ServiceName"`
		ServiceTags []interface{} `json:"ServiceTags"`
		Type        string        `json:"Type"`
		Definition  struct {
		} `json:"Definition"`
		CreateIndex int `json:"CreateIndex"`
		ModifyIndex int `json:"ModifyIndex"`
	} `json:"Checks"`
}

func (d *Diplomat) health(ctx context.Context, name string) (entries []*ServiceEntry, err error) {
	ch := d.group.DoChan(name, func() (interface{}, error) {
		return d.healthOnce(ctx, name)
	})
	select {
	case res := <-ch:
		if res.Err != nil {
			return nil, res.Err
		}
		return res.Val.([]*ServiceEntry), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (d *Diplomat) healthOnce(ctx context.Context, name string) (entries []*ServiceEntry, err error) {
	resp, err := d.req.Get(ctx, fmt.Sprintf("/v1/health/service/%s", name), request.Query{
		"stale":   "",
		"passing": "true",
	}, request.Headers{
		"x-telemtry-service": "consul",
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response code: %s %d %s", resp.Request.URL, resp.StatusCode, resp.String())
	}

	entriesResp := make([]*entryResponse, 0)
	err = resp.ToJSON(&entriesResp)
	if err != nil {
		return nil, err
	}
	for _, entryResp := range entriesResp {
		entries = append(entries, &ServiceEntry{
			Name: name,
			Host: entryResp.Service.Address,
			Port: entryResp.Service.Port,
		})
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("unexpected response body: %s %d %s", resp.Request.URL, resp.StatusCode, resp.String())
	}

	return entries, nil
}

func upper(str string) string {
	return strings.ToUpper(regexp.MustCompile("[^0-9a-zA-Z]+").ReplaceAllString(str, "_"))
}
