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
		id := c.Query("id")
		qtype := c.Query("operation")
		if id != "" && qtype == "addmodel" {
            sender := BuildSender(c)
			nId, _ := xid.FromString(id)
            base := NewBase("", sender)
			mm.addNewModel(&base, Id(nId))
		}
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
