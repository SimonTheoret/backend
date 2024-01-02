package back

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// Builds the query from a request. Returns a pointer to a FrontEndQuery object.
func BuildQuery(c *gin.Context, typeOfQuery queryType) FrontEndQuery {
	var reqContent Json
	c.BindJSON(reqContent)
	id := xid.New().String() // unique ID
	return FrontEndQuery{reqContent, typeOfQuery, id}
}
