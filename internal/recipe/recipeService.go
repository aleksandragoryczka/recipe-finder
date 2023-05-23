package recipe

import (
	"github.com/aleksandragoryczka/recipeFinder/internal/api"
	"github.com/aleksandragoryczka/recipeFinder/internal/database"
)

type Service struct {
	db *database.Database
}

func NewService() (*Service, error) {
	db, _ := database.NewDatabase()
	return &Service{db: db}, nil
}

func (s *Service) FindRecipeByIngredients(ingredients []string, numberOfRecipes int) ([]api.Recipe, error) {
	recipes, err := s.db.GetRecipeByIngredientsList(ingredients, numberOfRecipes)
	if err == nil && len(recipes) == numberOfRecipes {
		return recipes, nil
	} else if err == nil {
		hc := api.NewHttpClinet()
		recipes, err := hc.GetRecipes(ingredients, numberOfRecipes-len(recipes))
		if err != nil {
			return nil, err
		}
		//k := database.StringRecipe(recipes)
		return recipes, nil

	}
	return nil, nil

	//fmt.Println("hot:", recipes)

}
