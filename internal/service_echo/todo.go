package service_echo

import (
	"fmt"
	"github.com/harsha-aqfer/todo/pkg"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func getID(c echo.Context) (int64, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

func createTodo(c echo.Context) error {
	s := c.Get("service").(*Service)

	var req pkg.TodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := s.db.Todo.CreateTodo(&req); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func getTodo(c echo.Context) error {
	s := c.Get("service").(*Service)

	todoID, err := getID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	todo, err := s.db.Todo.GetTodo(todoID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todo)
}

func listTodos(c echo.Context) error {
	var (
		s   = c.Get("service").(*Service)
		all = c.QueryParam("all") == "true"
	)
	todos, err := s.db.Todo.ListTodos(all)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todos)
}

func updateTodo(c echo.Context) error {
	s := c.Get("service").(*Service)

	todoId, err := getID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var req pkg.TodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.IsZero() {
		return echo.NewHTTPError(http.StatusBadRequest, "empty body is not supported")
	}

	if err = s.db.Todo.UpdateTodo(todoId, &req); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func deleteTodo(c echo.Context) error {
	s := c.Get("service").(*Service)

	todoId, err := getID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = s.db.Todo.DeleteTodo(todoId); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
