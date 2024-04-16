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
	// "time"
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
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(form.File))
	// Define a channel to receive processed image information
	processedImages := make(chan map[string]interface{}, len(form.File)) // Adjust capacity based on number of file fields

	// Loop through each image field in the form
	for fieldName, files := range form.File {
		for _, file := range files {
			go func(fieldName string, file *multipart.FileHeader) {
				defer wg.Done() // Signal completion after processing
				ImageProcess(c, fieldName, file, processedImages)
			}(fieldName, file)
		}
	}

	
	// Collect and respond with processed image information
	wg.Wait()
	close(processedImages) // Close the channel after all sends

	var breakfastResult map[string]interface{}
	var lunchResult map[string]interface{}
	var dinnerResult map[string]interface{}

	for processedImage := range processedImages {
		fieldName := processedImage["fieldName"].(string)
		if fieldName == "breakfast_img" {
			breakfastResult = processedImage["data"].(map[string]interface{})
		} else if fieldName == "lunch_img" {
			lunchResult = processedImage["data"].(map[string]interface{})
		} else if fieldName == "dinner_img" {
			dinnerResult = processedImage["data"].(map[string]interface{})
		}
	}

	email := c.PostForm("email")
	date := c.PostForm("date")
	breakfast := c.PostForm("breakfast")
	lunch := c.PostForm("lunch")
	dinner := c.PostForm("dinner")
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
	fmt.Println(formData)
	formData["weight"] = weight

	err = DB.AppendMeals(formData)
	if err != nil {
		fmt.Println("Couldn't track calories", formData, err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": err})
		return

	}
	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

func ImageProcess(c *gin.Context, fieldName string, file *multipart.FileHeader, processedImages chan map[string]interface{}) {
	// fmt.Println("image processing")
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer openedFile.Close()

	// Read the entire file content into a byte slice
	imageData, err := ioutil.ReadAll(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Identify image type using magic bytes
	imageType := http.DetectContentType(imageData)

	switch imageType {
	case "image/png":
		imageType = "png"
	case "image/webp":
		imageType = "webp"
	case "image/jpeg":
		imageType = "png"
	default:
		imageType = "png"
	}

	result_map, err := AI.GenImageAI(imageData, imageType)
	if err != nil {
		fmt.Println("[Image processing]", err)
	}

	processedImages <- map[string]interface{}{
		"fieldName": fieldName,
		"data":      result_map, // Include the existing data map
	}
	
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
