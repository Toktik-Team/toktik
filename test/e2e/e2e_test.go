package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/auth"
	"toktik/kitex_gen/douyin/auth/authservice"

	kitexClient "github.com/cloudwego/kitex/client"

	"github.com/hashicorp/consul/api"
)

func TestPostgres(t *testing.T) {
	ctx := context.Background()

	const dbname = "test-db"
	const user = "postgres"
	const password = "password"

	port, err := nat.NewPort("tcp", "5432")
	require.NoError(t, err)

	container, err := setupPostgres(ctx,
		WithPort(port.Port()),
		WithInitialDatabase(user, password, dbname),
		WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	containerPort, err := container.MappedPort(ctx, port)
	assert.NoError(t, err)

	host, err := container.Host(ctx)
	assert.NoError(t, err)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, containerPort.Port(), user, password, dbname)

	// perform assertions
	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()

	result, err := db.Exec("CREATE TABLE IF NOT EXISTS test (id int, name varchar(255));")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	result, err = db.Exec("INSERT INTO test (id, name) VALUES (1, 'test');")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestContainerWithWaitForSQL(t *testing.T) {
	const dbname = "test-db"
	const user = "postgres"
	const password = "password"

	ctx := context.Background()

	var port = "5432/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", host, port.Port(), dbname)
	}

	t.Run("default query", func(t *testing.T) {
		container, err := setupPostgres(ctx, WithPort(port), WithInitialDatabase("postgres", "password", dbname), WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL)))
		require.NoError(t, err)
		require.NotNil(t, container)
	})
	t.Run("custom query", func(t *testing.T) {
		container, err := setupPostgres(
			ctx,
			WithPort(port),
			WithInitialDatabase(user, password, dbname),
			WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL).WithStartupTimeout(time.Second*5).WithQuery("SELECT 10")),
		)
		require.NoError(t, err)
		require.NotNil(t, container)
	})
	t.Run("custom bad query", func(t *testing.T) {
		container, err := setupPostgres(
			ctx,
			WithPort(port),
			WithInitialDatabase(user, password, dbname),
			WithWaitStrategy(wait.ForSQL(nat.Port(port), "postgres", dbURL).WithStartupTimeout(time.Second*5).WithQuery("SELECT 'a' from b")),
		)
		require.Error(t, err)
		require.Nil(t, container)
	})
}

func TestConsul(t *testing.T) {
	ctx := context.Background()

	container, err := setupConsul(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// perform assertions
	cfg := api.DefaultConfig()
	cfg.Address = container.endpoint
	client, err := api.NewClient(cfg)
	if nil != err {
		t.Fatal(err)
	}

	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	authClient, err := authservice.NewClient(config.AuthServiceName, kitexClient.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}

	authenticate, _ := authClient.Authenticate(ctx, &auth.AuthenticateRequest{
		Token: "123",
	})

	log.Printf("authenticate: %+v", authenticate)

	bs := []byte("apple")
	_, err = client.KV().Put(&api.KVPair{
		Key:   "fruit",
		Value: bs,
	}, nil)
	if nil != err {
		t.Fatal(err)
	}
	pair, _, err := client.KV().Get("fruit", nil)
	if err != nil {
		t.Fatal(err)
	}
	if pair.Key != "fruit" || !bytes.Equal(pair.Value, []byte("apple")) {
		t.Errorf("get KV: %v %s,expect them to be: 'fruit' and 'apple'\n", pair.Key, pair.Value)
	}
}

func TestBizContainer(t *testing.T) {
	ctx := context.Background()
	container, err := setupBiz(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("container: %+v", container)
}

func TestBiz(t *testing.T) {

}
