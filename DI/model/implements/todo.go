package implements

import (
	"github.com/hiroto1220/go-playground/di/model"
	"github.com/hiroto1220/go-playground/di/model/iface"
	"gorm.io/gorm"
)

type TodoModel struct {
	DB *gorm.DB
}

func NewTodoModel(db *gorm.DB) iface.TodoModeler {
	return &TodoModel{DB: db}
}

func (m *TodoModel) FetchTodos() ([]model.Todo, error) {
	var todos []model.Todo
	m.DB.Find(&todos)
	return todos, nil
}
