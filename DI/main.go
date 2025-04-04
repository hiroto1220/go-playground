package main

import (
	"github.com/hiroto1220/go-playground/di/controller"
	"github.com/hiroto1220/go-playground/di/model/implements"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// データベースの初期化
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//NewTodoModelでinterfaceを作成している
	todoModel := implements.NewTodoModel(db)
	//interface(依存性)をcontroller層へ注入している
	todoController := controller.NewTodoController(todoModel)
	todoController.FetchTodos()
}
