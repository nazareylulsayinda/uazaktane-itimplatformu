package main

import (
	"log"

	"golearn/config"
	"golearn/database"
	_ "golearn/docs" // swagger docs will be here after swag init
	"golearn/handlers"
	"golearn/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GoLearn LMS API
// @version 1.0
// @description GoLearn E-Learning Platform Backend API with Gin, GORM, SQLite.
// @host localhost:8090
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()
	database.Connect()

	r := gin.Default()

	// CORS Middleware (Simple version)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// IP Based Rate Limiting Middleware
	r.Use(middleware.RateLimit())

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected API Routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Course routes
		api.GET("/courses", handlers.GetCourses)
		api.GET("/courses/:id", handlers.GetCourse)

		api.POST("/courses", middleware.TeacherOnly(), handlers.CreateCourse)
		api.PUT("/courses/:id", middleware.TeacherOnly(), handlers.UpdateCourse)
		api.DELETE("/courses/:id", middleware.TeacherOnly(), handlers.DeleteCourse)

		// Lesson routes
		api.GET("/courses/:id/lessons", handlers.GetLessons)
		api.POST("/courses/:id/lessons", middleware.TeacherOnly(), handlers.CreateLesson)

		// Quiz routes
		api.GET("/lessons/:id/quiz", handlers.GetQuiz)
		api.POST("/lessons/:id/quiz", middleware.TeacherOnly(), handlers.CreateQuiz)
		api.POST("/quiz/:id/submit", handlers.SubmitQuiz)

		// Progress routes
		api.POST("/lessons/:id/complete", handlers.CompleteLesson)
		api.GET("/my/progress", handlers.GetMyProgress)
	}

	// Websocket for Classroom
	// Not adding AuthMiddleware to ws group here directly because web browsers natively don't send auth headers for wss easily, 
	// typically handled via query param ?token=xyz or custom headers setup in handler.
	ws := r.Group("/ws")
	ws.GET("/classroom/:courseId", handlers.WebSocketHandler)

	log.Printf("Server is starting at port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
