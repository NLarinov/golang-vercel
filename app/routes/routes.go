package routes

import (
	"database/sql"
	"golang-vercel/app/handler"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Main(app *gin.Engine, db *sql.DB) {
	cache := newCache(db)

	app.NoRoute(ErrRouter)

	app.GET("/ping", handler.Ping)

	route := app.Group("/api")
	{
		route.POST("/register", cache.Register)
		route.POST("/login", cache.Login)
		route.GET("/hello/:name", handler.Hello)
	}
}

func ErrRouter(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": "this page could not be found",
	})
}

func newCache(db *sql.DB) *handler.Cache {
	users := make(map[int]handler.User)
	rows, err := db.Query(`SELECT * from users`)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var u handler.User
		err = rows.Scan(&u.Name, &u.Surname, &u.University, &u.Faculty, &u.Email, &u.Password, &u.Id)
		if err != nil {
			log.Println(err)
		}
		users[u.Id] = u
	}

	return &handler.Cache{
		Users: users,
		Db:    db,
	}
}
