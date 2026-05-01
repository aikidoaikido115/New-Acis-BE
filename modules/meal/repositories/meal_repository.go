package repositories

import (
	"strings"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/meal/models"
	"gorm.io/gorm"
)

type GormMealRepository struct {
	db *gorm.DB
}

func NewGormMealRepository(db *gorm.DB) *GormMealRepository {
	return &GormMealRepository{
		db: db,
	}
}

type MealRepository interface {
	// Menu operations
	CreateMenu(meal *entities.Menu) (*entities.Menu, error)
	GetMenuByID(id string) (*entities.Menu, error)
	GetAllMenus() ([]*entities.Menu, error)
	UpdateMenu(meal *entities.Menu) (*entities.Menu, error)

	// MealPlan operations
	CreateMealPlan(mealPlan *entities.MealPlan) (*entities.MealPlan, error)
	GetMealPlanByID(id string) (*entities.MealPlan, error)
	GetAllMealPlans() ([]*entities.MealPlan, error)
	GetMealPlansByDate(date string) ([]*entities.MealPlan, error)
	GetMealHistory(date *string, mealType *string, search *string, page int, pageSize int) ([]models.MealHistoryItem, int64, error)
	DeleteMealPlansToday() ([]*entities.MealPlan, error)
	UpdateMealPlan(mealPlan *entities.MealPlan) (*entities.MealPlan, error)
}

func (r *GormMealRepository) CreateMenu(meal *entities.Menu) (*entities.Menu, error) {
	if err := r.db.Create(&meal).Error; err != nil {
		return nil, err
	}

	return r.GetMenuByID(meal.ID)
}

func (r *GormMealRepository) GetMenuByID(id string) (*entities.Menu, error) {
	var meal entities.Menu
	if err := r.db.Where("id = ?", id).First(&meal).Error; err != nil {
		return nil, err
	}

	return &meal, nil
}

func (r *GormMealRepository) GetAllMenus() ([]*entities.Menu, error) {
	var meals []*entities.Menu
	if err := r.db.Find(&meals).Error; err != nil {
		return nil, err
	}

	return meals, nil
}

func (r *GormMealRepository) UpdateMenu(meal *entities.Menu) (*entities.Menu, error) {
	if err := r.db.Save(&meal).Error; err != nil {
		return nil, err
	}

	return r.GetMenuByID(meal.ID)
}

func (r *GormMealRepository) CreateMealPlan(mealPlan *entities.MealPlan) (*entities.MealPlan, error) {
	if err := r.db.Omit("is_allergy").Create(&mealPlan).Error; err != nil {
		return nil, err
	}

	return r.GetMealPlanByID(mealPlan.ID)
}

func (r *GormMealRepository) GetMealPlanByID(id string) (*entities.MealPlan, error) {
	var mealPlan entities.MealPlan
	if err := r.db.Preload("Menu").Where("id = ?", id).First(&mealPlan).Error; err != nil {
		return nil, err
	}

	return &mealPlan, nil
}

func (r *GormMealRepository) GetAllMealPlans() ([]*entities.MealPlan, error) {
	var mealPlans []*entities.MealPlan
	if err := r.db.Preload("Menu").Find(&mealPlans).Error; err != nil {
		return nil, err
	}

	return mealPlans, nil
}

func (r *GormMealRepository) UpdateMealPlan(mealPlan *entities.MealPlan) (*entities.MealPlan, error) {
	if err := r.db.Model(&entities.MealPlan{}).
		Where("id = ?", mealPlan.ID).
		Omit("is_allergy").
		Updates(mealPlan).Error; err != nil {
		return nil, err
	}

	return r.GetMealPlanByID(mealPlan.ID)
}

func (r *GormMealRepository) GetMealPlansByDate(date string) ([]*entities.MealPlan, error) {
	// parse date in Asia/Bangkok location
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.FixedZone("Asia/Bangkok", 7*60*60) // fallback +7
	}

	parsed, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return nil, err
	}

	start := parsed
	end := start.Add(24 * time.Hour)

	var mealPlans []*entities.MealPlan
	if err := r.db.Preload("Menu").
		Where("created_at >= ? AND created_at < ?", start, end).
		Order("created_at DESC").
		Find(&mealPlans).Error; err != nil {
		return nil, err
	}

	return mealPlans, nil
}

func (r *GormMealRepository) DeleteMealPlansToday() ([]*entities.MealPlan, error) {
	var todayMealPlans []*entities.MealPlan

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("Menu").
			Where("created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
			Where("created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
			Order("created_at DESC").
			Find(&todayMealPlans).Error; err != nil {
			return err
		}

		if len(todayMealPlans) == 0 {
			return nil
		}

		if err := tx.Where("created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
			Where("created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
			Delete(&entities.MealPlan{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return todayMealPlans, nil
}

func (r *GormMealRepository) GetMealHistory(date *string, mealType *string, search *string, page int, pageSize int) ([]models.MealHistoryItem, int64, error) {
	buildQuery := func(withPagination bool) *gorm.DB {
		query := r.db.Table("meal_plans mp").
			Joins("LEFT JOIN menus m ON mp.menu_id = m.id").
			Joins("LEFT JOIN menus bm ON mp.back_up_menu_id = bm.id")

		if date != nil && strings.TrimSpace(*date) != "" {
			loc, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				loc = time.FixedZone("Asia/Bangkok", 7*60*60)
			}
			parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(*date), loc)
			if err == nil {
				start := parsed
				end := start.Add(24 * time.Hour)
				query = query.Where("mp.created_at >= ? AND mp.created_at < ?", start, end)
			}
		}

		if mealType != nil && strings.TrimSpace(*mealType) != "" {
			query = query.Where("LOWER(mp.meal_type) = ?", strings.ToLower(strings.TrimSpace(*mealType)))
		}

		if search != nil && strings.TrimSpace(*search) != "" {
			like := "%" + strings.TrimSpace(*search) + "%"
			query = query.Where(
				`m.menu_name ILIKE ? OR bm.menu_name ILIKE ? OR mp.created_by_staff_id ILIKE ? OR mp.staff_name ILIKE ?`,
				like, like, like, like,
			)
		}

		query = query.Order("mp.created_at DESC")

		if withPagination {
			offset := (page - 1) * pageSize
			query = query.Offset(offset).Limit(pageSize)
		}

		return query
	}

	var total int64
	countSubQuery := buildQuery(false).Select("mp.id")
	if err := r.db.Table("(?) AS filtered_meal_history", countSubQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := make([]models.MealHistoryItem, 0)
	err := buildQuery(true).
		Select(`
			TO_CHAR((mp.created_at AT TIME ZONE 'Asia/Bangkok')::date, 'YYYY-MM-DD') AS date,
			mp.meal_type,
			m.menu_name,
			mp.main_amount,
			bm.menu_name AS backup_menu_name,
			mp.back_up_amount AS backup_amount,
			COALESCE(NULLIF(TRIM(mp.staff_name), ''), mp.created_by_staff_id) AS staff_name
		`).
		Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
