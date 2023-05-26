package main 

import (
	"github.com/blessium/todo-sample/database"
	"github.com/blessium/todo-sample/handlers"
	"github.com/blessium/todo-sample/repositories"
	"github.com/blessium/todo-sample/services"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	db, err := database.NewGormDatabase()
	if err != nil {
		panic(err.Error())
	}

	todoRepo := repositories.NewGormTodoRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	todos := e.Group("todos")
	todos.POST("", todoHandler.AddTodo)
    todos.GET("", todoHandler.GetTodos)
    todos.PUT("", todoHandler.UpdateTodos)
    todos.DELETE("", todoHandler.DeleteTodos)

	todo := todos.Group("/:id")
	todo.GET("", todoHandler.GetTodo)
	todo.PUT("", todoHandler.UpdateTodo)
	todo.DELETE("", todoHandler.DeleteTodo)
}
