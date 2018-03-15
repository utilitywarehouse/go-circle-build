//go:generate protoc -I ../../proto  --go_out=plugins=grpc:${GOPATH}/src ../../proto/example.proto
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/utilitywarehouse/go-circle-build/examples/grpc-svc/internal/service"
	"github.com/utilitywarehouse/go-circle-build/examples/grpc-svc/pkg/pb/example"
	"github.com/utilitywarehouse/go-operational/op"
)

const (
	appName       = "grpc-svc"
	appDesc       = "ProjectTemplate"
	appOwner      = "telecom"
	appOwnerSlack = "#telecom"
	appURL        = "https://github.com/utilitywarehouse/go-stump"
)

var (
	gitHash string
)

func main() {
	app := cli.App(appName, appDesc)

	logLevel := app.String(cli.StringOpt{
		Name:   "log-level",
		Desc:   "log level [debug|info|warn|error]",
		EnvVar: "LOG_LEVEL",
		Value:  "info",
	})
	opPort := app.Int(cli.IntOpt{
		Name:   "op-port",
		Desc:   "The port to listen on for HTTP connections for Op Info",
		EnvVar: "OP_PORT",
		Value:  8081,
	})
	grpcPort := app.Int(cli.IntOpt{
		Name:   "grpc-port",
		Desc:   "GRPC port",
		Value:  8090,
		EnvVar: "GRPC_PORT",
	})

	app.Action = func() {
		grpc_prometheus.EnableHandlingTimeHistogram()
		logger := setUpLogger(*logLevel)
		grpc_logrus.ReplaceGrpcLogger(logger)

		svc := service.New("dummy")

		go func() {
			if errv := http.ListenAndServe(fmt.Sprintf(":%d", *opPort), op.NewHandler(opStatus())); errv != nil {
				logger.WithError(errv).Panicf("error listening on port: %d", *opPort)
			}
		}()

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		gSrv := grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_prometheus.UnaryServerInterceptor,
				grpc_recovery.UnaryServerInterceptor(),
				grpc_logrus.UnaryServerInterceptor(logger),
			)),
		)
		example.RegisterExampleServer(gSrv, svc)
		go waitForShutdown(func() {
			logger.Warn("shutdown")
			gSrv.GracefulStop()
		})
		reflection.Register(gSrv)
		if err := gSrv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
	log.Info("app starting")
	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("app stopped with error")
	}
}

func opStatus() *op.Status {
	return op.NewStatus(appName, appDesc).
		AddOwner(appOwner, appOwnerSlack).
		SetRevision(gitHash).
		ReadyUseHealthCheck().
		AddLink("VCS Repository", appURL)
}

func setUpLogger(level string) *log.Entry {
	l, err := log.ParseLevel(level)
	if err != nil {
		log.WithError(err).Panic("error parsing log level")
	}
	logger := log.Logger{
		Out:       os.Stderr,
		Formatter: &log.JSONFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     l,
	}
	return log.NewEntry(&logger)
}

func waitForShutdown(shutdown func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	shutdown()
}
