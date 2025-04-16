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

type MysqlContainer struct {
	mysqlContainer testcontainers.Container
	URI            string
}

func setupMysqlContainerWithModule(ctx context.Context, schemaPath string) (*MysqlContainer, error) {
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithScripts(schemaPath),

		// github.com/testcontainers/testcontainers-go/modules/mysqlを使ってデータベース名やユーザー名、パスワードを指定することも可能
		// mysql.WithDatabase("test"),
		// mysql.WithUsername("user"),
		// mysql.WithPassword("password"),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	// MySQLに接続
	dsn, err := mysqlContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %s", err)
	}

	mysqlC := &MysqlContainer{
		mysqlContainer: mysqlContainer,
		URI:            dsn,
	}

	return mysqlC, nil
}

func setupMysqlContainer(ctx context.Context, schemaPath string) (*MysqlContainer, error) {
	// MySQLコンテナの起動設定を作成（環境変数・ポート・初期SQLファイルの配置など）
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

	// コンテナ起動リクエストを作成
	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	// コンテナを起動
	mysqlContainer, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	// コンテナのホストとポートを取得
	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %s", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %s", err)
	}

	// MySQLに接続
	dsn := fmt.Sprintf("user:password@tcp(%s:%s)/test", host, port.Port())

	mysqlC := &MysqlContainer{
		mysqlContainer: mysqlContainer,
		URI:            dsn,
	}

	return mysqlC, nil
}

func TestUser_FetchUserByID(t *testing.T) {
	ctx := context.Background()

	// コンテナ内で実行したいschema.sqlのパスを指定
	schemaPath, err := filepath.Abs("../models/schema.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Testcontainersを使ってMySQLコンテナを起動させる関数を呼び出しコンテナを起動
	// 関数の実装については後述
	mysqlContainer, err := setupMysqlContainerWithModule(ctx, schemaPath)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %s", err)
	}

	// テスト終了後にコンテナをクリーンアップ
	testcontainers.CleanupContainer(t, mysqlContainer.mysqlContainer)

	db, err := sql.Open("mysql", mysqlContainer.URI)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	// テストケースを作成
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

	// テスト実行
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
