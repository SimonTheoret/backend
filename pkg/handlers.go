package back

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func HandlerModelPost(mm *modelMapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		qid := c.Query("id")
		if qid == "" {
			c.JSON(
				http.StatusBadRequest,
				Json{
					"content": fmt.Sprintf(
						"No model with matching id. Given id is empty: %s",
						qid,
					),
				})
		} else {
			id, err := xid.FromString(qid)
			if err != nil {
				c.JSON(
					http.StatusBadRequest,
					Json{
						"content": fmt.Sprintf("Not an id: %s", qid),
					})
			}
			inChan, ok := mm.InputChannels[Id(id)]
			if !ok {
				c.JSON(
					http.StatusBadRequest,
					Json{
						"content": fmt.Sprintf("id is not registered: %s", id),
					})
			}
			inChan <- BuildQuery(c, UnknownQuery)
		}
	}
}

func HandlerModelGet(mm *modelMapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		qid := c.Query("id")
		if qid == "" {
			c.JSON(
				http.StatusBadRequest,
				Json{
					"content": fmt.Sprintf("No model with matching id. Given id is empty: %s", qid),
				})
		} else {
			id, err := xid.FromString(qid)
			if err != nil {
				c.JSON(
					http.StatusBadRequest,
					Json{
						"content": fmt.Sprintf("Not an id: %s", qid),
					})
			}
			inChan, ok := mm.InputChannels[Id(id)]
			if !ok {
				c.JSON(
					http.StatusBadRequest,
					Json{
						"content": fmt.Sprintf("id is not registered: %s", id),
					})
			}
			inChan <- BuildQuery(c, UnknownQuery)
		}
	}
}

// Builds a sender based on the type given during the addmodel operation. It defaults to
// building a HttpModel with destination "127.0.0.1:8080".
func BuildSender(c *gin.Context) Sender {
	t := strings.ToLower(c.Query("type"))
	var sender Sender
	switch t {
	case "httpmodel":
		dest := c.Query("destination")
		sender = &HttpModel{dest}

	default:
		dest := "127.0.0.1:8080"
		sender = &HttpModel{dest}
	}
	return sender
}
