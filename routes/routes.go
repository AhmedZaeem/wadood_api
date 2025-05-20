package routes

import (
	"net/http"
	"wadood/auth"
	"github.com/gin-gonic/gin" 
)

func Setup() *gin.Engine {
    r := gin.Default()
    r.POST("/register", auth.Register)
    r.POST("/login", auth.Login)
    r.POST("/logout", auth.Logout)
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Welcome to Wadood API")
    })
    return r
}