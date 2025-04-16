package models

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUser_FetchUserByIDWithModule(t *testing.T) {
	// Testconainersを使用してMySQLコンテナを準備
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithScripts("./schema.sql"),

		// データベース名やユーザー名、パスワードを指定することも可能
		// mysql.WithDatabase("test"),
		// mysql.WithUsername("user"),
		// mysql.WithPassword("password"),
	)

	defer func() {
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	if err != nil {
		t.Fatalf("failed to start MySQL container: %s", err)
	}

	// MySQLに接続
	dsn, err := mysqlContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	// テストデータを作成
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
			name: "userIDが1の時にデータベースからuserID 1に対応したユーザーを取得できる",
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
		{
			name: "データベースに登録していないuserIDを指定した時にエラーが発生する",
			args: args{
				userID: "999",
			},
			want:    nil,
			wantErr: true,
		},
	}

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

func TestUser_FetchUserByID(t *testing.T) {
	// Testconainersを使用してMySQLコンテナを準備
	ctx := context.Background()

	schemaPath, err := filepath.Abs("../models/schema.sql")
	if err != nil {
		t.Fatal(err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0.36",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_USER":          "user",
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_PASSWORD":      "password",
			"MYSQL_DATABASE":      "test",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server"),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      schemaPath,
				ContainerFilePath: "/docker-entrypoint-initdb.d/schema.sql",
				FileMode:          0o644,
			},
		},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	mysqlContainer, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	if err != nil {
		t.Fatalf("failed to start MySQL container: %s", err)
	}

	// コンテナのホストとポートを取得
	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatal(err)
	}

	// MySQLに接続
	dsn := fmt.Sprintf("user:password@tcp(%s:%s)/test", host, port.Port())

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// データベースに接続できるか確認
	err = db.Ping()
	if err != nil {
		t.Fatalf("failed to connect to database: %s", err)
	} else {
		t.Log("データベースに接続できました")
	}

	r := NewUserRepository(db)

	// テストデータを作成
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
			name: "userIDが1の時にデータベースからuserID 1に対応したユーザーを取得できる",
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
		{
			name: "データベースに登録していないuserIDを指定した時にエラーが発生する",
			args: args{
				userID: "999",
			},
			want:    nil,
			wantErr: true,
		},
	}

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
