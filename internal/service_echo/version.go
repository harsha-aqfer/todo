package service_echo

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func getVersion(c echo.Context) error {
	s := c.Get("service").(*Service)
	version, err := s.db.Version.GetVersion()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"version": version})
}
