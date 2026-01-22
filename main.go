package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file - only in development
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ No .env file found, using environment variables")
		}
	}

	// Create Fiber app with production config
	app := fiber.New(fiber.Config{
		AppName:               "Baseet Company Registration",
		ReadTimeout:           40 * time.Second,
		WriteTimeout:          40 * time.Second,
		IdleTimeout:           120 * time.Second,
		DisableStartupMessage: os.Getenv("APP_ENV") == "production",
		EnablePrintRoutes:     os.Getenv("APP_ENV") != "production",
		CaseSensitive:         false,
		StrictRouting:         false,
		ProxyHeader:           fiber.HeaderXForwardedFor,
		ErrorHandler:          customErrorHandler,
	})

	// ========== MIDDLEWARE SETUP ==========

	// Security middleware (order is important!)
	app.Use(cors.New())
	app.Use(recover.New())

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:               100,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error":   "Too many requests",
				"message": "Please try again later",
			})
		},
	}))

	// Performance middleware
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Cache static assets
	app.Use(etag.New())

	// Logger (only in development)
	if os.Getenv("APP_ENV") != "production" {
		app.Use(logger.New(logger.Config{
			Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
			TimeFormat: "2006-01-02 15:04:05",
		}))
	}

	// ========== DATABASE ==========
	database.ConnectDB()

	// ========== ROUTES ==========
	routes.SetupAllRoutes(app)

	// ========== STATIC FILES ==========

	// 1. Serve uploaded files (important: before the client build!)
	uploadsPath := "./uploads"
	if _, err := os.Stat(uploadsPath); os.IsNotExist(err) {
		log.Printf("⚠️  Uploads directory not found, creating: %s", uploadsPath)
		if err := os.MkdirAll(uploadsPath, 0755); err != nil {
			log.Printf("❌ Failed to create uploads directory: %v", err)
		}
	}

	// Serve uploads with shorter cache (files might change)
	app.Static("/", uploadsPath, fiber.Static{
		Compress:      true,
		ByteRange:     true,
		CacheDuration: 1 * time.Hour,
		MaxAge:        3600,
		Browse:        false, // Important: disable directory listing for uploads
	})

	// 2. Serve React app build
	clientBuildPath := filepath.Join("..", "companies-registeration", "dist")

	if _, err := os.Stat(clientBuildPath); os.IsNotExist(err) {
		log.Printf("⚠️  Client build directory not found at: %s", clientBuildPath)
	} else {
		// Static files with cache headers
		app.Static("/", clientBuildPath, fiber.Static{
			Compress:      true,
			ByteRange:     true,
			CacheDuration: 24 * time.Hour,
			MaxAge:        86400,
			Index:         "index.html", // Specify index file
		})

		// Cache static assets aggressively (CSS, JS, images)
		app.Use("/static/*", cache.New(cache.Config{
			Expiration:   30 * 24 * time.Hour, // 30 days
			CacheControl: true,
		}))

		// Cache other assets like images, fonts
		app.Use(cache.New(cache.Config{
			Expiration:   7 * 24 * time.Hour, // 7 days
			CacheControl: true,
			KeyGenerator: func(c *fiber.Ctx) string {
				// Only cache asset files
				path := c.Path()
				if filepath.Ext(path) == ".js" ||
					filepath.Ext(path) == ".css" ||
					filepath.Ext(path) == ".png" ||
					filepath.Ext(path) == ".jpg" ||
					filepath.Ext(path) == ".jpeg" ||
					filepath.Ext(path) == ".gif" ||
					filepath.Ext(path) == ".svg" ||
					filepath.Ext(path) == ".woff" ||
					filepath.Ext(path) == ".woff2" ||
					filepath.Ext(path) == ".ttf" ||
					filepath.Ext(path) == ".eot" ||
					filepath.Ext(path) == ".ico" {
					return path
				}
				return ""
			},
		}))

		// SPA fallback - must be LAST after all other routes
		app.Get("*", func(c *fiber.Ctx) error {
			return c.SendFile(filepath.Join(clientBuildPath, "index.html"))
		})
	}

	// ========== ERROR HANDLERS ==========
	// 404 handler (will only be reached if no other route matches)
	app.Use(notFoundHandler)

	// ========== START SERVER ==========
	port := os.Getenv("PORT")
	if port == "" {
		port = "42069"
	}

	host := "0.0.0.0"
	if os.Getenv("APP_ENV") == "development" {
		host = "localhost"
	}

	log.Printf("🚀 Server starting on %s:%s in %s mode", host, port, os.Getenv("APP_ENV"))

	if err := app.Listen(host + ":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

// ========== HELPER FUNCTIONS ==========

func customErrorHandler(c *fiber.Ctx, err error) error {
	// Log the error
	log.Printf("Error: %v - Path: %s - Method: %s", err, c.Path(), c.Method())

	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	if os.Getenv("APP_ENV") == "production" && code == 500 {
		message = "Something went wrong"
	}

	if c.Get("Accept") == "application/json" {
		return c.Status(code).JSON(fiber.Map{
			"error":   true,
			"message": message,
			"code":    code,
		})
	}

	return c.Status(code).SendString(message)
}

func notFoundHandler(c *fiber.Ctx) error {
	// Check if it's an API request
	if c.Get("Accept") == "application/json" || c.Path()[:4] == "/api" {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Endpoint not found",
		})
	}

	// For SPA, try to serve index.html if it exists
	clientBuildPath := filepath.Join("..", "client", "build")
	indexPath := filepath.Join(clientBuildPath, "index.html")

	if _, err := os.Stat(indexPath); err == nil {
		return c.SendFile(indexPath)
	}

	return c.Status(404).SendString("Page not found")
}
