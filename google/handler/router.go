package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yaumianwar/go-oauth-google/google"
)

func Blueprint(routerGroup *gin.RouterGroup, db *gorm.DB) {
	svc := google.NewService(db)
	routerGroup.GET("/", IndexHandler)
	routerGroup.GET("/login", LoginHandler)
	routerGroup.GET("/GoogleCallback",AuthHandler(&svc))
}
