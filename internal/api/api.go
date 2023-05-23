package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	BaseUrl = "https://api.spoonacular.com/recipes"
	Key     = "baaadd57c0b04b8dbba77bf63d2e6adc"
)

type HttpClient struct {
	httpClient *http.Client
}

type Recipe struct {
	Id                int
	Title             string
	MissedIngredients []string
	UsedIngredients   []string
	Calories          float64
	Proteins          float64
	Carbs             float64
}

type RecipeInfo struct {
	Id                int          `json:"id"`
	Title             string       `json:"title"`
	MissedIngredients []Ingredient `json:"missedIngredients"`
	UsedIngredients   []Ingredient `json:"usedIngredients"`
}

type Ingredient struct {
	Name string `json:"name"`
}

type NutritionList struct {
	Nutrition Nutrition `json:"nutrition"`
}

type Nutrient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type Nutrition struct {
	Nutrients []Nutrient `json:"nutrients"`
}

func NewHttpClinet() *HttpClient {
	return &HttpClient{
		httpClient: &http.Client{},
	}
}

func (httpClient *HttpClient) GetRecipes(passedIngredients []string, maxNumberOfRecipes int) ([]Recipe, error) {
	endpoint := fmt.Sprintf("%s/findByIngredients", BaseUrl)
	params := url.Values{}
	params.Set("apiKey", Key)
	params.Set("ingredients", strings.Join(passedIngredients, ",+"))
	params.Set("number", fmt.Sprintf("%v", maxNumberOfRecipes))
	uri := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	unescapedUri, _ := url.QueryUnescape(uri)
	request, err := http.NewRequest(http.MethodGet, unescapedUri, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var recipesInfo []RecipeInfo

	err = json.NewDecoder(resp.Body).Decode(&recipesInfo)
	if err != nil {
		return nil, err
	}
	recipes := make([]Recipe, 0)
	for _, recipesInfo := range recipesInfo {
		nutritionList, err := httpClient.GetRecipeNutritionsInfo(recipesInfo.Id)
		if err != nil {
			return nil, err
		}
		recipe := Recipe{
			Id:                recipesInfo.Id,
			Title:             recipesInfo.Title,
			MissedIngredients: FormatToString(recipesInfo.MissedIngredients),
			UsedIngredients:   FormatToString(recipesInfo.UsedIngredients),
			Calories:          nutritionList[0],
			Proteins:          nutritionList[1],
			Carbs:             nutritionList[2],
		}

		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func FormatToString(ingredients []Ingredient) []string {
	var toReturn []string
	for _, i := range ingredients {
		toReturn = append(toReturn, i.Name)
	}
	return toReturn
}

func (httpClient *HttpClient) GetRecipeNutritionsInfo(recipeId int) ([]float64, error) {
	//fmt.Println(recipeId)
	endpoint := fmt.Sprintf("%s/%d/information", BaseUrl, recipeId)
	params := url.Values{}
	params.Set("apiKey", Key)
	params.Set("includeNutrition", "true")
	uri := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	unescapedUri, _ := url.QueryUnescape(uri)
	request, _ := http.NewRequest(http.MethodGet, unescapedUri, nil)

	request.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var nutritionsList NutritionList
	err = json.NewDecoder(resp.Body).Decode(&nutritionsList)
	if err != nil {
		return nil, err
	}

	var calories, proteins, carbs float64
	for _, nutrient := range nutritionsList.Nutrition.Nutrients {
		switch nutrient.Name {
		case "Calories":
			calories = nutrient.Amount
		case "Protein":
			proteins = nutrient.Amount
		case "Carbohydrates":
			carbs = nutrient.Amount
		}
	}
	return []float64{calories, proteins, carbs}, nil
}
