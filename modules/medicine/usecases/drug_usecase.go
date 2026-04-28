package usecases

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	audit_constants "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/constants"
	audit_repository "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/medicine/models"
	medicine_repository "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/repositories"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repository "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/google/uuid"
)

type DrugUsecase interface {
	CreateDrugMaster(req models.CreateDrugMasterRequest, userID string) (*entities.DrugMaster, error)
	GetDrugMasters(userID string) ([]*entities.DrugMaster, error)
	GetDrugMasterByID(id string, userID string) (*entities.DrugMaster, error)
	UpdateDrugMasterByID(id string, req models.UpdateDrugMasterRequest, userID string) (*entities.DrugMaster, error)
	DeleteDrugMasterByID(id string, userID string) error

	CreatePersonalDrug(req models.CreatePersonalDrugRequest, userID string) (*entities.PersonalDrug, error)
	GetPersonalDrugsOverview(req models.PersonalDrugOverviewQueryParams, userID string) (*models.PersonalDrugOverviewResponse, error)
	GetPersonalDrugsByResidentID(residentID string, userID string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugsByResidentIDToday(residentID string, userID string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugByID(id string, userID string) (*entities.PersonalDrug, error)
	UpdatePersonalDrugByID(id string, req models.UpdatePersonalDrugRequest, userID string) (*entities.PersonalDrug, error)
	DeletePersonalDrugByID(id string, userID string) error

	CreateDrugPlan(req models.CreateDrugPlanRequest, userID string) (*entities.DrugPlan, error)
	GetDrugPlansTodayResidentSummary(userID string) (*models.DrugPlanResidentSummaryResponse, error)
	GetDrugPlansToday(userID string) ([]*entities.DrugPlan, error)
	GetDrugPlansOverview(req models.DrugPlanOverviewQueryParams, userID string) (*models.DrugPlanOverviewResponse, error)
	GetDrugAdministrationHistory(req models.DrugAdministrationHistoryQueryParams, userID string) (*models.DrugAdministrationHistoryResponse, error)
	GetDrugPlansByResidentID(residentID string, userID string) ([]*entities.DrugPlan, error)
	GetDrugPlansByResidentIDToday(residentID string, userID string) ([]*entities.DrugPlan, error)
	GetDrugPlans(userID string) ([]*entities.DrugPlan, error)
	GetDrugPlanByID(id string, userID string) (*entities.DrugPlan, error)
	UpdateDrugPlanByID(id string, req models.UpdateDrugPlanRequest, userID string) (*entities.DrugPlan, error)
	ForceGenerateTodayDrugPlans(userID string) (*models.DrugPlanGenerationResponse, error)
	ForceGenerateTodayDrugPlansByResidentID(residentID string, userID string) (*models.DrugPlanGenerationResponse, error)
	TakeDrugPlanByID(id string, req models.TakeDrugPlanByIDRequest, userID string) (*entities.DrugPlan, error)
	OmitDrugPlanByID(id string, req models.OmitDrugPlanByIDRequest, userID string) (*entities.DrugPlan, error)
	TakeDrugPlansByResidentIDToday(residentID string, req models.TakeDrugPlansByResidentRequest, userID string) ([]*entities.DrugPlan, error)
	OmitDrugPlansByResidentIDToday(residentID string, req models.OmitDrugPlansByResidentRequest, userID string) ([]*entities.DrugPlan, error)
	DeleteDrugPlanByID(id string, userID string) error
}

type DrugUseCaseImpl struct {
	drugRepo     medicine_repository.DrugRepository
	auditLogRepo audit_repository.AuditLogRepository
	userRepo     user_repository.UserRepository
}

func NewDrugUseCase(
	drugRepo medicine_repository.DrugRepository,
	auditLogRepo audit_repository.AuditLogRepository,
	userRepo user_repository.UserRepository,
) *DrugUseCaseImpl {
	return &DrugUseCaseImpl{
		drugRepo:     drugRepo,
		auditLogRepo: auditLogRepo,
		userRepo:     userRepo,
	}
}

func (uc *DrugUseCaseImpl) ensureMedicalStaff(userID string) error {
	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userRepo.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can access personal drug data")
	}

	return nil
}

// func (uc *DrugUseCaseImpl) CreateDrugMaster(req models.CreateDrugMasterRequest, userID string) (*entities.DrugMaster, error) {
// 	if err := uc.ensureMedicalStaff(userID); err != nil {
// 		return nil, err
// 	}

// 	name := strings.TrimSpace(req.Name)
// 	dose := strings.TrimSpace(req.Dose)
// 	if name == "" || dose == "" {
// 		return nil, errors.New("name and dose are required")
// 	}

// // Dose format must be: <number> <unit>, e.g. 50 mg, 5 mL
// pattern := regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)\s*(mcg|mg|g|kg|ml|l|iu)$`)
// matches := pattern.FindStringSubmatch(dose)
func normalizeDose(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	pattern := regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)\s*(.+)$`)
	matches := pattern.FindStringSubmatch(trimmed)
	if len(matches) != 3 {
		return "", errors.New("invalid dose format: use '<number> <unit>'")
	}

	amount := matches[1]
	unitRaw := strings.TrimSpace(matches[2])
	if unitRaw == "" {
		return "", errors.New("invalid dose format: use '<number> <unit>'")
	}

	unitKey := strings.ToLower(unitRaw)
	unitMap := map[string]string{
		"mcg": "mcg",
		"mg":  "mg",
		"g":   "g",
		"kg":  "kg",
		"ml":  "mL",
		"l":   "L",
		"iu":  "IU",
	}
	// dose = amount + " " + unitMap[unit]
	if mapped, ok := unitMap[unitKey]; ok {
		unitRaw = mapped
	} else {
		unitRaw = strings.Join(strings.Fields(unitRaw), " ")
	}

	return amount + " " + unitRaw, nil
}

func (uc *DrugUseCaseImpl) CreateDrugMaster(req models.CreateDrugMasterRequest, userID string) (*entities.DrugMaster, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	dose := strings.TrimSpace(req.Dose)
	if name == "" || dose == "" {
		return nil, errors.New("name and dose are required")
	}

	normalizedDose, err := normalizeDose(dose)
	if err != nil {
		return nil, err
	}
	dose = normalizedDose

	exists, err := uc.drugRepo.DrugMasterExistsByNameAndDose(name, dose)
	if err != nil {
		return nil, errors.New("failed to verify drug master existence: " + err.Error())
	}
	if exists {
		return nil, errors.New("drug master already exists")
	}

	drugMaster := &entities.DrugMaster{
		ID:   uuid.New().String(),
		Name: name,
		Dose: dose,
	}

	created, err := uc.drugRepo.CreateDrugMaster(drugMaster)
	if err != nil {
		return nil, errors.New("failed to create drug master: " + err.Error())
	}

	newValue, _ := json.Marshal(created)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "drug_masters",
		RecordID:  created.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newValue),
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return created, nil
}

func (uc *DrugUseCaseImpl) GetDrugMasters(userID string) ([]*entities.DrugMaster, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetAllDrugMasters()
	if err != nil {
		return nil, errors.New("failed to get drug masters: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) GetDrugMasterByID(id string, userID string) (*entities.DrugMaster, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetDrugMasterByID(id)
	if err != nil {
		return nil, errors.New("drug master not found: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) UpdateDrugMasterByID(id string, req models.UpdateDrugMasterRequest, userID string) (*entities.DrugMaster, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	current, err := uc.drugRepo.GetDrugMasterByID(id)
	if err != nil {
		return nil, errors.New("drug master not found: " + err.Error())
	}

	oldValue, _ := json.Marshal(current)

	newName := current.Name
	newDose := current.Dose

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}
		newName = name
	}

	if req.Dose != nil {
		dose := strings.TrimSpace(*req.Dose)
		if dose == "" {
			return nil, errors.New("dose cannot be empty")
		}

		// pattern := regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)\s*(mcg|mg|g|kg|ml|l|iu)$`)
		// matches := pattern.FindStringSubmatch(dose)
		// if len(matches) != 3 {
		// 	return nil, errors.New("invalid dose format: use '<number> <unit>' and allowed units are mcg, mg, g, kg, mL, L, IU")
		// }

		// amount := matches[1]
		// unit := strings.ToLower(matches[2])
		// unitMap := map[string]string{
		// 	"mcg": "mcg",
		// 	"mg":  "mg",
		// 	"g":   "g",
		// 	"kg":  "kg",
		// 	"ml":  "mL",
		// 	"l":   "L",
		// 	"iu":  "IU",
		normalizedDose, err := normalizeDose(dose)
		if err != nil {
			return nil, err
		}
		// newDose = amount + " " + unitMap[unit]
		newDose = normalizedDose
	}

	if newName != current.Name || newDose != current.Dose {
		exists, err := uc.drugRepo.DrugMasterExistsByNameAndDose(newName, newDose)
		if err != nil {
			return nil, errors.New("failed to verify drug master existence: " + err.Error())
		}
		if exists {
			return nil, errors.New("drug master already exists")
		}
	}

	current.Name = newName
	current.Dose = newDose

	updated, err := uc.drugRepo.UpdateDrugMaster(current)
	if err != nil {
		return nil, errors.New("failed to update drug master: " + err.Error())
	}

	newValue, _ := json.Marshal(updated)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "drug_masters",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldValue),
		NewValue:  string(newValue),
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return updated, nil
}

func (uc *DrugUseCaseImpl) DeleteDrugMasterByID(id string, userID string) error {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return err
	}

	current, err := uc.drugRepo.GetDrugMasterByID(id)
	if err != nil {
		return errors.New("drug master not found: " + err.Error())
	}

	oldValue, _ := json.Marshal(current)

	if err := uc.drugRepo.DeleteDrugMaster(id); err != nil {
		return errors.New("failed to delete drug master: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "drug_masters",
		RecordID:  id,
		UserID:    userID,
		Action:    audit_constants.AuditActionDelete,
		OldValue:  string(oldValue),
		NewValue:  "",
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return nil
}

func normalizeEnumInput(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func validateTakeType(value string) error {
	v := normalizeEnumInput(value)
	if v != "regular" && v != "as_needed" {
		return errors.New("take_type must be 'regular' or 'as_needed'")
	}
	return nil
}

func validateTimeOfDay(value string) error {
	v := strings.TrimSpace(value)
	if v == "" {
		return errors.New("time_of_day is required")
	}

	allowed := map[string]struct{}{
		"เช้า":    {},
		"กลางวัน": {},
		"เย็น":    {},
		"ก่อนนอน": {},
		"morning": {},
		"noon":    {},
		"evening": {},
		"bedtime": {},
	}

	parts := strings.Split(v, ",")
	for _, part := range parts {
		token := normalizeEnumInput(part)
		if token == "" {
			return errors.New("time_of_day contains an empty value")
		}

		if _, ok := allowed[token]; !ok {
			return errors.New("time_of_day must contain only: เช้า, กลางวัน, เย็น, ก่อนนอน")
		}
	}

	return nil
}

func parseDateInput(value string) (time.Time, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return time.Time{}, errors.New("date cannot be empty")
	}

	parsed, err := time.Parse("2006-01-02", v)
	if err == nil {
		return parsed, nil
	}

	parsed, err = time.Parse(time.RFC3339, v)
	if err == nil {
		return parsed, nil
	}

	return time.Time{}, errors.New("invalid date format: use YYYY-MM-DD")
}

func parseOptionalDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}

	parsed, err := parseDateInput(v)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parseOptionalDateTime(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, errors.New("invalid datetime format: use RFC3339")
	}

	return &parsed, nil
}

func (uc *DrugUseCaseImpl) cleanupExpiredAsNeededPersonalDrugs(residentID *string) (int, error) {
	today := time.Now()
	expiredPersonalDrugs, err := uc.drugRepo.GetExpiredAsNeededPersonalDrugs(today, residentID)
	if err != nil {
		return 0, errors.New("failed to get expired as-needed personal drugs: " + err.Error())
	}

	deletedCount := 0
	for _, personalDrug := range expiredPersonalDrugs {
		if err := uc.drugRepo.DeleteDrugPlansByPdID(personalDrug.ID); err != nil {
			return deletedCount, errors.New("failed to delete expired as-needed drug plans: " + err.Error())
		}

		if err := uc.drugRepo.DeletePersonalDrug(personalDrug.ID); err != nil {
			return deletedCount, errors.New("failed to delete expired as-needed personal drug: " + err.Error())
		}

		deletedCount++
	}

	return deletedCount, nil
}

func (uc *DrugUseCaseImpl) ensureTodayDrugPlansWithResult(residentID *string) (*models.DrugPlanGenerationResponse, error) {
	deletedCount, err := uc.cleanupExpiredAsNeededPersonalDrugs(residentID)
	if err != nil {
		return nil, err
	}

	today := time.Now()
	personalDrugs, err := uc.drugRepo.GetActivePersonalDrugsForDate(today, residentID)
	if err != nil {
		return nil, errors.New("failed to get active personal drugs for lazy generation: " + err.Error())
	}

	generatedCount := 0
	skippedCount := 0

	for _, personalDrug := range personalDrugs {
		exists, existsErr := uc.drugRepo.HasDrugPlanForPersonalDrugOnDate(personalDrug.ID, today)
		if existsErr != nil {
			return nil, errors.New("failed to check existing daily drug plan: " + existsErr.Error())
		}
		if exists {
			skippedCount++
			continue
		}

		initialIsOmitted := false
		now := time.Now()
		_, createErr := uc.drugRepo.CreateDrugPlan(&entities.DrugPlan{
			ID:             uuid.New().String(),
			PdID:           personalDrug.ID,
			IsTaken:        false,
			TakenAt:        nil,
			GivenByStaffID: "",
			IsOmitted:      &initialIsOmitted,
			OmittedReason:  nil,
			Notes:          nil,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		if createErr != nil {
			return nil, errors.New("failed to create lazy daily drug plan: " + createErr.Error())
		}

		generatedCount++
	}

	response := &models.DrugPlanGenerationResponse{
		GeneratedCount:       generatedCount,
		SkippedExistingCount: skippedCount,
		ExpiredDeletedCount:  deletedCount,
		Scope:                "all",
	}

	if residentID != nil && *residentID != "" {
		response.Scope = "resident"
		response.ResidentID = *residentID
	}

	return response, nil
}

func (uc *DrugUseCaseImpl) ensureTodayDrugPlans(residentID *string) error {
	_, err := uc.ensureTodayDrugPlansWithResult(residentID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *DrugUseCaseImpl) ForceGenerateTodayDrugPlans(userID string) (*models.DrugPlanGenerationResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	return uc.ensureTodayDrugPlansWithResult(nil)
}

func (uc *DrugUseCaseImpl) ForceGenerateTodayDrugPlansByResidentID(residentID string, userID string) (*models.DrugPlanGenerationResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.drugRepo.ResidentExistsByID(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	return uc.ensureTodayDrugPlansWithResult(&residentID)
}

func normalizeOptionalText(value *string) *string {
	if value == nil {
		return nil
	}

	v := strings.TrimSpace(*value)
	if v == "" {
		return nil
	}

	return &v
}

func isTodayInBangkok(t time.Time) bool {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return false
	}

	now := time.Now().In(loc)
	target := t.In(loc)

	return now.Year() == target.Year() && now.YearDay() == target.YearDay()
}

func (uc *DrugUseCaseImpl) resolveMedicalStaffIDByName(firstName string, lastName string) (string, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if firstName == "" || lastName == "" {
		return "", errors.New("staff_first_name and staff_last_name are required")
	}

	users, err := uc.userRepo.GetUsersByFirstAndLastName(firstName, lastName)
	if err != nil {
		return "", errors.New("failed to find staff by name: " + err.Error())
	}

	if len(users) == 0 {
		return "", errors.New("staff not found by provided name")
	}

	staffIDs := make([]string, 0)
	for _, user := range users {
		if user == nil {
			continue
		}

		roleName := strings.TrimSpace(user.Role.Name)
		if roleName == "" {
			role, roleErr := uc.userRepo.GetRoleByID(user.RoleID)
			if roleErr != nil {
				continue
			}
			roleName = role.Name
		}

		if roleName != user_constants.RoleMedicalStaff {
			continue
		}

		staff, staffErr := uc.userRepo.GetStaffByUserID(user.ID)
		if staffErr != nil {
			continue
		}

		staffIDs = append(staffIDs, staff.ID)
	}

	if len(staffIDs) == 0 {
		return "", errors.New("provided staff name is not a medical staff")
	}

	if len(staffIDs) > 1 {
		return "", errors.New("multiple medical staff found with the same name")
	}

	return staffIDs[0], nil
}

func (uc *DrugUseCaseImpl) applyDrugPlanActionByID(id string, givenByStaffID string, isTaken bool, omittedReason *string, note *string, userID string) (*entities.DrugPlan, error) {
	current, err := uc.drugRepo.GetDrugPlanByID(id)
	if err != nil {
		return nil, errors.New("drug plan not found: " + err.Error())
	}

	if !isTodayInBangkok(current.CreatedAt) {
		return nil, errors.New("drug plan action is allowed for today records only")
	}

	oldValue, _ := json.Marshal(current)

	now := time.Now()
	isOmitted := !isTaken

	current.IsTaken = isTaken
	current.IsOmitted = &isOmitted
	current.TakenAt = &now
	current.GivenByStaffID = givenByStaffID
	current.Notes = normalizeOptionalText(note)
	current.UpdatedAt = now

	if isTaken {
		current.OmittedReason = nil
	} else {
		current.OmittedReason = normalizeOptionalText(omittedReason)
		if current.OmittedReason == nil {
			return nil, errors.New("omitted_reason is required when omitting a drug")
		}
	}

	updated, err := uc.drugRepo.UpdateDrugPlan(current)
	if err != nil {
		return nil, errors.New("failed to update drug plan: " + err.Error())
	}

	newValue, _ := json.Marshal(updated)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "drug_plans",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldValue),
		NewValue:  string(newValue),
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return updated, nil
}

func validateAsNeededDateRange(takeType string, startDate *time.Time, endDate *time.Time) error {
	if normalizeEnumInput(takeType) != "as_needed" {
		return nil
	}

	if startDate == nil || endDate == nil {
		return errors.New("start_date and end_date are required when take_type is 'as_needed'")
	}

	if endDate.Before(*startDate) {
		return errors.New("end_date must be greater than or equal to start_date")
	}

	return nil
}

func (uc *DrugUseCaseImpl) CreatePersonalDrug(req models.CreatePersonalDrugRequest, userID string) (*entities.PersonalDrug, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if err := validateTakeType(req.TakeType); err != nil {
		return nil, err
	}
	if err := validateTimeOfDay(req.TimeOfDay); err != nil {
		return nil, err
	}
	if req.Frequency <= 0 {
		return nil, errors.New("frequency must be greater than 0")
	}

	startDate, err := parseOptionalDate(req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date: " + err.Error())
	}
	endDate, err := parseOptionalDate(req.EndDate)
	if err != nil {
		return nil, errors.New("invalid end_date: " + err.Error())
	}
	if err := validateAsNeededDateRange(req.TakeType, startDate, endDate); err != nil {
		return nil, err
	}
	if normalizeEnumInput(req.TakeType) != "as_needed" {
		startDate = nil
		endDate = nil
	}

	residentExists, err := uc.drugRepo.ResidentExistsByID(req.ResidentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	if _, err := uc.drugRepo.GetDrugMasterByID(req.DmID); err != nil {
		return nil, errors.New("drug master not found: " + err.Error())
	}

	personalDrug := &entities.PersonalDrug{
		ID:          uuid.New().String(),
		ResidentID:  req.ResidentID,
		DmID:        req.DmID,
		Amount:      req.Amount,
		AmountUnit:  req.AmountUnit,
		Frequency:   req.Frequency,
		TimeOfDay:   req.TimeOfDay,
		Timing:      req.Timing,
		Description: "",
		TakeType:    normalizeEnumInput(req.TakeType),
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if req.Description != nil {
		personalDrug.Description = *req.Description
	}

	created, err := uc.drugRepo.CreatePersonalDrug(personalDrug)
	if err != nil {
		return nil, errors.New("failed to create personal drug: " + err.Error())
	}

	newValue, _ := json.Marshal(created)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "personal_drugs",
		RecordID:  created.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newValue),
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return created, nil
}

func (uc *DrugUseCaseImpl) GetPersonalDrugsOverview(req models.PersonalDrugOverviewQueryParams, userID string) (*models.PersonalDrugOverviewResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if _, err := uc.cleanupExpiredAsNeededPersonalDrugs(nil); err != nil {
		return nil, err
	}

	if req.TimeOfDay != nil && strings.TrimSpace(*req.TimeOfDay) != "" {
		if err := validateTimeOfDay(*req.TimeOfDay); err != nil {
			return nil, err
		}
	}

	if req.TakeType != nil && strings.TrimSpace(*req.TakeType) != "" {
		if err := validateTakeType(*req.TakeType); err != nil {
			return nil, err
		}
		v := normalizeEnumInput(*req.TakeType)
		req.TakeType = &v
	}

	page := 1
	if req.Page != nil {
		if *req.Page <= 0 {
			return nil, errors.New("page must be greater than 0")
		}
		page = *req.Page
	}

	pageSize := 20
	if req.PageSize != nil {
		if *req.PageSize <= 0 {
			return nil, errors.New("page_size must be greater than 0")
		}
		pageSize = *req.PageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}

	result, total, err := uc.drugRepo.GetPersonalDrugsTodayCustom(req.TimeOfDay, req.Search, req.TakeType, page, pageSize)
	if err != nil {
		return nil, errors.New("failed to get personal drugs overview: " + err.Error())
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &models.PersonalDrugOverviewResponse{
		Items: result,
		Pagination: models.DrugAdministrationHistoryPagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (uc *DrugUseCaseImpl) GetPersonalDrugsByResidentIDToday(residentID string, userID string) ([]*entities.PersonalDrug, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.drugRepo.ResidentExistsByID(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	if _, err := uc.cleanupExpiredAsNeededPersonalDrugs(&residentID); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetPersonalDrugsByResidentIDToday(residentID)
	if err != nil {
		return nil, errors.New("failed to get personal drugs by resident: " + err.Error())
	}
	return result, nil
}

func (uc *DrugUseCaseImpl) GetPersonalDrugsByResidentID(residentID string, userID string) ([]*entities.PersonalDrug, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.drugRepo.ResidentExistsByID(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	if _, err := uc.cleanupExpiredAsNeededPersonalDrugs(&residentID); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetPersonalDrugsByResidentID(residentID)
	if err != nil {
		return nil, errors.New("failed to get personal drugs by resident: " + err.Error())
	}
	return result, nil
}

func (uc *DrugUseCaseImpl) GetPersonalDrugByID(id string, userID string) (*entities.PersonalDrug, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if _, err := uc.cleanupExpiredAsNeededPersonalDrugs(nil); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetPersonalDrugByID(id)
	if err != nil {
		return nil, errors.New("personal drug not found: " + err.Error())
	}
	return result, nil
}

func (uc *DrugUseCaseImpl) UpdatePersonalDrugByID(id string, req models.UpdatePersonalDrugRequest, userID string) (*entities.PersonalDrug, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	current, err := uc.drugRepo.GetPersonalDrugByID(id)
	if err != nil {
		return nil, errors.New("personal drug not found: " + err.Error())
	}

	oldValue, _ := json.Marshal(current)

	if req.ResidentID != nil {
		residentExists, err := uc.drugRepo.ResidentExistsByID(*req.ResidentID)
		if err != nil {
			return nil, errors.New("failed to verify resident existence: " + err.Error())
		}
		if !residentExists {
			return nil, errors.New("resident not found")
		}
		current.ResidentID = *req.ResidentID
	}

	if req.DmID != nil {
		if _, err := uc.drugRepo.GetDrugMasterByID(*req.DmID); err != nil {
			return nil, errors.New("drug master not found: " + err.Error())
		}
		current.DmID = *req.DmID
	}
	if req.Amount != nil {
		current.Amount = *req.Amount
	}
	if req.AmountUnit != nil {
		current.AmountUnit = *req.AmountUnit
	}
	if req.Frequency != nil {
		if *req.Frequency <= 0 {
			return nil, errors.New("frequency must be greater than 0")
		}
		current.Frequency = *req.Frequency
	}
	if req.TimeOfDay != nil {
		if err := validateTimeOfDay(*req.TimeOfDay); err != nil {
			return nil, err
		}
		current.TimeOfDay = *req.TimeOfDay
	}
	if req.Timing != nil {
		current.Timing = *req.Timing
	}
	if req.Description != nil {
		current.Description = *req.Description
	}

	newTakeType := current.TakeType
	if req.TakeType != nil {
		if err := validateTakeType(*req.TakeType); err != nil {
			return nil, err
		}
		newTakeType = normalizeEnumInput(*req.TakeType)
	}

	newStartDate := current.StartDate
	if req.StartDate != nil {
		parsedStartDate, err := parseOptionalDate(req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start_date: " + err.Error())
		}
		newStartDate = parsedStartDate
	}

	newEndDate := current.EndDate
	if req.EndDate != nil {
		parsedEndDate, err := parseOptionalDate(req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date: " + err.Error())
		}
		newEndDate = parsedEndDate
	}

	if err := validateAsNeededDateRange(newTakeType, newStartDate, newEndDate); err != nil {
		return nil, err
	}
	if newTakeType != "as_needed" {
		newStartDate = nil
		newEndDate = nil
	}

	current.TakeType = newTakeType
	current.StartDate = newStartDate
	current.EndDate = newEndDate
	current.UpdatedAt = time.Now()

	updated, err := uc.drugRepo.UpdatePersonalDrug(current)
	if err != nil {
		return nil, errors.New("failed to update personal drug: " + err.Error())
	}
	newValue, _ := json.Marshal(updated)

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "personal_drugs",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldValue),
		NewValue:  string(newValue),
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return updated, nil
}

func (uc *DrugUseCaseImpl) DeletePersonalDrugByID(id string, userID string) error {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return err
	}

	current, err := uc.drugRepo.GetPersonalDrugByID(id)
	if err != nil {
		return errors.New("personal drug not found: " + err.Error())
	}
	oldValue, _ := json.Marshal(current)

	if err := uc.drugRepo.DeletePersonalDrug(id); err != nil {
		return errors.New("failed to delete personal drug: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "personal_drugs",
		RecordID:  id,
		UserID:    userID,
		Action:    audit_constants.AuditActionDelete,
		OldValue:  string(oldValue),
		NewValue:  "",
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return nil
}

func (uc *DrugUseCaseImpl) CreateDrugPlan(req models.CreateDrugPlanRequest, userID string) (*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if _, err := uc.drugRepo.GetPersonalDrugByID(req.PdID); err != nil {
		return nil, errors.New("personal drug not found: " + err.Error())
	}

	if _, err := uc.userRepo.GetStaffByID(req.GivenByStaffID); err != nil {
		return nil, errors.New("staff not found: " + err.Error())
	}

	initialIsOmitted := false

	now := time.Now()
	drugPlan := &entities.DrugPlan{
		ID:             uuid.New().String(),
		PdID:           req.PdID,
		IsTaken:        false,
		TakenAt:        nil,
		GivenByStaffID: req.GivenByStaffID,
		IsOmitted:      &initialIsOmitted,
		OmittedReason:  nil,
		Notes:          req.Notes,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	created, err := uc.drugRepo.CreateDrugPlan(drugPlan)
	if err != nil {
		return nil, errors.New("failed to create drug plan: " + err.Error())
	}

	newValue, _ := json.Marshal(created)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "drug_plans",
		RecordID:  created.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newValue),
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return created, nil
}

func (uc *DrugUseCaseImpl) GetDrugPlansTodayResidentSummary(userID string) (*models.DrugPlanResidentSummaryResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if err := uc.ensureTodayDrugPlans(nil); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetDrugPlansTodayResidentSummary()
	if err != nil {
		return nil, errors.New("failed to get drug plans summary: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) GetDrugPlansToday(userID string) ([]*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if err := uc.ensureTodayDrugPlans(nil); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetDrugPlansToday()
	if err != nil {
		return nil, errors.New("failed to get today's drug plans: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) GetDrugPlansOverview(req models.DrugPlanOverviewQueryParams, userID string) (*models.DrugPlanOverviewResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if err := uc.ensureTodayDrugPlans(nil); err != nil {
		return nil, err
	}

	if req.TimeOfDay != nil && strings.TrimSpace(*req.TimeOfDay) != "" {
		if err := validateTimeOfDay(*req.TimeOfDay); err != nil {
			return nil, err
		}
	}

	if req.TakeType != nil && strings.TrimSpace(*req.TakeType) != "" {
		if err := validateTakeType(*req.TakeType); err != nil {
			return nil, err
		}
		v := normalizeEnumInput(*req.TakeType)
		req.TakeType = &v
	}

	if len(req.LabelIDs) > 0 {
		expandedLabelIDs := make([]string, 0, len(req.LabelIDs))
		for _, id := range req.LabelIDs {
			for _, part := range strings.Split(id, ",") {
				part = strings.TrimSpace(part)
				if part != "" {
					expandedLabelIDs = append(expandedLabelIDs, part)
				}
			}
		}
		req.LabelIDs = expandedLabelIDs
	}

	page := 1
	if req.Page != nil {
		if *req.Page <= 0 {
			return nil, errors.New("page must be greater than 0")
		}
		page = *req.Page
	}

	pageSize := 20
	if req.PageSize != nil {
		if *req.PageSize <= 0 {
			return nil, errors.New("page_size must be greater than 0")
		}
		pageSize = *req.PageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}

	result, total, err := uc.drugRepo.GetDrugPlansTodayCustom(req.TimeOfDay, req.Search, req.TakeType, req.Floor, req.LabelIDs, page, pageSize)
	if err != nil {
		return nil, errors.New("failed to get drug plans overview: " + err.Error())
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &models.DrugPlanOverviewResponse{
		Items: result,
		Pagination: models.DrugAdministrationHistoryPagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (uc *DrugUseCaseImpl) GetDrugAdministrationHistory(req models.DrugAdministrationHistoryQueryParams, userID string) (*models.DrugAdministrationHistoryResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if req.Date != nil && strings.TrimSpace(*req.Date) != "" {
		if _, err := parseDateInput(*req.Date); err != nil {
			return nil, errors.New("invalid date: " + err.Error())
		}
		d := strings.TrimSpace(*req.Date)
		req.Date = &d
	}

	if req.TimeOfDay != nil && strings.TrimSpace(*req.TimeOfDay) != "" {
		if err := validateTimeOfDay(*req.TimeOfDay); err != nil {
			return nil, err
		}
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		status := strings.ToLower(strings.TrimSpace(*req.Status))
		if status != "taken" && status != "omitted" && status != "pending" {
			return nil, errors.New("status must be one of: taken, omitted, pending")
		}
		req.Status = &status
	}

	page := 1
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}

	pageSize := 20
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = *req.PageSize
	}
	if pageSize > 100 {
		pageSize = 100
	}

	shouldEnsureToday := true
	if req.Date != nil && strings.TrimSpace(*req.Date) != "" {
		parsedDate, _ := parseDateInput(*req.Date)
		shouldEnsureToday = isTodayInBangkok(parsedDate)
	}

	if shouldEnsureToday {
		if err := uc.ensureTodayDrugPlans(nil); err != nil {
			return nil, err
		}
	}

	items, total, err := uc.drugRepo.GetDrugAdministrationHistory(req, page, pageSize)
	if err != nil {
		return nil, errors.New("failed to get drug administration history: " + err.Error())
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &models.DrugAdministrationHistoryResponse{
		Items: items,
		Pagination: models.DrugAdministrationHistoryPagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (uc *DrugUseCaseImpl) GetDrugPlansByResidentID(residentID string, userID string) ([]*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.drugRepo.ResidentExistsByID(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	if err := uc.ensureTodayDrugPlans(&residentID); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetDrugPlansByResidentID(residentID)
	if err != nil {
		return nil, errors.New("failed to get drug plans by resident: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) GetDrugPlansByResidentIDToday(residentID string, userID string) ([]*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.drugRepo.ResidentExistsByID(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	if err := uc.ensureTodayDrugPlans(&residentID); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetDrugPlansByResidentIDToday(residentID)
	if err != nil {
		return nil, errors.New("failed to get today's drug plans by resident: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) GetDrugPlans(userID string) ([]*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if err := uc.ensureTodayDrugPlans(nil); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetAllDrugPlans()
	if err != nil {
		return nil, errors.New("failed to get drug plans: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) GetDrugPlanByID(id string, userID string) (*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if err := uc.ensureTodayDrugPlans(nil); err != nil {
		return nil, err
	}

	result, err := uc.drugRepo.GetDrugPlanByID(id)
	if err != nil {
		return nil, errors.New("drug plan not found: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) UpdateDrugPlanByID(id string, req models.UpdateDrugPlanRequest, userID string) (*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	current, err := uc.drugRepo.GetDrugPlanByID(id)
	if err != nil {
		return nil, errors.New("drug plan not found: " + err.Error())
	}

	oldValue, _ := json.Marshal(current)

	if req.PdID != nil {
		if _, err := uc.drugRepo.GetPersonalDrugByID(*req.PdID); err != nil {
			return nil, errors.New("personal drug not found: " + err.Error())
		}
		current.PdID = *req.PdID
	}

	if req.GivenByStaffID != nil {
		if _, err := uc.userRepo.GetStaffByID(*req.GivenByStaffID); err != nil {
			return nil, errors.New("staff not found: " + err.Error())
		}
		current.GivenByStaffID = *req.GivenByStaffID
	}

	if req.IsTaken != nil {
		current.IsTaken = *req.IsTaken
	}

	if req.TakenAt != nil {
		parsedTakenAt, err := parseOptionalDateTime(req.TakenAt)
		if err != nil {
			return nil, err
		}
		current.TakenAt = parsedTakenAt
	}

	if req.IsOmitted != nil {
		current.IsOmitted = req.IsOmitted
	}

	if req.OmittedReason != nil {
		if strings.TrimSpace(*req.OmittedReason) == "" {
			current.OmittedReason = nil
		} else {
			current.OmittedReason = req.OmittedReason
		}
	}

	if req.Notes != nil {
		if strings.TrimSpace(*req.Notes) == "" {
			current.Notes = nil
		} else {
			current.Notes = req.Notes
		}
	}

	if current.IsOmitted != nil && *current.IsOmitted && (current.OmittedReason == nil || strings.TrimSpace(*current.OmittedReason) == "") {
		return nil, errors.New("omitted_reason is required when is_omitted is true")
	}

	current.UpdatedAt = time.Now()

	updated, err := uc.drugRepo.UpdateDrugPlan(current)
	if err != nil {
		return nil, errors.New("failed to update drug plan: " + err.Error())
	}

	newValue, _ := json.Marshal(updated)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "drug_plans",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldValue),
		NewValue:  string(newValue),
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return updated, nil
}

func (uc *DrugUseCaseImpl) TakeDrugPlanByID(id string, req models.TakeDrugPlanByIDRequest, userID string) (*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	staffID, err := uc.resolveMedicalStaffIDByName(req.StaffFirstName, req.StaffLastName)
	if err != nil {
		return nil, err
	}

	return uc.applyDrugPlanActionByID(id, staffID, true, nil, req.Note, userID)
}

func (uc *DrugUseCaseImpl) OmitDrugPlanByID(id string, req models.OmitDrugPlanByIDRequest, userID string) (*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	staffID, err := uc.resolveMedicalStaffIDByName(req.StaffFirstName, req.StaffLastName)
	if err != nil {
		return nil, err
	}

	omittedReason := strings.TrimSpace(req.OmittedReason)
	if omittedReason == "" {
		return nil, errors.New("omitted_reason is required")
	}

	return uc.applyDrugPlanActionByID(id, staffID, false, &omittedReason, req.Note, userID)
}

func (uc *DrugUseCaseImpl) TakeDrugPlansByResidentIDToday(residentID string, req models.TakeDrugPlansByResidentRequest, userID string) ([]*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.drugRepo.ResidentExistsByID(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	staffID, err := uc.resolveMedicalStaffIDByName(req.StaffFirstName, req.StaffLastName)
	if err != nil {
		return nil, err
	}

	drugPlans, err := uc.drugRepo.GetDrugPlansByResidentIDToday(residentID)
	if err != nil {
		return nil, errors.New("failed to get today's drug plans by resident: " + err.Error())
	}

	updatedPlans := make([]*entities.DrugPlan, 0, len(drugPlans))
	for _, drugPlan := range drugPlans {
		updated, updateErr := uc.applyDrugPlanActionByID(drugPlan.ID, staffID, true, nil, req.Note, userID)
		if updateErr != nil {
			return nil, updateErr
		}
		updatedPlans = append(updatedPlans, updated)
	}

	return updatedPlans, nil
}

func (uc *DrugUseCaseImpl) OmitDrugPlansByResidentIDToday(residentID string, req models.OmitDrugPlansByResidentRequest, userID string) ([]*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.drugRepo.ResidentExistsByID(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	staffID, err := uc.resolveMedicalStaffIDByName(req.StaffFirstName, req.StaffLastName)
	if err != nil {
		return nil, err
	}

	omittedReason := strings.TrimSpace(req.OmittedReason)
	if omittedReason == "" {
		return nil, errors.New("omitted_reason is required")
	}

	drugPlans, err := uc.drugRepo.GetDrugPlansByResidentIDToday(residentID)
	if err != nil {
		return nil, errors.New("failed to get today's drug plans by resident: " + err.Error())
	}

	updatedPlans := make([]*entities.DrugPlan, 0, len(drugPlans))
	for _, drugPlan := range drugPlans {
		updated, updateErr := uc.applyDrugPlanActionByID(drugPlan.ID, staffID, false, &omittedReason, req.Note, userID)
		if updateErr != nil {
			return nil, updateErr
		}
		updatedPlans = append(updatedPlans, updated)
	}

	return updatedPlans, nil
}

func (uc *DrugUseCaseImpl) DeleteDrugPlanByID(id string, userID string) error {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return err
	}

	current, err := uc.drugRepo.GetDrugPlanByID(id)
	if err != nil {
		return errors.New("drug plan not found: " + err.Error())
	}

	oldValue, _ := json.Marshal(current)

	if err := uc.drugRepo.DeleteDrugPlan(id); err != nil {
		return errors.New("failed to delete drug plan: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "drug_plans",
		RecordID:  id,
		UserID:    userID,
		Action:    audit_constants.AuditActionDelete,
		OldValue:  string(oldValue),
		NewValue:  "",
	}
	_, _ = uc.auditLogRepo.CreateAuditLog(auditLog)

	return nil
}
