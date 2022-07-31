package main

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/golang-jwt/jwt/v4"
	_ "github.com/google/uuid"
	"github.com/inamandev/aspire-backend-test-2022/handlers/auth"
	"github.com/inamandev/aspire-backend-test-2022/handlers/loans"
	"github.com/inamandev/aspire-backend-test-2022/middlewares"
	_ "github.com/mattn/go-sqlite3"
)

var (
	app *fiber.App
)

func init() {
	log.SetFlags(log.Ltime | log.Llongfile)
	if ok := os.Setenv("JWT_SIGNING_KEY", "thisisjustasimeple"); ok != nil {
		log.Println(ok)
	}
}

/* func init() {
if ok := godotenv.Load(); ok != nil {
	log.Println(ok)
	panic(ok)
}
} */

// import _
var testLength = 3700000

// var testLength = 100000
var wg sync.WaitGroup

func main() {
	os.Remove("test.db")
	log.SetFlags(log.Ltime | log.Lshortfile)
	defer handlePanic()
	db, err := sql.Open("sqlite3", "test.db?_journal_mode=WAL")
	// db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Println("error connecting to database", err)
	}
	defer db.Close()
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	// db.SetConnMaxIdleTime(500 * time.Second)
	// db.SetConnMaxLifetime(500 * time.Second)
	// db.SetMaxIdleConns()
	if ok := db.Ping(); ok != nil {
		log.Println("unable to ping the database", ok)
	}
	start := time.Now()
	createTable(db)
	stmt := prepareQurey(db)
	defer stmt.Close()
	totalInserted := 0
	for i := 0; i < testLength; i++ {
		wg.Add(1)
		// log.Println("pringting test number", i)
		go insertQuery(stmt, i, &wg, &totalInserted)
	}
	log.Println("-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_WE FINISHED THE LOOP-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_")
	wg.Wait()
	log.Println("time for loop test", time.Since(start))
}

func prepareQurey(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare(`INSERT INTO person (id, name, age) VALUES(?,?,?)`)
	if err != nil {
		log.Println("unable to prepare insert query", err)
		panic(err)
	}
	return stmt
}

func insertQuery(stmt *sql.Stmt, id int, wg *sync.WaitGroup, total *int) {
	res, err := stmt.Exec(id, "Naman", 22)
	if err != nil {
		log.Println("unable to execute the query", err)
		panic(err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		log.Println("unable to get result", err)
		panic(err)
	}
	*total += 1
	log.Println(affectedRows, "person(s) inserted for id", id, "total inserted", *total)
	wg.Done()
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE person (
		id int,
		name text,
		age int
	);`)
	if err != nil {
		log.Println("unable to create table", err)
		return err
	}
	return nil
}

func startHttpApp() {
	app = fiber.New()
	app.Use(etag.New())
	app.Use(logger.New())
	app.Use(middlewares.Recover)
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/user/auth", auth.Login)
	v1.Post("/loans", middlewares.Auth, loans.Create)
	if ok := app.Listen(":3031"); ok != nil {
		log.Println("Unable to start the application due to below error")
		log.Println(ok)
		panic(ok)
	}
}

func handlePanic() {
	if ok := recover(); ok != nil {
		log.Println("recovering from panic", ok)
		app.Listen(":3031")
	}
}
