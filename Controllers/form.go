package Controllers

import (
	AI "blissfulbites/AI"
	DB "blissfulbites/DB"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"
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

	c.Redirect(http.StatusFound, "/dashboard")
}

func AppendMealsHandler(c *gin.Context) {

	var wg sync.WaitGroup

	breakfastChan := make(chan map[string]interface{})
	lunchChan := make(chan map[string]interface{})
	dinnerChan := make(chan map[string]interface{})

	wg.Add(3)

	go func() {
		fmt.Println("image processing 1 ")

		defer wg.Done()
		breakfastMap := ImageProcess(c, "breakfast_img")
		breakfastChan <- breakfastMap
	}()

	go func() {
		fmt.Println("image processing2")
		time.Sleep(3 * time.Second)
		defer wg.Done()
		lunchMap := ImageProcess(c, "lunch_img")
		lunchChan <- lunchMap
	}()

	go func() {
		fmt.Println("image processing3")
		time.Sleep(7 * time.Second)
		defer wg.Done()
		dinnerMap := ImageProcess(c, "dinner_img")
		dinnerChan <- dinnerMap
	}()
	fmt.Println("idhar ")

	breakfastResult := <-breakfastChan
	lunchResult := <-lunchChan
	dinnerResult := <-dinnerChan
	
	wg.Wait()
	fmt.Println("khatam")


	email := c.PostForm("email")
	date := c.PostForm("date")
	breakfast := c.PostForm("breakfast")
	lunch := c.PostForm("lunch")
	dinner := c.PostForm("dinner")
	fmt.Println(breakfast)
	weight := c.PostForm("weight")

	formData := make(map[string]interface{})
	formData["email"] = email
	formData["date"] = date
	if breakfast == "" && breakfastResult != nil {
		formData["breakfast"] = breakfastResult
	} else if breakfast != "" && breakfastResult == nil {
		formData["breakfast"] = breakfast
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Breakfast value empty."})
		return
	}

	if lunch == "" && lunchResult != nil {
		formData["lunch"] = lunchResult
	} else if lunch != "" && lunchResult == nil {
		formData["lunch"] = lunch
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lunch value empty."})
		return
	}

	if dinner == "" && dinnerResult != nil {
		formData["dinner"] = dinnerResult
	} else if lunch != "" && dinnerResult == nil {
		formData["dinner"] = dinner
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dinner value empty."})
		return
	}

	formData["weight"] = weight

	err := DB.AppendMeals(formData)
	if err != nil {
		fmt.Println("Couldn't track calories", formData, err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": err})
		return

	}
	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

func ImageProcess(c *gin.Context, id string) map[string]interface{} {
	// fmt.Println("image processing")
	file, header, err := c.Request.FormFile(id)
	if err != nil {
		fmt.Println("image processing phase2", id, err)

		return nil
	}
	defer file.Close()
	fmt.Println("image processing phase2")

	breakfast_imgType := GetMime(header)
	switch breakfast_imgType {
	case "image/png":
		breakfast_imgType = "png"
	case "image/webp":
		breakfast_imgType = "webp"
	case "image/jpeg":
		breakfast_imgType = "png"
	default:
		breakfast_imgType = "png"
	}

	breakfast_data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	breakfast_map, err := AI.GenImageAI(breakfast_data, breakfast_imgType)
	if err != nil {
		return nil

	}

	return breakfast_map
}

func GetMime(header *multipart.FileHeader) string {
	contentType := header.Header.Get("Content-Type")
	mime, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		fmt.Println("Error parsing content type:", err)
		return ""
	}

	fmt.Println("MIME type:", mime)
	return mime
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
		fmt.Println("[Update diet handler]", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "couldn't get updated"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "plan updated"})
}

func GenDietPlan(c *gin.Context) {
	var userData map[string]interface{}
	if err := c.ShouldBindJSON(&userData); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonBytes, _ := json.Marshal(userData)
	userString := string(jsonBytes)

	res, err := AI.GenAI(userString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "couldn't generate diet plan"})
		return
	}
	err = DB.UpdateDiet(userData["email"].(string), res, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "coudln't store in db"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"diet_plan": res})

}
