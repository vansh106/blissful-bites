package Controllers

import (
	DB "blissfulbites/DB"
	// "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ContactHandler(c *gin.Context) {
	email := c.PostForm("email")
	message := c.PostForm("message")

	err := DB.UpdateDM(email, message)
	if err != nil {
		fmt.Println("[Contact handler]",err)
		c.JSON(http.StatusInternalServerError, gin.H{"status":"message sent"})
	}

	c.JSON(http.StatusOK, gin.H{"status":"message sent"})

}
