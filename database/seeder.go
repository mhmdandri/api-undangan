package database

import (
	"api-undangan/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassword(plain string)(string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

func seedUsers(db *gorm.DB) {
	admin, err := hashPassword("admin123")
	if err != nil{
		log.Println("gagal hash password admin:", err)
		return
	}
	user, err := hashPassword("user123")
	if err != nil{
		log.Println("gagal hash password user:", err)
		return
	}
	users := []models.User{
        {Name: "Admin", Email: "admin@mohaproject.dev", Password: admin, Role: "admin"},
        {Name: "User", Email: "user@mohaproject.dev", Password: user, Role: "user"},
    }

    for _, u := range users {
        if err := db.Where("email = ?", u.Email).
            Attrs(models.User{
                Name:     u.Name,
                Password: u.Password,
                Role:     u.Role,
            }).
            FirstOrCreate(&u).Error; err != nil {
            log.Printf("gagal membuat user %s: %v\n", u.Email, err)
        }
    }
}