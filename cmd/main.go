package main

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/dimazhornyk/generic-proving-network/internal/logic/handlers"
	"github.com/dimazhornyk/generic-proving-network/internal/logic/sync"
	"github.com/dimazhornyk/generic-proving-network/internal/presenters"
	"go.uber.org/fx"
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
			common.InitGobModels,
			connectors.NewDocker,
			connectors.NewPrivateKey,
			connectors.NewHost,
			connectors.NewEthereum,
			logic.NewDHT,
			logic.NewConnectionHolder,
			logic.NewDiscovery,
			handlers.NewProvingRequestsHandler,
			handlers.NewVotingHandler,
			handlers.NewStatusUpdatesHandler,
			handlers.NewProofsHandler,
			connectors.NewPubSub,
			logic.NewGlobalMessaging,
			logic.NewNodesMap,
			logic.NewStorage,
			logic.NewService,
			sync.NewInitialSyncer,
			presenters.NewAPI,
			presenters.NewListener,
			// handles proofs generation, important for service to start first because it has to pull docker images
			fx.Invoke(func(ctx context.Context, service *logic.Service) error {
				return service.Start()
			}),
			// sync initial storage state
			fx.Invoke(func(ctx context.Context, syncer *sync.InitialSyncer) error {
				return syncer.Sync(ctx)
			}),
			// sends status updates
			fx.Invoke(func(ctx context.Context, cfg *common.Config, messaging *logic.StatusSharing) error {
				return messaging.Init(ctx, cfg.Consumers)
			}),
			// provides data for initial sync for others
			fx.Invoke(func(ctx context.Context, syncer *sync.InitialSyncer) {
				syncer.ProvideData()
			}),
			// listens to other's messages
			fx.Invoke(func(ctx context.Context, listener *presenters.Listener) error { // others listener
				return listener.Listen(ctx)
			}),
		),
	)
}
