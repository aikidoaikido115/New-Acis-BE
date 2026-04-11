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
	GetPersonalDrugsOverview(req models.PersonalDrugOverviewQueryParams, userID string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugsByResidentID(residentID string, userID string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugsByResidentIDToday(residentID string, userID string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugByID(id string, userID string) (*entities.PersonalDrug, error)
	UpdatePersonalDrugByID(id string, req models.UpdatePersonalDrugRequest, userID string) (*entities.PersonalDrug, error)
	DeletePersonalDrugByID(id string, userID string) error

	CreateDrugPlan(req models.CreateDrugPlanRequest, userID string) (*entities.DrugPlan, error)
	GetDrugPlansTodayResidentSummary(userID string) (*models.DrugPlanResidentSummaryResponse, error)
	GetDrugPlansToday(userID string) ([]*entities.DrugPlan, error)
	GetDrugPlansOverview(req models.DrugPlanOverviewQueryParams, userID string) ([]*entities.DrugPlan, error)
	GetDrugPlansByResidentID(residentID string, userID string) ([]*entities.DrugPlan, error)
	GetDrugPlansByResidentIDToday(residentID string, userID string) ([]*entities.DrugPlan, error)
	GetDrugPlans(userID string) ([]*entities.DrugPlan, error)
	GetDrugPlanByID(id string, userID string) (*entities.DrugPlan, error)
	UpdateDrugPlanByID(id string, req models.UpdateDrugPlanRequest, userID string) (*entities.DrugPlan, error)
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

	if userRole.Name != user_constants.RoleMedicalStaff {
		return errors.New("only users with 'Medical Staff' role can access personal drug data")
	}

	return nil
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

	// Dose format must be: <number> <unit>, e.g. 50 mg, 5 mL
	pattern := regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)\s*(mcg|mg|g|kg|ml|l|iu)$`)
	matches := pattern.FindStringSubmatch(dose)
	if len(matches) != 3 {
		return nil, errors.New("invalid dose format: use '<number> <unit>' and allowed units are mcg, mg, g, kg, mL, L, IU")
	}

	amount := matches[1]
	unit := strings.ToLower(matches[2])
	unitMap := map[string]string{
		"mcg": "mcg",
		"mg":  "mg",
		"g":   "g",
		"kg":  "kg",
		"ml":  "mL",
		"l":   "L",
		"iu":  "IU",
	}
	dose = amount + " " + unitMap[unit]

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

		pattern := regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)\s*(mcg|mg|g|kg|ml|l|iu)$`)
		matches := pattern.FindStringSubmatch(dose)
		if len(matches) != 3 {
			return nil, errors.New("invalid dose format: use '<number> <unit>' and allowed units are mcg, mg, g, kg, mL, L, IU")
		}

		amount := matches[1]
		unit := strings.ToLower(matches[2])
		unitMap := map[string]string{
			"mcg": "mcg",
			"mg":  "mg",
			"g":   "g",
			"kg":  "kg",
			"ml":  "mL",
			"l":   "L",
			"iu":  "IU",
		}
		newDose = amount + " " + unitMap[unit]
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

func (uc *DrugUseCaseImpl) GetPersonalDrugsOverview(req models.PersonalDrugOverviewQueryParams, userID string) ([]*entities.PersonalDrug, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
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

	result, err := uc.drugRepo.GetPersonalDrugsTodayCustom(req.TimeOfDay, req.Search, req.TakeType)
	if err != nil {
		return nil, errors.New("failed to get personal drugs overview: " + err.Error())
	}

	return result, nil
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

	initialIsOmmitted := false

	now := time.Now()
	drugPlan := &entities.DrugPlan{
		ID:             uuid.New().String(),
		PdID:           req.PdID,
		IsTaken:        false,
		TakenAt:        nil,
		GivenByStaffID: req.GivenByStaffID,
		IsOmmitted:     &initialIsOmmitted,
		OmmittedReason: nil,
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

	result, err := uc.drugRepo.GetDrugPlansToday()
	if err != nil {
		return nil, errors.New("failed to get today's drug plans: " + err.Error())
	}

	return result, nil
}

func (uc *DrugUseCaseImpl) GetDrugPlansOverview(req models.DrugPlanOverviewQueryParams, userID string) ([]*entities.DrugPlan, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
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

	result, err := uc.drugRepo.GetDrugPlansTodayCustom(req.TimeOfDay, req.Search, req.TakeType)
	if err != nil {
		return nil, errors.New("failed to get drug plans overview: " + err.Error())
	}

	return result, nil
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

	if req.IsOmmitted != nil {
		current.IsOmmitted = req.IsOmmitted
	}

	if req.OmmittedReason != nil {
		if strings.TrimSpace(*req.OmmittedReason) == "" {
			current.OmmittedReason = nil
		} else {
			current.OmmittedReason = req.OmmittedReason
		}
	}

	if req.Notes != nil {
		if strings.TrimSpace(*req.Notes) == "" {
			current.Notes = nil
		} else {
			current.Notes = req.Notes
		}
	}

	if current.IsOmmitted != nil && *current.IsOmmitted && (current.OmmittedReason == nil || strings.TrimSpace(*current.OmmittedReason) == "") {
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
