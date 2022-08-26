package pier

import (
	"math/rand"
	"strconv"
	"time"
)

type DomainScheduler interface {
	GetDomain(domains []string, token string) string
	SetDomains(domains []string)
	Method() string
}

type ModulusScheduler struct {
	Domains []string
}

func (m *ModulusScheduler) GetDomain(domains []string, token string) string {
	if len(token) > 5 {
		v, err := strconv.ParseInt(token[len(token)-5:], 16, 64)
		if err == nil {
			return domains[int(v) % len(domains)]
		}
	}

   return domains[rand.Intn(len(domains))]
}

func (m *ModulusScheduler) SetDomains(domains []string) {
	m.Domains = domains
}

func (m *ModulusScheduler) Method() string {
	return "modulus_hash"
}

type WeightRandomScheduler struct {
	Domains []string
}

func (m *WeightRandomScheduler) GetDomain(domains []string, token string) string {
	if len(m.Domains) > 0 {
		return m.Domains[rand.Intn(len(m.Domains))]
	}

	return domains[rand.Intn(len(domains))]
}

func (m *WeightRandomScheduler) SetDomains(domains []string) {
	m.Domains = domains
}

func (m *WeightRandomScheduler) Method() string {
	return "weight_random"
}

func getScheduler(method string) DomainScheduler {
	if method == "weight_random" {
		return &WeightRandomScheduler{}
	}
	return &ModulusScheduler{}
}

func init() {
	rand.Seed(time.Now().Unix())
}
