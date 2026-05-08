package seed

import (
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedMenus creates default menu data in the database
func SeedMenus(db *gorm.DB) {
	log.Println("Seeding menus...")

	menus := []entities.Menu{
		{
			ID:          uuid.New().String(),
			MenuName:    "คาเปรเซ่สลัด",
			Description: "มะเขือเทศสด, ใบโหระพา, มอสซาเรลล่ายี่ห้อ A, น้ำมันมะกอก, บัลซามิก",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "กระเพราะหมูสับ",
			Description: "หมูสับ, ใบกะเพรา, พริก, กระเทียม, ซอสปรุงรส",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ข้ามต้มกุ้ง",
			Description: "ข้าวหอมมะลิ, กุ้งสด, ขิงซอย, ต้นหอม, ขึ้นฉ่าย",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ต้มยำกุ้ง",
			Description: "กุ้ง, เห็ด, ตะไคร้, ใบมะกรูด, น้ำปลา, น้ำมะนาว",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ผัดไทยกุ้งสด",
			Description: "เส้นจันท์, กุ้ง, เต้าหู้, ไข่, ถั่วงอก, ใบกุยช่าย",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ซูชิแซลมอน",
			Description: "ข้าวซูชิ, ปลาแซลมอน, สาหร่าย, วาซาบิ, โชยุ",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "สปาเกตตีคาโบนารา",
			Description: "เส้นสปาเกตตี, เบคอน, ไข่แดง, พาร์เมซานชีส, พริกไทยดำ",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ซีซาร์สลัด",
			Description: "ผักโรมานี, ขนมปังกรอบ, พาร์เมซานชีส, น้ำสลัดซีซาร์",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "แกงเขียวหวานไก่",
			Description: "เนื้อไก่, พริกแกงเขียวหวาน, กะทิ, มะเขือเปราะ, ใบโหระพา, พริกชี้ฟ้า",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ส้มตำไทย",
			Description: "มะละกอสับ, มะเขือเทศ, ถั่วฝักยาว, ถั่วลิสงคั่ว, กุ้งแห้ง, พริก, กระเทียม, น้ำปลา, มะนาว",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "มัสมั่นเนื้อ",
			Description: "เนื้อวัวตุ๋น, พริกแกงมัสมั่น, กะทิ, มันฝรั่ง, หอมใหญ่, ถั่วลิสงคั่ว",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ข้าวผัดปู",
			Description: "ข้าวสวย, เนื้อปู, ไข่ไก่, ต้นหอมซอย, กระเทียม, ซอสปรุงรส",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ต้มข่าไก่",
			Description: "เนื้อไก่, กะทิ, เห็ดฟาง, ข่า, ตะไคร้, ใบมะกรูด, น้ำมะนาว",
		},
		{
			ID:          uuid.New().String(),
			MenuName:    "ยำวุ้นเส้นหมูสับ",
			Description: "วุ้นเส้น, หมูสับ, หอมใหญ่, มะเขือเทศ, ขึ้นฉ่าย, พริก, น้ำยำรสแซ่บ",
		},
	}

	for _, menu := range menus {
		var existingMenu entities.Menu
		result := db.Where("menu_name = ?", menu.MenuName).First(&existingMenu)

		if result.Error != nil {
			if err := db.Create(&menu).Error; err != nil {
				log.Printf("❌ Failed to seed menu '%s': %v", menu.MenuName, err)
			} else {
				log.Printf("✅ Seeded menu: %s (ID: %s)", menu.MenuName, menu.ID)
			}
		} else {
			log.Printf("⏭️  Menu already exists: %s (ID: %s)", existingMenu.MenuName, existingMenu.ID)
		}
	}

	log.Println("Menus seeding completed!")
}
