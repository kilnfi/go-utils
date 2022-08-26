package docker

import (
	"context"
	"fmt"
	"time"

	dockerref "github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockernetwork "github.com/docker/docker/api/types/network"
	dockervolume "github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
)

type Compose struct {
	cfg *ComposeConfig

	logger logrus.FieldLogger

	dockerc docker.CommonAPIClient

	services []*service
	volumes  []*volume
	networks []*network
}

func NewCompose(cfg *ComposeConfig) (*Compose, error) {
	dockerc, err := NewClient(cfg.Client)
	if err != nil {
		return nil, err
	}

	c := &Compose{
		cfg:     cfg,
		dockerc: dockerc,
	}
	c.SetLogger(logrus.StandardLogger())

	return c, nil
}

func (c *Compose) SetLogger(logger logrus.FieldLogger) {
	c.logger = logger.WithField("component", "compose")
}

func (c *Compose) Up(ctx context.Context) error {
	if err := c.createNetworks(ctx); err != nil {
		return err
	}

	if err := c.createVolumes(ctx); err != nil {
		return err
	}

	if err := c.pullImages(ctx); err != nil {
		return err
	}

	if err := c.createContainers(ctx); err != nil {
		return err
	}

	if err := c.startContainers(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Compose) Down(ctx context.Context) error {
	var rErr error
	if err := c.stopContainers(ctx); rErr == nil && err != nil {
		rErr = err
	}

	if err := c.removeContainers(ctx); err != nil {
		rErr = err
	}

	if err := c.removeVolumes(ctx); rErr == nil && err != nil {
		rErr = err
	}

	if err := c.removeNetworks(ctx); rErr == nil && err != nil {
		rErr = err
	}

	return rErr
}

func (c *Compose) Name(name string) string {
	return c.name(name)
}

func (c *Compose) name(name string) string {
	return fmt.Sprintf("%v_%v", c.cfg.Namespace, name)
}

func (c *Compose) RegisterNetwork(name string, cfg dockertypes.NetworkCreate) {
	c.networks = append(c.networks, &network{
		name: c.name(name),
		cfg:  cfg,
	})
}

type network struct {
	name string
	cfg  dockertypes.NetworkCreate

	id  string
	err error
}

func (c *Compose) createNetworks(ctx context.Context) error {
	for _, ntwrk := range c.networks {
		err := c.createNetwork(ctx, ntwrk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compose) createNetwork(ctx context.Context, ntwrk *network) error {
	logger := c.logger.WithField("network.name", ntwrk.name)

	logger.Infof("create network...")
	resp, err := c.dockerc.NetworkCreate(ctx, ntwrk.name, ntwrk.cfg)
	if err != nil {
		ntwrk.err = err
		logger.WithError(err).Infof("failed to create network...")
		return err
	}

	ntwrk.id = resp.ID
	logger = logger.WithField("network.id", ntwrk.id)
	if resp.Warning != "" {
		logger.Warn(resp.Warning)
	}

	logger.Infof("network created")

	return nil
}

func (c *Compose) removeNetworks(ctx context.Context) error {
	var rErr error
	for _, ntwrk := range c.networks {
		err := c.removeNetwork(ctx, ntwrk)
		if rErr == nil && err != nil {
			rErr = err
		}
	}
	return rErr
}

func (c *Compose) removeNetwork(ctx context.Context, ntwrk *network) error {
	if ntwrk.id != "" {
		logger := c.logger.WithField("network.name", ntwrk.name)
		logger.Infof("remove network...")
		err := c.dockerc.NetworkRemove(ctx, ntwrk.id)
		if err != nil {
			logger.WithError(err).Errorf("failed to remove network")
			return err
		}
		logger.Infof("network removed")
	}

	return nil
}

type volume struct {
	name string
	cfg  dockervolume.VolumeCreateBody

	err error
}

func (c *Compose) RegisterVolume(name string, cfg dockervolume.VolumeCreateBody) {
	name = c.name(name)
	cfg.Name = name
	c.volumes = append(c.volumes, &volume{
		name: name,
		cfg:  cfg,
	})
}

func (c *Compose) createVolumes(ctx context.Context) error {
	for _, vol := range c.volumes {
		err := c.createVolume(ctx, vol)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compose) createVolume(ctx context.Context, vol *volume) error {
	logger := c.logger.WithField("volume.name", vol.name)

	logger.Infof("create volume...")
	_, err := c.dockerc.VolumeCreate(ctx, vol.cfg)
	if err != nil {
		vol.err = err
		logger.WithError(err).Infof("failed to create volume...")
		return err
	}

	logger.Infof("volume created")

	return nil
}

func (c *Compose) removeVolumes(ctx context.Context) error {
	var rErr error
	for _, vol := range c.volumes {
		err := c.removeVolume(ctx, vol)
		if rErr == nil && err != nil {
			rErr = err
		}
	}
	return rErr
}

func (c *Compose) removeVolume(ctx context.Context, vol *volume) error {
	if vol.err == nil {
		logger := c.logger.WithField("volume.name", vol.name)
		logger.Infof("remove volume...")
		err := c.dockerc.VolumeRemove(ctx, vol.name, true)
		if err != nil {
			logger.WithError(err).Errorf("failed to remove volume")
			return err
		}
		logger.Infof("volume removed")
	}

	return nil
}

type ServiceConfig struct {
	Name       string
	Container  *dockercontainer.Config
	Host       *dockercontainer.HostConfig
	Networking *dockernetwork.NetworkingConfig
	IsReady    func(context.Context, *dockertypes.ContainerJSON) error
	DependsOn  []string
}

type service struct {
	name string
	cfg  *ServiceConfig

	id  string
	err error
}

func (c *Compose) RegisterService(name string, cfg *ServiceConfig) {
	c.services = append(c.services, &service{
		name: c.name(name),
		cfg:  cfg,
	})
}

func (c *Compose) GetContainer(ctx context.Context, name string) (*dockertypes.ContainerJSON, error) {
	svc, err := c.getService(name)
	if err != nil {
		return nil, err
	}

	return c.getContainer(ctx, svc)
}

func (c *Compose) getContainer(ctx context.Context, svc *service) (*dockertypes.ContainerJSON, error) {
	if svc.err != nil {
		return nil, svc.err
	}

	container, err := c.dockerc.ContainerInspect(ctx, svc.id)
	if err != nil {
		return nil, err
	}

	return &container, nil
}

func (c *Compose) getService(name string) (*service, error) {
	name = c.name(name)
	for _, svc := range c.services {
		if svc.name == name {
			return svc, nil
		}
	}

	return nil, fmt.Errorf("no service %q", name)
}

func (c *Compose) pullImages(ctx context.Context) error {
	for _, svc := range c.services {
		err := c.pullImage(ctx, svc.cfg.Container.Image)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compose) pullImage(ctx context.Context, image string) error {
	logger := c.logger.WithField("image", image)
	logger.Infof("pull image...")

	_, err := dockerref.ParseNormalizedNamed(image)
	if err != nil {
		logger.WithError(err).Errorf("invalid image name")
		return err
	}

	options := types.ImagePullOptions{
		RegistryAuth: "", // TODO: deal with docker registry authentication
	}

	respBody, err := c.dockerc.ImagePull(ctx, image, options)
	if err != nil {
		return err
	}
	defer respBody.Close()

	return jsonmessage.DisplayJSONMessagesStream(
		respBody,
		logger.Writer(),
		0,
		false,
		nil,
	)
}

func (c *Compose) createContainers(ctx context.Context) error {
	for _, svc := range c.services {
		err := c.createContainer(ctx, svc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compose) createContainer(ctx context.Context, svc *service) error {
	logger := c.logger.WithField("container.name", svc.name)
	logger.Infof("create container...")

	resp, err := c.dockerc.ContainerCreate(
		ctx,
		svc.cfg.Container,
		svc.cfg.Host,
		svc.cfg.Networking,
		nil,
		svc.name,
	)
	if err != nil {
		logger.WithError(err).Errorf("failed to crate container")
		svc.err = err
		return err
	}

	svc.id = resp.ID
	logger = logger.WithField("container.id", svc.id)

	for _, warning := range resp.Warnings {
		logger.Warn(warning)
	}

	logger.Infof("container created")

	return nil
}

func (c *Compose) removeContainers(ctx context.Context) error {
	var rErr error
	for _, svc := range c.services {
		err := c.removeContainer(ctx, svc)
		if rErr == nil && err != nil {
			rErr = err
		}
	}
	return rErr
}

func (c *Compose) removeContainer(ctx context.Context, svc *service) error {
	if svc.id != "" {
		logger := c.logger.WithField("container.name", svc.name)
		logger.Infof("remove container...")
		err := c.dockerc.ContainerRemove(
			ctx,
			svc.id,
			dockertypes.ContainerRemoveOptions{},
		)
		if err != nil {
			logger.WithError(err).Errorf("failed to remove container")
			return err
		}
		logger.Infof("container removed")
	}

	return nil
}

func (c *Compose) startContainers(ctx context.Context) error {
	var rErr error
	for _, svc := range c.services {
		err := c.startContainer(ctx, svc)
		if rErr == nil && err != nil {
			rErr = err
		}
	}
	return rErr
}

func (c *Compose) startContainer(ctx context.Context, svc *service) error {
	if svc.id != "" {
		for _, dep := range svc.cfg.DependsOn {
			err := c.WaitContainer(ctx, dep, 10*time.Second)
			if err != nil {
				return err
			}
		}
		logger := c.logger.WithField("container.name", svc.name)
		logger.Infof("start container...")
		err := c.dockerc.ContainerStart(
			ctx,
			svc.id,
			dockertypes.ContainerStartOptions{},
		)
		if err != nil {
			svc.err = err
			logger.WithError(err).Errorf("failed to start container")
			return err
		}
		logger.Infof("container started")

		err = c.isContainerReady(ctx, svc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compose) isContainerReady(ctx context.Context, svc *service) error {
	if svc.err != nil {
		return svc.err
	}

	if svc.id == "" {
		return fmt.Errorf("container not started")
	}

	return nil
}

func (c *Compose) stopContainers(ctx context.Context) error {
	var rErr error
	for _, svc := range c.services {
		err := c.stopContainer(ctx, svc)
		if rErr == nil && err != nil {
			rErr = err
		}
	}
	return rErr
}

func (c *Compose) stopContainer(ctx context.Context, svc *service) error {
	if svc.id != "" && svc.err == nil {
		logger := c.logger.WithField("container.name", svc.name)
		logger.Infof("stop container...")
		err := c.dockerc.ContainerStop(
			ctx,
			svc.id,
			nil,
		)
		if err != nil {
			svc.err = err
			logger.WithError(err).Errorf("failed to stop container")
			return err
		}
		logger.Infof("container stopped")
	}

	return nil
}

func (c *Compose) WaitContainer(ctx context.Context, name string, timeout time.Duration) error {
	svc, err := c.getService(name)
	if err != nil {
		return err
	}

	return c.waitContainer(ctx, svc, timeout)
}

func (c *Compose) waitContainer(ctx context.Context, svc *service, timeout time.Duration) error {
	logger := c.logger.WithField("container.name", svc.name)
	logger.Infof("wait for container to be ready...")

	container, err := c.getContainer(ctx, svc)
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-timeoutCtx.Done():
			logger.WithError(timeoutCtx.Err()).Errorf("container never got ready...")
			return timeoutCtx.Err()
		case <-ticker.C:
			if svc.err != nil {
				return svc.err
			}

			if err := svc.cfg.IsReady(timeoutCtx, container); err == nil {
				logger.Infof("container is ready")
				return nil
			} else {
				logger.WithError(err).Warnf("container not yet ready...")
			}
		}
	}
}
