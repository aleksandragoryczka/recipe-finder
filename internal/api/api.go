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
	Key     = "e70bfeb0be7146aa951be98ce9083f08"
)

type HttpClient struct {
	httpClient *http.Client
}

type Recipe struct {
	Id                int
	Title             string
	MissedIngredients []Ingredient
	UsedIngredients   []Ingredient
	Calories          Nutrient
	Proteins          Nutrient
	Carbs             Nutrient
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
	recipes := make([]Recipe, len(recipesInfo))
	for _, recipesInfo := range recipesInfo {
		nutritionList, err := httpClient.GetRecipeNutritionsInfo(recipesInfo.Id)
		if err != nil {
			return nil, err
		}
		recipe := Recipe{
			Id:                recipesInfo.Id,
			Title:             recipesInfo.Title,
			MissedIngredients: recipesInfo.MissedIngredients,
			UsedIngredients:   recipesInfo.UsedIngredients,
			Calories:          nutritionList[0],
			Proteins:          nutritionList[1],
			Carbs:             nutritionList[2],
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (httpClient *HttpClient) GetRecipeNutritionsInfo(recipeId int) ([]Nutrient, error) {
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

	var calories, proteins, carbs Nutrient
	for _, nutrient := range nutritionsList.Nutrition.Nutrients {
		switch nutrient.Name {
		case "Calories":
			calories.Name = "Calories"
			calories.Amount = nutrient.Amount
		case "Protein":
			proteins.Name = "Proteins"
			proteins.Amount = nutrient.Amount
		case "Carbohydrates":
			carbs.Name = "Carbs"
			carbs.Amount = nutrient.Amount
		}
	}
	return []Nutrient{calories, proteins, carbs}, nil
}
