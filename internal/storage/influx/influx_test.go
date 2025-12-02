package influx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInflux(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container, hostPort, err := runGenericInfluxV3(ctx)
	require.NoError(t, err)
	require.NotZero(t, hostPort)
	require.NotNil(t, container)

	defer container.Terminate(ctx)
}

func runGenericInfluxV3(ctx context.Context) (testcontainers.Container, string, error) {
	//portWithProto := fmt.Sprintf("%s/tcp", exposedPort)

	req := testcontainers.ContainerRequest{
		Image:        "influxdb:3-core",
		Cmd:          []string{"--without-auth"},
		ExposedPorts: []string{"8181/tcp", "8086/tcp"},
		WaitingFor:   wait.ForListeningPort("8181/tcp").WithStartupTimeout(2 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, "", err
	}

	hostPort, err := container.MappedPort(ctx, "8181")

	return container, hostPort.Port(), err
}
