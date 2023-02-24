package main

import (
	"ginchat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:1234@tcp(127.0.0.1:3307)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.UserBasic{})

	// Create
	user := &models.UserBasic{}
	user.Name = "yanming"
	db.Create(user)

	// Read
	// fmt.Println(db.First(user, 1))

	// Update - update product's price to 200
	db.Model(user).Update("Password", "1234")
	// Update - update multiple fields
	// db.Model(user).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(user).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	// db.Delete(user, 1)
}
