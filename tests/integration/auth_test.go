package integration

import (
	"context"
	"github.com/bubalync/uni-auth/internal/app"
	"github.com/bubalync/uni-auth/internal/config"
	"github.com/bubalync/uni-auth/tests/integration/containers"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	configPath = "./../../config/local.yaml"
	host       = "http://localhost:8080"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Containers
	postgres := containers.NewPostgresContainer(ctx)
	defer postgres.Close(ctx)

	connectionString, _ := postgres.Container.ConnectionString(ctx)

	log.Println("connectionString = ", connectionString)
	cfg := config.NewConfigByPath(configPath)
	cfg.PG.Url = connectionString

	go func() {
		app.Run(cfg)
	}()

	time.Sleep(5 * time.Second)

	os.Exit(m.Run())
}

func TestRegistration(t *testing.T) {
	client := resty.New()

	const registrationPath = host + "/api/v1/users/register"

	type regReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type regResp struct {
		Id string `json:"id"`
	}

	var resp regResp
	r, err := client.R().
		SetBody(regReq{Email: "test@example.com", Password: "password"}).
		SetResult(&resp).
		Post(registrationPath)

	log.Println(resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, r.StatusCode())
	assert.NoError(t, uuid.Validate(resp.Id))
}
