package migration

import (
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	u := &entity.User{}
	p := &entity.Product{}
	pc := &entity.ProductCategory{}
	o := &entity.Order{}
	oi := &entity.OrderItem{}

	_ = db.Migrator().DropTable(u, p, pc, o, oi)
	_ = db.AutoMigrate(u, p, pc, o, oi)
}
