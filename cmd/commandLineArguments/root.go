package commandlinearguments

import (
	"fmt"
	"github.com/aleksandragoryczka/recipeFinder/internal/recipe"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
)

var ingredients string
var numberOfRecipes int
var rootCmd = &cobra.Command{
	Use:   "recipeFinder",
	Short: "recipeFinder - program generating meals from ingredients in your fridge",
	Long: `recipeFinder is a program that helps people avoid wasting food 
			while maintaining a healthy lifestyle by generating a list of meals
			that can be prepared from ingredients in your fridge 
			with minimal number of missing ingredients`,
	Run: func(cmd *cobra.Command, args []string) {
		i := strings.Split(ingredients, ",")
		sort.Strings(i)
		inputIngredients := strings.Join(i[:], ",")

		recipeService := recipe.NewService()

		recipes := recipeService.FindRecipeByIngredients(inputIngredients, numberOfRecipes)
		for _, recipe := range recipes {
			fmt.Printf("Name: %s\n", recipe.Title)
			fmt.Printf("Used Ingredients: %s\n", strings.Join(recipe.UsedIngredients, ", "))
			fmt.Printf("Missing Ingredients: %s\n", strings.Join(recipe.MissedIngredients, ", "))
			fmt.Printf("Calories: %.2f\n", recipe.Calories)
			fmt.Printf("Proteins: %.2f\n", recipe.Proteins)
			fmt.Printf("Carbs: %.2f\n", recipe.Carbs)
			fmt.Println(strings.Repeat("-", 70) + "\n")
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&ingredients, "ingredients", "", "Ingredients list")
	rootCmd.Flags().IntVar(&numberOfRecipes, "numberOfRecipes", 5, "Number of recipes")
	err := rootCmd.MarkFlagRequired("ingredients")
	if err != nil {
		return
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, "Error while executing your CLI: '%s", err)
		os.Exit(1)
	}

}
