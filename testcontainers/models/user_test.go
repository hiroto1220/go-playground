package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUser_FetchUserByID(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				userID: "1",
			},
			want: &User{
				ID:    "1",
				Name:  "yamada taro",
				Email: "taro.yamada@example.com",
			},
			wantErr: false,
		},
	}

	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("test"),
		mysql.WithUsername("user"),
		mysql.WithPassword("password"),
		mysql.WithScripts("./schema.sql"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("3306/tcp").WithStartupTimeout(5*time.Minute)),
	)

	defer func() {
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatal(err)
	}

	// ← ここで接続文字列を動的に組み立てる
	dsn := fmt.Sprintf("user:password@tcp(%s:%s)/test", host, port.Port())

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer db.Close()

	r := NewUserRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FetchUserByID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.FetchUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.FetchUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
