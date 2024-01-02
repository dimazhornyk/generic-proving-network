package main

import (
	"context"
	"fmt"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/dimazhornyk/generic-proving-network/internal/logic/handlers"
	"github.com/dimazhornyk/generic-proving-network/internal/logic/sync"
	"github.com/dimazhornyk/generic-proving-network/internal/presenters"
	"github.com/dimazhornyk/generic-proving-network/proto"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"time"
)

const port = 5050

func main() {
	ctx := context.Background()
	app := buildApp(ctx)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}

func buildApp(ctx context.Context) *fx.App {
	return fx.New(
		fx.Provide(
			func() context.Context { return ctx },
			common.NewConfig,
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
			logic.NewNetworkParticipants,
			logic.NewGlobalMessaging,
			logic.NewStatusMap,
			logic.NewStorage,
			logic.NewService,
			sync.NewInitialSyncer,
			presenters.NewAPI,
			presenters.NewListener,
		),
		fx.Invoke(common.InitGobModels),
		// handles proofs generation, important for service to start first because it has to pull docker images
		fx.Invoke(func(ctx context.Context, service *logic.Service) error {
			return service.Start()
		}),
		// sync initial storage state
		fx.Invoke(func(ctx context.Context, syncer *sync.InitialSyncer) error {
			return syncer.Sync(ctx)
		}),
		// listens to others' messages
		fx.Invoke(func(ctx context.Context, listener *presenters.Listener) { // others listener
			go listener.Listen(ctx)
			time.Sleep(time.Second * 1) // wait for the listeners to start
		}),
		// sends status updates
		fx.Invoke(func(ctx context.Context, cfg *common.Config, messaging *logic.StatusSharing) error {
			return messaging.Init(ctx, cfg.Consumers)
		}),
		// provides data for initial sync for others
		fx.Invoke(func(ctx context.Context, syncer *sync.InitialSyncer) {
			syncer.ProvideData()
		}),
		// starts grpc server
		fx.Invoke(func(ctx context.Context, api *presenters.API) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				return fmt.Errorf("error starting a tcp listener: %w", err)
			}

			grpcServer := grpc.NewServer()
			proto.RegisterProvingNetworkServiceServer(grpcServer, api)
			slog.Info("grpc server started", "port", port)

			return grpcServer.Serve(listener)
		}),
	)
}
