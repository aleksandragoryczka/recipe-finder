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

func (s *Service) FindRecipeByIngredients(inputIngredients string, numberOfRecipes int) []api.Recipe {
	recipes, err := s.db.GetRecipeByIngredientsList(inputIngredients, numberOfRecipes)
	if err == nil && len(recipes) == numberOfRecipes {
		err := s.db.CloseDatabaseConnection()
		if err != nil {
			fmt.Println("Error closing DB connection: ", err)
		}
		return recipes
	} else if err == nil {
		hc := api.NewHttpClient()
		apiRecipes := hc.GetRecipes(inputIngredients, numberOfRecipes, len(recipes))
		if err != nil {
			fmt.Println("Error gettingRecipes from API: ", err)
			return nil
		}
		s.db.InsertTransaction(apiRecipes, inputIngredients)
		err := s.db.CloseDatabaseConnection()
		if err != nil {
			fmt.Println("Error closing DB connection: ", err)
		}
		return append(recipes, apiRecipes...)

	} else {
		fmt.Println("Error in FindRecipeByIngredients: ", err)
	}
	return nil
}
