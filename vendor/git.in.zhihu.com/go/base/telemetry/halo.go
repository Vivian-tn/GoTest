package telemetry

import (
	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"git.in.zhihu.com/go/base/telemetry/statsd"
)

var (
	globalHaloClient statsd.Client
)

func init() {
	var err error
	globalHaloClient, err = statsd.New("span")
	if err != nil {
		panic(err)
	}
}

type Error = halo.Error

func NewClassErr(class string) Error {
	return halo.NewClassErr(class)
}

func WrapErr(err error, class string) Error {
	return halo.WrapErr(err, class)
}

func WrapErrWithUnknownClass(err error) Error {
	return halo.WrapErrWithUnknownClass(err)
}

func WrapErrWithStack(err error) Error {
	return halo.WrapErrWithStack(err)
}
