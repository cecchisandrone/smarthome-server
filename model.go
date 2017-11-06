package main

import "github.com/jinzhu/gorm"

type (
	// todoModel describes a todoModel type
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	// transformedTodo represents a formatted todo
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)
