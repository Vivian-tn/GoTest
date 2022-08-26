package telemetry

import (
	"context"
	"fmt"
	"time"

	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/base/zae"
	"github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
)

// https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/trace/semantic_conventions/database.md

type DatabaseSystem = string

const (
	DatabaseMySQL  DatabaseSystem = "MySQL"
	DatabaseSQLite DatabaseSystem = "SQLite"
	DatabaseRedis  DatabaseSystem = "Redis"
)

const (
	databaseSegmentKey ctxKey = "database-segment-key"
)

func DatabaseSegmentFromContext(ctx context.Context) *DatabaseSegment {
	val := ctx.Value(databaseSegmentKey)
	if sp, ok := val.(*DatabaseSegment); ok {
		return sp
	}
	return nil
}

func StartDatabaseSegment(ctx context.Context, ds *DatabaseSegment) (*DatabaseSegment, context.Context, error) {
	openSpan, ctx := StartChildSpanWithContext(ctx, "")
	ds.openSpan = openSpan
	ds.start = time.Now()
	ctx = context.WithValue(ctx, databaseSegmentKey, ds)
	return ds, ctx, nil
}

type DatabaseSegment struct {
	openSpan      opentracing.Span
	haloSpan      *halo.ClientSpan
	start         time.Time
	SlowThreshold time.Duration
	System        DatabaseSystem
	Host          string
	PortOrPath    string
	DatabaseName  string
	Collection    string
	Operation     string
	Query         string
	Parameters    map[string]interface{}
}

func (s *DatabaseSegment) Create(ctx context.Context) *DatabaseSegment {
	return &DatabaseSegment{
		haloSpan: &halo.ClientSpan{
			Client:        globalHaloClient,
			SlowThreshold: s.SlowThreshold,
			Service:       zae.Service(),
			Method:        MethodFromContext(ctx),
		},
		System:       s.System,
		Host:         s.Host,
		PortOrPath:   s.PortOrPath,
		DatabaseName: s.DatabaseName,
		Collection:   s.Collection,
		Operation:    s.Operation,
		Query:        s.Query,
	}
}

func (s *DatabaseSegment) Name() string {
	if s.Collection == "" {
		return fmt.Sprintf("%s.%s/%s", s.System, s.DatabaseName, s.Operation)
	}
	return fmt.Sprintf("%s.%s.%s/%s", s.System, s.DatabaseName, s.Collection, s.Operation)
}

func (s *DatabaseSegment) Address() string {
	if s.Host != "" && s.PortOrPath != "" {
		return s.Host + ":" + s.PortOrPath
	}
	if s.Host != "" {
		return s.Host
	}
	return "unknown"
}

func (s *DatabaseSegment) End(ctx context.Context, input Error) {
	if s.System == "" {
		s.System = "unknown"
	}
	if s.DatabaseName == "" {
		s.DatabaseName = "unknown"
	}
	if s.Operation == "" {
		s.Operation = "unknown"
	}

	if s.Query == "" {
		collection := s.Collection
		if collection == "" {
			collection = "unknown"
		}
		s.Query = fmt.Sprintf(`'%s' on '%s' using '%s'`, s.Operation, collection, s.System)
	}

	elapsed := time.Since(s.start)

	{
		s.haloSpan.TargetService = fmt.Sprintf("%s_%s", s.System, s.DatabaseName)
		s.haloSpan.TargetMethod = s.Operation

		s.haloSpan.End(elapsed, input)
	}

	{
		if s.openSpan != nil {
			s.openSpan.SetOperationName(s.Name())

			s.openSpan.SetTag("span.kind", "client")
			s.openSpan.SetTag("db.system", s.System)
			s.openSpan.SetTag("db.connection_string", s.Address())
			s.openSpan.SetTag("db.name", s.DatabaseName)
			s.openSpan.SetTag("db.statement", s.Query)
			s.openSpan.SetTag("db.parameters", s.Parameters)

			if input != nil {
				s.openSpan.SetTag("error", true)
				s.openSpan.LogFields(spanlog.String("message", input.Error()))
			}

			s.openSpan.Finish()
		}
	}

	{
		if input != nil {
			fields := log.Fields{
				"db.system":            s.System,
				"db.connection_string": s.Address(),
				"db.name":              s.DatabaseName,
				"db.statement":         s.Query,
				"elapsed":              elapsed.String(),
				"error.class":          input.Class(),
			}
			log.WithFields(ctx, fields).WithError(input).Error(input.Error())
		}
	}
}
