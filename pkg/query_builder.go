package back

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// Builds the query from a request. Returns a pointer to a message[queryType] object.
func BuildQuery(c *gin.Context, typeOfQuery string, receiverId Id) *message[queryType] {
	var reqContent Json
	var q queryType
	c.BindJSON(reqContent)
	switch typeOfQuery {
	case "predict":
		q = Predict

	case "getlogs":
		q = GetLogs

	case "cleanlogs":
		q = CleanLogs
	}
	ip := c.ClientIP()

	id := Id(xid.New()) // unique ID
	return &message[queryType]{
		content:         reqContent,
		messageType:     q,
		id:              id,
		creationTime:    time.Now(),
		receiverId:      id,
		sender:          ip,
		queryOptions:    nil,
		responseOptions: nil,
	}
}
