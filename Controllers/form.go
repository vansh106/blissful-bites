package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	DB "blissfulbites/DB"
)

func FormHandler(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}
	formData := make(map[string]interface{})
	for key, values := range c.Request.PostForm {
		if len(values) == 1 {
			formData[key] = values[0]
		} else {
			formData[key] = values
		}
	}

	err := DB.InsertUserData(formData)
	if err != nil {
		fmt.Println("[FORM handler]", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error to store form"})
	}

	c.Redirect(http.StatusSeeOther, "/dashboard")
}

func AppendMealsHandler(c *gin.Context) {
	email := c.PostForm("email")
	date := c.PostForm("date")
	breakfast := c.PostForm("breakfast")
	lunch := c.PostForm("lunch")
	dinner := c.PostForm("dinner")
	weight := c.PostForm("weight")

	formData := make(map[string]interface{})
	formData["email"]=email
	formData["date"]=date
	formData["breakfast"]=breakfast
	formData["lunch"]=lunch
	formData["dinner"]=dinner
	formData["weight"]=weight

	err := DB.AppendMeals(formData)
	if err != nil {
		fmt.Println("Couldn't track calories", formData, err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": err})
		return 

	}
	c.JSON(http.StatusOK, gin.H{"status":"message sent"})


}
