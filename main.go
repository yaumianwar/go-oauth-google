package main

import (
	"github.com/gin-gonic/gin"
	google "github.com/yaumianwar/go-oauth-google/google/handler"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/spf13/viper"
	"github.com/jinzhu/gorm"
	"github.com/yaumianwar/go-oauth-google/middleware"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}

func main() {
	db, err := gorm.Open("postgres", viper.GetString("datasource"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)

	router := gin.Default()
	store := sessions.NewCookieStore([]byte(google.RandToken(64)))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("goquestsession", store))
	router.Static("/css", "./static/css")
	router.Static("/img", "./static/img")
	router.LoadHTMLGlob("templates/*")

	google.Blueprint(router.Group(""), db)

	authorized := router.Group("/battle")
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/field", google.FieldHandler)
	}

	router.Run("127.0.0.1:3000")
}