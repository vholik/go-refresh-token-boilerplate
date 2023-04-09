package router

import (
	"net/http"
	"ocenakademik/internal/user"
	"ocenakademik/util"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, jwtSecret string) {
	r = gin.Default()

	authMiddleware := util.NewAuthMiddleware(jwtSecret).TokenAuthMiddleware()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/signup", userHandler.SignUp)
	r.GET("/signin", userHandler.SignIn)
	r.GET("/refresh", userHandler.Refresh)
	r.GET("/hello", authMiddleware, func(c *gin.Context) {
		id := c.Request.Context().Value("userId")
		c.JSON(http.StatusOK, gin.H{"message": id})
	})
}

func Start(addr string) error {
	return r.Run(addr)
}
