package telemetry

import (
	"context"

	"git.in.zhihu.com/go/base/telemetry/sentry"
)

func Exec(jobName string, f func(ctx context.Context) Error) Error {
	txn, ctx, e := StartTransaction(context.Background(), &Transaction{
		System: TransactionExec,
		Method: jobName,
	})
	if e != nil {
		return WrapErrWithUnknownClass(e)
	}

	var err Error
	sentry.Recover(func() {
		err = f(ctx)
	}, func(e error) {
		err = WrapErrWithUnknownClass(e)
	})
	txn.End(ctx, err)

	return err
}
