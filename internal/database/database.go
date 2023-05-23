package database

import (
	"database/sql"
	"fmt"
	"github.com/aleksandragoryczka/recipeFinder/internal/api"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
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
            carbs REAL);`)
	if err != nil {
		fmt.Println("Error creating recipes table: ", err)
		return nil, nil
	}
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS used_ingredients (
            id SERIAL PRIMARY KEY,
            key_f INT,
            used_ingredient VARCHAR(255),
            FOREIGN KEY (key_f) REFERENCES recipes (id) ON DELETE CASCADE
        );`)
	if err != nil {
		fmt.Println("Error creating used_ingredients table: ", err)
		return nil, nil
	}
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS missing_ingredients (
            id SERIAL PRIMARY KEY,
            key_f INT,
            missing_ingredient VARCHAR(255),
            FOREIGN KEY (key_f) REFERENCES recipes (id) ON DELETE CASCADE
        );`)
	if err != nil {
		fmt.Println("Error creating missing_ingredients table: ", err)
		return nil, nil
	}

	return &Database{db: db}, nil
}

func (db *Database) GetRecipeByIngredientsList(ingredients []string, numberOfRecipes int) ([]api.Recipe, error) {
	q := `SELECT
    	r.id_recipe,
        r.title,
        ARRAY_AGG(DISTINCT ui.used_ingredient) AS used_ingredients,
        ARRAY_AGG(DISTINCT mi.missing_ingredient) AS missing_ingredients,
        r.calories,
        r.proteins,
        r.carbs
      FROM
        recipes r
      JOIN
        used_ingredients ui ON r.id = ui.key_f
      JOIN
        missing_ingredients mi ON r.id = mi.key_f
      WHERE
        ui.used_ingredient = ANY($1)
      GROUP BY
        r.id
      LIMIT
        $2;`

	rows, err := db.db.Query(q, pq.Array(ingredients), numberOfRecipes)
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
		fmt.Println(recipe.Id)
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

func (db *Database) InsertTransaction(recipes []api.Recipe) {
	tx, err := db.db.Begin()
	if err != nil {
		fmt.Println("Error beginning transaction: ", err)
	}

	for _, recipe := range recipes {
		err := db.InsertRecipe(recipe)
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

func (db *Database) InsertRecipe(recipe api.Recipe) error {

	var id int
	err := db.db.QueryRow(`INSERT INTO recipes (id_recipe, title, calories, proteins, carbs) 
							VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		recipe.Id, recipe.Title, recipe.Calories, recipe.Proteins, recipe.Carbs).Scan(&id)
	if err != nil {
		fmt.Println("Error inserting row into recipes: ", err)
		return nil
	}

	for _, usedIngredient := range recipe.UsedIngredients {
		_, err = db.db.Exec("INSERT INTO used_ingredients (key_f, used_ingredient) VALUES ($1, $2)",
			id, usedIngredient)
		if err != nil {
			fmt.Println("Error inserting row into used_ingredients: ", err)
			return nil
		}
	}

	for _, missingIngredient := range recipe.MissedIngredients {
		_, err = db.db.Exec("INSERT INTO missing_ingredients (key_f, missing_ingredient) VALUES ($1, $2)",
			id, missingIngredient)
		if err != nil {
			fmt.Println("Error inserting row into missing_ingredients: ", err)
			return nil
		}
	}
	return nil
}

func (db *Database) CloseDatabaseConnection() error {
	return db.db.Close()
}
