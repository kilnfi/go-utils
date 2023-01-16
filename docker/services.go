package docker

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/stdlib" // imported so pgx sql driver is registered
	kilnsql "github.com/kilnfi/go-utils/sql"
)

type PostgresServiceOpts struct {
	Version, Port, User, Password string
}

var pgDialect = "postgres"

func (opts *PostgresServiceOpts) SetDefault() *PostgresServiceOpts {
	if opts.Version == "" {
		opts.Version = "14.4"
	}

	if opts.Port == "" {
		opts.Port = "0" // expose on random available port
	}

	if opts.User == "" {
		opts.User = pgDialect
	}

	if opts.Password == "" {
		opts.Password = pgDialect
	}

	return opts
}

// SQLConfig returns an SQL config to connect to the postgres container from the host
func (opts *PostgresServiceOpts) SQLConfig(container *dockertypes.ContainerJSON) (*kilnsql.Config, error) {
	portBindings, err := GetPortBindings("5432", container)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portBindings[0].HostPort)
	if err != nil {
		return nil, err
	}

	return (&kilnsql.Config{
		Host:     portBindings[0].HostIP,
		Port:     uint16(port),
		User:     opts.User,
		Password: opts.Password,
		Dialect:  pgDialect,
		DBName:   pgDialect,
	}).SetDefault(), nil
}

func NewPostgresServiceConfig(opts *PostgresServiceOpts) (*ServiceConfig, error) {
	ports, portBindings, err := nat.ParsePortSpecs([]string{
		fmt.Sprintf("%v:5432", opts.Port),
	})
	if err != nil {
		return nil, err
	}

	containerCfg := &dockercontainer.Config{
		Image:        fmt.Sprintf("postgres:%v", opts.Version),
		ExposedPorts: ports,
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%v", opts.Password),
			fmt.Sprintf("POSTGRES_USER=%v", opts.User),
		},
	}

	hostCfg := &dockercontainer.HostConfig{
		PortBindings: portBindings,
	}

	isReady := func(ctx context.Context, container *dockertypes.ContainerJSON) error {
		sqlCfg, err := opts.SQLConfig(container)
		if err != nil {
			return err
		}

		db, err := sql.Open("pgx", sqlCfg.DSN().String())
		if err != nil {
			return err
		}

		return db.PingContext(ctx)
	}

	return &ServiceConfig{
		Container: containerCfg,
		Host:      hostCfg,
		IsReady:   isReady,
	}, nil
}

type TraefikServiceOpts struct {
	Version, Port string
}

func (opts *TraefikServiceOpts) SetDefault() *TraefikServiceOpts {
	if opts.Version == "" {
		opts.Version = "v2.8.0"
	}

	if opts.Port == "" {
		opts.Port = "0" // expose on random available port
	}

	return opts
}

// Addr returns address to connect to the traefik container from the host
func (opts *TraefikServiceOpts) Addr(container *dockertypes.ContainerJSON) (string, error) {
	portBindings, err := GetPortBindings("80", container)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%v:%v", portBindings[0].HostIP, portBindings[0].HostPort), nil
}

func NewTreafikServiceConfig(opts *TraefikServiceOpts) (*ServiceConfig, error) {
	ports, portBindings, err := nat.ParsePortSpecs([]string{
		fmt.Sprintf("%v:80", opts.Port),
	})
	if err != nil {
		return nil, err
	}

	volumes, binds, err := ParseVolumes("/var/run/docker.sock:/var/run/docker.sock:ro")
	if err != nil {
		return nil, err
	}

	containerCfg := &dockercontainer.Config{
		Image: fmt.Sprintf("traefik:%v", opts.Version),
		Cmd: []string{
			"--api.insecure=true",
			"--providers.docker=true",
			"--providers.docker.exposedbydefault=false",
			"--entrypoints.web.address=:80",
			"--ping=true",
			"--ping.entryPoint=web",
		},
		ExposedPorts: ports,
		Volumes:      volumes,
	}

	hostCfg := &dockercontainer.HostConfig{
		PortBindings: portBindings,
		Binds:        binds,
	}

	isReady := func(ctx context.Context, container *dockertypes.ContainerJSON) error {
		addr, err := opts.Addr(container)
		if err != nil {
			return err
		}

		req, _ := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%v/ping", addr),
			http.NoBody,
		)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("invalid status: %v", resp.Status)
		}

		return nil
	}

	return &ServiceConfig{
		Container: containerCfg,
		Host:      hostCfg,
		IsReady:   isReady,
	}, nil
}
