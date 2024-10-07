package drivers

import (
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	storagemart "github.com/mi4r/gophermart/internal/storage/gophermart"
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
		storage.isTest = true
		if err := storage.Open(); err != nil {
			return err
		}
		return storage.Ping()
	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// if err := storage.autoMigrate(); err != nil {
	// 	pool.Purge(resource)
	// 	log.Fatal(err)
	// }

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
		user    storagemart.User
		wantErr bool
	}{
		{
			name: "create_Admin",
			user: storagemart.User{
				Creds: storagemart.Creds{
					Login:    "admin",
					Password: "admin",
				},
			},
			wantErr: false,
		},
		{
			name: "create_User1",
			user: storagemart.User{
				Creds: storagemart.Creds{
					Login:    "user1",
					Password: "user1",
				},
			},
			wantErr: false,
		},
		{
			name: "create_User2",
			user: storagemart.User{
				Creds: storagemart.Creds{
					Login:    "user2",
					Password: "user2",
				},
			},
			wantErr: false,
		},
		{
			name: "create_OneMoreUser1",
			user: storagemart.User{
				Creds: storagemart.Creds{
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

func TestUsersReadOne(t *testing.T) {
	tests := []struct {
		name  string
		login string
	}{
		{
			name:  "read_Admin",
			login: "admin",
		},
		{
			name:  "create_User1",
			login: "user1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := storage.UserReadOne(tt.login)
			if err != nil {
				t.Error(err)
			}
			if user.Login != tt.login {
				t.Errorf("want login %s, get login %s", user.Login, tt.login)
			}
		})
	}
}

func TestOrdersCreate(t *testing.T) {
	tests := []struct {
		name    string
		number  string
		login   string
		wantErr bool
	}{
		{
			name:    "create_AdminOrder1",
			number:  "123",
			login:   "admin",
			wantErr: false,
		},
		{
			name:    "create_AdminOrder2",
			number:  "1234",
			login:   "admin",
			wantErr: false,
		},
		{
			name:    "create_User1Order1",
			number:  "123",
			login:   "user1",
			wantErr: true,
		},
		{
			name:    "create_User1Order2",
			number:  "1234",
			login:   "user1",
			wantErr: true,
		},
		{
			name:    "create_User1Order3",
			number:  "12345",
			login:   "user1",
			wantErr: false,
		},
		{
			name:    "create_User2Order3",
			number:  "12345",
			login:   "user2",
			wantErr: true,
		},
		{
			name:    "create_UserOrder4",
			number:  "123456",
			login:   "user2",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.UserOrderCreate(tt.login, tt.number); err != nil {
				if tt.wantErr {
					t.Skipf("want err. number %s already exists. err: %s", tt.number, err.Error())
				} else {
					t.Errorf("not created order login %s number %s", tt.login, tt.number)
					t.Error(err)
				}
			}
		})
	}
}

func TestOrderReadOne(t *testing.T) {
	tests := []struct {
		name   string
		number string
	}{
		{
			name:   "read_order_123",
			number: "123",
		},
		{
			name:   "read_order_1234",
			number: "1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := storage.UserOrderReadOne(tt.number)
			if err != nil {
				t.Error(err)
			}
			if user.Number != tt.number {
				t.Errorf("want number %s, get login %s", user.Number, tt.number)
			}
		})
	}
}

func TestOrderReadAllByUser(t *testing.T) {
	tests := []struct {
		name      string
		userLogin string
	}{
		{
			name:      "read_orders_admin",
			userLogin: "admin",
		},
		{
			name:      "read_orders_user1",
			userLogin: "user1",
		},
		{
			name:      "read_orders_user2",
			userLogin: "user2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders, err := storage.UserOrdersReadByLogin(tt.userLogin)
			if err != nil {
				t.Error(err)
			}
			if len(orders) == 0 {
				t.Errorf("Orders not found. storagemart.User login: %s", tt.userLogin)
			}
			t.Logf("Order user %s %+v", tt.userLogin, orders)
		})
	}
}
