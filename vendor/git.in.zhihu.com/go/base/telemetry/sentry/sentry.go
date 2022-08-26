package sentry

import (
	"context"
	"errors"
	"fmt"

	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"git.in.zhihu.com/go/base/telemetry/log"
)

type Tags = log.Fields

// SetIncludePaths
func SetIncludePaths(p []string) {
	log.SetIncludePaths(p)
}

func Recover(f func(), r func(error)) {
	defer func() {
		var rvalErr error
		err := recover()
		switch rval := err.(type) {
		case nil:
			return
		case error:
			rvalErr = rval
		default:
			rvalErr = errors.New(fmt.Sprint(rval))
		}
		r(halo.WrapErrWithStack(rvalErr))
	}()

	f()
}

func CapturePanic(ctx context.Context, tags Tags, f func()) {
	Recover(f, func(err error) {
		log.WithFields(ctx, tags).Error(err)
	})
}
