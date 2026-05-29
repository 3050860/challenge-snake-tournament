package tests

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

var (
	//cfg                    *config.Server
	Env *dockerenv.TestEnvironment
	//ChallengesWorkerConfig worker.Config
)

type Config struct {
}
type TestEnvironment struct {
	Mongodb testcontainers.Container
	Net     *testcontainers.DockerNetwork
	Config  Config
}

type testEnvironmentRequest struct {
	db bool
}

type TestEnvironmentOption func(req *testEnvironmentRequest)

func NewTestEnvironment(ctx context.Context, modifiers ...TestEnvironmentOption) (*TestEnvironment, error) {

}

func CreateNetwork(ctx context.Context) (*testcontainers.DockerNetwork, error) {
	net, err := network.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create network:%w", err)
	}
	fmt.Printf("Created network: %s\n", net.ID)

	return net, nil
}
