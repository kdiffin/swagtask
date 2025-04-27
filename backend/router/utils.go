package router

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// funny right
func getIdAsStr(c echo.Context) (int, error) {
	idStr := c.Param("id")
	id, errConv := strconv.Atoi(idStr)

	return id, errConv
}