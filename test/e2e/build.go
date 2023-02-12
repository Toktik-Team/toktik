package main

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
)

type bizContainer struct {
	testcontainers.Container
	endpoint string
}

func setupBiz(ctx context.Context) (*bizContainer, error) {

	// print current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("current working directory: %s", dir)
	req := testcontainers.ContainerRequest{
		ExposedPorts: []string{"40126/tcp"},
		FromDockerfile: testcontainers.FromDockerfile{
			PrintBuildLog: true,
			Context:       "../../",
			Dockerfile:    "Dockerfile",
		},
	}
	log.Printf("req: %+v", req)

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	mappedPort, err := container.MappedPort(ctx, "40126")
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &bizContainer{Container: container, endpoint: fmt.Sprintf("%s:%s", host, mappedPort.Port())}, nil
}
