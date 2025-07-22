package seeders

import (
	"rowjak/website-inspektorat/app/models"

	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *UserSeeder) Signature() string {
	return "UserSeeder"
}

// Run executes the seeder logic.
func (s *UserSeeder) Run() error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Inspektorat2018."), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		NamaLengkap: "Inspektorat",
		Email:       "inspektoratpekalongankab@gmail.com",
		Password:    string(hashedPassword),
		Role:        "admin",
	}

	// Insert user data
	return facades.Orm().Query().Create(&user)
}

// go run . artisan migrate:fresh --seed --seeder=UserSeeder
