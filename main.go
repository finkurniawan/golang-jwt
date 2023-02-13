package main

import (
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"golang-jwt/src/api/v1/app"
	"golang-jwt/src/api/v1/controller"
	"golang-jwt/src/api/v1/exception"
	"golang-jwt/src/api/v1/helper"
	"golang-jwt/src/api/v1/middleware"
	"golang-jwt/src/api/v1/repository"
	"golang-jwt/src/api/v1/service"
	"net/http"
)

func main() {
	//load environment
	envErr := godotenv.Load(".env")
	helper.PanicError(envErr)

	//connection to database
	db := app.Database()

	//validator
	validate := *validator.New()

	//	repository
	userRepository := repository.NewUserRepositoryImpl()

	//service
	userService := service.NewUserServiceImpl(
		userRepository,
		db,
		validate,
	)

	//controller
	userController := controller.NewUserControllerImpl(userService)

	//router
	router := httprouter.New()

	//[USER]
	router.POST("/api/v1/user", userController.Create)
	router.POST("/api/v1/auth", userController.Auth)
	router.GET("/api/v1/refresh-token", userController.CreateWithRefreshToken)
	router.PUT("/api/v1/user/:user_id", userController.Update)
	router.DELETE("/api/v1/user/:user_id", userController.Delete)
	router.GET("/api/v1/user/:user_id", userController.FindById)
	router.GET("/api/v1/user", userController.FindAll)

	//error handler
	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()

	helper.PanicError(err)

}
