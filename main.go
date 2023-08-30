package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowCredentials: true,
		// 在使用 JWT 的时候，因为我们使用了 Authorization 的头部，所以要加上
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 为了 JWT
		ExposeHeaders: []string{"X-Jwt-Token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://") {
				// 开发环境
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	})
}

func main() {
	db := initDB()
	server := initServer()
	user := initUser(db)
	user.RegisterRoutes(server)
	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initUser(db *gorm.DB) *web.UserHandler {
	userDao := dao.NewUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepository)
	return web.NewUserHandler(userService)
}

func initServer() *gin.Engine {
	server := gin.Default()
	server.Use(corsHandler())

	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))

	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())

	return server
}
