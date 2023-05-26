package handlers

import (
	old_err "errors"
	"github.com/blessium/todo-sample/errors"
	"github.com/blessium/todo-sample/services"
	"github.com/blessium/todo-sample/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var (
	HandlerRequiredTitleError = "todo schema: \"title\" field is required"
	HandlerFormatDueTimeError = "todo schema: Title field is required"
)

type TodoHandler struct {
	todoService services.ITodoService
}

func NewTodoHandler(s services.ITodoService) TodoHandler {
	return TodoHandler{
		todoService: s,
	}
}

type TodoAddRequest struct {
	Title   string    `json:"title"`
	Notes   string    `json:"notes"`
	DueDate time.Time `json:"due_date"`
}

func (t TodoAddRequest) Validate() error {
	if t.Title == "" {
		return old_err.New(HandlerRequiredTitleError)
	}
	return nil
}

func (t TodoAddRequest) mapToService() services.Todo {
	return services.Todo{
		Title:   t.Title,
		Notes:   t.Notes,
		DueDate: t.DueDate,
	}
}

type TodoUpdateRequest struct {
	Title     string    `json:"title"`
	Notes     string    `json:"notes"`
	DueDate   time.Time `json:"due_date"`
	Completed bool      `json:"completed"`
}

type TodosUpdateRequest struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Notes     string    `json:"notes"`
	DueDate   time.Time `json:"due_date"`
	Completed bool      `json:"completed"`
}

func (t TodosUpdateRequest) mapToService() services.Todo {
	return services.Todo{
		ID:        t.ID,
		Title:     t.Title,
		Notes:     t.Notes,
		DueDate:   t.DueDate,
		Completed: t.Completed,
	}
}

func (t TodoUpdateRequest) mapToService() services.Todo {
	return services.Todo{
		Title:     t.Title,
		Notes:     t.Notes,
		DueDate:   t.DueDate,
		Completed: t.Completed,
	}
}

type TodoFullResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Notes        string    `json:"notes"`
	CreationDate time.Time `json:"creation_date"`
	DueDate      time.Time `json:"due_date"`
	Completed    bool      `json:"completed"`
}

func FullResponseFromService(t services.Todo) TodoFullResponse {
	return TodoFullResponse{
		ID:           t.ID,
		Title:        t.Title,
		Notes:        t.Notes,
		CreationDate: t.CreationDate,
		DueDate:      t.DueDate,
		Completed:    t.Completed,
	}
}

func handleError(e error, path string) (int, errors.HttpErrorResponse) {
	var status int
	e, ok := e.(errors.Error)
	if !ok {
		status = http.StatusInternalServerError
		return status, errors.NewHttpErrorResponse(uint(status), "internal error", "internal server error", path)
	}

	if errors.IsType(e, errors.ErrInternal) {
		status = http.StatusInternalServerError
	} else if errors.IsType(e, errors.ErrValidation) {
		status = http.StatusBadRequest
	} else if errors.IsType(e, errors.ErrNotExist) {
		status = http.StatusNotFound
	}

	return status, errors.NewHttpErrorResponse(uint(status), e.(errors.Error).Type.String(), e.Error(), path)
}

func (t TodoHandler) AddTodo(c echo.Context) error {
	var r TodoAddRequest
	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := r.Validate(); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	todo := r.mapToService()
	todo, err := t.todoService.AddTodo(todo)
	if err != nil {
		return c.JSON(handleError(err, "/todos"))
	}

	return c.JSON(http.StatusCreated, FullResponseFromService(todo))
}

func (t TodoHandler) DeleteTodo(c echo.Context) error {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := t.todoService.DeleteTodo(id); err != nil {
		return c.JSON(handleError(err, "/todos/{id}"))
	}

	return c.NoContent(http.StatusNoContent)
}

func (t TodoHandler) DeleteTodos(c echo.Context) error {
	if err := t.todoService.DeleteTodos(); err != nil {
		return c.JSON(handleError(err, "/todos"))
	}
	return c.NoContent(http.StatusNoContent)
}

func (t TodoHandler) UpdateTodo(c echo.Context) error {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var r TodoUpdateRequest
	if err := c.Bind(&r); err != nil {
		return err
	}

	todo, err := t.todoService.UpdateTodo(id, r.mapToService())
	if err != nil {
		return c.JSON(handleError(err, "/todos/2"))
	}

	return c.JSON(http.StatusOK, FullResponseFromService(todo))
}

func (t TodoHandler) UpdateTodos(c echo.Context) error {
	var r []TodoUpdateRequest
	if err := c.Bind(&r); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var ser_todos []services.Todo
	for _, todo := range r {
		ser_todos = append(ser_todos, todo.mapToService())
	}

	todos, err := t.todoService.UpdateTodos(ser_todos)
	if err != nil {
		return c.JSON(handleError(err, "/todos"))
	}

	var result_todos []TodoFullResponse
	for _, todo := range todos {
		result_todos = append(result_todos, FullResponseFromService(todo))
	}

	return c.JSON(http.StatusOK, result_todos)
}

func (t TodoHandler) GetTodo(c echo.Context) error {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	todo, err := t.todoService.GetTodo(id)
	if err != nil {
		return c.JSON(handleError(err, "/todos/{id}"))
	}

	return c.JSON(http.StatusFound, FullResponseFromService(todo))
}

func (t TodoHandler) GetTodos(c echo.Context) error {
	todos, err := t.todoService.GetTodos()
	if err != nil {
		return c.JSON(handleError(err, "/todos"))
	}
	var r_todos []TodoFullResponse
	for _, todo := range todos {
		r_todos = append(r_todos, FullResponseFromService(todo))
	}

	return c.JSON(http.StatusFound, r_todos)
}
