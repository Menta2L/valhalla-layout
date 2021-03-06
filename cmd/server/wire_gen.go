// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/menta2l/valhalla-layout/internal/biz"
	"github.com/menta2l/valhalla-layout/internal/conf"
	"github.com/menta2l/valhalla-layout/internal/data"
	"github.com/menta2l/valhalla-layout/internal/server"
	"github.com/menta2l/valhalla-layout/internal/service"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, registry *conf.Registry, confData *conf.Data, auth *conf.Auth, logger log.Logger, tracerProvider *trace.TracerProvider) (*kratos.App, func(), error) {
	client := data.NewEntClient(confData, logger)
	dataData, cleanup, err := data.NewData(confData, client, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	registrar := data.NewRegister(registry)
	app := newApp(logger, httpServer, grpcServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
