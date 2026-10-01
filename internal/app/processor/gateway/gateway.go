package pgateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/RonIT-401/catalog-service/internal/app/config/section"
	"github.com/RonIT-401/catalog-service/internal/app/processor"
	catalogv1 "github.com/RonIT-401/catalog-service/internal/pkg/grpc/gen/catalog/v1"
)

const (
	gatewayReadHeaderTimeout = 5 * time.Second
	gatewayReadTimeout       = 30 * time.Second
	gatewayWriteTimeout      = 30 * time.Second
	gatewayIdleTimeout       = 120 * time.Second
	gatewayShutdownTimeout   = 5 * time.Second
)

type gatewayProc struct {
	server   *http.Server
	addr     string
	grpcAddr string
}

func NewGateway(
	cfgGateway section.ProcessorGateway,
	cfgGrpc section.ProcessorGrpc,
) processor.Processor {
	httpAddr := fmt.Sprintf(":%d", cfgGateway.ListenPort)

	grpcAddr := net.JoinHostPort("localhost", strconv.Itoa(int(cfgGrpc.ListenPort)))

	server := http.Server{
		Addr:              httpAddr,
		ReadHeaderTimeout: gatewayReadHeaderTimeout,
		ReadTimeout:       gatewayReadTimeout,
		WriteTimeout:      gatewayWriteTimeout,
		IdleTimeout:       gatewayIdleTimeout,
	}

	return &gatewayProc{
		server:   &server,
		addr:     httpAddr,
		grpcAddr: grpcAddr,
	}
}

func (p *gatewayProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := catalogv1.RegisterCatalogServiceHandlerFromEndpoint(ctx, mux, p.grpcAddr, opts)
	if err != nil {
		log.Fatal()
		return
	}

	p.server.Handler = mux

	lc := net.ListenConfig{}

	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal()
		return
	}

	log.Info().Msgf("gRPC gateway started on %s", p.addr)

	go func() {
		_ = p.server.Serve(l)
	}()

	processor.WatchForShutdown(ctx, wg, processor.CloserFunc(func() error {
		return l.Close()
	}))

	processor.WatchForShutdown(
		ctx,
		wg,
		processor.NewCloserContextFunc(p.server.Shutdown, ctx, gatewayShutdownTimeout),
	)
}
