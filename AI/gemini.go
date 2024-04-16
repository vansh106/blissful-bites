package AI

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var client *genai.Client

// var ctx *context.Context

func InitializeModel(key string) {
	ctx := context.Background()
	client, _ = genai.NewClient(ctx, option.WithAPIKey(key))
	fmt.Println("gemini initialized")
}

func GenImageAI(imgData []byte, imgType string) (map[string]interface{}, error) {
	vmodel := client.GenerativeModel("gemini-pro-vision")
	ctx := context.Background()

	prompt := `
You are an expert in nutritionist where you need to see the food items from the image
               and calculate the total calories, also provide the details of every food items with calories intake
               is below JSON format. 
			   NOTE: Don't include anything apart from below format, like number of food items or irrelevant text )
			   NOTE: Don't use bold words
			   NOTE: Don't include "'''JSON" in the output
			   {
               	"Food Item 1": no of calories (integer)
                "Food Item 2": no of calories (integer)
				"Total calories": (integer)
			   }
`
	resp, err := vmodel.GenerateContent(ctx, genai.Text(prompt), genai.ImageData(imgType, imgData))
	if err != nil {
		return nil,err

	}
	fmt.Println(resp.Candidates[0].Content.Parts[0])
	var jsonData []byte
	if jsonData, err = json.Marshal(resp); err != nil {
		return nil,err
	}
	stringData := string(jsonData)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(stringData), &data); err != nil {
		fmt.Println("Error:", err)
		return nil,err

	}
	value := data["Candidates"].([]interface{})[0].(map[string]interface{})["Content"].(map[string]interface{})["Parts"].([]interface{})[0].(string)
	
	var meal map[string]interface{}
	if err := json.Unmarshal([]byte(value), &meal); err != nil {
		fmt.Println("Error:", err)
		return nil,err
	}
	return meal, nil
}

func GenAI(userData string) (string, error) {
	model := client.GenerativeModel("gemini-pro")
	ctx := context.Background()

	prompt := fmt.Sprintf(`
You are an expert in nutritionist, named "Mannat", where you assess users data ( along with the calories if they have been tracking) and suggest very very short diet plans ( in a nutshell, in few lines~) accordingly.
%s
`, userData)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "",err

	}
	// fmt.Println(resp.Candidates[0].Content.Parts[0])
	var jsonData []byte
	if jsonData, err = json.Marshal(resp); err != nil {
		return "",err
	}
	stringData := string(jsonData)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(stringData), &data); err != nil {
		fmt.Println("Error:", err)
		return "",err

	}
	value := data["Candidates"].([]interface{})[0].(map[string]interface{})["Content"].(map[string]interface{})["Parts"].([]interface{})[0].(string)
	return value, nil
}
