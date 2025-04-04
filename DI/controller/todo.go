package controller

import (
	"fmt"
	"log"

	"github.com/hiroto1220/go-playground/di/model"
	"github.com/hiroto1220/go-playground/di/model/iface"
)

type TodoController struct {
	Model iface.TodoModeler
}

// model.TodoModeler(interface)に依存している。つまり、model層の実装は知らず、interfaceを実装していれば引数に取れる。
func NewTodoController(m iface.TodoModeler) *TodoController {
	return &TodoController{Model: m}
}

func (c *TodoController) FetchTodos() ([]model.Todo, error) {
	posts, err := c.Model.FetchTodos()
	if err != nil {
		log.Println("Failed to get posts: ", err)
		return nil, err
	}
	fmt.Println(posts)
	return posts, nil
}
