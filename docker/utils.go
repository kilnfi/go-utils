package docker

import (
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/compose/loader"
	dockertypes "github.com/docker/docker/api/types"
	dockermount "github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

func ParseVolumes(specs ...string) (volumes map[string]struct{}, binds []string, err error) {
	volumes = make(map[string]struct{})
	for _, spec := range specs {
		parsed, err := loader.ParseVolume(spec)
		if err != nil {
			return nil, nil, err
		}

		if parsed.Source != "" {
			toBind := spec

			if parsed.Type == string(dockermount.TypeBind) {
				if arr := strings.SplitN(spec, ":", 2); len(arr) == 2 {
					hostPart := arr[0]
					if strings.HasPrefix(hostPart, "."+string(filepath.Separator)) || hostPart == "." {
						if absHostPart, err := filepath.Abs(hostPart); err == nil {
							hostPart = absHostPart
						}
					}
					toBind = hostPart + ":" + arr[1]
				}
			}

			// after creating the bind mount we want to delete it from the copts.volumes values because
			// we do not want bind mounts being committed to image configs
			binds = append(binds, toBind)
		} else {
			volumes[spec] = struct{}{}
		}
	}

	return
}

// GetPortBindings returns the list of port bindingto access container on given port
func GetPortBindings(port string, container *dockertypes.ContainerJSON) ([]nat.PortBinding, error) {
	proto := "tcp"
	parts := strings.SplitN(port, "/", 2)

	if len(parts) == 2 && len(parts[1]) != 0 {
		port = parts[0]
		proto = parts[1]
	}

	natPort := port + "/" + proto
	newP, err := nat.NewPort(proto, port)
	if err != nil {
		return nil, err
	}

	if frontends, exists := container.NetworkSettings.Ports[newP]; exists && frontends != nil {
		return frontends, nil
	}

	return nil, errors.Errorf("Error: No public port '%s' published for container %s", natPort, container.ID)
}
