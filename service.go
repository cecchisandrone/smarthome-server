package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// createTodo add a new todo
func createTodo(c *gin.Context) {
	var todoJSON transformedTodo
	if err := c.ShouldBindJSON(&todoJSON); err == nil {
		todo := todoModel{Title: todoJSON.Title, Completed: todoJSON.Completed}
		db.Save(&todo)
		todoJSON.ID = todo.ID
		c.JSON(http.StatusCreated, todoJSON)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// fetchSingleTodo fetch a single todo
func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}
	c.JSON(http.StatusOK, _todo)
}

// fetchAllTodo fetch all todos
func fetchAllTodos(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo
	db.Find(&todos)
	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	//transforms the todos for building a good response
	for _, item := range todos {
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: item.Completed})
	}
	c.JSON(http.StatusOK, _todos)
}
