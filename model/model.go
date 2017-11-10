package model

import "github.com/jinzhu/gorm"

type (
	// todoModel describes a todoModel type
	TodoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	// transformedTodo represents a formatted todo
	TransformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)
