package influx

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const dbName = "currate-test"

var (
	container     testcontainers.Container
	hostPort      string
	influxClient  *Client
	testRatesMap  map[string][]domain.Rate
	ratesQuantity = uint(10)
	err           error
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// prepare influxdb3 container
	container, hostPort, err = runGenericInfluxV3(ctx, dbName)
	if err != nil {
		panic(err)
	}
	defer container.Terminate(ctx)

	// prepare client
	influxClient, err = NewClient(
		"http://"+net.JoinHostPort("127.0.0.1", hostPort),
		"Bearer",
		dbName,
	)
	if err != nil {
		panic(err)
	}

	testRatesMap = map[string][]domain.Rate{
		"BTC": generateTestRates("BTC", ratesQuantity),
		"ETH": generateTestRates("ETH", ratesQuantity),
		"TRX": generateTestRates("TRX", ratesQuantity),
	}

	if code := m.Run(); code != 0 {
		os.Exit(code)
	}
}

func TestHealthCheck(t *testing.T) {
	require.Eventually(t, func() bool {
		return influxClient.HealthCheck(context.Background()) == nil
	}, time.Second*5, time.Second)
}

func TestInfluxInsertAndGetCurrency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// insert test data
	for _, rates := range testRatesMap {
		err = influxClient.Insert(ctx, rates...)
		require.NoError(t, err)
	}

	for _, v := range []uint{10, 5, 15} {
		t.Run(fmt.Sprintf("limit=%d", v), func(t *testing.T) {
			testInfluxInsertAndGetCurrency(t, ctx, v)
		})
	}
}

func testInfluxInsertAndGetCurrency(t *testing.T, ctx context.Context, limit uint) {
	for currency, sourceRates := range testRatesMap {
		ratesFromInflux, err := influxClient.Get(ctx, currency, limit)
		require.NoError(t, err)
		validateLastNRates(t, ratesFromInflux, sourceRates, int(limit))
	}
}
func TestInfluxInsertAndGetAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// validate GetAll
	for _, v := range []uint{10, 5, 15} {
		t.Run(fmt.Sprintf("limit=%d", v), func(t *testing.T) {
			testInfluxInsertAndGetAll(t, ctx, v)
		})
	}
}

func testInfluxInsertAndGetAll(t *testing.T, ctx context.Context, limit uint) {
	influxRatesMap, err := influxClient.GetAll(ctx, limit)
	require.NoError(t, err)
	require.Len(t, influxRatesMap, len(testRatesMap))

	for currency, rates := range influxRatesMap {
		sourceRates, ok := testRatesMap[currency]
		require.True(t, ok)
		validateLastNRates(t, rates, sourceRates, int(limit))
	}
}

func runGenericInfluxV3(ctx context.Context, dbName string) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image: "influxdb:3-core",
		Cmd: []string{
			"influxdb3",
			"serve",
			"--node-id=node0", // note: the node-id flag is required
			"--without-auth",
			"--bucket", dbName,
		},
		Env: map[string]string{
			"INFLUXDB3_HOST_URL": net.JoinHostPort("127.0.0.1", "8181"),
		},

		ExposedPorts: []string{"8181/tcp"}, // 8181 is the default exposed port for influxdb3
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

func generateTestRates(currency string, qty uint) []domain.Rate {
	res := make([]domain.Rate, qty)

	for i := 0; i < int(qty); i++ {
		res[i] = domain.Rate{
			Currency:  currency,
			Quote:     "USD",
			Provider:  "https://min-api.cryptocompare.com",
			Value:     rand.Float64(),
			Timestamp: time.Now().UTC(),
		}
		time.Sleep(time.Millisecond * 10) // a little delay to avoid timestamp collisions
	}

	return res
}

func validateLastNRates(t *testing.T, ratesFromInflux, sourceRates []domain.Rate, limit int) {
	if limit > len(sourceRates) {
		require.Equal(t, len(ratesFromInflux), len(sourceRates))
	} else {
		require.Equal(t, limit, len(ratesFromInflux))
	}

	for i, rateFromInflux := range ratesFromInflux {
		require.Equal(t, sourceRates[len(sourceRates)-i-1].Value, rateFromInflux.Value)
		require.Equal(t, sourceRates[len(sourceRates)-i-1].Timestamp, rateFromInflux.Timestamp)
		require.Equal(t, sourceRates[len(sourceRates)-i-1].Provider, rateFromInflux.Provider)
		require.Equal(t, sourceRates[len(sourceRates)-i-1].Quote, rateFromInflux.Quote)
		require.Equal(t, sourceRates[len(sourceRates)-i-1].Currency, rateFromInflux.Currency)
	}
}
