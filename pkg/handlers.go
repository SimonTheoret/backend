package back

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlerModelPost(mm *modelMapper) gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(
				http.StatusBadRequest,
				Json{
					"content": fmt.Sprintf("No model with matching id. Given id is empty: %s", id)})
		} else {
			inChan := mm.InputChannels[Id(id)]
			inChan <- BuildQuery(c, UnknownQuery)
		}
	}
}

func HandlerModelGet(mm *modelMapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(
				http.StatusBadRequest,
				Json{
					"content": fmt.Sprintf("No model with matching id. Given id is empty: %s", id)})
		} else {
			inChan := mm.InputChannels[Id(id)]
			inChan <- BuildQuery(c, GetLogs)
		}
	}
}
