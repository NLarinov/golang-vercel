package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"golang-vercel/app/routes"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

var (
	app *gin.Engine
)

func init() {
	app = gin.New()

	database := openDB()

	routes.Main(app, database)
}

func openDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Established a successful connection!")

	return db
}

// Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
