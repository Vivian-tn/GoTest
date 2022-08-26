package tzone

import "time"

type Option func(*Client)

func TargetName(name string) Option {
	return Option(func(c *Client) {
		c.targetName = name
	})
}

func Timeout(t time.Duration) Option {
	return func(c *Client) {
		c.timeout = t
	}
}

func HostPort(host string, port string) Option {
	return func(c *Client) {
		c.hostPort = host + ":" + port
	}
}
