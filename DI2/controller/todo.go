package controller

import (
	"log"

	"github.com/hiroto1220/go-playground/di2/model"
)

type TodoController struct {
	Model *model.TodoModel
}

// model層でgorm.DBを使っていることを知っている。つまりmodel層に依存している。
func NewTodoController(m *model.TodoModel) *TodoController {
	return &TodoController{Model: m}
}

func (c *TodoController) FetchTodos() ([]model.Todo, error) {
	posts, err := c.Model.FetchTodos()
	if err != nil {
		log.Println("Failed to get todos: ", err)
		return nil, err
	}
	return posts, nil
}
