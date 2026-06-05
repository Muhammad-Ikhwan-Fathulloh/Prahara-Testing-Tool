package main

import (
	"context"
	"net/http"
	"os"
	"prahara-api/auth"
	"prahara-api/models"
	"prahara-api/services"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("prahara.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.TestScript{}, &models.TestRun{}, &models.TestingURL{})

	influxURL := os.Getenv("INFLUX_URL")
	if influxURL == "" {
		influxURL = "http://influxdb:8086"
	}
	k6Service := services.NewK6Service(db, influxURL, "prahara-token", "prahara", "prahara")

	r := gin.Default()
	r.Use(cors.Default())

	// Auth Routes
	r.POST("/api/auth/register", func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hashed, _ := auth.HashPassword(input.Password)
		user := models.User{Username: input.Username, Password: hashed, Role: input.Role}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username exists"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username})
	})

	r.POST("/api/auth/login", func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if !auth.CheckPasswordHash(input.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		token, _ := auth.GenerateToken(user.ID, user.Username, user.Role)
		c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "username": user.Username, "role": user.Role}})
	})

	// Test Routes
	protected := r.Group("/api")
	// Simplified auth for demo, usually you'd add a middleware here
	protected.GET("/scripts", func(c *gin.Context) {
		var scripts []models.TestScript
		db.Find(&scripts)
		c.JSON(http.StatusOK, scripts)
	})

	protected.POST("/scripts", func(c *gin.Context) {
		var script models.TestScript
		c.ShouldBindJSON(&script)
		db.Create(&script)
		c.JSON(http.StatusCreated, script)
	})

	protected.POST("/scripts/:id/run", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		run, err := k6Service.RunTest(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, run)
	})

	protected.GET("/metrics", func(c *gin.Context) {
		timeRange := c.DefaultQuery("range", "-1h")
		metrics, err := k6Service.GetMetrics(context.Background(), timeRange)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, metrics)
	})

	// URL Management Routes
	protected.GET("/urls", func(c *gin.Context) {
		category := c.Query("category")
		var urls []models.TestingURL
		query := db
		if category != "" {
			query = query.Where("category = ?", category)
		}
		query.Find(&urls)
		c.JSON(http.StatusOK, urls)
	})

	protected.POST("/urls", func(c *gin.Context) {
		var testingURL models.TestingURL
		if err := c.ShouldBindJSON(&testingURL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&testingURL)
		c.JSON(http.StatusCreated, testingURL)
	})

	protected.POST("/run-dynamic", func(c *gin.Context) {
		var input struct {
			URL      string `json:"url"`
			Method   string `json:"method"`
			VUs      int    `json:"vus"`
			Duration string `json:"duration"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		run, err := k6Service.RunDynamicTest(input.URL, input.Method, input.VUs, input.Duration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, run)
	})

	// Serve static files from Vue dist
	// In production, this will be in /app/web/dist
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	r.NoRoute(func(c *gin.Context) {
		// If the request is for /api, don't serve index.html
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}
		c.File("./web/dist/index.html")
	})

	r.Run(":3000")
}
