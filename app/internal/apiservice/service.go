package apiservice

import (
	"context"
	"time"

	"github.com/go-fuego/fuego"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApiService struct {
	cluster client.Client

	server *fuego.Server
}

type RouteRepository interface {
	Register(*fuego.Server)
}

func New(cluster client.Client) *ApiService {
	return &ApiService{
		cluster: cluster,
	}
}

func (apiService *ApiService) Run(ctx context.Context) error {
	server := fuego.NewServer(
		fuego.WithAddr(":9999"),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				DisableLocalSave: true,
			}),
		),
	)

	apiService.server = server

	var routeRepositories = []RouteRepository{
		NewAppAgentRepository(apiService.cluster),
	}

	for _, routeRepository := range routeRepositories {
		routeRepository.Register(server)
	}

	var errorChannel = make(chan error)
	go func() {
		if err := server.RunTLS("/tmp/k8s-extension-server/extension-certs/tls.crt", "/tmp/k8s-extension-server/extension-certs/tls.key"); err != nil {
			select {
			case <-ctx.Done():
				return
			case errorChannel <- err:
			}
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errorChannel:
		if err != nil {
			return err
		}
	}

	return nil
}
