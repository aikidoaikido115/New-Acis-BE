package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	audit_constants "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/constants"
	audit_repo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	emr_repo "github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	meal_constants "github.com/aikidoaikido115/New-Acis-BE/modules/meal/constants"
	"github.com/aikidoaikido115/New-Acis-BE/modules/meal/models"
	meal_repo "github.com/aikidoaikido115/New-Acis-BE/modules/meal/repositories"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	aiinfra "github.com/aikidoaikido115/New-Acis-BE/pkg/ai"
	"github.com/google/uuid"
)

type AllergyAIClient interface {
	CheckAllergy(ctx context.Context, payload aiinfra.CheckAllergyRequest) ([]byte, error)
}

type MealUsecase interface {
	// Menu operations
	CreateMenu(menu *entities.Menu, userID string) (*entities.Menu, error)
	GetMenuByID(id string, userID string) (*entities.Menu, error)
	GetAllMenus(userID string) ([]*entities.Menu, error)
	UpdateMenu(menuID string, req models.UpdateMenuRequest, userID string) (*entities.Menu, error)

	// MealPlan operations
	CreateMealPlan(mealPlan *entities.MealPlan, userID string, humanInTheLoop bool) (*entities.MealPlan, *models.AllergyCheckSummary, error)
	GetMealPlanByID(id string, userID string) (*entities.MealPlan, error)
	GetAllMealPlans(userID string) ([]*entities.MealPlan, error)
	GetMealPlansToday(userID string) ([]*entities.MealPlan, error)
	UpdateMealPlan(mealPlanID string, mealPlan *entities.MealPlan, userID string) (*entities.MealPlan, error)
}

type MealUseCaseImpl struct {
	repo         meal_repo.MealRepository
	emrrepo      emr_repo.EmrRepository
	auditlogrepo audit_repo.AuditLogRepository
	userrepo     user_repo.UserRepository
	allergyAI    AllergyAIClient
}

func NewMealUseCase(
	repo meal_repo.MealRepository,
	emrrepo emr_repo.EmrRepository,
	auditlogrepo audit_repo.AuditLogRepository,
	userrepo user_repo.UserRepository,
	allergyAI AllergyAIClient,
) MealUsecase {
	return &MealUseCaseImpl{
		repo:         repo,
		emrrepo:      emrrepo,
		auditlogrepo: auditlogrepo,
		userrepo:     userrepo,
		allergyAI:    allergyAI,
	}
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
	if err := validateMenuDescriptionFormat(menu.Description); err != nil {
		return nil, err
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

func (uc *MealUseCaseImpl) UpdateMenu(menuID string, req models.UpdateMenuRequest, userID string) (*entities.Menu, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	if req.MenuName == nil && req.Description == nil {
		return nil, errors.New("at least one field must be provided")
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

	if req.MenuName != nil {
		name := strings.TrimSpace(*req.MenuName)
		if name == "" {
			return nil, errors.New("menu_name cannot be empty")
		}
		existingMenu.MenuName = name
	}
	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		if desc == "" {
			return nil, errors.New("description cannot be empty")
		}
		if err := validateMenuDescriptionFormat(desc); err != nil {
			return nil, err
		}
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

func (uc *MealUseCaseImpl) CreateMealPlan(mealPlan *entities.MealPlan, userID string, humanInTheLoop bool) (*entities.MealPlan, *models.AllergyCheckSummary, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, nil, err
	}

	if mealPlan == nil {
		return nil, nil, errors.New("meal plan payload is required")
	}

	if err := uc.validateMealPlan(mealPlan); err != nil {
		return nil, nil, err
	}

	menu, err := uc.repo.GetMenuByID(mealPlan.MenuID)
	if err != nil {
		return nil, nil, errors.New("menu does not exist")
	}

	if mealPlan.BackUpMenuID == nil || strings.TrimSpace(*mealPlan.BackUpMenuID) == "" {
		return nil, nil, errors.New("backup_menu_id is required")
	}

	backupMenuID := strings.TrimSpace(*mealPlan.BackUpMenuID)
	backupMenu, err := uc.repo.GetMenuByID(backupMenuID)
	if err != nil {
		return nil, nil, errors.New("backup menu does not exist")
	}
	mealPlan.BackUpMenuID = &backupMenuID

	allergyStats, err := uc.emrrepo.GetResidentAllergyStatsDashboard()
	if err != nil {
		return nil, nil, errors.New("failed to get resident allergy stats: " + err.Error())
	}

	allergyDetails := make([]aiinfra.AllergyDetail, 0, len(allergyStats.AllergyDetails))
	for _, detail := range allergyStats.AllergyDetails {
		allergyDetails = append(allergyDetails, aiinfra.AllergyDetail{
			AllergyID:   detail.AllergyID,
			AllergyName: detail.AllergyName,
			Count:       int(detail.ResidentCount),
		})
	}

	var warningSummary *models.AllergyCheckSummary

	// If human_in_the_loop is true, skip AI check and mark as reviewed by human.
	// Otherwise, check allergy for both main and backup menus in separate rounds.
	if !humanInTheLoop {
		mainResult, mainPassed, err := uc.checkMenuAllergy(menu, allergyDetails)
		if err != nil {
			return nil, nil, err
		}

		backupResult, backupPassed, err := uc.checkMenuAllergy(backupMenu, allergyDetails)
		if err != nil {
			return nil, nil, err
		}

		summary := &models.AllergyCheckSummary{
			MainMenuPassed:   mainPassed,
			BackupMenuPassed: backupPassed,
			MainMenuResult:   mainResult,
			BackupMenuResult: backupResult,
		}

		switch {
		case mainPassed && backupPassed:
			// Both passed. Save meal plan.
		case !mainPassed && backupPassed:
			// Main failed but backup passed. Save with warning.
			warningSummary = summary
		case mainPassed && !backupPassed:
			allergyErr := models.NewAllergyCheckError(
				"backup menu is not safe for allergy group: "+backupResult.Reason,
				backupResult.Status,
				summary,
			)
			return nil, nil, allergyErr
		default:
			status := meal_constants.AllergyCheckStatusAllergyWarn
			if mainResult.Status == meal_constants.AllergyCheckStatusManualReview || backupResult.Status == meal_constants.AllergyCheckStatusManualReview {
				status = meal_constants.AllergyCheckStatusManualReview
			}

			allergyErr := models.NewAllergyCheckError(
				"both main and backup menus are not safe for allergy group",
				status,
				summary,
			)
			return nil, nil, allergyErr
		}
	}

	deletedMealPlans, err := uc.repo.DeleteMealPlansToday()
	if err != nil {
		return nil, nil, errors.New("failed to replace today's meal plans: " + err.Error())
	}

	for _, deletedMealPlan := range deletedMealPlans {
		oldMealPlanData, _ := json.Marshal(map[string]interface{}{
			"menu_id":        deletedMealPlan.MenuID,
			"backup_menu_id": deletedMealPlan.BackUpMenuID,
			"main_amount":    deletedMealPlan.MainAmount,
			"backup_amount":  deletedMealPlan.BackUpAmount,
			"meal_type":      deletedMealPlan.MealType,
		})
		uc.createAuditLog(userID, audit_constants.AuditActionDelete, "meal_plans", deletedMealPlan.ID, string(oldMealPlanData), "")
	}

	mealPlan.ID = uuid.New().String()
	createdMealPlan, err := uc.repo.CreateMealPlan(mealPlan)
	if err != nil {
		return nil, nil, err
	}

	newMealPlanData, _ := json.Marshal(map[string]interface{}{
		"menu_id":        createdMealPlan.MenuID,
		"backup_menu_id": createdMealPlan.BackUpMenuID,
		"main_amount":    createdMealPlan.MainAmount,
		"backup_amount":  createdMealPlan.BackUpAmount,
		"meal_type":      createdMealPlan.MealType,
	})
	uc.createAuditLog(userID, audit_constants.AuditActionInsert, "meal_plans", createdMealPlan.ID, "", string(newMealPlanData))

	return createdMealPlan, warningSummary, nil
}

func (uc *MealUseCaseImpl) checkMenuAllergy(menu *entities.Menu, allergyDetails []aiinfra.AllergyDetail) (*aiinfra.CheckAllergyResponse, bool, error) {
	aiResponse, err := uc.allergyAI.CheckAllergy(context.Background(), aiinfra.CheckAllergyRequest{
		MenuData: aiinfra.MenuData{
			MenuName:        menu.MenuName,
			MenuDescription: menu.Description,
		},
		AllergyDetails: allergyDetails,
	})
	if err != nil {
		return nil, false, errors.New("failed to check allergy by ai: " + err.Error())
	}

	var aiCheckResult aiinfra.CheckAllergyResponse
	if err := json.Unmarshal(aiResponse, &aiCheckResult); err != nil {
		return nil, false, errors.New("failed to parse ai response: " + err.Error())
	}

	switch aiCheckResult.Status {
	case meal_constants.AllergyCheckStatusSafe:
		return &aiCheckResult, true, nil
	case meal_constants.AllergyCheckStatusManualReview, meal_constants.AllergyCheckStatusAllergyWarn:
		return &aiCheckResult, false, nil
	default:
		return nil, false, errors.New("unknown ai status: " + aiCheckResult.Status)
	}
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

func (uc *MealUseCaseImpl) GetMealPlansToday(userID string) ([]*entities.MealPlan, error) {
	if err := uc.ensureKitchenStaff(userID); err != nil {
		return nil, err
	}

	return uc.repo.GetMealPlansToday()
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

func validateMenuDescriptionFormat(description string) error {
	parts := strings.Split(description, ",")
	if len(parts) < 2 {
		return errors.New("description must be comma-separated, e.g. \"ingredient1, ingredient2\"")
	}

	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return errors.New("description contains an empty ingredient")
		}
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

	if role.Name != user_constants.RoleKitchenStaff && role.Name != user_constants.RoleSuperUser && role.Name != user_constants.RoleAdmin {
		return errors.New("only users with 'Kitchen Staff', 'Super User', or 'Admin' role can manage meals")
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
