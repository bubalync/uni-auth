package containers

import (
	"context"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName     = "uni-auth"
	dbUser     = "user"
	dbPassword = "password"
)

type PostgresContainer struct {
	Container *postgres.PostgresContainer
}

func NewPostgresContainer(ctx context.Context) *PostgresContainer {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("..", "..", "migrations", "000001_create_users.up.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Fatalf("Failed to start postgres container: %s", err)
		return nil
	}

	return &PostgresContainer{
		postgresContainer,
	}
}

func (c *PostgresContainer) Close(ctx context.Context) {
	func() {
		if err := c.Container.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %s", err)
		}
	}()
}
