package handler

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	University string `json:"university"`
	Faculty    string `json:"faculty"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Id         int    `json:"id"`
}

type Cache struct {
	Users map[int]User
	Db    *sql.DB
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping": "pong"})
}

func Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello %v", c.Param("name"))
}

func (cc *Cache) Register(c *gin.Context) {
	var u User

	defer func() {
		if recover() != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error while parsing",
			})
		}
	}()

	if c.BindJSON(&u) == nil {
		if cc.checkLogin(u.Email) {

			token, err := generateToken(cc.write(u))
			if err != nil {
				log.Panic(err)
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "User registered",
				"token":   token,
			})
		} else {
			c.JSON(http.StatusConflict, gin.H{
				"message": "User already exists",
			})
		}
	} else {
		log.Panic("error while parsing")
	}
}

func (cc *Cache) checkLogin(login string) bool {
	for _, v := range cc.Users {
		if v.Email == login {
			return false
		}
	}
	return true
}

func generateToken(id int) (string, error) {
	payload := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return t, nil
}

func (cc *Cache) write(u User) int {
	u.Id = int(time.Now().Unix())
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.Password = string(hash)
	cc.Users[u.Id] = u

	sqlStatement := `INSERT INTO users (name, surname, university, faculty, email, password, id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = cc.Db.Exec(sqlStatement, u.Name, u.Surname, u.University, u.Faculty, u.Email, u.Password, u.Id)
	if err != nil {
		panic(err)
	}

	cc.Users[u.Id] = u

	return u.Id
}

func (cc *Cache) Login(c *gin.Context) {
	var u User

	defer func() {
		if recover() != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error while parsing",
			})
		}
	}()

	if c.BindJSON(&u) == nil {
		if u.Id == cc.Users[u.Id].Id {

			token, err := generateToken(u.Id)
			if err != nil {
				log.Panic(err)
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "User logged in",
				"token":   token,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
		}
	} else {
		panic("error while parsing")
	}
}
