package iface

import (
	"github.com/hiroto1220/go-playground/di/model"
)

type TodoModeler interface {
	FetchTodos() ([]model.Todo, error)
}
