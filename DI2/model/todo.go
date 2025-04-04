package model

import "gorm.io/gorm"

type Todo struct {
	ID    int
	Title string
}

type TodoModel struct {
	DB *gorm.DB
}

func NewTodoModel(db *gorm.DB) *TodoModel {
	return &TodoModel{DB: db}
}

func (m *TodoModel) FetchTodos() ([]Todo, error) {
	var todos []Todo
	if err := m.DB.Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}
