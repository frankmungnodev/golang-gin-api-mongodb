package routes

import (
	"franky/go-api-gin/controllers"
	"franky/go-api-gin/middlewares"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (ar *AuthRouter) DefineV1Routes(v1Router *gin.RouterGroup) {
	authcontroller := &controllers.AuthController{}

	routerProtector := &middlewares.ProtectRoute{}

	authRouter := v1Router.Group("/auth")
	authRouter.POST("/register", authcontroller.Register)
	authRouter.POST("/login", authcontroller.Login)
	authRouter.POST("/logout", authcontroller.Logout)
	authRouter.GET("/me", routerProtector.Protect, authcontroller.GetMe)
}
