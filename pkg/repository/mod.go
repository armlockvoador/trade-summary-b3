package repository

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewTradeRepository, fx.As(new(TradeRepository))),
	),
)
