package migration

import (
	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/hasher"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	users := []*entity.User{
		{Email: "alice@example.com", Username: "alice", Password: hashPassword("Alice12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "bob@example.com", Username: "bob", Password: hashPassword("Bob12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "charlie@example.com", Username: "charlie", Password: hashPassword("Charlie12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "doni@example.com", Username: "doni", Password: hashPassword("Doni12345"), Role: entity.RoleUser, IsVerified: true},
		{Email: "david@example.com", Username: "david", Password: hashPassword("David12345"), Role: entity.RoleUser, IsVerified: true},
	}
	productCategory := []*entity.ProductCategory{
		{Name: "Makanan"},
		{Name: "Minuman"},
		{Name: "Pakaian"},
	}
	products := []*entity.Product{
		{Name: "Aqua", ProductCategoryId: 2, ProductCode: "MIN-001", Stock: 6, Price: decimal.NewFromInt(3000)},
		{Name: "Le Mineral", ProductCategoryId: 2, ProductCode: "MIN-002", Stock: 5, Price: decimal.NewFromInt(4000)},
		{Name: "Oreo", ProductCategoryId: 1, ProductCode: "MIN-001", Stock: 4, Price: decimal.NewFromInt(6000)},
		{Name: "Baju Biru", ProductCategoryId: 3, ProductCode: "PAK-001", Stock: 2, Price: decimal.NewFromInt(7000)},
	}
	db.Create(users)
	db.Create(productCategory)
	db.Create(products)
}

func hashPassword(text string) string {
	h := hasher.NewHasher()
	hashedText, err := h.Hash(text)
	if err != nil {
		return ""
	}
	return string(hashedText)
}
