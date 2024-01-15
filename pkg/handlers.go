package back

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func HandlerModelPost(mm *modelMapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
        qtype:= c.Query("operation")
		if id == "" {
			c.JSON(
				http.StatusBadRequest,
				Json{
					"content": fmt.Sprintf("No model with matching id. Given id is empty: %s", id),
				})
		} else {
			nId, _ := xid.FromString(id)
			inChan := mm.InputChannels[Id(nId)]
			inChan <- BuildQuery(c, qtype, Id(nId))
		}
	}
}

func HandlerModelGet(mm *modelMapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
        qtype := c.Query("operation")
		if id == "" {
			c.JSON(
				http.StatusBadRequest,
				Json{
					"content": fmt.Sprintf("No model with matching id. Given id is empty: %s", id),
				})
		} else {
			nId, _ := xid.FromString(id)
			inChan := mm.InputChannels[Id(nId)]
			inChan <- BuildQuery(c, qtype, Id(nId))
		}
	}
}
