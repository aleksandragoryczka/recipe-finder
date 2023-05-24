package database

import (
	"database/sql"
	"fmt"
	"github.com/aleksandragoryczka/recipeFinder/internal/api"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "user"
	dbname   = "recipeFinderDb"
)

type StringRecipe struct {
	Id                int
	Title             string
	MissedIngredients []string
	UsedIngredients   []string
	Calories          float64
	Proteins          float64
	Carbs             float64
}

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		fmt.Println("Error connecting to db: ", err)
		return nil, nil
	}

	_, err = db.Query(`CREATE TABLE IF NOT EXISTS recipes (
    		id SERIAL PRIMARY KEY,
            id_recipe INT,
            title VARCHAR(255),
            calories REAL,
            proteins REAL,
            carbs REAL,
    		used_ingredients VARCHAR(500),
    		missing_ingredients VARCHAR(500),
    		input_ingredients VARCHAR(255));`)
	if err != nil {
		fmt.Println("Error creating recipes table: ", err)
		return nil, nil
	}

	return &Database{db: db}, nil
}

func (db *Database) GetRecipeByIngredientsList(ingredients string, numberOfRecipes int) ([]api.Recipe, error) {
	q := `SELECT
    	r.id_recipe,
        r.title,
        STRING_TO_ARRAY(r.used_ingredients, ',') AS used_ingredients,
        STRING_TO_ARRAY(r.missing_ingredients, ',') AS missing_ingredients,
        r.calories,
        r.proteins,
        r.carbs
      FROM
        recipes r
      WHERE
          r.input_ingredients = $1
      LIMIT
        $2;`

	rows, err := db.db.Query(q, ingredients, numberOfRecipes)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	recipes := make([]api.Recipe, 0)
	defer rows.Close()

	for rows.Next() {
		var recipe api.Recipe
		err := rows.Scan(
			&recipe.Id,
			&recipe.Title,
			pq.Array(&recipe.UsedIngredients),
			pq.Array(&recipe.MissedIngredients),
			&recipe.Calories,
			&recipe.Proteins,
			&recipe.Carbs,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return recipes, nil
}

func (db *Database) InsertTransaction(recipes []api.Recipe, inputIngredients string) {
	tx, err := db.db.Begin()
	if err != nil {
		fmt.Println("Error beginning transaction: ", err)
	}

	for _, recipe := range recipes {
		err := db.InsertRecipe(recipe, inputIngredients)
		if err != nil {
			fmt.Println("Error inserting single Recipe in db: ", err)
			tx.Rollback()
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error committing transaction: ", err)
		tx.Rollback()
	}
}

func (db *Database) InsertRecipe(recipe api.Recipe, inputIngredients string) error {

	_, err := db.db.Query(`INSERT INTO recipes (id_recipe, title, calories, proteins, carbs, 
                     		used_ingredients, missing_ingredients, input_ingredients) 
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		recipe.Id, recipe.Title, recipe.Calories, recipe.Proteins, recipe.Carbs,
		strings.Join(recipe.UsedIngredients, ","), strings.Join(recipe.MissedIngredients, ","),
		inputIngredients)
	if err != nil {
		fmt.Println("Error inserting row into recipes: ", err)
		return nil
	}
	return nil
}

func (db *Database) CloseDatabaseConnection() error {
	return db.db.Close()
}
