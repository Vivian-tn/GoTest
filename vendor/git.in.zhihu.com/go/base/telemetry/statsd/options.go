package statsd

type options struct {
	addrs      []string
	globalTags map[string]string
	metadata   bool
}

type option func(o *options)

func WithAddrs(addrs ...string) option {
	return func(o *options) {
		o.addrs = addrs
	}
}

func WithGlobalTags(tags map[string]string) option {
	return func(o *options) {
		o.globalTags = tags
	}
}

func WithMetadata() option {
	return func(o *options) {
		o.metadata = true
	}
}

func WithoutMetadata() option {
	return func(o *options) {
		o.metadata = false
	}
}
