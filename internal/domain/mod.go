package domain

import (
	"go.uber.org/fx"
	"negotiation-history-B3/internal/domain/trade"
	"negotiation-history-B3/internal/domain/trade/finder"
	"negotiation-history-B3/internal/domain/trade/processor"
)

var Common = fx.Options(
	fx.Provide(
		fx.Annotate(processor.New, fx.As(new(trade.Processor))),
		fx.Annotate(finder.New, fx.As(new(trade.Finder))),
	),
)
