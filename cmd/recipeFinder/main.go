package main

import (
	"fmt"
	"github.com/aleksandragoryczka/recipeFinder/cmd/commandLineArguments"
	"github.com/aleksandragoryczka/recipeFinder/internal/api"
)

func main() {
	commandlinearguments.Execute()

	ingredients := []string{"tomatoes", "eggs"}
	maxRecipes := 1

	hc := api.NewHttpClinet()

	recipes, err := hc.GetRecipes(ingredients, maxRecipes)
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Println(recipes)
	}

}
