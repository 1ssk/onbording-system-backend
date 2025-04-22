package controllers

import (
	"net/http"
	"os"
	"regexp"
	"time"
	"log"

	"github.com/1ssk/go-jwt/initializers"
	"github.com/1ssk/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

func Signup(c *gin.Context) {
	var body models.User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if !emailRegex.MatchString(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	if len(body.Password) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 4 characters long",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash), Role: "user"}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user, email may already exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user_id": user.ID,
	})
}

func Login(c *gin.Context) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("Invalid input: %v", err)
        c.JSON(400, gin.H{"error": "Invalid input"})
        return
    }

    var user models.User
    if err := initializers.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        log.Printf("User not found: %v", err)
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        log.Printf("Invalid password: %v", err)
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub":  user.ID,
        "role": user.Role,
        "exp":  time.Now().Add(time.Hour * 24).Unix(),
    })

    secret := os.Getenv("SECRET")
    if secret == "" {
        log.Printf("SECRET is not set")
        c.JSON(500, gin.H{"error": "Server configuration error"})
        return
    }

    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        log.Printf("Failed to create token: %v", err)
        c.JSON(500, gin.H{"error": "Failed to create token"})
        return
    }

    log.Printf("Generated token: %s", tokenString[:10]+"...")
    c.SetCookie("Authorization", tokenString, 3600*24, "/", "localhost", false, true)
    log.Printf("Set cookie: Authorization=%s", tokenString[:10]+"...")

    c.JSON(200, gin.H{
        "message": "Login successful",
        "user_id": user.ID,
        "token":   tokenString, 
    })
}

func UpdateAdmin(c *gin.Context) {
	id := c.Param("id")

	var roleAdmin struct {
		Role string
	}

	c.Bind(&roleAdmin)

	var updateAdmin models.User
	initializers.DB.First(&updateAdmin, id)

	initializers.DB.Model(&updateAdmin).Updates(models.User{
		Role: roleAdmin.Role,
	})

	c.JSON(200, gin.H{
		"Role": "update role",
	})
}

func Validate(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User validated",
		"user":    user,
	})
}
