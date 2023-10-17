package connectors

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"io"
	"os"
)

type Docker struct {
	client *client.Client
}

func NewDocker() (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Docker{client: cli}, nil
}

func (d Docker) Pull(image string) error {
	out, err := d.client.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(os.Stdout, out)

	return err
}

func (d Docker) HasImage(image string) (bool, error) {
	list, err := d.client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false, err
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

func (d Docker) CreateNewContainer(image string) (string, error) {
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8000",
	}
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		panic("Unable to get the port")
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
		panic(err)
	}

	if err := d.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "error starting a container")
	}

	fmt.Printf("Container %s is started", cont.ID)
	return cont.ID, nil
}
