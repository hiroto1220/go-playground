package main

import (
	"fmt"

	"github.com/hiroto1220/go-playground/di2/controller"
	"github.com/hiroto1220/go-playground/di2/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// データベースの初期化
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// モデル層の初期化
	todoModel := model.NewTodoModel(db)

	// コントローラ層の初期化
	todoController := controller.NewTodoController(todoModel.DB)

	// コントローラのメソッドを呼び出し
	todos, err := todoController.FetchTodos()
	if err != nil {
		fmt.Printf("Error fetching todos: %v\n", err)
		return
	}

	fmt.Println("Fetched Todos:")
	for _, todo := range todos {
		fmt.Printf("- ID: %d, Title: %s\n", todo.ID, todo.Title)
	}
}
