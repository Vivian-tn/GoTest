package sentry

import (
	"runtime"
	"strings"

	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"github.com/getsentry/raven-go"
)

var includePaths []string

func SetIncludePaths(p []string) {
	includePaths = p
}

type Error struct {
	message string
	frames  []runtime.Frame
}

func (r Error) Error() string {
	return r.message
}

func (r Error) Caller() *runtime.Frame {
	for _, fr := range r.frames {
		packageName, _ := halo.ParsePackage(fr.Function)
		for _, prefix := range includePaths {
			if strings.HasPrefix(packageName, prefix) && !strings.Contains(packageName, "vendor") && !strings.Contains(packageName, "third_party") {
				return &fr
			}
		}
	}
	return &r.frames[0]
}

func (r Error) Stacktrace() *raven.Stacktrace {
	frames := make([]*raven.StacktraceFrame, 0, len(r.frames))
	for _, fr := range r.frames {
		frame := raven.NewStacktraceFrame(
			fr.PC, fr.Function, fr.File, fr.Line, 3, includePaths,
		)
		if frame != nil {
			frames = append(frames, frame)
		}
	}
	return &raven.Stacktrace{Frames: frames}
}
