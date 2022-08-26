package docker

import (
	"net/http"

	docker "github.com/docker/docker/client"
	dockersockets "github.com/docker/go-connections/sockets"
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

func NewClient(cfg *ClientConfig) (*docker.Client, error) {
	host, err := docker.ParseHostURL(cfg.Host)
	if err != nil {
		return nil, err
	}

	httpc, err := kilnhttp.NewClient(cfg.Client)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{}
	err = dockersockets.ConfigureTransport(transport, host.Scheme, host.Host)
	if err != nil {
		return nil, err
	}
	httpc.Transport = transport

	opts := []docker.Opt{
		docker.WithHTTPClient(httpc),
		docker.WithHost(cfg.Host),
		docker.WithVersion(cfg.Version),
	}

	return docker.NewClientWithOpts(opts...)
}
