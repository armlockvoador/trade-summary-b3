package app

import (
	"context"
	"go.uber.org/fx"
	"negotiation-history-B3/internal/domain"
	"negotiation-history-B3/pkg/infra"
	"negotiation-history-B3/pkg/repository"
)

func NewApp(invokeFuncs ...interface{}) *fx.App {
	var fxInvokes []fx.Option
	for _, f := range invokeFuncs {
		fxInvokes = append(fxInvokes, fx.Invoke(f))
	}

	return fx.New(
		domain.Common,
		infra.Module,
		repository.Module,
		fx.Options(fxInvokes...),
	)
}

func StartApp(ctx context.Context, invokeFuncs ...interface{}) error {
	app := NewApp(invokeFuncs...)
	if err := app.Start(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	return app.Stop(context.Background())
}
