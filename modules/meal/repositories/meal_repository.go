package repositories

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
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
	GetMealPlansToday() ([]*entities.MealPlan, error)
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

func (r *GormMealRepository) GetMealPlansToday() ([]*entities.MealPlan, error) {
	var mealPlans []*entities.MealPlan
	if err := r.db.Preload("Menu").
		Where("created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
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
