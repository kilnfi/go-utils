package docker

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v4/stdlib" // imported so pgx sql driver is registered
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
		portBindings, err := GetPortBindings("5432", container)
		if err != nil {
			return err
		}

		port, err := strconv.Atoi(portBindings[0].HostPort)
		if err != nil {
			return err
		}

		cfg := (&kilnsql.Config{
			Host:     portBindings[0].HostIP,
			Port:     uint16(port),
			User:     opts.User,
			Password: opts.Password,
			Dialect:  pgDialect,
			DBName:   pgDialect,
		}).SetDefault()

		db, err := sql.Open("pgx", cfg.DSN().String())
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
