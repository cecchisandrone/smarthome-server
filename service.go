package main

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Db *gorm.DB `inject:""`
}

// createTodo add a new todo
func (s *Service) createTodo(c *gin.Context) {
	var todoJSON model.TransformedTodo
	if err := c.ShouldBindWith(&todoJSON, binding.JSON); err == nil {
		todo := model.TodoModel{Title: todoJSON.Title, Completed: todoJSON.Completed}
		s.Db.Save(&todo)
		todoJSON.ID = todo.ID
		c.JSON(http.StatusCreated, todoJSON)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// fetchSingleTodo fetch a single todo
func (s *Service) fetchSingleTodo(c *gin.Context) {
	var todo model.TodoModel
	todoID := c.Param("id")
	s.Db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	_todo := model.TransformedTodo{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}
	c.JSON(http.StatusOK, _todo)
}

// fetchAllTodo fetch all todos
func (s *Service) fetchAllTodos(c *gin.Context) {
	var todos []model.TodoModel
	var _todos []model.TransformedTodo
	s.Db.Find(&todos)
	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	//transforms the todos for building a good response
	for _, item := range todos {
		_todos = append(_todos, model.TransformedTodo{ID: item.ID, Title: item.Title, Completed: item.Completed})
	}
	c.JSON(http.StatusOK, _todos)
}

func (s *Service) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}
