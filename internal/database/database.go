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
		fmt.Println("laczenie")
		return nil, nil
	}

	_, err = db.Query("CREATE TABLE  IF NOT EXISTS recipes (" +
		"id INT PRIMARY KEY, " +
		"title VARCHAR(255), " +
		"calories REAL, " +
		"proteins REAL, " +
		"carbs REAL);")
	if err != nil {
		//fmt.Println("tworzenie query")
		return nil, nil
	}
	_, err = db.Query("CREATE TABLE IF NOT EXISTS used_ingredients(" +
		"id INT PRIMARY KEY," +
		"recipe_id INT NOT NULL," +
		"used_ingredient VARCHAR(255) NOT NULL," +
		"FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE" +
		");")
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	_, err = db.Query("CREATE TABLE IF NOT EXISTS missing_ingredients(" +
		"id INT PRIMARY KEY," +
		"recipe_id INT NOT NULL," +
		"missing_ingredient VARCHAR(255) NOT NULL," +
		"FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE" +
		");")
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	//defer db.Close()
	return &Database{db: db}, nil
}

var (
	id int
)

func (db *Database) GetRecipeByIngredientsList(ingredients []string, numberOfRecipes int) ([]api.Recipe, error) {
	q := `SELECT
        r.title,
        ARRAY_AGG(DISTINCT ui.used_ingredient) AS used_ingredients,
        ARRAY_AGG(DISTINCT mi.missing_ingredient) AS missing_ingredients,
        r.calories,
        r.proteins,
        r.carbs
      FROM
        recipes r
      INNER JOIN
        used_ingredients ui ON r.id = ui.recipe_id
      LEFT JOIN
        missing_ingredients mi ON r.id = mi.recipe_id
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
	defer rows.Close()

	return recipes, nil
}
