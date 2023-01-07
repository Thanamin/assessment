package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
	// _ "github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func createNewExpense(c echo.Context) error {
	e := Expense{}
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	err = row.Scan(&e.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}

func getExpenseByID(c echo.Context) error {
	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	e := Expense{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}

func updateExpenseByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	e := Expense{}
	e.ID = id

	err = c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses statment:" + err.Error()})
	}

	if _, err := stmt.Exec(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(e.Tags)); err != nil {
		log.Fatal("error execute update ", err)
	}

	return c.JSON(http.StatusOK, e)
}

func getAllExpense(c echo.Context) error {

	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	rows, err := stmt.Query()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	eTotal := []Expense{}
	for rows.Next() {
		e := Expense{}
		err = rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "Can't Scan Expense:" + err.Error()})
		}
		eTotal = append(eTotal, e)
	}

	return c.JSON(http.StatusOK, eTotal)
}

var db *sql.DB

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	//Connect Database
	url := "postgres://glvrfamy:3otHFqRv3zLOpqIuSSeq4OS6XhClBm_X@john.db.elephantsql.com/glvrfamy"
	var err error
	// url := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount INT, note TEXT, tags TEXT[]);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	// log.Println("Okay")

	//

	apiPort := ":2565"
	// apiPort := os.Getenv("PORT")

	e := echo.New()

	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == "apidesign" || password == "45678" {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.POST("/expenses", createNewExpense)
	e.GET("/expenses/:id", getExpenseByID)
	e.PUT("/expenses/:id", updateExpenseByID)
	e.GET("/expenses", getAllExpense)

	log.Fatal(e.Start(apiPort))
}
