package main

import (
	"context"
	"go.uber.org/fx"
	"multi-proving-client/internal/common"
	"multi-proving-client/internal/connectors"
	"multi-proving-client/internal/logic"
	"multi-proving-client/internal/presenters"
)

func main() {
	app := buildApp()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}

	//chat.Main()
}

func buildApp() *fx.App {
	return fx.New(
		fx.Provide(
			func() context.Context { return context.Background() },
			common.NewConfig,
			connectors.NewDocker,
			connectors.NewHost,
			logic.NewDHT,
			logic.NewDiscovery,
			connectors.NewPubSub,
			logic.NewGlobalMessaging,
			logic.NewNodesMap,
			logic.NewState,
			logic.NewService,
			presenters.NewAPI,
			presenters.NewListener,
			// sends status updates
			fx.Invoke(func(ctx context.Context, cfg *common.Config, messaging *logic.GlobalMessaging) error {
				return messaging.Init(ctx, cfg.Consumers)
			}),
			// listens to other's messages
			fx.Invoke(func(ctx context.Context, listener *presenters.Listener) error { // others listener
				return listener.Listen(ctx)
			}),
			// handles proofs generation
			fx.Invoke(func(ctx context.Context, service *logic.Service) error { // service
				return service.Start(ctx)
			}),
		),
	)
}
