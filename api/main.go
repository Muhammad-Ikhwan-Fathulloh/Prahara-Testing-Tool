package main

import (
	"context"
	"net/http"
	"os"
	"prahara-api/auth"
	"prahara-api/models"
	"prahara-api/services"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Ensure data directory exists
	os.MkdirAll("./data", 0755)

	db, err := gorm.Open(sqlite.Open("./data/prahara.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate
	db.AutoMigrate(&models.User{}, &models.TestScript{}, &models.TestRun{}, &models.TestingURL{})

	// Initialize K6 Service
	k6Service := &services.K6Service{
		DB:     db,
		Bucket: "prahara-metrics",
		Org:    "prahara-org",
		Token:  "prahara-token-1234567890",
	}

	r := gin.Default()
	r.Use(cors.Default())

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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

		hashedPassword, _ := auth.HashPassword(input.Password)
		user := models.User{
			Username: input.Username,
			Password: hashedPassword,
			Role:     input.Role,
		}

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
		now := time.Now()
		user.LastLogin = &now
		db.Save(&user)

		token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}
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
		category := c.Query("category")
		metrics, err := k6Service.GetMetrics(c.Request.Context(), timeRange, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, metrics)
	})

	// Test Runs History
	protected.GET("/runs", func(c *gin.Context) {
		var runs []models.TestRun
		db.Order("created_at desc").Limit(20).Find(&runs)
		c.JSON(http.StatusOK, runs)
	})

	protected.GET("/runs/:id/metrics", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		metrics, err := k6Service.GetRunMetrics(context.Background(), uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, metrics)
	})

	protected.GET("/runs/:id/logs", func(c *gin.Context) {
		id := c.Param("id")
		var run models.TestRun
		if err := db.First(&run, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Run not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"logs": run.Logs})
	})

	// User Management (Admin Only)
	protected.GET("/users", func(c *gin.Context) {
		if c.GetString("role") != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		var users []models.User
		db.Find(&users)
		c.JSON(http.StatusOK, users)
	})

	protected.POST("/users", func(c *gin.Context) {
		if c.GetString("role") != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user := models.User{
			Username: input.Username,
			Password: string(hashedPassword),
			Role:     input.Role,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	})

	protected.DELETE("/users/:id", func(c *gin.Context) {
		if c.GetString("role") != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		id := c.Param("id")
		db.Delete(&models.User{}, id)
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	})

	// URL Management Routes
	protected.GET("/urls", func(c *gin.Context) {
		category := c.Query("category")
		var urls []models.TestingURL
		if category != "" && category != "ALL" {
			db.Where("category = ?", category).Find(&urls)
		} else {
			db.Find(&urls)
		}
		c.JSON(http.StatusOK, urls)
	})

	protected.POST("/urls", func(c *gin.Context) {
		var input struct {
			Name     string `json:"name"`
			URL      string `json:"url"`
			Category string `json:"category"`
			Script   string `json:"script"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		url := models.TestingURL{
			Name:     input.Name,
			URL:      input.URL,
			Category: input.Category,
			Script:   input.Script,
		}
		db.Create(&url)
		c.JSON(http.StatusCreated, url)
	})

	protected.PUT("/urls/:id", func(c *gin.Context) {
		id := c.Param("id")
		var input struct {
			Name     string `json:"name"`
			URL      string `json:"url"`
			Category string `json:"category"`
			Script   string `json:"script"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var url models.TestingURL
		if err := db.First(&url, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		db.Model(&url).Updates(models.TestingURL{
			Name:     input.Name,
			URL:      input.URL,
			Category: input.Category,
			Script:   input.Script,
		})
		c.JSON(http.StatusOK, url)
	})

	protected.DELETE("/urls/:id", func(c *gin.Context) {
		id := c.Param("id")
		db.Delete(&models.TestingURL{}, id)
		c.JSON(http.StatusOK, gin.H{"message": "URL deleted"})
	})

	protected.POST("/run-dynamic", func(c *gin.Context) {
		var input struct {
			URL           string `json:"url"`
			Method        string `json:"method"`
			VUs           int    `json:"vus"`
			Duration      string `json:"duration"`
			Category      string `json:"category"`
			ScriptID      uint   `json:"script_id"`
			ScriptContent string `json:"script_content"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		run, err := k6Service.RunDynamicTest(db, input.URL, input.Method, input.VUs, input.Duration, input.Category, input.ScriptID, input.ScriptContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, run)
	})

	r.Run(":8080")
}
