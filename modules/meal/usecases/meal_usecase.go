package usecases

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	audit_constants "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/constants"
	audit_repo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	meal_repo "github.com/aikidoaikido115/New-Acis-BE/modules/meal/repositories"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/google/uuid"
)

type MealUsecase interface {
	// Menu operations
	CreateMenu(menu *entities.Menu, userID string) (*entities.Menu, error)
	GetMenuByID(id string, userID string) (*entities.Menu, error)
	GetAllMenus(userID string) ([]*entities.Menu, error)
	UpdateMenu(menuID string, menu *entities.Menu, userID string) (*entities.Menu, error)

	// MealPlan operations
	CreateMealPlan(mealPlan *entities.MealPlan, userID string) (*entities.MealPlan, error)
	GetMealPlanByID(id string, userID string) (*entities.MealPlan, error)
	GetAllMealPlans(userID string) ([]*entities.MealPlan, error)
	UpdateMealPlan(mealPlanID string, mealPlan *entities.MealPlan, userID string) (*entities.MealPlan, error)
}

type MealUseCaseImpl struct {
	repo         meal_repo.MealRepository
	auditlogrepo audit_repo.AuditLogRepository
	userrepo     user_repo.UserRepository
}

func NewMealUseCase(repo meal_repo.MealRepository, auditlogrepo audit_repo.AuditLogRepository, userrepo user_repo.UserRepository) MealUsecase {
	return &MealUseCaseImpl{repo: repo, auditlogrepo: auditlogrepo, userrepo: userrepo}
}

func (uc *MealUseCaseImpl) CreateMenu(menu *entities.Menu, userID string) (*entities.Menu, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	if menu == nil {
		return nil, errors.New("menu payload is required")
	}

	menu.MenuName = strings.TrimSpace(menu.MenuName)
	menu.Description = strings.TrimSpace(menu.Description)

	if menu.MenuName == "" {
		return nil, errors.New("menu_name is required")
	}

	if menu.Description == "" {
		return nil, errors.New("description is required")
	}

	menu.ID = uuid.New().String()
	createdMenu, err := uc.repo.CreateMenu(menu)
	if err != nil {
		return nil, err
	}

	newMenuData, _ := json.Marshal(map[string]interface{}{
		"menu_name":   createdMenu.MenuName,
		"description": createdMenu.Description,
	})
	uc.createAuditLog(userID, audit_constants.AuditActionInsert, "menus", createdMenu.ID, "", string(newMenuData))

	return createdMenu, nil
}

func (uc *MealUseCaseImpl) GetMenuByID(id string, userID string) (*entities.Menu, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("menu id is required")
	}

	return uc.repo.GetMenuByID(id)
}

func (uc *MealUseCaseImpl) GetAllMenus(userID string) ([]*entities.Menu, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	return uc.repo.GetAllMenus()
}

func (uc *MealUseCaseImpl) UpdateMenu(menuID string, menu *entities.Menu, userID string) (*entities.Menu, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	if menu == nil {
		return nil, errors.New("menu payload is required")
	}

	menuID = strings.TrimSpace(menuID)
	if menuID == "" {
		return nil, errors.New("menu id is required")
	}

	existingMenu, err := uc.repo.GetMenuByID(menuID)
	if err != nil {
		return nil, errors.New("menu not found")
	}

	oldMenuData, _ := json.Marshal(map[string]interface{}{
		"menu_name":   existingMenu.MenuName,
		"description": existingMenu.Description,
	})

	if name := strings.TrimSpace(menu.MenuName); name != "" {
		existingMenu.MenuName = name
	}
	if desc := strings.TrimSpace(menu.Description); desc != "" {
		existingMenu.Description = desc
	}

	updatedMenu, err := uc.repo.UpdateMenu(existingMenu)
	if err != nil {
		return nil, err
	}

	newMenuData, _ := json.Marshal(map[string]interface{}{
		"menu_name":   updatedMenu.MenuName,
		"description": updatedMenu.Description,
	})
	uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "menus", updatedMenu.ID, string(oldMenuData), string(newMenuData))

	return updatedMenu, nil
}

func (uc *MealUseCaseImpl) CreateMealPlan(mealPlan *entities.MealPlan, userID string) (*entities.MealPlan, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	if mealPlan == nil {
		return nil, errors.New("meal plan payload is required")
	}

	if err := uc.validateMealPlan(mealPlan); err != nil {
		return nil, err
	}

	mealPlan.ID = uuid.New().String()
	createdMealPlan, err := uc.repo.CreateMealPlan(mealPlan)
	if err != nil {
		return nil, err
	}

	newMealPlanData, _ := json.Marshal(map[string]interface{}{
		"menu_id":        createdMealPlan.MenuID,
		"backup_menu_id": createdMealPlan.BackUpMenuID,
		"main_amount":    createdMealPlan.MainAmount,
		"backup_amount":  createdMealPlan.BackUpAmount,
		"meal_type":      createdMealPlan.MealType,
	})
	uc.createAuditLog(userID, audit_constants.AuditActionInsert, "meal_plans", createdMealPlan.ID, "", string(newMealPlanData))

	return createdMealPlan, nil
}

func (uc *MealUseCaseImpl) GetMealPlanByID(id string, userID string) (*entities.MealPlan, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("meal plan id is required")
	}

	return uc.repo.GetMealPlanByID(id)
}

func (uc *MealUseCaseImpl) GetAllMealPlans(userID string) ([]*entities.MealPlan, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	return uc.repo.GetAllMealPlans()
}

func (uc *MealUseCaseImpl) UpdateMealPlan(mealPlanID string, mealPlan *entities.MealPlan, userID string) (*entities.MealPlan, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	if mealPlan == nil {
		return nil, errors.New("meal plan payload is required")
	}

	mealPlanID = strings.TrimSpace(mealPlanID)
	if mealPlanID == "" {
		return nil, errors.New("meal plan id is required")
	}

	existingMealPlan, err := uc.repo.GetMealPlanByID(mealPlanID)
	if err != nil {
		return nil, errors.New("meal plan not found")
	}

	oldMealPlanData, _ := json.Marshal(map[string]interface{}{
		"menu_id":        existingMealPlan.MenuID,
		"backup_menu_id": existingMealPlan.BackUpMenuID,
		"main_amount":    existingMealPlan.MainAmount,
		"backup_amount":  existingMealPlan.BackUpAmount,
		"meal_type":      existingMealPlan.MealType,
	})

	if menuID := strings.TrimSpace(mealPlan.MenuID); menuID != "" {
		if _, err := uc.repo.GetMenuByID(menuID); err != nil {
			return nil, errors.New("menu does not exist")
		}
		existingMealPlan.MenuID = menuID
	}

	if mealPlan.MainAmount > 0 {
		existingMealPlan.MainAmount = mealPlan.MainAmount
	}

	if mealPlan.BackUpAmount != nil {
		if *mealPlan.BackUpAmount < 0 {
			return nil, errors.New("backup_amount cannot be negative")
		}
		existingMealPlan.BackUpAmount = mealPlan.BackUpAmount
	}

	if mealType := strings.ToLower(strings.TrimSpace(mealPlan.MealType)); mealType != "" {
		if mealType != "breakfast" && mealType != "lunch" && mealType != "dinner" {
			return nil, errors.New("meal_type must be one of breakfast, lunch, dinner")
		}
		existingMealPlan.MealType = mealType
	}

	if mealPlan.BackUpMenuID != nil {
		backupID := strings.TrimSpace(*mealPlan.BackUpMenuID)
		if backupID == "" {
			existingMealPlan.BackUpMenuID = nil
		} else {
			if backupID == existingMealPlan.MenuID {
				return nil, errors.New("backup_menu_id must be different from menu_id")
			}
			if _, err := uc.repo.GetMenuByID(backupID); err != nil {
				return nil, errors.New("backup menu does not exist")
			}
			existingMealPlan.BackUpMenuID = &backupID
		}
	}

	if existingMealPlan.MainAmount <= 0 {
		return nil, errors.New("main_amount must be greater than 0")
	}

	if existingMealPlan.BackUpAmount != nil && *existingMealPlan.BackUpAmount < 0 {
		return nil, errors.New("backup_amount cannot be negative")
	}

	if existingMealPlan.MealType != "breakfast" && existingMealPlan.MealType != "lunch" && existingMealPlan.MealType != "dinner" {
		return nil, errors.New("meal_type must be one of breakfast, lunch, dinner")
	}

	if existingMealPlan.BackUpMenuID != nil && strings.TrimSpace(*existingMealPlan.BackUpMenuID) == "" {
		existingMealPlan.BackUpMenuID = nil
	}

	updatedMealPlan, err := uc.repo.UpdateMealPlan(existingMealPlan)
	if err != nil {
		return nil, err
	}

	newMealPlanData, _ := json.Marshal(map[string]interface{}{
		"menu_id":        updatedMealPlan.MenuID,
		"backup_menu_id": updatedMealPlan.BackUpMenuID,
		"main_amount":    updatedMealPlan.MainAmount,
		"backup_amount":  updatedMealPlan.BackUpAmount,
		"meal_type":      updatedMealPlan.MealType,
	})
	uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "meal_plans", updatedMealPlan.ID, string(oldMealPlanData), string(newMealPlanData))

	return updatedMealPlan, nil
}

func (uc *MealUseCaseImpl) validateMealPlan(mealPlan *entities.MealPlan) error {
	mealPlan.MenuID = strings.TrimSpace(mealPlan.MenuID)
	mealPlan.MealType = strings.ToLower(strings.TrimSpace(mealPlan.MealType))

	if mealPlan.MenuID == "" {
		return errors.New("menu_id is required")
	}

	if _, err := uc.repo.GetMenuByID(mealPlan.MenuID); err != nil {
		return errors.New("menu does not exist")
	}

	if mealPlan.MainAmount <= 0 {
		return errors.New("main_amount must be greater than 0")
	}

	if mealPlan.BackUpAmount != nil && *mealPlan.BackUpAmount < 0 {
		return errors.New("backup_amount cannot be negative")
	}

	if mealPlan.MealType != "breakfast" && mealPlan.MealType != "lunch" && mealPlan.MealType != "dinner" {
		return errors.New("meal_type must be one of breakfast, lunch, dinner")
	}

	if mealPlan.BackUpMenuID != nil {
		backupID := strings.TrimSpace(*mealPlan.BackUpMenuID)
		if backupID == "" {
			return errors.New("backup_menu_id cannot be empty")
		}
		if backupID == mealPlan.MenuID {
			return errors.New("backup_menu_id must be different from menu_id")
		}
		if _, err := uc.repo.GetMenuByID(backupID); err != nil {
			return errors.New("backup menu does not exist")
		}
		mealPlan.BackUpMenuID = &backupID
	}

	return nil
}

func (uc *MealUseCaseImpl) ensureKitchenStaff(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	role, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}

	if role.Name != user_constants.RoleKitchenStaff {
		return errors.New("only users with 'Kitchen Staff' role can manage meals")
	}

	return nil
}

func (uc *MealUseCaseImpl) createAuditLog(userID string, action string, tableName string, recordID string, oldValue string, newValue string) {
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: tableName,
		RecordID:  recordID,
		UserID:    userID,
		Action:    action,
		OldValue:  oldValue,
		NewValue:  newValue,
	}

	if _, err := uc.auditlogrepo.CreateAuditLog(auditLog); err != nil {
		log.Printf("[ERROR] Failed to create audit log for %s %s: %v", tableName, recordID, err)
	}
}
