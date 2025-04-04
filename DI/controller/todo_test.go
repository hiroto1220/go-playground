package controller

import (
	"reflect"
	"testing"

	"github.com/hiroto1220/go-playground/di/model"
)

type MockTodoModel struct {
	models []model.Todo
	err    error
}

func (m *MockTodoModel) FetchTodos() ([]model.Todo, error) {
	return m.models, m.err
}

// テストコード
func TestTodoController(t *testing.T) {
	// モックのセットアップ
	mockModel := &MockTodoModel{
		models: []model.Todo{
			{ID: 1, Title: "Task1"},
			{ID: 2, Title: "Task2"},
		},
		err: nil,
	}

	// コントローラのセットアップ
	controller := NewTodoController(mockModel)

	// メソッドの実行
	todos, err := controller.FetchTodos()
	if err != nil {
		t.Error("Expected no error, got ", err)
	}

	expected := []model.Todo{
		{ID: 1, Title: "Task1"},
		{ID: 2, Title: "Task2"},
	}

	if !reflect.DeepEqual(todos, expected) {
		t.Errorf("Expected %v, got %v", expected, todos)
	}
}
