package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Recipe struct {
	Id               int     `json:"id"`
	Title            string  `json:"title"`
	Steps            string  `json:"steps"`
	GramsCoffee      float64 `json:"gramsCoffee"`
	GramsWater       float64 `json:"gramsWater"`
	WaterTempCelsius float64 `json:"waterTempCelsius"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type CreateRecipeRequest struct {
	Title            string  `json:"title"`
	Steps            string  `json:"steps"`
	GramsCoffee      float64 `json:"gramsCoffee"`
	GramsWater       float64 `json:"gramsWater"`
	WaterTempCelsius float64 `json:"waterTempCelsius"`
}

type UpdateRecipeRequest struct {
	Title            *string  `json:"title,omitempty"`
	Steps            *string  `json:"steps,omitempty"`
	GramsCoffee      *float64 `json:"gramsCoffee,omitempty"`
	GramsWater       *float64 `json:"gramsWater,omitempty"`
	WaterTempCelsius *float64 `json:"waterTempCelsius,omitempty"`
}

var db *sql.DB

func deleteRecipe(id int) error {
	_, err := db.Exec("DELETE FROM recipes WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}

func handleDeleteRecipe(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	if err := deleteRecipe(id); err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func buildUpdateRecipeQuery(recipe *UpdateRecipeRequest, id int) (string, []interface{}) {
	query := `UPDATE recipes SET `

	queryUpdates := make([]string, 0, 6)
	queryArgs := make([]interface{}, 0, 6)

	if recipe.Title != nil {
		queryUpdates = append(queryUpdates, `title = ?`)
		queryArgs = append(queryArgs, recipe.Title)
	}

	if recipe.Steps != nil {
		queryUpdates = append(queryUpdates, `steps = ?`)
		queryArgs = append(queryArgs, recipe.Steps)
	}

	if recipe.GramsCoffee != nil {
		queryUpdates = append(queryUpdates, `gramsCoffee = ?`)
		queryArgs = append(queryArgs, recipe.GramsCoffee)
	}

	if recipe.GramsWater != nil {
		queryUpdates = append(queryUpdates, `gramsWater = ?`)
		queryArgs = append(queryArgs, recipe.GramsWater)
	}

	if recipe.WaterTempCelsius != nil {
		queryUpdates = append(queryUpdates, `waterTempCelsius = ?`)
		queryArgs = append(queryArgs, recipe.WaterTempCelsius)
	}

	query += strings.Join(queryUpdates, ",") + ` WHERE id = ?`
	queryArgs = append(queryArgs, id)

	return query, queryArgs
}

func updateRecipe(recipe *UpdateRecipeRequest, id int) error {
	// TODO: Check that id exists in db
	query, queryArgs := buildUpdateRecipeQuery(recipe, id)

	_, err := db.Exec(query, queryArgs...)

	if err != nil {
		return err
	}

	return nil
}

func handleUpdateRecipe(c *gin.Context) {
	req := new(UpdateRecipeRequest)
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	if err := updateRecipe(req, id); err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func createRecipe(recipe *CreateRecipeRequest) error {
	query := `INSERT INTO recipes 
	(title, steps, gramsCoffee, gramsWater, waterTempCelsius)
	VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(
		query,
		recipe.Title,
		recipe.Steps,
		recipe.GramsCoffee,
		recipe.GramsWater,
		recipe.WaterTempCelsius,
	)

	if err != nil {
		return err
	}

	return nil
}

func handleCreateRecipe(c *gin.Context) {
	req := new(CreateRecipeRequest)
	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	if err := createRecipe(req); err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func handleGetRecipesById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	var recipe Recipe
	row := db.QueryRow("SELECT * FROM recipes WHERE id = ?", id)

	if err := row.Scan(&recipe.Id, &recipe.Title, &recipe.Steps, &recipe.GramsCoffee, &recipe.GramsWater, recipe.WaterTempCelsius); err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, recipe)
}

func handleGetRecipes(c *gin.Context) {
	var recipes []Recipe
	rows, err := db.Query("SELECT * FROM recipes")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, ResponseError{Error: "Failed to find recipes!"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var recipe Recipe
		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Steps, &recipe.GramsCoffee, &recipe.GramsWater, recipe.WaterTempCelsius); err != nil {
			c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
			return
		}
		recipes = append(recipes, recipe)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, &ResponseError{Error: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, recipes)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConfig := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST"),
		DBName: os.Getenv("DB_NAME"),
	}

	var err error
	db, err = sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		log.Fatalf("Could not connect to the database!")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected!")

	router := gin.Default()
	router.GET("/coffee/recipes", handleGetRecipes)
	router.GET("/coffee/recipes/:id", handleGetRecipesById)
	router.PUT("/coffee/recipes/:id", handleUpdateRecipe)
	router.DELETE("/coffee/recipes/:id", handleDeleteRecipe)
	router.POST("/coffee/recipes", handleCreateRecipe)
	router.Run(":8080")
}
