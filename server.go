// main.go

package main

import (
	DB "blissfulbites/DB"
	Controllers "blissfulbites/Controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)


func main() {

	_= DB.ConnectPsql()
	_= DB.CreateTableUserDetails()

	r := gin.Default()

	// r.LoadHTMLGlob("static/client/*")
    r.LoadHTMLGlob("static/**/*.html")
	r.Static("/css", "static/css")
	r.Static("/logos", "static/logos")

    //client endpoints
	r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "home.html", gin.H{})
    })

    r.GET("/login", func(c *gin.Context) {
        c.HTML(http.StatusOK, "auth.html", gin.H{})
    })

    r.GET("/signup", func(c *gin.Context) {
        c.HTML(http.StatusOK, "signup.html", gin.H{})
    })

	r.GET("/dashboard", func(c *gin.Context) {
        c.HTML(http.StatusOK, "dashboard.html", gin.H{})
    })

	r.GET("/form", func(c *gin.Context) {
        c.HTML(http.StatusOK, "form.html", gin.H{})
    })

    r.GET("/track", func(c *gin.Context) {
        c.HTML(http.StatusOK, "track.html", gin.H{})
    })

    r.GET("/contact", func(c *gin.Context) {
        c.HTML(http.StatusOK, "contact.html", gin.H{})
    })

    r.POST("/contactUs", func(c *gin.Context) {
        Controllers.ContactHandler(c)
    })

	r.POST("/userFormDetails", func(c *gin.Context) {
		Controllers.FormHandler(c)
    })

    r.POST("/trackMeal", func(c *gin.Context) {
		Controllers.AppendMealsHandler(c)
    })

    r.GET("/userDetails", func(c *gin.Context) {
		Controllers.UserDataHandler(c)
    })

	r.GET("/firstlogin", func(c *gin.Context) {
        Controllers.FirstLoginHandler(c)
    })

    // admin endpoints
    
    r.GET("/admin", func(c *gin.Context) {
        Controllers.AllUsersDataHandler(c)
    })

    r.GET("/dm", func(c *gin.Context) {
        Controllers.DmHandler(c)
    })

    r.GET("/user", func(c *gin.Context) {
        c.HTML(http.StatusOK, "user.html", gin.H{})
    })
	
    r.POST("/updateDiet", func(c *gin.Context) {
        Controllers.UpdateDietHandler(c)
    })
	
	// r.Static("/static", "./static")
	r.Run(":8080")
}
