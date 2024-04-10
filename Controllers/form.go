package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	DB "blissfulbites/DB"
)

func FormHandler(c *gin.Context) {
	number := c.PostForm("age")
	fmt.Println(number)
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

	c.Redirect(http.StatusFound, "/dashboard")
}

func AppendMealsHandler(c *gin.Context) {
	email := c.PostForm("email")
	date := c.PostForm("date")
	breakfast := c.PostForm("breakfast")
	lunch := c.PostForm("lunch")
	dinner := c.PostForm("dinner")
	weight := c.PostForm("weight")

	formData := make(map[string]interface{})
	formData["email"] = email
	formData["date"] = date
	formData["breakfast"] = breakfast
	formData["lunch"] = lunch
	formData["dinner"] = dinner
	formData["weight"] = weight

	err := DB.AppendMeals(formData)
	if err != nil {
		fmt.Println("Couldn't track calories", formData, err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": err})
		return

	}
	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

func UpdateDietHandler(c *gin.Context) {
	email := c.PostForm("email")
	diet := c.PostForm("diet_plan")
	healthscore := c.PostForm("healthscore")
	hs, err := strconv.Atoi(healthscore)
	if err != nil {
		hs = 0 
	}
	fmt.Println(email, diet, hs)
	err = DB.UpdateDiet(email, diet, hs)
	if err != nil {
		fmt.Println("[Update diet handler]",err)
		c.JSON(http.StatusInternalServerError, gin.H{"status":"couldn't get updated"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status":"plan updated"})
}
