package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"scau-daily/internal/config"
	"scau-daily/internal/database"
	"scau-daily/internal/handler"
	"scau-daily/internal/middleware"
	"scau-daily/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize database connections
	database.Init()
	defer database.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "SCAU Daily API",
		ErrorHandler: errorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(fiberlog.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Initialize services
	authSvc := &service.AuthService{}
	scheduleSvc := &service.ScheduleService{}
	todoSvc := &service.TodoService{}
	todaySvc := &service.TodayService{}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc)
	scheduleHandler := handler.NewScheduleHandler(scheduleSvc)
	todoHandler := handler.NewTodoHandler(todoSvc)
	todayHandler := handler.NewTodayHandler(todaySvc)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/wx-login", authHandler.WxLogin)
	auth.Post("/refresh", authHandler.Refresh)

	// Auth routes (protected)
	authProtected := api.Group("/auth", middleware.JWTAuth())
	authProtected.Get("/me", authHandler.GetMe)
	authProtected.Post("/bind-student", authHandler.BindStudent)

	// Schedule routes (protected)
	schedule := api.Group("/schedule", middleware.JWTAuth())
	schedule.Post("/sync", scheduleHandler.Sync)
	schedule.Get("/today", scheduleHandler.Today)
	schedule.Get("/week", scheduleHandler.Week)
	schedule.Get("/courses", scheduleHandler.Courses)
	schedule.Get("/free-slots", scheduleHandler.FreeSlots)

	// Todo routes (protected)
	todos := api.Group("/todos", middleware.JWTAuth())
	todos.Get("/", todoHandler.List)
	todos.Post("/", todoHandler.Create)
	todos.Patch("/:id", todoHandler.Update)
	todos.Delete("/:id", todoHandler.Delete)

	// Today routes (protected)
	today := api.Group("/today", middleware.JWTAuth())
	today.Get("/overview", todayHandler.Overview)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		_ = app.Shutdown()
	}()

	// Start server
	addr := ":" + config.AppConfig.Port
	log.Printf("Starting server on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error":   "internal_error",
		"message": err.Error(),
	})
}
