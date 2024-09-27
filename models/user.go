package models

import (
	"fmt"
	"golang_rest_api/db"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID       uuid.UUID
	Email    string
	Password string
}

// HashPassword hashes the user's password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks if the provided password matches the hashed one
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Save inserts a new user into the database
func (u *User) Save() error {
	u.ID = uuid.New() // Generate a new UUID for the user

	query := `INSERT INTO users (id, email, password) VALUES ($1, $2, $3);`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Fatalf("Error preparing query: %v", err)
		return err
	}

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
		return err
	}

	_, err = stmt.Exec(u.ID, u.Email, hashedPassword) // Hash the password before storing it
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
		return err
	}

	log.Printf("User saved with ID: %v", u.ID)
	return nil
}

var JwtSecretKey = []byte("MIHuAgEAMBAGByqGSM49AgEGBSuBBAAjBIHWMIHTAgEBBEIAuOQ7wchpXzehkzl6ywqmayBl+Am+TtQ5AkMZDOVPq6AkRN4w1YhLfyCCRWmFgPzfvuQqImN7Ry5bfeqMEN2HVguhgYkDgYYABAGINhz1Z/yOtijOrsw/DDJo+hV2PnzgxzJUsDOUOxAUs+azx6T1TnMzPMSI0mzdZzqdeYMIJooM2euM+ZbrJgmgGgC+r2/kSRIzr4sq4X9X4hGWhRITFX+WiVS+OE0VH549rbzD+prixzON/Ta6scJgc5JjRKO09EwO6Zp/8W/dl7GHJQ==")

func (u *User) GenerateJWT() (string, error) {
	claims := jwt.MapClaims{
		"user_id": u.ID.String(),
		"email":   u.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Expiration time of 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecretKey)
}

func (u *User) Login(email, password string) (string, error) {
	var storedUser User

	// Fetch user from the database by email
	query := `SELECT id, email, password FROM users WHERE email = $1;`
	err := db.DB.QueryRow(query, email).Scan(&storedUser.ID, &storedUser.Email, &storedUser.Password)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Compare the provided password with the hashed password
	if !CheckPasswordHash(password, storedUser.Password) {
		return "", fmt.Errorf("invalid password")
	}

	// If password matches, generate a JWT token
	token, err := storedUser.GenerateJWT()
	if err != nil {
		return "", fmt.Errorf("could not generate JWT: %v", err)
	}

	return token, nil
}
