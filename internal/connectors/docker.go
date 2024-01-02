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
	"slices"
	"strconv"
)

const privatePort = 3000

type Docker struct {
	client     *client.Client
	usedImages []string
}

func NewDocker() (*Docker, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "error creating a docker client")
	}

	return &Docker{
		client:     c,
		usedImages: make([]string, 0),
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
	containers, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return "", errors.Wrap(err, "error getting a list of containers")
	}

	for _, c := range containers {
		if c.Image == image {
			for _, port := range c.Ports {
				if port.PrivatePort == privatePort {
					return strconv.Itoa(int(port.PublicPort)), nil
				}
			}
		}
	}

	return "", errors.Errorf("container with image %s is not started", image)
}

func (d Docker) StartContainers(images []string) error {
	containers, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return errors.Wrap(err, "error getting a list of containers")
	}

	for _, image := range images {
		for _, c := range containers {
			if c.Image == image {
				slog.Warn("container is already started", slog.String("image", image))

				continue
			}
		}

		slog.Info("pulling an image", slog.String("image", image))
		if err := d.Pull(image); err != nil {
			return errors.Wrap(err, "error pulling an image")
		}

		c, err := d.CreateNewContainer(image)
		if err != nil {
			return errors.Wrap(err, "error starting a container")
		}

		slog.Info("container is started", slog.String("id", c.ID), slog.String("image", image))
		d.usedImages = append(d.usedImages, image)
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

	containerPort, err := nat.NewPort("tcp", strconv.Itoa(privatePort))
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

	containers, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return errors.Wrap(err, "error getting a list of containers")
	}

	for _, c := range containers {
		if slices.Contains(d.usedImages, c.Image) {
			id := c.ID

			eg.Go(func() error {
				return d.client.ContainerStop(context.Background(), id, container.StopOptions{})
			})
		}
	}

	return errors.Wrap(eg.Wait(), "error stopping containers")
}
