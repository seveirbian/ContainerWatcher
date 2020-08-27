package docker

import (
	"context"
	"os"

	"github.com/docker/docker/client"
)

// Ctx context
var Ctx context.Context

// Client a client to communicate with dockerd
var Client *client.Client

func init() {
	cli, err := client.NewEnvClient()
	if err != nil {
		os.Exit(1)
	}
	Client = cli

	Ctx = context.Background()
}
