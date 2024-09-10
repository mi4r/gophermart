package storage

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mi4r/gophermart/lib/logger"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var storage *pgxDriver

func TestMain(m *testing.M) {

	logger.InitLogger("debug")
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Запуск контейнера с PostgreSQL
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=test",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// Автоматическое удаление контейнера после завершения
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://test:test@%s/test?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseURL)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		storage = NewPgxDriver(databaseURL)
		if err := storage.Open(); err != nil {
			return err
		}
		return storage.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err := storage.autoMigrate(); err != nil {
		pool.Purge(resource)
		log.Fatal(err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
		storage.Close()
	}()

	// run tests
	m.Run()
}

func TestDatabaseConnection(t *testing.T) {
	// Пример теста: проверка соединения с базой данных
	err := storage.Ping()
	if err != nil {
		t.Fatalf("Could not ping database: %s", err)
	}
}

func TestUsersCreate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "create_Admin",
			user: User{
				Creds: Creds{
					Login:    "admin",
					Password: "admin",
				},
			},
			wantErr: false,
		},
		{
			name: "create_User1",
			user: User{
				Creds: Creds{
					Login:    "user1",
					Password: "user1",
				},
			},
			wantErr: false,
		},
		{
			name: "create_User2",
			user: User{
				Creds: Creds{
					Login:    "user2",
					Password: "user2",
				},
			},
			wantErr: false,
		},
		{
			name: "create_OneMoreUser1",
			user: User{
				Creds: Creds{
					Login:    "user1",
					Password: "user1",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.UserCreate(tt.user); err != nil {
				if tt.wantErr {
					t.Log(err)
				} else {
					t.Errorf("Not created user: %+v", tt.user)
					t.Error(err)
				}
			}
		})
	}
}

// func TestAutoMigrate(t *testing.T) {
// 	storage.autoMigrate()
// }
