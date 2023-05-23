package recipe

import (
	"fmt"
	"github.com/aleksandragoryczka/recipeFinder/internal/api"
	"github.com/aleksandragoryczka/recipeFinder/internal/database"
)

type Service struct {
	db *database.Database
}

func NewService() *Service {
	db, _ := database.NewDatabase()
	return &Service{db: db}
}

func (s *Service) FindRecipeByIngredients(ingredients []string, numberOfRecipes int) []api.Recipe {
	recipes, err := s.db.GetRecipeByIngredientsList(ingredients, numberOfRecipes)
	if err == nil && len(recipes) > 0 {
		err := s.db.CloseDatabaseConnection()
		if err != nil {
			fmt.Println("Error closing DB connection: ", err)
		}
		return recipes
	} else if err == nil {
		hc := api.NewHttpClient()
		recipes := hc.GetRecipes(ingredients, numberOfRecipes)
		if err != nil {
			fmt.Println("Error gettingRecipes from API: ", err)
			return nil
		}
		s.db.InsertTransaction(recipes)
		err := s.db.CloseDatabaseConnection()
		if err != nil {
			fmt.Println("Error closing DB connection: ", err)
		}
		return recipes

	} else {
		fmt.Println("Error in FindRecipeByIngredients: ", err)
	}
	return nil
}
