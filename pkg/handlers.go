package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ModelPredict(c *gin.Context) {
	cCp := c.Copy()
	id := cCp.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Json{"content": fmt.Sprintf("No model with matching id. Given id is empty", id)})
	}
}
