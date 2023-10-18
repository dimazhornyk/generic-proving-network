package connectors

import (
	"context"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"log/slog"
	"os"
)

type Docker struct {
	client     *client.Client
	containers map[string]common.Container
}

func NewDocker() (*Docker, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error creating a docker client")
	}

	return &Docker{
		client:     c,
		containers: make(map[string]common.Container),
	}, nil
}

func (d Docker) Pull(image string) error {
	out, err := d.client.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error pulling an image")
	}

	defer out.Close()

	if _, err = io.Copy(os.Stdout, out); err != nil {
		return errors.Wrap(err, "error copying output")
	}

	return nil
}

func (d Docker) HasImage(image string) (bool, error) {
	list, err := d.client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false, errors.Wrap(err, "error getting a list of images")
	}

	for _, img := range list {
		for _, imgFullName := range img.RepoTags {
			if imgFullName == image {
				return true, nil
			}
		}
	}

	return false, nil
}

func (d Docker) GetContainerPort(image string) (string, error) {
	if c, ok := d.containers[image]; ok {
		return c.SourcePort, nil
	}

	return "", errors.Errorf("container with image %s is not started", image)
}

func (d Docker) StartContainers(images []string) error {
	for _, image := range images {
		if _, ok := d.containers[image]; ok {
			slog.Warn("container is already started", slog.String("image", image))

			continue
		}

		if err := d.Pull(image); err != nil {
			return errors.Wrap(err, "error pulling an image")
		}

		c, err := d.CreateNewContainer(image)
		if err != nil {
			return errors.Wrap(err, "error starting a container")
		}

		d.containers[image] = c
	}

	return nil
}

func (d Docker) CreateNewContainer(image string) (common.Container, error) {
	port, err := common.AvailablePort()
	if err != nil {
		return common.Container{}, errors.Wrap(err, "error getting an available port")
	}

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: port,
	}

	containerPort, err := nat.NewPort("tcp", "3000")
	if err != nil {
		return common.Container{}, errors.Wrap(err, "error creating a container port")
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	cont, err := d.client.ContainerCreate(
		context.Background(),
		&container.Config{Image: image},
		&container.HostConfig{PortBindings: portBinding},
		nil,
		nil,
		"",
	)
	if err != nil {
		return common.Container{}, errors.Wrap(err, "error creating a container")
	}

	if err := d.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		return common.Container{}, errors.Wrap(err, "error starting a container")
	}

	slog.Info("Container is started", slog.String("id", cont.ID))

	return common.Container{
		ID:         cont.ID,
		SourcePort: port,
	}, nil
}

func (d Docker) StopContainers() error {
	eg := errgroup.Group{}

	for _, c := range d.containers {
		cont := c

		eg.Go(func() error {
			return d.client.ContainerStop(context.Background(), cont.ID, container.StopOptions{})
		})
	}

	return errors.Wrap(eg.Wait(), "error stopping containers")
}
