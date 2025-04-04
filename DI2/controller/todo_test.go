package controller

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hiroto1220/go-playground/di2/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFetchTodos(t *testing.T) {
	// モックDBの作成
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// sqlite_version()を設定
	// SQLiteを使っていることを知ってしまっている
	mock.ExpectQuery(`select sqlite_version\(\)`).WillReturnRows(
		sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.36.0"),
	)

	// gorm.DB を初期化
	gormDB, err := gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create gorm.DB: %v", err)
	}

	// モックのクエリと返り値の設定
	mock.ExpectQuery("SELECT \\* FROM `todos`").WillReturnRows(
		sqlmock.NewRows([]string{"id", "title"}).
			AddRow(1, "Task1").
			AddRow(2, "Task2"),
	)

	// TodoControllerのセットアップ
	m := model.NewTodoModel(gormDB)
	controller := NewTodoController(m)

	// メソッドをテスト
	todos, err := controller.FetchTodos()

	// エラーがないことを確認
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// 期待する結果
	expected := []model.Todo{
		{ID: 1, Title: "Task1"},
		{ID: 2, Title: "Task2"},
	}

	// reflect.DeepEqual で結果を検証
	if !reflect.DeepEqual(todos, expected) {
		t.Errorf("Expected %v, got %v", expected, todos)
	}

	// モックが期待通りに呼び出されたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %v", err)
	}
}
