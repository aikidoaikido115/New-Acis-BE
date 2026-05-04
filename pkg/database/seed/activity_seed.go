package seed

import (
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedActivities creates default activity data in the database.
func SeedActivities(db *gorm.DB) {
	log.Println("Seeding activities...")

	var seedStaff entities.Staff
	if err := db.First(&seedStaff).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("⏭️  Skip activities seeding: no staff found (staff_id is required)")
			return
		}

		log.Printf("❌ Failed to prepare activities seeding: %v", err)
		return
	}

	activities := []entities.Activity{
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "เกมฝึกความจำ ต่อคำ",
			ActivityType: "กิจกรรมกระตุ้นสมอง",
			Description:  strPtr("เล่นต่อคำง่ายๆ เป็นกลุ่มเล็ก เพื่อกระตุ้นการคิดและการจำ"),
			Location:     strPtr("โถงกิจกรรมชั้น 1"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "ปริศนาอักษรไขว้",
			ActivityType: "กิจกรรมกระตุ้นสมอง",
			Description:  strPtr("ฝึกสมองด้วยโจทย์คำศัพท์แบบง่าย เหมาะกับผู้สูงอายุ"),
			Location:     strPtr("มุมอ่านหนังสือ"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "เล่านิทานและแชร์ประสบการณ์ชีวิต",
			ActivityType: "กิจกรรมกระตุ้นสมอง",
			Description:  strPtr("สลับกันเล่าเรื่องในอดีต เพื่อฝึกการสื่อสารและความจำระยะยาว"),
			Location:     strPtr("โถงกลาง"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "ระบายสีภาพธรรมชาติ",
			ActivityType: "กิจกรรมสร้างสรรค์",
			Description:  strPtr("ใช้สีไม้และสีเทียน ฝึกสมาธิและกล้ามเนื้อมัดเล็ก"),
			Location:     strPtr("ห้องศิลปะ"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "พับกระดาษและงานฝีมือ",
			ActivityType: "กิจกรรมสร้างสรรค์",
			Description:  strPtr("ทำดอกไม้กระดาษและงานประดิษฐ์ง่ายๆ"),
			Location:     strPtr("ห้องกิจกรรม 2"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "ปลูกต้นไม้กระถางเล็ก",
			ActivityType: "กิจกรรมสร้างสรรค์",
			Description:  strPtr("ปลูกพืชสมุนไพรและไม้ประดับ ดูแลง่าย"),
			Location:     strPtr("สวนหย่อมด้านหลัง"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "กายบริหารเบาๆ ยืดเส้น",
			ActivityType: "กิจกรรมทางกาย",
			Description:  strPtr("ยืดเหยียด 20 นาที โดยมีเจ้าหน้าที่ดูแลใกล้ชิด"),
			Location:     strPtr("ลานเอนกประสงค์"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "โยคะผู้สูงอายุ",
			ActivityType: "กิจกรรมทางกาย",
			Description:  strPtr("ท่าง่าย เน้นการหายใจและการทรงตัว"),
			Location:     strPtr("โถงชั้น 2"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "เดินเล่นในสวน",
			ActivityType: "กิจกรรมทางกาย",
			Description:  strPtr("เดินช้าๆ รอบสวนพร้อมวัดชีพจรเบื้องต้นก่อนและหลัง"),
			Location:     strPtr("ทางเดินสวนในบ้านพัก"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "พูดคุยกลุ่มเช้า",
			ActivityType: "กิจกรรมสังคม",
			Description:  strPtr("วงสนทนาหัวข้อชีวิตประจำวัน เพื่อสร้างปฏิสัมพันธ์"),
			Location:     strPtr("โถงรับแขก"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "เกมกระดานโดมิโน",
			ActivityType: "กิจกรรมสังคม",
			Description:  strPtr("เล่นโดมิโนและเกมกระดานเพื่อฝึกคิดเชิงกลยุทธ์"),
			Location:     strPtr("ห้องสันทนาการ"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "สวดมนต์และนั่งสมาธิ",
			ActivityType: "กิจกรรมด้านจิตใจ/ศาสนา",
			Description:  strPtr("ช่วงเช้า 30 นาที เพื่อความสงบและผ่อนคลาย"),
			Location:     strPtr("ห้องพระ"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "ฟังธรรม",
			ActivityType: "กิจกรรมด้านจิตใจ/ศาสนา",
			Description:  strPtr("เปิดเสียงบรรยายธรรมะระดับเข้าใจง่าย"),
			Location:     strPtr("ห้องประชุมเล็ก"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "ร้องเพลงคาราโอเกะ",
			ActivityType: "กิจกรรมบันเทิง",
			Description:  strPtr("เพลงลูกกรุงและเพลงเก่ายอดนิยม"),
			Location:     strPtr("โถงอเนกประสงค์"),
		},
		{
			ID:           uuid.New().String(),
			StaffID:      seedStaff.ID,
			ActivityName: "ดูหนังคลาสสิกและฟังเพลงเก่า",
			ActivityType: "กิจกรรมบันเทิง",
			Description:  strPtr("ฉายภาพยนตร์เก่าและเปิดเพลงย้อนยุคหลังมื้อเย็น"),
			Location:     strPtr("ห้องดูทีวีรวม"),
		},
	}

	for _, activity := range activities {
		var existingActivity entities.Activity
		result := db.Where("activity_name = ? AND activity_type = ?", activity.ActivityName, activity.ActivityType).First(&existingActivity)

		if result.Error != nil {
			if err := db.Create(&activity).Error; err != nil {
				log.Printf("❌ Failed to seed activity '%s': %v", activity.ActivityName, err)
			} else {
				log.Printf("✅ Seeded activity: %s (ID: %s)", activity.ActivityName, activity.ID)
			}
		} else {
			log.Printf("⏭️  Activity already exists: %s (ID: %s)", existingActivity.ActivityName, existingActivity.ID)
		}
	}

	log.Println("Activities seeding completed!")
}

func strPtr(value string) *string {
	return &value
}
