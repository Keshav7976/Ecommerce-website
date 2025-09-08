//config/db.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/keshav7976/ecommerce/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Example: root:password@tcp(127.0.0.1:3306)/ecommerce?parseTime=true
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/ecommerce?parseTime=true"
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto-migrate in correct order
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("Migration failed for User:", err)
	}

	err = DB.AutoMigrate(&models.Category{})
	if err != nil {
		fmt.Println("Migration failed for Category:", err)
	}

	err = DB.AutoMigrate(&models.Item{})
	if err != nil {
		fmt.Println("Migration failed for Item:", err)
	}

	err = DB.AutoMigrate(&models.CartItem{})
	if err != nil {
		fmt.Println("Migration failed for CartItem:", err)
	}

	// Predefined categories (seed if not present)
	categories := []models.Category{
		{Name: "Electronics"},
		{Name: "Books"},
		{Name: "Fashion"},
		{Name: "Home & Kitchen"},
		{Name: "Sports"},
		{Name: "Toys"},
		{Name: "Groceries"},
	}

	for _, c := range categories {
		var existing models.Category
		err := DB.Where("name = ?", c.Name).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				DB.Create(&c)
			} else {
				fmt.Println("Error checking category:", err)
			}
		}
	}

	fmt.Println("Database connected and categories seeded successfully!")
}