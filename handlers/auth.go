package handlers

import (
	"time"
	"waweb-v2/config"
	"waweb-v2/database"
	"waweb-v2/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Index(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token != "" {
		return c.Redirect("/dashboard")
	}
	return c.Redirect("/login")
}

func LoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}

func Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, exists := database.GetUserByEmail(email)
	if !exists {
		return c.Status(401).SendString(`
			<div class="alert alert-error" role="alert">
				Invalid email or password
			</div>
		`)
	}

	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.Status(401).SendString(`
			<div class="alert alert-error" role="alert">
				Invalid email or password
			</div>
		`)
	}

	// Create JWT token
	cfg := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return c.Status(500).SendString("Error generating token")
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})

	// Return HTMX redirect
	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(200)
}

func RegisterPage(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title": "Register",
	})
}

func Register(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	name := c.FormValue("name")

	// Check if user exists
	if _, exists := database.GetUserByEmail(email); exists {
		return c.Status(400).SendString(`
			<div class="alert alert-error" role="alert">
				Email already registered
			</div>
		`)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).SendString("Error creating account")
	}

	// Create user
	user := &models.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := database.CreateUser(user); err != nil {
		return c.Status(500).SendString("Error creating account")
	}

	// Create JWT token
	cfg := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return c.Status(500).SendString("Error generating token")
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})

	// Return HTMX redirect
	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(200)
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("token")
	return c.Redirect("/login")
}
