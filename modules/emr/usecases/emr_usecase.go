package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"strconv"

	// "mime/multipart"
	// "os"
	// "sync"
	"strings"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	audit_constants "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/constants"
	audit_repo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	emr_constants "github.com/aikidoaikido115/New-Acis-BE/modules/emr/constants"
	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	// "golang.org/x/text/unicode/norm"
)

type EmrUsecase interface {

	// Resident operations
	CreateResident(resident *entities.Resident, userID string) (*entities.Resident, error)
	GetResidentByID(id string, userID string) (*entities.Resident, error)
	GetResidentByRoomID(roomID string, userID string) ([]*entities.Resident, error)
	GetAllResidents(userID string) ([]*entities.Resident, error)
	GetResidentOverview(req models.ResidentQueryParams, userID string) (*models.ResidentOverviewListResponse, error)
	UpdateResidentByID(residentID string, data models.UpdateResidentRequest, userID string) (*entities.Resident, error)

	// Dashboard operations
	GetNumberOfResidentsDashboard(userID string) (models.NumberOfResidentsDashboardResponse, error)
	GetResidentGenderStatsDashboard(userID string) (models.ResidentGenderStatsDashboardResponse, error)
	GetResidentAllergyStatsDashboard(userID string) (models.ResidentAllergyStatsDashboardResponse, error)
	GetResidentDrugAllergyStatsDashboard(userID string) (models.ResidentDrugAllergyStatsDashboardResponse, error)

	// Room operations
	GetRoomByID(id string, userID string) (*entities.Room, error)
	GetAllRooms(userID string) ([]*entities.Room, error)
	CreateRoom(room *entities.Room, userID string) (*entities.Room, error)
	UpdateRoomByID(roomID string, data models.UpdateRoomRequest, userID string) (*entities.Room, error)

	// IntakeLabel operations
	CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error)
	// GetIntakeLabelByID(id string) (*entities.IntakeLabels, error)
	GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error)
	GetAllIntakeLabels(userID string) ([]*entities.IntakeLabels, error)
	// Allergy operations
	CreateAllergy(allergy *entities.Allergy) (*entities.Allergy, error)
	GetAllergyByName(allergyName string) (*entities.Allergy, error)
	GetAllAllergies(userID string) ([]*entities.Allergy, error)
	// DrugAllergy operations
	CreateDrugAllergy(drugAllergy *entities.DrugAllergy) (*entities.DrugAllergy, error)
	GetDrugAllergyByName(allergyName string) (*entities.DrugAllergy, error)
	GetAllDrugAllergies(userID string) ([]*entities.DrugAllergy, error)

	// ResidentLabel operations (many-to-many)
	CreateIntakeLabelByResidentID(residentID string, labels []models.IntakeLabelRequest, userID string) ([]*entities.ResidentLabels, error)
	GetResidentLabelsByResidentID(residentID string, userID string) ([]*entities.ResidentLabels, error)
	// ResidentAllergy operations (many-to-many)
	CreateAllergyByResidentID(residentID string, allergies []models.AllergyRequest, userID string) ([]*entities.ResidentAllergies, error)
	GetResidentAllergiesByResidentID(residentID string, userID string) ([]*entities.ResidentAllergies, error)
	GetAllResidentAllergies(userID string) ([]*models.ResidentAllergyListResponse, error)
	// ResidentDrugAllergy operations (many-to-many)
	CreateDrugAllergyByResidentID(residentID string, drugAllergies []models.DrugAllergyRequest, userID string) ([]*entities.ResidentDA, error)
	GetResidentDrugAllergiesByResidentID(residentID string, userID string) ([]*entities.ResidentDA, error)
	GetAllResidentDrugAllergies(userID string) ([]*models.ResidentDrugAllergyListResponse, error)

	// VitalSign operations
	CreateVitalSign(vitalSign *entities.VitalSign, dateInput string, userID string) (*entities.VitalSign, error)

	GetVitalSignsOverview(req models.VitalSignQueryParams, userID string) (*models.VitalSignsOverviewResponse, error)
	GetVitalSignsByResident(residentID string, dateInput string, isLatest string, userID string) ([]*entities.VitalSign, error)
	GetRoomVitalSigns(roomID string, isLatest string, userID string) ([]*entities.VitalSign, error)
	GetVitalSignsHistory(residentID string, userID string) ([]*entities.VitalSign, error)

	UpdateVitalSignByID(vitalSignID string, vitalSign *entities.VitalSign, userID string) (*entities.VitalSign, error)

	// LaboratoryValue operations
	CreateLaboratoryValue(laboratoryValue *entities.LaboratoryValue, dateInput string, timeOfDayInput string, userID string) (*entities.LaboratoryValue, error)
	GetLaboratoryValuesOverview(req models.LaboratoryValueQueryParams, userID string) (*models.LaboratoryValuesOverviewResponse, error)
	GetLaboratoryValuesByResident(residentID string, dateInput string, isLatest string, userID string) ([]*entities.LaboratoryValue, error)
	GetRoomLaboratoryValues(roomID string, isLatest string, userID string) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesHistory(residentID string, userID string) ([]*entities.LaboratoryValue, error)
	GetUrineOutputSumByResidentID(residentID string, req models.LaboratoryValueQueryParams, userID string) (*models.UrineOutputSummaryByResidentResponse, error)

	UpdateLaboratoryValueByID(laboratoryValueID string, laboratoryValue *entities.LaboratoryValue, userID string) (*entities.LaboratoryValue, error)

	// NurseNote operations
	CreateNurseNote(note *entities.NurseNote, userID string) (*entities.NurseNote, error)
	GetNurseNotesOverview(userID string) ([]*entities.NurseNote, error)
	GetNurseNotesByResidentID(residentID string, userID string) ([]*entities.NurseNote, error)
	UpdateNurseNoteByID(noteID string, note *entities.NurseNote, userID string) (*entities.NurseNote, error)
	DeleteNurseNoteByID(noteID string, userID string) error

	// WoundCareNote operations
	CreateWoundCareNote(note *entities.WoundCareNote, userID string, imageFile multipart.File) (*entities.WoundCareNote, error)
	GetWoundCareNotesOverview(userID string) ([]*entities.WoundCareNote, error)
	GetWoundCareNotesByResidentID(residentID string, userID string) ([]*entities.WoundCareNote, error)
	UpdateWoundCareNoteByID(noteID string, note *entities.WoundCareNote, userID string, imageFile multipart.File) (*entities.WoundCareNote, error)
	DeleteWoundCareNoteByID(noteID string, userID string) error

	// RelativeNote operations
	CreateRelativeNote(note *entities.RelativeNote, userID string) (*entities.RelativeNote, error)
	GetRelativeNotesOverview(userID string) ([]*entities.RelativeNote, error)
	GetRelativeNotesByResidentID(residentID string, userID string) ([]*entities.RelativeNote, error)
	UpdateRelativeNoteByID(noteID string, note *entities.RelativeNote, userID string) (*entities.RelativeNote, error)
	DeleteRelativeNoteByID(noteID string, userID string) error

	// Relative portal operations
	IssueRelativeMagicLink(residentID string, userID string) (*models.RelativeMagicLinkResponse, error)
	GetRelativeMagicLink(residentID string, userID string) (*models.RelativeMagicLinkResponse, error)
	RelativePortalLogin(req models.RelativePortalLoginRequest) (*models.RelativePortalLoginResponse, error)
	GetRelativeDashboard(userID string, dateInput string) (*models.RelativeDashboardResponse, error)
	GetRelativePatientInfo(userID string) (*models.RelativePatientInfoResponse, error)

	// DoctorOrder operations
	CreateDoctorOrder(order *entities.DoctorOrder, userID string) (*entities.DoctorOrder, error)
	GetDoctorOrdersOverview(userID string) ([]*entities.DoctorOrder, error)
	GetDoctorOrdersByResidentID(residentID string, userID string) ([]*entities.DoctorOrder, error)
	UpdateDoctorOrderByID(orderID string, order *entities.DoctorOrder, userID string) (*entities.DoctorOrder, error)
	DeleteDoctorOrderByID(orderID string, userID string) error
	//todo search resident by like sql
	//todo overview resident
}

type EmrUseCaseImpl struct {
	emrrepo      repositories.EmrRepository
	auditlogrepo audit_repo.AuditLogRepository
	userrepo     user_repo.UserRepository
	supa         configs.Supabase
	jwtSecret    string
}

func NewEmrUseCase(
	emrrepo repositories.EmrRepository,
	auditlogrepo audit_repo.AuditLogRepository,
	userrepo user_repo.UserRepository,
	supa configs.Supabase,
	jwtConfig configs.JWT) EmrUsecase {
	return &EmrUseCaseImpl{
		emrrepo:      emrrepo,
		auditlogrepo: auditlogrepo,
		userrepo:     userrepo,
		supa:         supa,
		jwtSecret:    jwtConfig.Secret,
	}
}

func (uc *EmrUseCaseImpl) uploadWoundCareImage(file multipart.File) (*string, error) {
	if file == nil {
		return nil, nil
	}

	fileExtension, err := utils.DetectFileType(file)
	if err != nil {
		return nil, errors.New("invalid file: " + err.Error())
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return nil, errors.New("failed to reset file pointer: " + err.Error())
	}

	fileName := uuid.New().String() + fileExtension
	imageURL, err := utils.UploadFile2Supa(file, fileName, "wound_care/", uc.supa)
	if err != nil {
		return nil, errors.New("failed to upload wound care image: " + err.Error())
	}

	return &imageURL, nil
}

func (uc *EmrUseCaseImpl) ensureMedicalStaff(userID string) error {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can access EMR")
	}

	return nil
}

func (uc *EmrUseCaseImpl) ensureRelative(userID string) error {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleRelative {
		return errors.New("only users with 'Relative' role can access relative portal")
	}

	return nil
}

func (uc *EmrUseCaseImpl) buildThaiDOBPassword(dateOfBirth time.Time) string {
	dd := dateOfBirth.Day()
	mm := int(dateOfBirth.Month())
	yyyyThai := dateOfBirth.Year() + 543
	return fmt.Sprintf("%02d%02d%04d", dd, mm, yyyyThai)
}

func (uc *EmrUseCaseImpl) ensureUniqueRelativeIdentity(base string) (string, string, error) {
	for i := 0; i < 1000; i++ {
		suffix := ""
		if i > 0 {
			suffix = fmt.Sprintf("_%d", i)
		}
		username := fmt.Sprintf("%s%s", base, suffix)
		email := fmt.Sprintf("%s%s@relative.local", base, suffix)

		usernameExists, err := uc.userrepo.UsernameExists(username)
		if err != nil {
			return "", "", err
		}
		if usernameExists {
			continue
		}

		emailExists, err := uc.userrepo.EmailExists(email)
		if err != nil {
			return "", "", err
		}
		if emailExists {
			continue
		}

		return username, email, nil
	}

	return "", "", errors.New("failed to allocate unique relative identity")
}

func (uc *EmrUseCaseImpl) ensureRelativeAccountForResident(resident *entities.Resident) (*entities.Relative, error) {
	existing, err := uc.emrrepo.GetRelativeByResidentID(resident.ID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to get relative by resident: " + err.Error())
	}

	role, err := uc.userrepo.GetRoleByName(user_constants.RoleRelative)
	if err != nil {
		return nil, errors.New("failed to get relative role: " + err.Error())
	}

	plainPassword := uc.buildThaiDOBPassword(resident.DateOfBirth)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash relative password: " + err.Error())
	}

	usernameBase := fmt.Sprintf("relative_%s", resident.ID)
	username, email, err := uc.ensureUniqueRelativeIdentity(usernameBase)
	if err != nil {
		return nil, errors.New("failed to build relative credentials: " + err.Error())
	}

	createdUser, err := uc.userrepo.CreateUser(&entities.User{
		ID:        uuid.New().String(),
		RoleID:    role.ID,
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		IsApprove: true,
		FirstName: "ญาติ",
		LastName:  resident.LastName,
		Nickname:  "ญาติ",
	})
	if err != nil {
		return nil, errors.New("failed to create relative user: " + err.Error())
	}

	createdRelative, err := uc.emrrepo.CreateRelative(&entities.Relative{
		ID:               uuid.New().String(),
		UserID:           createdUser.ID,
		ResidentID:       resident.ID,
		RelativePassword: string(hashedPassword),
		Relation:         "ญาติ",
		Phone:            "",
	})
	if err != nil {
		return nil, errors.New("failed to create relative profile: " + err.Error())
	}

	return createdRelative, nil
}

func (uc *EmrUseCaseImpl) buildRelativeMagicLink(residentID, token string) string {
	return "/relative/login?resident_id=" + url.QueryEscape(residentID) + "&token=" + url.QueryEscape(token)
}

func splitTextList(value *string) []string {
	if value == nil {
		return []string{}
	}

	normalized := strings.ReplaceAll(*value, "\r\n", "\n")
	parts := strings.FieldsFunc(normalized, func(r rune) bool {
		return r == '\n' || r == ',' || r == ';'
	})

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (uc *EmrUseCaseImpl) resolveStaffDisplayName(staffID string, cache map[string]string) string {
	if staffID == "" {
		return ""
	}

	if cache != nil {
		if name, ok := cache[staffID]; ok {
			return name
		}
	}

	staff, err := uc.userrepo.GetStaffByID(staffID)
	if err != nil {
		if cache != nil {
			cache[staffID] = staffID
		}
		return staffID
	}

	fullName := strings.TrimSpace(strings.TrimSpace(staff.User.FirstName) + " " + strings.TrimSpace(staff.User.LastName))
	if fullName == "" {
		fullName = strings.TrimSpace(staff.User.Nickname)
	}
	if fullName == "" {
		fullName = strings.TrimSpace(staff.User.Username)
	}
	if fullName == "" {
		fullName = staffID
	}

	if cache != nil {
		cache[staffID] = fullName
	}

	return fullName
}

func (uc *EmrUseCaseImpl) populateNurseNotesStaffNames(notes []*entities.NurseNote) {
	cache := map[string]string{}
	for _, note := range notes {
		if note == nil {
			continue
		}
		note.CreatedByStaffName = uc.resolveStaffDisplayName(note.CreatedByStaffID, cache)
	}
}

func (uc *EmrUseCaseImpl) populateNurseNoteStaffName(note *entities.NurseNote) {
	if note == nil {
		return
	}
	note.CreatedByStaffName = uc.resolveStaffDisplayName(note.CreatedByStaffID, nil)
}

func (uc *EmrUseCaseImpl) populateWoundCareNotesStaffNames(notes []*entities.WoundCareNote) {
	cache := map[string]string{}
	for _, note := range notes {
		if note == nil {
			continue
		}
		note.CreatedByStaffName = uc.resolveStaffDisplayName(note.CreatedByStaffID, cache)
	}
}

func (uc *EmrUseCaseImpl) populateWoundCareNoteStaffName(note *entities.WoundCareNote) {
	if note == nil {
		return
	}
	note.CreatedByStaffName = uc.resolveStaffDisplayName(note.CreatedByStaffID, nil)
}

func (uc *EmrUseCaseImpl) populateRelativeNotesStaffNames(notes []*entities.RelativeNote) {
	cache := map[string]string{}
	for _, note := range notes {
		if note == nil {
			continue
		}
		note.CreatedByStaffName = uc.resolveStaffDisplayName(note.CreatedByStaffID, cache)
	}
}

func (uc *EmrUseCaseImpl) populateRelativeNoteStaffName(note *entities.RelativeNote) {
	if note == nil {
		return
	}
	note.CreatedByStaffName = uc.resolveStaffDisplayName(note.CreatedByStaffID, nil)
}

func (uc *EmrUseCaseImpl) populateDoctorOrdersStaffNames(orders []*entities.DoctorOrder) {
	cache := map[string]string{}
	for _, order := range orders {
		if order == nil {
			continue
		}
		order.CreatedByStaffName = uc.resolveStaffDisplayName(order.CreatedByStaffID, cache)
	}
}

func (uc *EmrUseCaseImpl) populateDoctorOrderStaffName(order *entities.DoctorOrder) {
	if order == nil {
		return
	}
	order.CreatedByStaffName = uc.resolveStaffDisplayName(order.CreatedByStaffID, nil)
}

func (uc *EmrUseCaseImpl) CreateResident(resident *entities.Resident, userID string) (*entities.Resident, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	if resident.DateOfBirth.IsZero() {
		return nil, errors.New("date_of_birth is required")
	}

	if resident.DateOfBirth.After(time.Now()) {
		return nil, errors.New("date_of_birth cannot be in the future")
	}

	roomExists, err := uc.emrrepo.RoomExists(resident.RoomID)
	if err != nil {
		return nil, errors.New("failed to verify room existence: " + err.Error())
	}
	if !roomExists {
		return nil, errors.New("room does not exist")
	}

	resident.Gender = strings.ToLower(strings.TrimSpace(resident.Gender))
	if resident.Gender != "male" && resident.Gender != "female" && resident.Gender != "other" {
		return nil, errors.New("gender must be either 'male', 'female', or 'other'")
	}

	resident.IdCardNumber = strings.TrimSpace(resident.IdCardNumber)
	if len(resident.IdCardNumber) != 13 {
		return nil, errors.New("ID card number must be 13 characters long or not missing")
	}

	idCardExists, err := uc.emrrepo.IdCardNumberExists(resident.IdCardNumber)
	if err != nil {
		return nil, errors.New("failed to verify ID card number existence: " + err.Error())
	}
	if idCardExists {
		return nil, errors.New("ID card number already exists")
	}

	if resident.ExpectedCheckOutDate != nil && resident.ExpectedCheckOutDate.Before(resident.CheckInDate) {
		return nil, errors.New("expected check-out date cannot be before check-in date")
	}

	if resident.Status != emr_constants.Active && resident.Status != emr_constants.InActive {
		return nil, errors.New("status must be either 'active' or 'inactive' or not missing")
	}

	if resident.ResucitationStatus != nil && *resident.ResucitationStatus != emr_constants.CPR && *resident.ResucitationStatus != emr_constants.DNR {
		return nil, errors.New("resuscitation status must be either 'CPR' or 'DNR'")
	}

	if resident.EmergencyHospitalPhone != nil && len(*resident.EmergencyHospitalPhone) != 10 {
		return nil, errors.New("emergency hospital phone must be 10 characters long")
	}

	resident.ID = uuid.New().String()

	createdResident, err := uc.emrrepo.CreateResident(resident)
	if err != nil {
		return nil, errors.New("failed to create resident: " + err.Error())
	}

	if _, err := uc.ensureRelativeAccountForResident(createdResident); err != nil {
		return nil, err
	}

	newResidentData, _ := json.Marshal(map[string]interface{}{
		"first_name":                    createdResident.FirstName,
		"last_name":                     createdResident.LastName,
		"date_of_birth":                 createdResident.DateOfBirth,
		"gender":                        createdResident.Gender,
		"nickname":                      createdResident.Nickname,
		"id_card_number":                createdResident.IdCardNumber,
		"purpose_of_stay":               createdResident.PurposeOfStay,
		"check_in_date":                 createdResident.CheckInDate,
		"expected_check_out_date":       createdResident.ExpectedCheckOutDate,
		"status":                        createdResident.Status,
		"pre_existing_conditions":       createdResident.PreExistingConditions,
		"pre_existing_conditions_notes": createdResident.PreExistingConditionsNotes,
		"resuscitation_status":          createdResident.ResucitationStatus,
		"surgical_history":              createdResident.SugicalHistory,
		"preferred_emergency_hospital":  createdResident.PreferredEmergencyHospital,
		"emergency_hospital_phone":      createdResident.EmergencyHospitalPhone,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "residents",
		RecordID:  createdResident.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newResidentData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for new resident %s: %v", createdResident.ID, err)
	}

	return createdResident, nil
}

func (uc *EmrUseCaseImpl) GetResidentByID(id string, userID string) (*entities.Resident, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	resident, err := uc.emrrepo.GetResidentByID(id)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}
	return resident, nil
}

func (uc *EmrUseCaseImpl) GetResidentByRoomID(roomID string, userID string) ([]*entities.Resident, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residents, err := uc.emrrepo.GetResidentByRoomID(roomID)
	if err != nil {
		return nil, errors.New("failed to get residents by room ID: " + err.Error())
	}
	return residents, nil
}

func (uc *EmrUseCaseImpl) GetAllResidents(userID string) ([]*entities.Resident, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residents, err := uc.emrrepo.GetAllResidents()
	if err != nil {
		return nil, errors.New("failed to get all residents: " + err.Error())
	}
	return residents, nil
}

func (uc *EmrUseCaseImpl) GetResidentOverview(req models.ResidentQueryParams, userID string) (*models.ResidentOverviewListResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	var (
		residents []*entities.Resident
		total     int64
		err       error
	)

	// Swagger UI ส่ง label_ids เป็น comma-separated string เดียว เช่น "id1,id2"
	// ต้อง split ก่อนใช้งาน
	var expandedLabelIDs []string
	for _, id := range req.LabelIDs {
		for _, part := range strings.Split(id, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				expandedLabelIDs = append(expandedLabelIDs, part)
			}
		}
	}
	req.LabelIDs = expandedLabelIDs

	hasFilter := req.Floor != nil || len(req.LabelIDs) > 0 ||
		(req.Search != nil && *req.Search != "") ||
		(req.Status != nil && *req.Status != "")

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
	} else if req.Limit > 0 {
		pageSize = req.Limit
	}
	if pageSize > 100 {
		pageSize = 100
	}

	hasPagination := req.Page != nil || req.PageSize != nil || req.Limit > 0 || req.Offset > 0
	if req.Page == nil && req.Offset > 0 {
		page = (req.Offset / pageSize) + 1
	}

	req.Limit = pageSize
	req.Offset = (page - 1) * pageSize

	if req.Status != nil && *req.Status != "" && *req.Status != emr_constants.Active && *req.Status != emr_constants.InActive {
		return nil, errors.New("status must be 'active' or 'inactive'")
	}

	log.Printf("เข้ามาใน overview")
	if hasFilter || hasPagination {
		residents, total, err = uc.emrrepo.GetResidentsCustom(req)
		log.Printf("มีการใช้ filter ใน GetResidentOverview: floor=%v, label_ids=%v, search=%v, status=%v | จำนวน residents ที่ได้จาก custom query: %d", req.Floor, req.LabelIDs, req.Search, req.Status, len(residents))
	} else {
		residents, err = uc.emrrepo.GetAllResidents()
		total = int64(len(residents))
		log.Printf("ไม่มีการใช้ filter ใน GetResidentOverview | จำนวน residents ที่ได้จาก GetAllResidents: %d", len(residents))
	}
	if err != nil {
		return nil, err
	}

	response := make([]*models.ResidentOverviewResponse, 0)
	for _, r := range residents {
		labelNames := make([]string, len(r.ResidentLabels))
		for i, rl := range r.ResidentLabels {
			labelNames[i] = rl.IntakeLabel.LabelName
		}

		response = append(response, &models.ResidentOverviewResponse{
			ResidentID:   r.ID,
			FirstName:    r.FirstName,
			LastName:     r.LastName,
			Nickname:     r.Nickname,
			RoomNumber:   r.Room.RoomNumber,
			IntakeLabels: labelNames,
		})
	}
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &models.ResidentOverviewListResponse{
		Items: response,
		Pagination: models.OverviewPagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (uc *EmrUseCaseImpl) UpdateResidentByID(residentID string, data models.UpdateResidentRequest, userID string) (*entities.Resident, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	resident, err := uc.emrrepo.GetResidentByID(residentID)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}

	oldResidentData, _ := json.Marshal(map[string]interface{}{
		"room_id":                       resident.RoomID,
		"first_name":                    resident.FirstName,
		"last_name":                     resident.LastName,
		"date_of_birth":                 resident.DateOfBirth,
		"gender":                        resident.Gender,
		"nickname":                      resident.Nickname,
		"id_card_number":                resident.IdCardNumber,
		"purpose_of_stay":               resident.PurposeOfStay,
		"check_in_date":                 resident.CheckInDate,
		"expected_check_out_date":       resident.ExpectedCheckOutDate,
		"status":                        resident.Status,
		"pre_existing_conditions":       resident.PreExistingConditions,
		"pre_existing_conditions_notes": resident.PreExistingConditionsNotes,
		"resuscitation_status":          resident.ResucitationStatus,
		"surgical_history":              resident.SugicalHistory,
		"preferred_emergency_hospital":  resident.PreferredEmergencyHospital,
		"emergency_hospital_phone":      resident.EmergencyHospitalPhone,
	})

	if data.RoomID != nil {
		roomExists, err := uc.emrrepo.RoomExists(*data.RoomID)
		if err != nil {
			return nil, errors.New("failed to verify room existence: " + err.Error())
		}
		if !roomExists {
			return nil, errors.New("room does not exist")
		}
		resident.RoomID = *data.RoomID
	}

	if data.FirstName != nil {
		resident.FirstName = *data.FirstName
	}

	if data.LastName != nil {
		resident.LastName = *data.LastName
	}

	if data.DateOfBirth != nil {
		if data.DateOfBirth.IsZero() {
			return nil, errors.New("date_of_birth cannot be zero")
		}
		if data.DateOfBirth.After(time.Now()) {
			return nil, errors.New("date_of_birth cannot be in the future")
		}
		resident.DateOfBirth = *data.DateOfBirth
	}

	if data.Gender != nil {
		gender := strings.ToLower(strings.TrimSpace(*data.Gender))
		if gender != "male" && gender != "female" && gender != "other" {
			return nil, errors.New("gender must be either 'male', 'female', or 'other'")
		}
		resident.Gender = gender
	}

	if data.Nickname != nil {
		resident.Nickname = data.Nickname
	}

	if data.IdCardNumber != nil {
		idCard := strings.TrimSpace(*data.IdCardNumber)
		if len(idCard) != 13 {
			return nil, errors.New("ID card number must be 13 characters long")
		}
		if idCard != resident.IdCardNumber {
			idCardExists, err := uc.emrrepo.IdCardNumberExists(idCard)
			if err != nil {
				return nil, errors.New("failed to verify ID card number existence: " + err.Error())
			}
			if idCardExists {
				return nil, errors.New("ID card number already exists")
			}
		}
		resident.IdCardNumber = idCard
	}

	if data.PurposeOfStay != nil {
		resident.PurposeOfStay = data.PurposeOfStay
	}

	if data.CheckInDate != nil {
		resident.CheckInDate = *data.CheckInDate
	}

	if data.ExpectedCheckOutDate != nil {
		if data.ExpectedCheckOutDate.Before(resident.CheckInDate) {
			return nil, errors.New("expected check-out date cannot be before check-in date")
		}
		resident.ExpectedCheckOutDate = data.ExpectedCheckOutDate
	}

	if data.Status != nil {
		if *data.Status != emr_constants.Active && *data.Status != emr_constants.InActive {
			return nil, errors.New("status must be either 'active' or 'inactive'")
		}
		resident.Status = *data.Status
	}

	if data.PreExistingConditions != nil {
		resident.PreExistingConditions = data.PreExistingConditions
	}

	if data.PreExistingConditionsNotes != nil {
		resident.PreExistingConditionsNotes = data.PreExistingConditionsNotes
	}

	if data.ResucitationStatus != nil {
		if *data.ResucitationStatus != emr_constants.CPR && *data.ResucitationStatus != emr_constants.DNR {
			return nil, errors.New("resuscitation status must be either 'CPR' or 'DNR'")
		}
		resident.ResucitationStatus = data.ResucitationStatus
	}

	if data.SugicalHistory != nil {
		resident.SugicalHistory = data.SugicalHistory
	}

	if data.PreferredEmergencyHospital != nil {
		resident.PreferredEmergencyHospital = data.PreferredEmergencyHospital
	}

	if data.EmergencyHospitalPhone != nil {
		if len(*data.EmergencyHospitalPhone) != 10 {
			return nil, errors.New("emergency hospital phone must be 10 characters long")
		}
		resident.EmergencyHospitalPhone = data.EmergencyHospitalPhone
	}

	updatedResident, err := uc.emrrepo.UpdateResident(resident)
	if err != nil {
		return nil, errors.New("failed to update resident: " + err.Error())
	}

	newResidentData, _ := json.Marshal(map[string]interface{}{
		"first_name":                    updatedResident.FirstName,
		"last_name":                     updatedResident.LastName,
		"date_of_birth":                 updatedResident.DateOfBirth,
		"gender":                        updatedResident.Gender,
		"nickname":                      updatedResident.Nickname,
		"id_card_number":                updatedResident.IdCardNumber,
		"purpose_of_stay":               updatedResident.PurposeOfStay,
		"check_in_date":                 updatedResident.CheckInDate,
		"expected_check_out_date":       updatedResident.ExpectedCheckOutDate,
		"status":                        updatedResident.Status,
		"pre_existing_conditions":       updatedResident.PreExistingConditions,
		"pre_existing_conditions_notes": updatedResident.PreExistingConditionsNotes,
		"resuscitation_status":          updatedResident.ResucitationStatus,
		"surgical_history":              updatedResident.SugicalHistory,
		"preferred_emergency_hospital":  updatedResident.PreferredEmergencyHospital,
		"emergency_hospital_phone":      updatedResident.EmergencyHospitalPhone,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "residents",
		RecordID:  updatedResident.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldResidentData),
		NewValue:  string(newResidentData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for resident %s: %v", updatedResident.ID, err)
	}

	if data.Labels != nil {
		existingLabels, _ := uc.emrrepo.GetResidentLabelsByResidentID(resident.ID)

		err = uc.emrrepo.DeleteResidentLabelsByResidentID(resident.ID)
		if err != nil {
			return nil, errors.New("failed to delete existing labels: " + err.Error())
		}

		if len(existingLabels) > 0 {
			oldLabelsData, _ := json.Marshal(existingLabels)
			auditLog := &entities.AuditLogs{
				ID:        uuid.New().String(),
				TableName: "resident_labels",
				RecordID:  resident.ID,
				UserID:    userID,
				Action:    audit_constants.AuditActionDelete,
				OldValue:  string(oldLabelsData),
				NewValue:  "",
			}
			_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
			if err != nil {
				log.Printf("[ERROR] Failed to create audit log for deleting resident labels %s: %v", resident.ID, err)
			}
		}

		seenLabelIDs := make(map[string]bool)

		for _, label := range data.Labels {
			if len(strings.TrimSpace(label.LabelName)) == 0 {
				return nil, errors.New("label name cannot be empty or whitespace")
			}

			labelID, err := uc.getOrCreateLabelID(label.LabelName)
			if err != nil {
				return nil, errors.New("failed to get or create label: " + err.Error())
			}

			if seenLabelIDs[labelID] {
				continue
			}
			seenLabelIDs[labelID] = true

			residentLabel := &entities.ResidentLabels{
				ResidentID: resident.ID,
				LabelID:    labelID,
				NoteText:   label.NoteText,
				NotedAt:    time.Now(), // เวลาปัจจุบัน
			}
			createdLabel, err := uc.emrrepo.CreateIntakeLabelByResidentID(residentLabel)
			if err != nil {
				return nil, errors.New("failed to create resident label: " + err.Error())
			}

			// Create audit log for each new label
			newLabelData, _ := json.Marshal(map[string]interface{}{
				"resident_id": createdLabel.ResidentID,
				"label_id":    createdLabel.LabelID,
				"label_name":  label.LabelName,
				"note_text":   label.NoteText,
				"noted_at":    createdLabel.NotedAt,
			})
			auditLog := &entities.AuditLogs{
				ID:        uuid.New().String(),
				TableName: "resident_labels",
				RecordID:  createdLabel.ResidentID + "-" + createdLabel.LabelID,
				UserID:    userID,
				Action:    audit_constants.AuditActionInsert,
				OldValue:  "",
				NewValue:  string(newLabelData),
			}
			_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
			if err != nil {
				log.Printf("[ERROR] Failed to create audit log for resident label %s-%s: %v", createdLabel.ResidentID, createdLabel.LabelID, err)
			}
		}

		// Fetch resident again to get updated labels
		finalResident, err := uc.emrrepo.GetResidentByID(resident.ID)
		if err != nil {
			return nil, errors.New("failed to fetch updated resident: " + err.Error())
		}
		return finalResident, nil
	}

	return updatedResident, nil
}

func (uc *EmrUseCaseImpl) GetNumberOfResidentsDashboard(userID string) (models.NumberOfResidentsDashboardResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return models.NumberOfResidentsDashboardResponse{}, err
	}

	response, err := uc.emrrepo.GetNumberOfResidentsDashboard()
	if err != nil {
		return models.NumberOfResidentsDashboardResponse{}, errors.New("failed to get dashboard data: " + err.Error())
	}
	return response, nil
}

func (uc *EmrUseCaseImpl) GetResidentGenderStatsDashboard(userID string) (models.ResidentGenderStatsDashboardResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return models.ResidentGenderStatsDashboardResponse{}, err
	}

	response, err := uc.emrrepo.GetNumberOfResidentGender()
	if err != nil {
		return models.ResidentGenderStatsDashboardResponse{}, errors.New("failed to get resident gender stats: " + err.Error())
	}

	if response.TotalResidents > 0 {
		response.MalePercentage = (float32(response.SumOfMale) / float32(response.TotalResidents)) * 100
		response.FemalePercentage = (float32(response.SumOfFemale) / float32(response.TotalResidents)) * 100
	} else {
		response.MalePercentage = 0
		response.FemalePercentage = 0
	}

	return response, nil
}

func (uc *EmrUseCaseImpl) GetResidentAllergyStatsDashboard(userID string) (models.ResidentAllergyStatsDashboardResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return models.ResidentAllergyStatsDashboardResponse{}, err
	}

	response, err := uc.emrrepo.GetResidentAllergyStatsDashboard()
	if err != nil {
		return models.ResidentAllergyStatsDashboardResponse{}, errors.New("failed to get resident allergy stats: " + err.Error())
	}

	return response, nil
}

func (uc *EmrUseCaseImpl) GetResidentDrugAllergyStatsDashboard(userID string) (models.ResidentDrugAllergyStatsDashboardResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return models.ResidentDrugAllergyStatsDashboardResponse{}, err
	}

	response, err := uc.emrrepo.GetResidentDrugAllergyStatsDashboard()
	if err != nil {
		return models.ResidentDrugAllergyStatsDashboardResponse{}, errors.New("failed to get resident drug allergy stats: " + err.Error())
	}

	return response, nil
}

func (uc *EmrUseCaseImpl) GetRoomByID(id string, userID string) (*entities.Room, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	room, err := uc.emrrepo.GetRoomByID(id)
	if err != nil {
		return nil, errors.New("room not found: " + err.Error())
	}
	return room, nil
}

func (uc *EmrUseCaseImpl) GetAllRooms(userID string) ([]*entities.Room, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	rooms, err := uc.emrrepo.GetAllRooms()
	if err != nil {
		return nil, errors.New("failed to get all rooms: " + err.Error())
	}
	return rooms, nil
}

func (uc *EmrUseCaseImpl) CreateRoom(room *entities.Room, userID string) (*entities.Room, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	roomNumberExists, err := uc.emrrepo.RoomNumberExists(room.RoomNumber)
	if err != nil {
		return nil, errors.New("failed to verify room number existence: " + err.Error())
	}
	if roomNumberExists {
		return nil, errors.New("room number already exists")
	}
	room.ID = uuid.New().String()
	createdRoom, err := uc.emrrepo.CreateRoom(room)
	if err != nil {
		return nil, errors.New("failed to create room: " + err.Error())
	}

	// Create audit log
	newRoomData, _ := json.Marshal(map[string]interface{}{
		"room_number": createdRoom.RoomNumber,
		"floor":       createdRoom.Floor,
		"staff_id":    createdRoom.StaffID,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "rooms",
		RecordID:  createdRoom.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newRoomData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for new room %s: %v", createdRoom.ID, err)
	}

	return createdRoom, nil
}

func (uc *EmrUseCaseImpl) UpdateRoomByID(roomID string, data models.UpdateRoomRequest, userID string) (*entities.Room, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	room, err := uc.emrrepo.GetRoomByID(roomID)
	if err != nil {
		return nil, errors.New("room not found: " + err.Error())
	}

	oldRoomData, _ := json.Marshal(map[string]interface{}{
		"staff_id":    room.StaffID,
		"floor":       room.Floor,
		"room_number": room.RoomNumber,
	})

	if data.StaffID != nil {
		room.StaffID = data.StaffID
	}

	updatedRoom, err := uc.emrrepo.UpdateRoom(room)
	if err != nil {
		return nil, errors.New("failed to update room: " + err.Error())
	}

	// Create audit log
	newRoomData, _ := json.Marshal(map[string]interface{}{
		"staff_id":    updatedRoom.StaffID,
		"floor":       updatedRoom.Floor,
		"room_number": updatedRoom.RoomNumber,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "rooms",
		RecordID:  updatedRoom.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldRoomData),
		NewValue:  string(newRoomData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for room %s: %v", updatedRoom.ID, err)
	}

	return updatedRoom, nil
}

func (uc *EmrUseCaseImpl) CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error) {
	labelExists, err := uc.emrrepo.LabelExists(label.LabelName)
	if err != nil {
		return nil, errors.New("failed to verify label existence: " + err.Error())
	}
	if labelExists {
		return nil, errors.New("label already exists")
	}

	label.ID = uuid.New().String()
	createdLabel, err := uc.emrrepo.CreateIntakeLabel(label)
	if err != nil {
		return nil, errors.New("failed to create intake label: " + err.Error())
	}
	return createdLabel, nil
}

// func (uc *EmrUseCaseImpl) GetIntakeLabelByID(id string) (*entities.IntakeLabels, error) {
// 	label, err := uc.emrrepo.GetIntakeLabelByID(id)
// 	if err != nil {
// 		return nil, errors.New("intake label not found: " + err.Error())
// 	}
// 	return label, nil
// }

func (uc *EmrUseCaseImpl) GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error) {
	label, err := uc.emrrepo.GetIntakeLabelByName(labelName)
	if err != nil {
		return nil, errors.New("intake label not found: " + err.Error())
	}
	return label, nil
}

func (uc *EmrUseCaseImpl) GetAllIntakeLabels(userID string) ([]*entities.IntakeLabels, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	labels, err := uc.emrrepo.GetAllIntakeLabels()
	if err != nil {
		return nil, errors.New("failed to get all intake labels: " + err.Error())
	}
	return labels, nil
}

func (uc *EmrUseCaseImpl) CreateAllergy(allergy *entities.Allergy) (*entities.Allergy, error) {
	allergyExists, err := uc.emrrepo.AllergyExists(allergy.AllergyName)
	if err != nil {
		return nil, errors.New("failed to verify allergy existence: " + err.Error())
	}
	if allergyExists {
		return nil, errors.New("allergy already exists")
	}

	allergy.ID = uuid.New().String()
	createdAllergy, err := uc.emrrepo.CreateAllergy(allergy)
	if err != nil {
		return nil, errors.New("failed to create allergy: " + err.Error())
	}
	return createdAllergy, nil
}

func (uc *EmrUseCaseImpl) GetAllergyByName(allergyName string) (*entities.Allergy, error) {
	allergy, err := uc.emrrepo.GetAllergyByName(allergyName)
	if err != nil {
		return nil, errors.New("allergy not found: " + err.Error())
	}
	return allergy, nil
}

func (uc *EmrUseCaseImpl) GetAllAllergies(userID string) ([]*entities.Allergy, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	allergies, err := uc.emrrepo.GetAllAllergies()
	if err != nil {
		return nil, errors.New("failed to get all allergies: " + err.Error())
	}
	return allergies, nil
}

func (uc *EmrUseCaseImpl) CreateDrugAllergy(drugAllergy *entities.DrugAllergy) (*entities.DrugAllergy, error) {
	drugAllergyExists, err := uc.emrrepo.DrugAllergyExists(drugAllergy.AllergyName)
	if err != nil {
		return nil, errors.New("failed to verify drug allergy existence: " + err.Error())
	}
	if drugAllergyExists {
		return nil, errors.New("drug allergy already exists")
	}

	drugAllergy.ID = uuid.New().String()
	createdDrugAllergy, err := uc.emrrepo.CreateDrugAllergy(drugAllergy)
	if err != nil {
		return nil, errors.New("failed to create drug allergy: " + err.Error())
	}
	return createdDrugAllergy, nil
}

func (uc *EmrUseCaseImpl) GetDrugAllergyByName(allergyName string) (*entities.DrugAllergy, error) {
	drugAllergy, err := uc.emrrepo.GetDrugAllergyByName(allergyName)
	if err != nil {
		return nil, errors.New("drug allergy not found: " + err.Error())
	}
	return drugAllergy, nil
}

func (uc *EmrUseCaseImpl) GetAllDrugAllergies(userID string) ([]*entities.DrugAllergy, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	drugAllergies, err := uc.emrrepo.GetAllDrugAllergies()
	if err != nil {
		return nil, errors.New("failed to get all drug allergies: " + err.Error())
	}
	return drugAllergies, nil
}

// Helper function
func (uc *EmrUseCaseImpl) getOrCreateLabelID(labelName string) (string, error) {
	labelExists, err := uc.emrrepo.LabelExists(labelName)
	if err != nil {
		return "", err
	}

	if labelExists {
		intakeLabel, err := uc.emrrepo.GetIntakeLabelByName(labelName)
		return intakeLabel.ID, err
	}

	newLabel, err := uc.emrrepo.CreateIntakeLabel(&entities.IntakeLabels{
		ID:        uuid.New().String(),
		LabelName: labelName,
	})
	return newLabel.ID, err
}

func (uc *EmrUseCaseImpl) getOrCreateAllergyID(allergyName string) (string, error) {
	allergyExists, err := uc.emrrepo.AllergyExists(allergyName)
	if err != nil {
		return "", err
	}

	if allergyExists {
		allergy, err := uc.emrrepo.GetAllergyByName(allergyName)
		return allergy.ID, err
	}

	newAllergy, err := uc.emrrepo.CreateAllergy(&entities.Allergy{
		ID:          uuid.New().String(),
		AllergyName: allergyName,
	})
	if err != nil {
		return "", err
	}

	return newAllergy.ID, nil
}

func (uc *EmrUseCaseImpl) getOrCreateDrugAllergyID(allergyName string) (string, error) {
	drugAllergyExists, err := uc.emrrepo.DrugAllergyExists(allergyName)
	if err != nil {
		return "", err
	}

	if drugAllergyExists {
		drugAllergy, err := uc.emrrepo.GetDrugAllergyByName(allergyName)
		return drugAllergy.ID, err
	}

	newDrugAllergy, err := uc.emrrepo.CreateDrugAllergy(&entities.DrugAllergy{
		ID:          uuid.New().String(),
		AllergyName: allergyName,
	})
	if err != nil {
		return "", err
	}

	return newDrugAllergy.ID, nil
}

func (uc *EmrUseCaseImpl) CreateIntakeLabelByResidentID(residentID string, labels []models.IntakeLabelRequest, userID string) ([]*entities.ResidentLabels, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	resident, err := uc.emrrepo.GetResidentByID(residentID)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}

	if len(labels) == 0 {
		return nil, errors.New("labels cannot be empty")
	}

	for _, label := range labels {

		if len(strings.TrimSpace(label.LabelName)) == 0 {
			return nil, errors.New("label name cannot be empty or whitespace")
		}

		labelID, err := uc.getOrCreateLabelID(label.LabelName)
		if err != nil {
			return nil, errors.New("failed to get or create label: " + err.Error())
		}

		residentLabelExists, err := uc.emrrepo.ResidentLabelExists(resident.ID, labelID)
		if err != nil {
			return nil, errors.New("failed to verify resident label existence: " + err.Error())
		}
		if residentLabelExists {
			continue
		}

		residentLabel := &entities.ResidentLabels{
			ResidentID: resident.ID,
			LabelID:    labelID,
			NoteText:   label.NoteText, // Optional note nullable
			NotedAt:    time.Now(),
		}
		createdLabel, err := uc.emrrepo.CreateIntakeLabelByResidentID(residentLabel)
		if err != nil {
			return nil, errors.New("failed to create resident label: " + err.Error())
		}

		// Create audit log for each label
		newLabelData, _ := json.Marshal(map[string]interface{}{
			"resident_id": createdLabel.ResidentID,
			"label_id":    createdLabel.LabelID,
			"label_name":  label.LabelName,
			"note_text":   label.NoteText,
			"noted_at":    createdLabel.NotedAt,
		})
		auditLog := &entities.AuditLogs{
			ID:        uuid.New().String(),
			TableName: "resident_labels",
			RecordID:  createdLabel.ResidentID + "-" + createdLabel.LabelID,
			UserID:    userID,
			Action:    audit_constants.AuditActionInsert,
			OldValue:  "",
			NewValue:  string(newLabelData),
		}
		_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
		if err != nil {
			log.Printf("[ERROR] Failed to create audit log for resident label %s-%s: %v", createdLabel.ResidentID, createdLabel.LabelID, err)
		}
	}

	residentLabels, err := uc.emrrepo.GetResidentLabelsByResidentID(resident.ID)
	if err != nil {
		return nil, errors.New("failed to get resident labels: " + err.Error())
	}

	return residentLabels, nil
}

func (uc *EmrUseCaseImpl) GetResidentLabelsByResidentID(residentID string, userID string) ([]*entities.ResidentLabels, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentLabels, err := uc.emrrepo.GetResidentLabelsByResidentID(residentID)
	if err != nil {
		return nil, errors.New("failed to get resident labels: " + err.Error())
	}
	return residentLabels, nil
}

func (uc *EmrUseCaseImpl) CreateAllergyByResidentID(residentID string, allergies []models.AllergyRequest, userID string) ([]*entities.ResidentAllergies, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	resident, err := uc.emrrepo.GetResidentByID(residentID)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}

	if len(allergies) == 0 {
		return nil, errors.New("allergies cannot be empty")
	}

	for _, allergy := range allergies {
		if len(strings.TrimSpace(allergy.AllergyName)) == 0 {
			return nil, errors.New("allergy name cannot be empty or whitespace")
		}

		allergyID, err := uc.getOrCreateAllergyID(allergy.AllergyName)
		if err != nil {
			return nil, errors.New("failed to get or create allergy: " + err.Error())
		}

		residentAllergyExists, err := uc.emrrepo.ResidentAllergyExists(resident.ID, allergyID)
		if err != nil {
			return nil, errors.New("failed to verify resident allergy existence: " + err.Error())
		}
		if residentAllergyExists {
			continue
		}

		residentAllergy := &entities.ResidentAllergies{
			ResidentID: resident.ID,
			AllergyID:  allergyID,
			NoteText:   allergy.NoteText,
			NotedAt:    time.Now(),
		}
		createdAllergy, err := uc.emrrepo.CreateAllergyByResidentID(residentAllergy)
		if err != nil {
			return nil, errors.New("failed to create resident allergy: " + err.Error())
		}

		newAllergyData, _ := json.Marshal(map[string]interface{}{
			"resident_id":  createdAllergy.ResidentID,
			"allergy_id":   createdAllergy.AllergyID,
			"allergy_name": allergy.AllergyName,
			"note_text":    allergy.NoteText,
			"noted_at":     createdAllergy.NotedAt,
		})
		auditLog := &entities.AuditLogs{
			ID:        uuid.New().String(),
			TableName: "resident_allergies",
			RecordID:  createdAllergy.ResidentID + "-" + createdAllergy.AllergyID,
			UserID:    userID,
			Action:    audit_constants.AuditActionInsert,
			OldValue:  "",
			NewValue:  string(newAllergyData),
		}
		_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
		if err != nil {
			log.Printf("[ERROR] Failed to create audit log for resident allergy %s-%s: %v", createdAllergy.ResidentID, createdAllergy.AllergyID, err)
		}
	}

	residentAllergies, err := uc.emrrepo.GetResidentAllergiesByResidentID(resident.ID)
	if err != nil {
		return nil, errors.New("failed to get resident allergies: " + err.Error())
	}

	return residentAllergies, nil
}

func (uc *EmrUseCaseImpl) GetResidentAllergiesByResidentID(residentID string, userID string) ([]*entities.ResidentAllergies, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentAllergies, err := uc.emrrepo.GetResidentAllergiesByResidentID(residentID)
	if err != nil {
		return nil, errors.New("failed to get resident allergies: " + err.Error())
	}
	return residentAllergies, nil
}

func (uc *EmrUseCaseImpl) GetAllResidentAllergies(userID string) ([]*models.ResidentAllergyListResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentAllergies, err := uc.emrrepo.GetAllResidentAllergies()
	if err != nil {
		return nil, errors.New("failed to get all resident allergies: " + err.Error())
	}
	return residentAllergies, nil
}

func (uc *EmrUseCaseImpl) CreateDrugAllergyByResidentID(residentID string, drugAllergies []models.DrugAllergyRequest, userID string) ([]*entities.ResidentDA, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	resident, err := uc.emrrepo.GetResidentByID(residentID)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}

	if len(drugAllergies) == 0 {
		return nil, errors.New("drug allergies cannot be empty")
	}

	for _, drugAllergy := range drugAllergies {
		if len(strings.TrimSpace(drugAllergy.AllergyName)) == 0 {
			return nil, errors.New("drug allergy name cannot be empty or whitespace")
		}

		drugAllergyID, err := uc.getOrCreateDrugAllergyID(drugAllergy.AllergyName)
		if err != nil {
			return nil, errors.New("failed to get or create drug allergy: " + err.Error())
		}

		residentDrugAllergyExists, err := uc.emrrepo.ResidentDrugAllergyExists(resident.ID, drugAllergyID)
		if err != nil {
			return nil, errors.New("failed to verify resident drug allergy existence: " + err.Error())
		}
		if residentDrugAllergyExists {
			continue
		}

		residentDA := &entities.ResidentDA{
			ResidentID:    resident.ID,
			DrugAllergyID: drugAllergyID,
			NoteText:      drugAllergy.NoteText,
			NotedAt:       time.Now(),
		}
		createdDrugAllergy, err := uc.emrrepo.CreateDrugAllergyByResidentID(residentDA)
		if err != nil {
			return nil, errors.New("failed to create resident drug allergy: " + err.Error())
		}

		newDrugAllergyData, _ := json.Marshal(map[string]interface{}{
			"resident_id":     createdDrugAllergy.ResidentID,
			"drug_allergy_id": createdDrugAllergy.DrugAllergyID,
			"allergy_name":    drugAllergy.AllergyName,
			"note_text":       drugAllergy.NoteText,
			"noted_at":        createdDrugAllergy.NotedAt,
		})
		auditLog := &entities.AuditLogs{
			ID:        uuid.New().String(),
			TableName: "resident_das",
			RecordID:  createdDrugAllergy.ResidentID + "-" + createdDrugAllergy.DrugAllergyID,
			UserID:    userID,
			Action:    audit_constants.AuditActionInsert,
			OldValue:  "",
			NewValue:  string(newDrugAllergyData),
		}
		_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
		if err != nil {
			log.Printf("[ERROR] Failed to create audit log for resident drug allergy %s-%s: %v", createdDrugAllergy.ResidentID, createdDrugAllergy.DrugAllergyID, err)
		}
	}

	residentDrugAllergies, err := uc.emrrepo.GetResidentDrugAllergiesByResidentID(resident.ID)
	if err != nil {
		return nil, errors.New("failed to get resident drug allergies: " + err.Error())
	}

	return residentDrugAllergies, nil
}

func (uc *EmrUseCaseImpl) GetResidentDrugAllergiesByResidentID(residentID string, userID string) ([]*entities.ResidentDA, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentDrugAllergies, err := uc.emrrepo.GetResidentDrugAllergiesByResidentID(residentID)
	if err != nil {
		return nil, errors.New("failed to get resident drug allergies: " + err.Error())
	}
	return residentDrugAllergies, nil
}

func (uc *EmrUseCaseImpl) GetAllResidentDrugAllergies(userID string) ([]*models.ResidentDrugAllergyListResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentDrugAllergies, err := uc.emrrepo.GetAllResidentDrugAllergies()
	if err != nil {
		return nil, errors.New("failed to get all resident drug allergies: " + err.Error())
	}
	return residentDrugAllergies, nil
}

func parseVitalSignDateInput(value string) (time.Time, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return time.Time{}, errors.New("date is required")
	}

	layouts := []string{"2006-01-02", "02-01-2006", "02/01/2006"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, v)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, errors.New("date must be in YYYY-MM-DD, DD-MM-YYYY, or DD/MM/YYYY format")
}

func normalizeVitalSignTimeOfDay(value string) (string, error) {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return "", errors.New("time_of_day is required")
	}

	allowed := map[string]string{
		"เช้า":         "เช้า",
		"morning":      "เช้า",
		"สาย":          "สาย",
		"late_morning": "สาย",
		"midmorning":   "สาย",
		"noon":         "สาย",
		"บ่าย":         "บ่าย",
		"afternoon":    "บ่าย",
		"เย็น":         "เย็น",
		"evening":      "เย็น",
		"กลางคืน":      "กลางคืน",
		"night":        "กลางคืน",
	}

	if normalized, ok := allowed[v]; ok {
		return normalized, nil
	}

	return "", errors.New("time_of_day must be one of: เช้า, สาย, บ่าย, เย็น, กลางคืน")
}

func (uc *EmrUseCaseImpl) CreateVitalSign(vitalSign *entities.VitalSign, dateInput string, userID string) (*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can create vital signs")
	}

	staff, err := uc.userrepo.GetStaffByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get staff ID: " + err.Error())
	}

	residentExists, err := uc.emrrepo.ResidentExists(vitalSign.ResidentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident does not exist or missing resident ID")
	}

	measurementDate, err := parseVitalSignDateInput(dateInput)
	if err != nil {
		return nil, err
	}
	normalizedTimeOfDay, err := normalizeVitalSignTimeOfDay(vitalSign.TimeOfDay)
	if err != nil {
		return nil, err
	}
	vitalSign.MeasurementDate = measurementDate
	vitalSign.TimeOfDay = normalizedTimeOfDay

	slotExists, err := uc.emrrepo.VitalSignSlotExists(vitalSign.ResidentID, measurementDate, normalizedTimeOfDay)
	if err != nil {
		return nil, errors.New("failed to validate existing vital sign slot: " + err.Error())
	}
	if slotExists {
		return nil, errors.New("vital sign already exists for this resident, date, and time_of_day")
	}

	if vitalSign.Temperature == nil &&
		vitalSign.HeartRate == nil &&
		vitalSign.BreathingRate == nil &&
		vitalSign.BloodPressureSystolic == nil &&
		vitalSign.BloodPressureDiastolic == nil &&
		vitalSign.OxygenSaturation == nil {
		return nil, errors.New("at least one vital sign must be provided")
	}

	if (vitalSign.BloodPressureSystolic != nil && vitalSign.BloodPressureDiastolic == nil) ||
		(vitalSign.BloodPressureSystolic == nil && vitalSign.BloodPressureDiastolic != nil) {
		return nil, errors.New("blood pressure systolic and diastolic must be provided together")
	}

	if vitalSign.Temperature != nil && (*vitalSign.Temperature < emr_constants.MinTemperature || *vitalSign.Temperature > emr_constants.MaxTemperature) {
		return nil, errors.New("temperature must be between 30 and 45 Celsius")
	}
	if vitalSign.HeartRate != nil && (*vitalSign.HeartRate < emr_constants.MinHeartRate || *vitalSign.HeartRate > emr_constants.MaxHeartRate) {
		return nil, errors.New("heart rate must be between 20 and 250 bpm")
	}
	if vitalSign.BreathingRate != nil && (*vitalSign.BreathingRate < emr_constants.MinBreathingRate || *vitalSign.BreathingRate > emr_constants.MaxBreathingRate) {
		return nil, errors.New("breathing rate must be between 5 and 60 breaths per minute")
	}

	if vitalSign.BloodPressureSystolic != nil && (*vitalSign.BloodPressureSystolic < emr_constants.MinBloodPressureSystolic || *vitalSign.BloodPressureSystolic > emr_constants.MaxBloodPressureSystolic) {
		return nil, errors.New("blood pressure systolic must be between 50 and 300 mmHg")
	}
	if vitalSign.BloodPressureDiastolic != nil && (*vitalSign.BloodPressureDiastolic < emr_constants.MinBloodPressureDiastolic || *vitalSign.BloodPressureDiastolic > emr_constants.MaxBloodPressureDiastolic) {
		return nil, errors.New("blood pressure diastolic must be between 30 and 200 mmHg")
	}

	if vitalSign.BloodPressureSystolic != nil && vitalSign.BloodPressureDiastolic != nil && *vitalSign.BloodPressureSystolic < *vitalSign.BloodPressureDiastolic {
		return nil, errors.New("blood pressure systolic cannot be less than diastolic")
	}

	if vitalSign.OxygenSaturation != nil && (*vitalSign.OxygenSaturation < emr_constants.MinOxygenSaturation || *vitalSign.OxygenSaturation > emr_constants.MaxOxygenSaturation) {
		return nil, errors.New("oxygen saturation must be between 50 and 100 percent")
	}

	vitalSign.ID = uuid.New().String()
	vitalSign.CreatedByStaffID = staff.ID

	createdVitalSign, err := uc.emrrepo.CreateVitalSign(vitalSign)
	if err != nil {
		return nil, errors.New("failed to create vital sign: " + err.Error())
	}

	newVitalSignData, _ := json.Marshal(map[string]interface{}{
		"resident_id":              createdVitalSign.ResidentID,
		"measurement_date":         createdVitalSign.MeasurementDate,
		"time_of_day":              createdVitalSign.TimeOfDay,
		"temperature":              createdVitalSign.Temperature,
		"heart_rate":               createdVitalSign.HeartRate,
		"breathing_rate":           createdVitalSign.BreathingRate,
		"blood_pressure_systolic":  createdVitalSign.BloodPressureSystolic,
		"blood_pressure_diastolic": createdVitalSign.BloodPressureDiastolic,
		"oxygen_saturation":        createdVitalSign.OxygenSaturation,
		"created_by_staff_id":      createdVitalSign.CreatedByStaffID,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "vital_signs",
		RecordID:  createdVitalSign.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newVitalSignData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for new vital sign %s: %v", createdVitalSign.ID, err)
	}

	return createdVitalSign, nil
}

// filterVitalSignsByStatus filters vital signs by status: "all", "normal", or "abnormal".
// Default (empty string or "all") returns all.
func buildVitalSignFieldStatuses(vs *entities.VitalSign) ([]string, []string, []models.VitalSignFieldStatus) {
	fieldStatuses := make([]models.VitalSignFieldStatus, 0, 6)
	normalList := make([]string, 0, 6)
	abnormalList := make([]string, 0, 6)

	appendStatus := func(key, label string, isAbnormal bool) {
		fieldStatuses = append(fieldStatuses, models.VitalSignFieldStatus{Key: key, Label: label, IsAbnormal: isAbnormal})
		if isAbnormal {
			abnormalList = append(abnormalList, label)
			return
		}
		normalList = append(normalList, label)
	}

	if vs.Temperature != nil {
		appendStatus("temperature", "อุณหภูมิ", *vs.Temperature < emr_constants.NormalTempLow || *vs.Temperature > emr_constants.NormalTempHigh)
	}
	if vs.HeartRate != nil {
		appendStatus("heart_rate", "ชีพจร", *vs.HeartRate < emr_constants.NormalHeartRateLow || *vs.HeartRate > emr_constants.NormalHeartRateHigh)
	}
	if vs.BreathingRate != nil {
		appendStatus("breathing_rate", "อัตราการหายใจ", *vs.BreathingRate < emr_constants.NormalBreathingRateLow || *vs.BreathingRate > emr_constants.NormalBreathingRateHigh)
	}
	if vs.BloodPressureSystolic != nil {
		appendStatus("blood_pressure_systolic", "ความดันตัวบน", *vs.BloodPressureSystolic < emr_constants.NormalSystolicLow || *vs.BloodPressureSystolic > emr_constants.NormalSystolicHigh)
	}
	if vs.BloodPressureDiastolic != nil {
		appendStatus("blood_pressure_diastolic", "ความดันตัวล่าง", *vs.BloodPressureDiastolic < emr_constants.NormalDiastolicLow || *vs.BloodPressureDiastolic > emr_constants.NormalDiastolicHigh)
	}
	if vs.OxygenSaturation != nil {
		appendStatus("oxygen_saturation", "ออกซิเจนในเลือด", *vs.OxygenSaturation < emr_constants.NormalOxygenSaturationLow)
	}

	return normalList, abnormalList, fieldStatuses
}

func (uc *EmrUseCaseImpl) GetVitalSignsOverview(req models.VitalSignQueryParams, userID string) (*models.VitalSignsOverviewResponse, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view vital signs overview")
	}

	if req.VitalSignStatus != "" && req.VitalSignStatus != "all" && req.VitalSignStatus != "normal" && req.VitalSignStatus != "abnormal" {
		return nil, errors.New("vitalsign_status must be 'all', 'normal', or 'abnormal'")
	}

	if req.Date == nil || strings.TrimSpace(*req.Date) == "" {
		return nil, errors.New("date is required")
	}

	selectedDate, err := parseVitalSignDateInput(*req.Date)
	if err != nil {
		return nil, err
	}

	var normalizedTimeOfDay *string
	if req.TimeOfDay != nil && strings.TrimSpace(*req.TimeOfDay) != "" {
		normalized, normalizeErr := normalizeVitalSignTimeOfDay(*req.TimeOfDay)
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		normalizedTimeOfDay = &normalized
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
	} else if req.Limit > 0 {
		pageSize = req.Limit
	}
	if pageSize > 100 {
		pageSize = 100
	}

	hasPagination := req.Page != nil || req.PageSize != nil || req.Limit > 0 || req.Offset > 0
	if req.Page == nil && req.Offset > 0 {
		page = (req.Offset / pageSize) + 1
	}

	limit := pageSize
	offset := (page - 1) * pageSize

	// กรณีธรรมดา: ทั้งหมด วันนี้
	var vitalSigns []*entities.VitalSign
	var total int64

	if req.Floor == nil && len(req.LabelIDs) == 0 && !hasPagination && normalizedTimeOfDay == nil {
		vitalSigns, err = uc.emrrepo.GetVitalSignsOnDate(selectedDate, false)
		total = int64(len(vitalSigns))
	} else {
		// กรณีมี filter: ใช้ Custom
		selectedDateStr := selectedDate.Format("2006-01-02")
		params := models.VitalSignQueryParams{
			Date:      &selectedDateStr,
			TimeOfDay: normalizedTimeOfDay,
			Floor:     req.Floor,
			LabelIDs:  req.LabelIDs,
			IsLatest:  false,
			Limit:     limit,
			Offset:    offset,
		}
		vitalSigns, total, err = uc.emrrepo.GetVitalSignsCustom(params)
	}
	if err != nil {
		return nil, err
	}
	overviewItems := make([]*models.VitalSignsOverviewItemResponse, 0, len(vitalSigns))
	for _, vitalSign := range vitalSigns {
		normalList, abnormalList, fieldStatuses := buildVitalSignFieldStatuses(vitalSign)
		overviewItems = append(overviewItems, &models.VitalSignsOverviewItemResponse{
			VitalSign:     vitalSign,
			NormalList:    normalList,
			AbnormalList:  abnormalList,
			FieldStatuses: fieldStatuses,
		})
	}

	if req.VitalSignStatus == "normal" {
		filteredItems := make([]*models.VitalSignsOverviewItemResponse, 0, len(overviewItems))
		for _, item := range overviewItems {
			if len(item.AbnormalList) == 0 {
				filteredItems = append(filteredItems, item)
			}
		}
		overviewItems = filteredItems
	} else if req.VitalSignStatus == "abnormal" {
		filteredItems := make([]*models.VitalSignsOverviewItemResponse, 0, len(overviewItems))
		for _, item := range overviewItems {
			if len(item.AbnormalList) > 0 {
				filteredItems = append(filteredItems, item)
			}
		}
		overviewItems = filteredItems
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &models.VitalSignsOverviewResponse{
		Items: overviewItems,
		Pagination: models.OverviewPagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (uc *EmrUseCaseImpl) GetVitalSignsByResident(residentID string, dateInput string, isLatest string, userID string) ([]*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view vital signs by resident")
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	selectedDate, err := parseVitalSignDateInput(dateInput)
	if err != nil {
		return nil, err
	}

	isLatestBool, err := strconv.ParseBool(isLatest)
	if err != nil {
		return nil, errors.New("invalid isLatest parameter you must provide a boolean value: " + err.Error())
	}

	vitalSigns, err := uc.emrrepo.GetVitalSignsByResidentIDOnDate(residentID, selectedDate, isLatestBool)
	if err != nil {
		return nil, errors.New("failed to get vital signs: " + err.Error())
	}
	return vitalSigns, nil
}

func (uc *EmrUseCaseImpl) GetRoomVitalSigns(roomID string, isLatest string, userID string) ([]*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view vital signs by room")
	}

	roomExists, err := uc.emrrepo.RoomExists(roomID)
	if err != nil {
		return nil, errors.New("failed to verify room existence: " + err.Error())
	}
	if !roomExists {
		return nil, errors.New("room not found")
	}

	isLatestBool, err := strconv.ParseBool(isLatest)
	if err != nil {
		return nil, errors.New("invalid isLatest parameter you must provide a boolean value: " + err.Error())
	}

	vitalSigns, err := uc.emrrepo.GetVitalSignsByRoomIDToday(roomID, isLatestBool)
	if err != nil {
		return nil, errors.New("failed to get vital signs: " + err.Error())
	}
	return vitalSigns, nil
}

func (uc *EmrUseCaseImpl) GetVitalSignsHistory(residentID string, userID string) ([]*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view vital signs history")
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	vitalSigns, err := uc.emrrepo.GetVitalSignsHistory(residentID)
	if err != nil {
		return nil, errors.New("failed to get vital signs history: " + err.Error())
	}
	return vitalSigns, nil
}

func (uc *EmrUseCaseImpl) UpdateVitalSignByID(vitalSignID string, vitalSign *entities.VitalSign, userID string) (*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can update vital signs")
	}

	existingVitalSign, err := uc.emrrepo.GetVitalSignByID(vitalSignID)
	if err != nil {
		return nil, errors.New("failed to get existing vital sign: " + err.Error())
	}

	oldVitalSignData, _ := json.Marshal(map[string]interface{}{
		"resident_id":              existingVitalSign.ResidentID,
		"temperature":              existingVitalSign.Temperature,
		"heart_rate":               existingVitalSign.HeartRate,
		"breathing_rate":           existingVitalSign.BreathingRate,
		"blood_pressure_systolic":  existingVitalSign.BloodPressureSystolic,
		"blood_pressure_diastolic": existingVitalSign.BloodPressureDiastolic,
		"oxygen_saturation":        existingVitalSign.OxygenSaturation,
	})

	if vitalSign.Temperature == nil &&
		vitalSign.HeartRate == nil &&
		vitalSign.BreathingRate == nil &&
		vitalSign.BloodPressureSystolic == nil &&
		vitalSign.BloodPressureDiastolic == nil &&
		vitalSign.OxygenSaturation == nil {
		return nil, errors.New("at least one vital sign must be provided")
	}

	if (vitalSign.BloodPressureSystolic != nil && vitalSign.BloodPressureDiastolic == nil) ||
		(vitalSign.BloodPressureSystolic == nil && vitalSign.BloodPressureDiastolic != nil) {
		return nil, errors.New("blood pressure systolic and diastolic must be provided together")
	}

	if vitalSign.Temperature != nil && (*vitalSign.Temperature < emr_constants.MinTemperature || *vitalSign.Temperature > emr_constants.MaxTemperature) {
		return nil, errors.New("temperature must be between 30 and 45 Celsius")
	}
	if vitalSign.HeartRate != nil && (*vitalSign.HeartRate < emr_constants.MinHeartRate || *vitalSign.HeartRate > emr_constants.MaxHeartRate) {
		return nil, errors.New("heart rate must be between 20 and 250 bpm")
	}
	if vitalSign.BreathingRate != nil && (*vitalSign.BreathingRate < emr_constants.MinBreathingRate || *vitalSign.BreathingRate > emr_constants.MaxBreathingRate) {
		return nil, errors.New("breathing rate must be between 5 and 60 breaths per minute")
	}

	if vitalSign.BloodPressureSystolic != nil && (*vitalSign.BloodPressureSystolic < emr_constants.MinBloodPressureSystolic || *vitalSign.BloodPressureSystolic > emr_constants.MaxBloodPressureSystolic) {
		return nil, errors.New("blood pressure systolic must be between 50 and 300 mmHg")
	}
	if vitalSign.BloodPressureDiastolic != nil && (*vitalSign.BloodPressureDiastolic < emr_constants.MinBloodPressureDiastolic || *vitalSign.BloodPressureDiastolic > emr_constants.MaxBloodPressureDiastolic) {
		return nil, errors.New("blood pressure diastolic must be between 30 and 200 mmHg")
	}

	if vitalSign.BloodPressureSystolic != nil && vitalSign.BloodPressureDiastolic != nil && *vitalSign.BloodPressureSystolic < *vitalSign.BloodPressureDiastolic {
		return nil, errors.New("blood pressure systolic cannot be less than diastolic")
	}

	if vitalSign.OxygenSaturation != nil && (*vitalSign.OxygenSaturation < emr_constants.MinOxygenSaturation || *vitalSign.OxygenSaturation > emr_constants.MaxOxygenSaturation) {
		return nil, errors.New("oxygen saturation must be between 50 and 100 percent")
	}

	if vitalSign.Temperature != nil {
		existingVitalSign.Temperature = vitalSign.Temperature
	}
	if vitalSign.HeartRate != nil {
		existingVitalSign.HeartRate = vitalSign.HeartRate
	}
	if vitalSign.BreathingRate != nil {
		existingVitalSign.BreathingRate = vitalSign.BreathingRate
	}
	if vitalSign.BloodPressureSystolic != nil {
		existingVitalSign.BloodPressureSystolic = vitalSign.BloodPressureSystolic
	}
	if vitalSign.BloodPressureDiastolic != nil {
		existingVitalSign.BloodPressureDiastolic = vitalSign.BloodPressureDiastolic
	}
	if vitalSign.OxygenSaturation != nil {
		existingVitalSign.OxygenSaturation = vitalSign.OxygenSaturation
	}

	updatedVitalSign, err := uc.emrrepo.UpdateVitalSignByID(existingVitalSign)
	if err != nil {
		return nil, errors.New("failed to update vital sign: " + err.Error())
	}

	newVitalSignData, _ := json.Marshal(map[string]interface{}{
		"resident_id":              updatedVitalSign.ResidentID,
		"temperature":              updatedVitalSign.Temperature,
		"heart_rate":               updatedVitalSign.HeartRate,
		"breathing_rate":           updatedVitalSign.BreathingRate,
		"blood_pressure_systolic":  updatedVitalSign.BloodPressureSystolic,
		"blood_pressure_diastolic": updatedVitalSign.BloodPressureDiastolic,
		"oxygen_saturation":        updatedVitalSign.OxygenSaturation,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "vital_signs",
		RecordID:  updatedVitalSign.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldVitalSignData),
		NewValue:  string(newVitalSignData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for vital sign %s: %v", updatedVitalSign.ID, err)
	}

	return updatedVitalSign, nil
}

func (uc *EmrUseCaseImpl) CreateLaboratoryValue(laboratoryValue *entities.LaboratoryValue, dateInput string, timeOfDayInput string, userID string) (*entities.LaboratoryValue, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can create laboratory values")
	}

	staff, err := uc.userrepo.GetStaffByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get staff ID: " + err.Error())
	}

	residentExists, err := uc.emrrepo.ResidentExists(laboratoryValue.ResidentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident does not exist or missing resident ID")
	}

	// Parse measurement date
	measurementDate, err := parseLaboratoryValueDateInput(dateInput)
	if err != nil {
		return nil, errors.New("invalid date format: " + err.Error())
	}
	laboratoryValue.MeasurementDate = measurementDate

	// Normalize time of day
	normalizedTimeOfDay, err := normalizeLaboratoryValueTimeOfDay(timeOfDayInput)
	if err != nil {
		return nil, errors.New("invalid time of day: " + err.Error())
	}
	laboratoryValue.TimeOfDay = normalizedTimeOfDay

	// Check if slot already exists
	slotExists, err := uc.emrrepo.LaboratoryValueSlotExists(laboratoryValue.ResidentID, measurementDate, normalizedTimeOfDay)
	if err != nil {
		return nil, errors.New("failed to check slot existence: " + err.Error())
	}
	if slotExists {
		return nil, errors.New("laboratory value for this resident on this date and time slot already exists")
	}

	if laboratoryValue.BloodGlucose == nil &&
		laboratoryValue.FluidIn == nil &&
		laboratoryValue.FluidOut == nil &&
		laboratoryValue.UrineOutput == nil &&
		laboratoryValue.Stool == nil &&
		laboratoryValue.DiaperChange == nil {
		return nil, errors.New("at least one laboratory value must be provided")
	}

	// UrineOutput และ UrineType ต้องมาคู่กัน
	if (laboratoryValue.UrineOutput != nil && laboratoryValue.UrineType == nil) ||
		(laboratoryValue.UrineOutput == nil && laboratoryValue.UrineType != nil) {
		return nil, errors.New("urine_output and urine_type must be provided together")
	}

	if laboratoryValue.UrineType != nil {
		urineType := *laboratoryValue.UrineType
		if urineType != emr_constants.UrineTypeML && urineType != emr_constants.UrineTypeTimes {
			return nil, errors.New("urine_type must be either 'ml' or 'times'")
		}
	}

	if laboratoryValue.BloodGlucose != nil && (*laboratoryValue.BloodGlucose < emr_constants.MinBloodGlucose || *laboratoryValue.BloodGlucose > emr_constants.MaxBloodGlucose) {
		return nil, errors.New("blood_glucose must be between 1 and 1000 mg/dL")
	}

	if laboratoryValue.FluidIn != nil && (*laboratoryValue.FluidIn < emr_constants.MinFluidIn || *laboratoryValue.FluidIn > emr_constants.MaxFluidIn) {
		return nil, errors.New("fluid_in must be between 0 and 10000 mL")
	}

	if laboratoryValue.FluidOut != nil && (*laboratoryValue.FluidOut < emr_constants.MinFluidOut || *laboratoryValue.FluidOut > emr_constants.MaxFluidOut) {
		return nil, errors.New("fluid_out must be between 0 and 10000 mL")
	}

	if laboratoryValue.UrineOutput != nil && laboratoryValue.UrineType != nil {
		if *laboratoryValue.UrineType == emr_constants.UrineTypeML {
			if *laboratoryValue.UrineOutput < emr_constants.MinUrineOutputML || *laboratoryValue.UrineOutput > emr_constants.MaxUrineOutputML {
				return nil, errors.New("urine_output (ml) must be between 0 and 5000 mL")
			}
		} else {
			if *laboratoryValue.UrineOutput < emr_constants.MinUrineOutputTimes || *laboratoryValue.UrineOutput > emr_constants.MaxUrineOutputTimes {
				return nil, errors.New("urine_output (times) must be between 0 and 50")
			}
		}
	}

	if laboratoryValue.Stool != nil && (*laboratoryValue.Stool < emr_constants.MinStool || *laboratoryValue.Stool > emr_constants.MaxStool) {
		return nil, errors.New("stool must be between 0 and 30 times")
	}

	if laboratoryValue.DiaperChange != nil && (*laboratoryValue.DiaperChange < emr_constants.MinDiaperChange || *laboratoryValue.DiaperChange > emr_constants.MaxDiaperChange) {
		return nil, errors.New("diaper_change must be between 0 and 30 times")
	}

	laboratoryValue.ID = uuid.New().String()
	laboratoryValue.CreatedByStaffID = staff.ID

	createdLaboratoryValue, err := uc.emrrepo.CreateLaboratoryValue(laboratoryValue)
	if err != nil {
		return nil, errors.New("failed to create laboratory value: " + err.Error())
	}

	newLabValueData, _ := json.Marshal(map[string]interface{}{
		"resident_id":         createdLaboratoryValue.ResidentID,
		"blood_glucose":       createdLaboratoryValue.BloodGlucose,
		"fluid_in":            createdLaboratoryValue.FluidIn,
		"fluid_out":           createdLaboratoryValue.FluidOut,
		"urine_output":        createdLaboratoryValue.UrineOutput,
		"urine_type":          createdLaboratoryValue.UrineType,
		"stool":               createdLaboratoryValue.Stool,
		"diaper_change":       createdLaboratoryValue.DiaperChange,
		"created_by_staff_id": createdLaboratoryValue.CreatedByStaffID,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "laboratory_values",
		RecordID:  createdLaboratoryValue.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newLabValueData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for new laboratory value %s: %v", createdLaboratoryValue.ID, err)
	}

	return createdLaboratoryValue, nil
}

func parseLaboratoryValueDateInput(value string) (time.Time, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return time.Time{}, errors.New("date is required")
	}

	layouts := []string{"2006-01-02", "02-01-2006", "02/01/2006"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, v)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, errors.New("date must be in YYYY-MM-DD, DD-MM-YYYY, or DD/MM/YYYY format")
}

func normalizeLaboratoryValueTimeOfDay(value string) (string, error) {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return "", errors.New("time_of_day is required")
	}

	allowed := map[string]string{
		"เช้า":         "เช้า",
		"morning":      "เช้า",
		"สาย":          "สาย",
		"late_morning": "สาย",
		"midmorning":   "สาย",
		"noon":         "สาย",
		"บ่าย":         "บ่าย",
		"afternoon":    "บ่าย",
		"เย็น":         "เย็น",
		"evening":      "เย็น",
		"กลางคืน":      "กลางคืน",
		"night":        "กลางคืน",
	}

	if normalized, ok := allowed[v]; ok {
		return normalized, nil
	}

	return "", errors.New("time_of_day must be one of: เช้า, สาย, บ่าย, เย็น, กลางคืน")
}

func buildLaboratoryValueFieldStatuses(lab *entities.LaboratoryValue) ([]string, []string, []models.LaboratoryValueFieldStatus) {
	fieldStatuses := make([]models.LaboratoryValueFieldStatus, 0, 7)
	normalList := make([]string, 0, 7)
	abnormalList := make([]string, 0, 7)

	appendStatus := func(key, label string, isAbnormal bool) {
		fieldStatuses = append(fieldStatuses, models.LaboratoryValueFieldStatus{Key: key, Label: label, IsAbnormal: isAbnormal})
		if isAbnormal {
			abnormalList = append(abnormalList, label)
			return
		}
		normalList = append(normalList, label)
	}

	if lab.BloodGlucose != nil {
		appendStatus("blood_glucose", "ระดับน้ำตาลในเลือด", *lab.BloodGlucose < emr_constants.NormalBloodGlucoseLow || *lab.BloodGlucose > emr_constants.NormalBloodGlucoseHigh)
	}
	if lab.FluidIn != nil {
		appendStatus("fluid_in", "ปริมาณน้ำเข้า", *lab.FluidIn < emr_constants.NormalFluidInLow || *lab.FluidIn > emr_constants.NormalFluidInHigh)
	}
	if lab.FluidOut != nil {
		appendStatus("fluid_out", "ปริมาณน้ำออก", *lab.FluidOut < emr_constants.NormalFluidOutLow || *lab.FluidOut > emr_constants.NormalFluidOutHigh)
	}
	if lab.UrineOutput != nil && lab.UrineType != nil {
		isAbnormal := false
		if *lab.UrineType == emr_constants.UrineTypeML {
			isAbnormal = *lab.UrineOutput < emr_constants.NormalUrineOutputMLLow || *lab.UrineOutput > emr_constants.NormalUrineOutputMLHigh
		} else {
			isAbnormal = *lab.UrineOutput < emr_constants.NormalUrineOutputTimesLow || *lab.UrineOutput > emr_constants.NormalUrineOutputTimesHigh
		}
		appendStatus("urine_output", "ปริมาณปัสสาวะ", isAbnormal)
	}
	if lab.Stool != nil {
		appendStatus("stool", "อุจจาระ", *lab.Stool > emr_constants.NormalStoolHigh)
	}
	if lab.DiaperChange != nil {
		appendStatus("diaper_change", "การเปลี่ยนผ้าอ้อม", *lab.DiaperChange > emr_constants.NormalDiaperChangeHigh)
	}

	return normalList, abnormalList, fieldStatuses
}

// filterLaboratoryValuesByStatus filters laboratory values by status: "all", "normal", or "abnormal".
// func filterLaboratoryValuesByStatus(labs []*entities.LaboratoryValue, status string) []*entities.LaboratoryValue {
// 	if status == "" || status == "all" {
// 		return labs
// 	}

// 	wantAbnormal := status == "abnormal"
// 	result := make([]*entities.LaboratoryValue, 0)
// 	for _, lab := range labs {
// 		isAbnormal := false

// 		if lab.BloodGlucose != nil && (*lab.BloodGlucose < emr_constants.NormalBloodGlucoseLow || *lab.BloodGlucose > emr_constants.NormalBloodGlucoseHigh) {
// 			isAbnormal = true
// 		}
// 		if lab.FluidIn != nil && (*lab.FluidIn < emr_constants.NormalFluidInLow || *lab.FluidIn > emr_constants.NormalFluidInHigh) {
// 			isAbnormal = true
// 		}
// 		if lab.FluidOut != nil && (*lab.FluidOut < emr_constants.NormalFluidOutLow || *lab.FluidOut > emr_constants.NormalFluidOutHigh) {
// 			isAbnormal = true
// 		}
// 		if lab.UrineOutput != nil && lab.UrineType != nil {
// 			if *lab.UrineType == emr_constants.UrineTypeML {
// 				if *lab.UrineOutput < emr_constants.NormalUrineOutputMLLow || *lab.UrineOutput > emr_constants.NormalUrineOutputMLHigh {
// 					isAbnormal = true
// 				}
// 			} else {
// 				if *lab.UrineOutput < emr_constants.NormalUrineOutputTimesLow || *lab.UrineOutput > emr_constants.NormalUrineOutputTimesHigh {
// 					isAbnormal = true
// 				}
// 			}
// 		}
// 		if lab.Stool != nil && *lab.Stool > emr_constants.NormalStoolHigh {
// 			isAbnormal = true
// 		}
// 		if lab.DiaperChange != nil && *lab.DiaperChange > emr_constants.NormalDiaperChangeHigh {
// 			isAbnormal = true
// 		}

// 		if isAbnormal == wantAbnormal {
// 			result = append(result, lab)
// 		}
// 	}
// 	return result
// }

func (uc *EmrUseCaseImpl) GetLaboratoryValuesOverview(req models.LaboratoryValueQueryParams, userID string) (*models.LaboratoryValuesOverviewResponse, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view laboratory values overview")
	}

	if req.LaboratoryValueStatus != "" && req.LaboratoryValueStatus != "all" && req.LaboratoryValueStatus != "normal" && req.LaboratoryValueStatus != "abnormal" {
		return nil, errors.New("laboratory_value_status must be 'all', 'normal', or 'abnormal'")
	}

	// Date is required
	if req.Date == nil || strings.TrimSpace(*req.Date) == "" {
		return nil, errors.New("date is required")
	}

	selectedDate, err := parseLaboratoryValueDateInput(*req.Date)
	if err != nil {
		return nil, err
	}

	var normalizedTimeOfDay *string
	if req.TimeOfDay != nil && strings.TrimSpace(*req.TimeOfDay) != "" {
		normalized, normalizeErr := normalizeLaboratoryValueTimeOfDay(*req.TimeOfDay)
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		normalizedTimeOfDay = &normalized
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
	} else if req.Limit > 0 {
		pageSize = req.Limit
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Get laboratory values for the specified date
	labs, err := uc.emrrepo.GetLaboratoryValuesOnDate(selectedDate, false)
	if err != nil {
		return nil, err
	}

	// Filter by time_of_day if provided
	if normalizedTimeOfDay != nil {
		filtered := make([]*entities.LaboratoryValue, 0)
		for _, lab := range labs {
			if strings.EqualFold(lab.TimeOfDay, *normalizedTimeOfDay) {
				filtered = append(filtered, lab)
			}
		}
		labs = filtered
	}

	// Apply floor and label filters
	if req.Floor != nil || len(req.LabelIDs) > 0 {
		// If filters are present, re-query with custom query
		params := models.LaboratoryValueQueryParams{
			Date:      req.Date,
			TimeOfDay: req.TimeOfDay,
			Floor:     req.Floor,
			LabelIDs:  req.LabelIDs,
		}
		labs, _, err = uc.emrrepo.GetLaboratoryValuesCustom(params)
		if err != nil {
			return nil, err
		}
	}

	// Build response items with field_statuses
	items := make([]*models.LaboratoryValuesOverviewItemResponse, 0, len(labs))
	for _, lab := range labs {
		normalList, abnormalList, fieldStatuses := buildLaboratoryValueFieldStatuses(lab)

		// Apply status filter
		if req.LaboratoryValueStatus == "normal" && len(abnormalList) > 0 {
			continue
		}
		if req.LaboratoryValueStatus == "abnormal" && len(abnormalList) == 0 {
			continue
		}

		items = append(items, &models.LaboratoryValuesOverviewItemResponse{
			LaboratoryValue: lab,
			NormalList:      normalList,
			AbnormalList:    abnormalList,
			FieldStatuses:   fieldStatuses,
		})
	}

	// Handle pagination
	totalItems := len(items)
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize
	if startIdx > totalItems {
		startIdx = totalItems
		endIdx = totalItems
	}
	if endIdx > totalItems {
		endIdx = totalItems
	}

	paginatedItems := make([]*models.LaboratoryValuesOverviewItemResponse, 0)
	if startIdx < totalItems {
		paginatedItems = items[startIdx:endIdx]
	}

	return &models.LaboratoryValuesOverviewResponse{
		Items: paginatedItems,
		Pagination: models.OverviewPagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}

func (uc *EmrUseCaseImpl) GetLaboratoryValuesByResident(residentID string, dateInput string, isLatest string, userID string) ([]*entities.LaboratoryValue, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view laboratory values by resident")
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	// Parse measurement date
	measurementDate, err := parseLaboratoryValueDateInput(dateInput)
	if err != nil {
		return nil, errors.New("invalid date format: " + err.Error())
	}

	isLatestBool, err := strconv.ParseBool(isLatest)
	if err != nil {
		return nil, errors.New("invalid isLatest parameter you must provide a boolean value: " + err.Error())
	}

	labs, err := uc.emrrepo.GetLaboratoryValuesByResidentIDOnDate(residentID, measurementDate, isLatestBool)
	if err != nil {
		return nil, errors.New("failed to get laboratory values: " + err.Error())
	}
	return labs, nil
}

func (uc *EmrUseCaseImpl) GetRoomLaboratoryValues(roomID string, isLatest string, userID string) ([]*entities.LaboratoryValue, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view laboratory values by room")
	}

	roomExists, err := uc.emrrepo.RoomExists(roomID)
	if err != nil {
		return nil, errors.New("failed to verify room existence: " + err.Error())
	}
	if !roomExists {
		return nil, errors.New("room not found")
	}

	isLatestBool, err := strconv.ParseBool(isLatest)
	if err != nil {
		return nil, errors.New("invalid isLatest parameter you must provide a boolean value: " + err.Error())
	}

	labs, err := uc.emrrepo.GetLaboratoryValuesByRoomIDToday(roomID, isLatestBool)
	if err != nil {
		return nil, errors.New("failed to get laboratory values: " + err.Error())
	}
	return labs, nil
}

func (uc *EmrUseCaseImpl) GetLaboratoryValuesHistory(residentID string, userID string) ([]*entities.LaboratoryValue, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view laboratory values history")
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	labs, err := uc.emrrepo.GetLaboratoryValuesHistory(residentID)
	if err != nil {
		return nil, errors.New("failed to get laboratory values history: " + err.Error())
	}
	return labs, nil
}

func (uc *EmrUseCaseImpl) GetUrineOutputSumByResidentID(residentID string, req models.LaboratoryValueQueryParams, userID string) (*models.UrineOutputSummaryByResidentResponse, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can view urine output summary")
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	req.ResidentID = &residentID

	mlResult, err := uc.emrrepo.GetUrineOutputSumGroupByResident(req, emr_constants.UrineTypeML)
	if err != nil {
		return nil, errors.New("failed to get urine ml sum: " + err.Error())
	}
	timesResult, err := uc.emrrepo.GetUrineOutputSumGroupByResident(req, emr_constants.UrineTypeTimes)
	if err != nil {
		return nil, errors.New("failed to get urine times sum: " + err.Error())
	}

	summary := &models.UrineOutputSummaryByResidentResponse{
		ResidentID: residentID,
		TotalML:    mlResult.TotalAmount,
		TotalTimes: timesResult.TotalAmount,
	}
	return summary, nil
}

func (uc *EmrUseCaseImpl) UpdateLaboratoryValueByID(laboratoryValueID string, laboratoryValue *entities.LaboratoryValue, userID string) (*entities.LaboratoryValue, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff && userRole.Name != user_constants.RoleSuperUser && userRole.Name != user_constants.RoleAdmin {
		return nil, errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can update laboratory values")
	}

	existing, err := uc.emrrepo.GetLaboratoryValueByID(laboratoryValueID)
	if err != nil {
		return nil, errors.New("laboratory value not found: " + err.Error())
	}

	oldData, _ := json.Marshal(map[string]interface{}{
		"resident_id":   existing.ResidentID,
		"blood_glucose": existing.BloodGlucose,
		"fluid_in":      existing.FluidIn,
		"fluid_out":     existing.FluidOut,
		"urine_output":  existing.UrineOutput,
		"urine_type":    existing.UrineType,
		"stool":         existing.Stool,
		"diaper_change": existing.DiaperChange,
	})

	if laboratoryValue.BloodGlucose == nil &&
		laboratoryValue.FluidIn == nil &&
		laboratoryValue.FluidOut == nil &&
		laboratoryValue.UrineOutput == nil &&
		laboratoryValue.UrineType == nil &&
		laboratoryValue.Stool == nil &&
		laboratoryValue.DiaperChange == nil {
		return nil, errors.New("at least one laboratory value must be provided")
	}

	if (laboratoryValue.UrineOutput != nil && laboratoryValue.UrineType == nil) ||
		(laboratoryValue.UrineOutput == nil && laboratoryValue.UrineType != nil) {
		return nil, errors.New("urine_output and urine_type must be provided together")
	}

	if laboratoryValue.UrineType != nil {
		if *laboratoryValue.UrineType != emr_constants.UrineTypeML && *laboratoryValue.UrineType != emr_constants.UrineTypeTimes {
			return nil, errors.New("urine_type must be either 'ml' or 'times'")
		}
	}

	if laboratoryValue.BloodGlucose != nil && (*laboratoryValue.BloodGlucose < emr_constants.MinBloodGlucose || *laboratoryValue.BloodGlucose > emr_constants.MaxBloodGlucose) {
		return nil, errors.New("blood_glucose must be between 1 and 1000 mg/dL")
	}
	if laboratoryValue.FluidIn != nil && (*laboratoryValue.FluidIn < emr_constants.MinFluidIn || *laboratoryValue.FluidIn > emr_constants.MaxFluidIn) {
		return nil, errors.New("fluid_in must be between 0 and 10000 mL")
	}
	if laboratoryValue.FluidOut != nil && (*laboratoryValue.FluidOut < emr_constants.MinFluidOut || *laboratoryValue.FluidOut > emr_constants.MaxFluidOut) {
		return nil, errors.New("fluid_out must be between 0 and 10000 mL")
	}
	if laboratoryValue.UrineOutput != nil && laboratoryValue.UrineType != nil {
		if *laboratoryValue.UrineType == emr_constants.UrineTypeML {
			if *laboratoryValue.UrineOutput < emr_constants.MinUrineOutputML || *laboratoryValue.UrineOutput > emr_constants.MaxUrineOutputML {
				return nil, errors.New("urine_output (ml) must be between 0 and 5000 mL")
			}
		} else {
			if *laboratoryValue.UrineOutput < emr_constants.MinUrineOutputTimes || *laboratoryValue.UrineOutput > emr_constants.MaxUrineOutputTimes {
				return nil, errors.New("urine_output (times) must be between 0 and 50")
			}
		}
	}
	if laboratoryValue.Stool != nil && (*laboratoryValue.Stool < emr_constants.MinStool || *laboratoryValue.Stool > emr_constants.MaxStool) {
		return nil, errors.New("stool must be between 0 and 30 times")
	}
	if laboratoryValue.DiaperChange != nil && (*laboratoryValue.DiaperChange < emr_constants.MinDiaperChange || *laboratoryValue.DiaperChange > emr_constants.MaxDiaperChange) {
		return nil, errors.New("diaper_change must be between 0 and 30 times")
	}

	if laboratoryValue.BloodGlucose != nil {
		existing.BloodGlucose = laboratoryValue.BloodGlucose
	}
	if laboratoryValue.FluidIn != nil {
		existing.FluidIn = laboratoryValue.FluidIn
	}
	if laboratoryValue.FluidOut != nil {
		existing.FluidOut = laboratoryValue.FluidOut
	}
	if laboratoryValue.UrineOutput != nil {
		existing.UrineOutput = laboratoryValue.UrineOutput
	}
	if laboratoryValue.UrineType != nil {
		existing.UrineType = laboratoryValue.UrineType
	}
	if laboratoryValue.Stool != nil {
		existing.Stool = laboratoryValue.Stool
	}
	if laboratoryValue.DiaperChange != nil {
		existing.DiaperChange = laboratoryValue.DiaperChange
	}

	updated, err := uc.emrrepo.UpdateLaboratoryValueByID(existing)
	if err != nil {
		return nil, errors.New("failed to update laboratory value: " + err.Error())
	}

	newData, _ := json.Marshal(map[string]interface{}{
		"resident_id":   updated.ResidentID,
		"blood_glucose": updated.BloodGlucose,
		"fluid_in":      updated.FluidIn,
		"fluid_out":     updated.FluidOut,
		"urine_output":  updated.UrineOutput,
		"urine_type":    updated.UrineType,
		"stool":         updated.Stool,
		"diaper_change": updated.DiaperChange,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "laboratory_values",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldData),
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for laboratory value %s: %v", updated.ID, err)
	}

	return updated, nil
}

func (uc *EmrUseCaseImpl) CreateNurseNote(note *entities.NurseNote, userID string) (*entities.NurseNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	staff, err := uc.userrepo.GetStaffByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get staff ID: " + err.Error())
	}

	residentExists, err := uc.emrrepo.ResidentExists(note.ResidentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident does not exist")
	}

	if strings.TrimSpace(note.Content) == "" {
		return nil, errors.New("content is required")
	}
	if strings.TrimSpace(note.Category) == "" {
		return nil, errors.New("category is required")
	}

	note.Priority = strings.ToLower(strings.TrimSpace(note.Priority))
	if note.Priority != "normal" && note.Priority != "urgent" {
		return nil, errors.New("priority must be either 'normal' or 'urgent'")
	}

	note.ID = uuid.New().String()
	note.CreatedByStaffID = staff.ID

	created, err := uc.emrrepo.CreateNurseNote(note)
	if err != nil {
		return nil, errors.New("failed to create nurse note: " + err.Error())
	}

	newData, _ := json.Marshal(created)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "nurse_notes",
		RecordID:  created.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for nurse note %s: %v", created.ID, err)
	}

	uc.populateNurseNoteStaffName(created)

	return created, nil
}

func (uc *EmrUseCaseImpl) GetNurseNotesOverview(userID string) ([]*entities.NurseNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	notes, err := uc.emrrepo.GetNurseNotesOverview()
	if err != nil {
		return nil, err
	}

	uc.populateNurseNotesStaffNames(notes)
	return notes, nil
}

func (uc *EmrUseCaseImpl) GetNurseNotesByResidentID(residentID string, userID string) ([]*entities.NurseNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	notes, err := uc.emrrepo.GetNurseNotesByResidentID(residentID)
	if err != nil {
		return nil, err
	}

	uc.populateNurseNotesStaffNames(notes)
	return notes, nil
}

func (uc *EmrUseCaseImpl) UpdateNurseNoteByID(noteID string, note *entities.NurseNote, userID string) (*entities.NurseNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	existing, err := uc.emrrepo.GetNurseNoteByID(noteID)
	if err != nil {
		return nil, errors.New("nurse note not found: " + err.Error())
	}

	oldData, _ := json.Marshal(existing)

	if note.Category != "" {
		existing.Category = note.Category
	}
	if note.Content != "" {
		existing.Content = note.Content
	}
	if note.Priority != "" {
		nextPriority := strings.ToLower(strings.TrimSpace(note.Priority))
		if nextPriority != "normal" && nextPriority != "urgent" {
			return nil, errors.New("priority must be either 'normal' or 'urgent'")
		}
		existing.Priority = nextPriority
	}
	if note.SendNote != existing.SendNote {
		existing.SendNote = note.SendNote
	}

	updated, err := uc.emrrepo.UpdateNurseNoteByID(existing)
	if err != nil {
		return nil, errors.New("failed to update nurse note: " + err.Error())
	}

	newData, _ := json.Marshal(updated)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "nurse_notes",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldData),
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for nurse note %s: %v", updated.ID, err)
	}

	uc.populateNurseNoteStaffName(updated)

	return updated, nil
}

func (uc *EmrUseCaseImpl) DeleteNurseNoteByID(noteID string, userID string) error {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return err
	}

	existing, err := uc.emrrepo.GetNurseNoteByID(noteID)
	if err != nil {
		return errors.New("nurse note not found: " + err.Error())
	}

	oldData, _ := json.Marshal(existing)
	if err := uc.emrrepo.DeleteNurseNoteByID(noteID); err != nil {
		return errors.New("failed to delete nurse note: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "nurse_notes",
		RecordID:  noteID,
		UserID:    userID,
		Action:    audit_constants.AuditActionDelete,
		OldValue:  string(oldData),
		NewValue:  "",
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for deleted nurse note %s: %v", noteID, err)
	}

	return nil
}

func (uc *EmrUseCaseImpl) CreateWoundCareNote(note *entities.WoundCareNote, userID string, imageFile multipart.File) (*entities.WoundCareNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	staff, err := uc.userrepo.GetStaffByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get staff ID: " + err.Error())
	}

	residentExists, err := uc.emrrepo.ResidentExists(note.ResidentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident does not exist")
	}

	if strings.TrimSpace(note.Location) == "" {
		return nil, errors.New("location is required")
	}
	if strings.TrimSpace(note.WoundType) == "" {
		return nil, errors.New("wound_type is required")
	}

	if imageFile != nil {
		imageURL, err := uc.uploadWoundCareImage(imageFile)
		if err != nil {
			return nil, err
		}
		note.ImageURL = imageURL
	}

	note.ID = uuid.New().String()
	note.CreatedByStaffID = staff.ID

	created, err := uc.emrrepo.CreateWoundCareNote(note)
	if err != nil {
		return nil, errors.New("failed to create wound care note: " + err.Error())
	}

	newData, _ := json.Marshal(created)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "wound_care_notes",
		RecordID:  created.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for wound care note %s: %v", created.ID, err)
	}

	uc.populateWoundCareNoteStaffName(created)

	return created, nil
}

func (uc *EmrUseCaseImpl) GetWoundCareNotesOverview(userID string) ([]*entities.WoundCareNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	notes, err := uc.emrrepo.GetWoundCareNotesOverview()
	if err != nil {
		return nil, err
	}

	uc.populateWoundCareNotesStaffNames(notes)
	return notes, nil
}

func (uc *EmrUseCaseImpl) GetWoundCareNotesByResidentID(residentID string, userID string) ([]*entities.WoundCareNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	notes, err := uc.emrrepo.GetWoundCareNotesByResidentID(residentID)
	if err != nil {
		return nil, err
	}

	uc.populateWoundCareNotesStaffNames(notes)
	return notes, nil
}

func (uc *EmrUseCaseImpl) UpdateWoundCareNoteByID(noteID string, note *entities.WoundCareNote, userID string, imageFile multipart.File) (*entities.WoundCareNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	existing, err := uc.emrrepo.GetWoundCareNoteByID(noteID)
	if err != nil {
		return nil, errors.New("wound care note not found: " + err.Error())
	}

	oldData, _ := json.Marshal(existing)

	if note.Location != "" {
		existing.Location = note.Location
	}
	if note.WoundType != "" {
		existing.WoundType = note.WoundType
	}
	if note.Size != nil {
		existing.Size = note.Size
	}
	if note.Treatment != nil {
		existing.Treatment = note.Treatment
	}
	if note.Supplies != nil {
		existing.Supplies = note.Supplies
	}
	if note.Status != nil {
		existing.Status = note.Status
	}
	if imageFile != nil {
		imageURL, err := uc.uploadWoundCareImage(imageFile)
		if err != nil {
			return nil, err
		}
		existing.ImageURL = imageURL
	} else if note.ImageURL != nil {
		existing.ImageURL = note.ImageURL
	}
	if note.Note != nil {
		existing.Note = note.Note
	}

	updated, err := uc.emrrepo.UpdateWoundCareNoteByID(existing)
	if err != nil {
		return nil, errors.New("failed to update wound care note: " + err.Error())
	}

	newData, _ := json.Marshal(updated)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "wound_care_notes",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldData),
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for wound care note %s: %v", updated.ID, err)
	}

	uc.populateWoundCareNoteStaffName(updated)

	return updated, nil
}

func (uc *EmrUseCaseImpl) DeleteWoundCareNoteByID(noteID string, userID string) error {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return err
	}

	existing, err := uc.emrrepo.GetWoundCareNoteByID(noteID)
	if err != nil {
		return errors.New("wound care note not found: " + err.Error())
	}

	oldData, _ := json.Marshal(existing)
	if err := uc.emrrepo.DeleteWoundCareNoteByID(noteID); err != nil {
		return errors.New("failed to delete wound care note: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "wound_care_notes",
		RecordID:  noteID,
		UserID:    userID,
		Action:    audit_constants.AuditActionDelete,
		OldValue:  string(oldData),
		NewValue:  "",
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for deleted wound care note %s: %v", noteID, err)
	}

	return nil
}

func (uc *EmrUseCaseImpl) CreateRelativeNote(note *entities.RelativeNote, userID string) (*entities.RelativeNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	staff, err := uc.userrepo.GetStaffByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get staff ID: " + err.Error())
	}

	residentExists, err := uc.emrrepo.ResidentExists(note.ResidentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident does not exist")
	}

	if strings.TrimSpace(note.Relation) == "" {
		return nil, errors.New("relation is required")
	}
	if strings.TrimSpace(note.Content) == "" {
		return nil, errors.New("content is required")
	}

	note.ID = uuid.New().String()
	note.CreatedByStaffID = staff.ID

	created, err := uc.emrrepo.CreateRelativeNote(note)
	if err != nil {
		return nil, errors.New("failed to create relative note: " + err.Error())
	}

	newData, _ := json.Marshal(created)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "relative_notes",
		RecordID:  created.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for relative note %s: %v", created.ID, err)
	}

	uc.populateRelativeNoteStaffName(created)

	return created, nil
}

func (uc *EmrUseCaseImpl) GetRelativeNotesOverview(userID string) ([]*entities.RelativeNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	notes, err := uc.emrrepo.GetRelativeNotesOverview()
	if err != nil {
		return nil, err
	}

	uc.populateRelativeNotesStaffNames(notes)
	return notes, nil
}

func (uc *EmrUseCaseImpl) GetRelativeNotesByResidentID(residentID string, userID string) ([]*entities.RelativeNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	notes, err := uc.emrrepo.GetRelativeNotesByResidentID(residentID)
	if err != nil {
		return nil, err
	}

	uc.populateRelativeNotesStaffNames(notes)
	return notes, nil
}

func (uc *EmrUseCaseImpl) UpdateRelativeNoteByID(noteID string, note *entities.RelativeNote, userID string) (*entities.RelativeNote, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	existing, err := uc.emrrepo.GetRelativeNoteByID(noteID)
	if err != nil {
		return nil, errors.New("relative note not found: " + err.Error())
	}

	oldData, _ := json.Marshal(existing)

	if note.Relation != "" {
		existing.Relation = note.Relation
	}
	if note.Content != "" {
		existing.Content = note.Content
	}
	if note.SendNote != existing.SendNote {
		existing.SendNote = note.SendNote
	}

	updated, err := uc.emrrepo.UpdateRelativeNoteByID(existing)
	if err != nil {
		return nil, errors.New("failed to update relative note: " + err.Error())
	}

	newData, _ := json.Marshal(updated)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "relative_notes",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldData),
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for relative note %s: %v", updated.ID, err)
	}

	uc.populateRelativeNoteStaffName(updated)

	return updated, nil
}

func (uc *EmrUseCaseImpl) DeleteRelativeNoteByID(noteID string, userID string) error {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return err
	}

	existing, err := uc.emrrepo.GetRelativeNoteByID(noteID)
	if err != nil {
		return errors.New("relative note not found: " + err.Error())
	}

	oldData, _ := json.Marshal(existing)
	if err := uc.emrrepo.DeleteRelativeNoteByID(noteID); err != nil {
		return errors.New("failed to delete relative note: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "relative_notes",
		RecordID:  noteID,
		UserID:    userID,
		Action:    audit_constants.AuditActionDelete,
		OldValue:  string(oldData),
		NewValue:  "",
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for deleted relative note %s: %v", noteID, err)
	}

	return nil
}

func (uc *EmrUseCaseImpl) IssueRelativeMagicLink(residentID string, userID string) (*models.RelativeMagicLinkResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	resident, err := uc.emrrepo.GetResidentByID(residentID)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}

	relative, err := uc.ensureRelativeAccountForResident(resident)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiredAt := now.Add(30 * 24 * time.Hour)
	tokenText := uuid.New().String()

	createdToken, err := uc.emrrepo.CreateRelativeMagicLinkToken(&entities.RelativeMagicLinkToken{
		ID:              uuid.New().String(),
		RelativeID:      relative.ID,
		ResidentID:      resident.ID,
		Token:           tokenText,
		ExpiresAt:       expiredAt,
		LastAccessedAt:  nil,
		CreatedByUserID: userID,
	})
	if err != nil {
		return nil, errors.New("failed to issue magic link token: " + err.Error())
	}

	return &models.RelativeMagicLinkResponse{
		ResidentID: resident.ID,
		RelativeID: relative.ID,
		Token:      createdToken.Token,
		MagicLink:  uc.buildRelativeMagicLink(resident.ID, createdToken.Token),
		ExpiresAt:  createdToken.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (uc *EmrUseCaseImpl) GetRelativeMagicLink(residentID string, userID string) (*models.RelativeMagicLinkResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	resident, err := uc.emrrepo.GetResidentByID(residentID)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}

	relative, err := uc.ensureRelativeAccountForResident(resident)
	if err != nil {
		return nil, err
	}

	token, err := uc.emrrepo.GetLatestValidRelativeMagicLinkTokenByResidentID(residentID, time.Now())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uc.IssueRelativeMagicLink(residentID, userID)
		}
		return nil, errors.New("failed to get magic link token: " + err.Error())
	}

	return &models.RelativeMagicLinkResponse{
		ResidentID: resident.ID,
		RelativeID: relative.ID,
		Token:      token.Token,
		MagicLink:  uc.buildRelativeMagicLink(resident.ID, token.Token),
		ExpiresAt:  token.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (uc *EmrUseCaseImpl) RelativePortalLogin(req models.RelativePortalLoginRequest) (*models.RelativePortalLoginResponse, error) {
	password := strings.TrimSpace(req.Password)
	if password == "" {
		return nil, errors.New("password is required")
	}
	if len(password) != 8 {
		return nil, errors.New("password must be in DDMMYYYY format")
	}
	for _, ch := range password {
		if ch < '0' || ch > '9' {
			return nil, errors.New("password must be in DDMMYYYY format")
		}
	}

	var (
		relative *entities.Relative
		resident *entities.Resident
		tokenRec *entities.RelativeMagicLinkToken
		err      error
	)

	now := time.Now()

	if strings.TrimSpace(req.Token) != "" {
		tokenRec, err = uc.emrrepo.GetRelativeMagicLinkTokenByToken(strings.TrimSpace(req.Token))
		if err != nil {
			return nil, errors.New("invalid magic link")
		}
		if !tokenRec.ExpiresAt.After(now) {
			return nil, errors.New("magic link expired")
		}

		relative, err = uc.emrrepo.GetRelativeByID(tokenRec.RelativeID)
		if err != nil {
			return nil, errors.New("relative account not found")
		}
		resident, err = uc.emrrepo.GetResidentByID(tokenRec.ResidentID)
		if err != nil {
			return nil, errors.New("resident not found")
		}
		if strings.TrimSpace(req.ResidentID) != "" && strings.TrimSpace(req.ResidentID) != resident.ID {
			return nil, errors.New("resident_id does not match magic link")
		}
	} else {
		residentID := strings.TrimSpace(req.ResidentID)
		if residentID == "" {
			return nil, errors.New("resident_id is required")
		}
		resident, err = uc.emrrepo.GetResidentByID(residentID)
		if err != nil {
			return nil, errors.New("resident not found")
		}
		relative, err = uc.emrrepo.GetRelativeByResidentID(residentID)
		if err != nil {
			return nil, errors.New("relative account not found")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(relative.RelativePassword), []byte(password)); err != nil {
		return nil, errors.New("invalid resident birthday password")
	}

	user, err := uc.userrepo.GetUserByID(relative.UserID)
	if err != nil {
		return nil, errors.New("relative user not found")
	}

	expiryDuration := 30 * time.Minute
	if req.Remember {
		expiryDuration = 48 * time.Hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     now.Add(expiryDuration).Unix(),
		"iat":     now.Unix(),
		"jti":     uuid.New().String(),
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return nil, errors.New("failed to generate token: " + err.Error())
	}

	if tokenRec != nil {
		_ = uc.emrrepo.TouchRelativeMagicLinkTokenLastAccessed(tokenRec.ID, now)
	}

	return &models.RelativePortalLoginResponse{
		Token:      tokenString,
		UserID:     user.ID,
		Username:   user.Username,
		Email:      user.Email,
		RoleName:   "relative",
		ResidentID: resident.ID,
	}, nil
}

func (uc *EmrUseCaseImpl) GetRelativeDashboard(userID string, dateInput string) (*models.RelativeDashboardResponse, error) {
	if err := uc.ensureRelative(userID); err != nil {
		return nil, err
	}

	relative, err := uc.emrrepo.GetRelativeByUserID(userID)
	if err != nil {
		return nil, errors.New("relative profile not found")
	}

	resident, err := uc.emrrepo.GetResidentByID(relative.ResidentID)
	if err != nil {
		return nil, errors.New("resident not found")
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	selectedDate := time.Now().In(loc)
	if strings.TrimSpace(dateInput) != "" {
		parsed, parseErr := time.ParseInLocation("2006-01-02", strings.TrimSpace(dateInput), loc)
		if parseErr != nil {
			return nil, errors.New("date must be in YYYY-MM-DD format")
		}
		selectedDate = parsed
	}

	notes, err := uc.emrrepo.GetRelativeNotesByResidentID(relative.ResidentID)
	if err != nil {
		return nil, errors.New("failed to get relative notes: " + err.Error())
	}

	resultNotes := make([]models.RelativeDashboardNote, 0)
	var lastUpdated *time.Time
	for _, note := range notes {
		if note == nil {
			continue
		}
		noteTime := note.CreatedAt.In(loc)
		if noteTime.Year() != selectedDate.Year() || noteTime.Month() != selectedDate.Month() || noteTime.Day() != selectedDate.Day() {
			continue
		}
		if !note.SendNote {
			continue
		}

		resultNotes = append(resultNotes, models.RelativeDashboardNote{
			ID:        note.ID,
			Content:   note.Content,
			CreatedAt: note.CreatedAt.Format(time.RFC3339),
		})

		if lastUpdated == nil || note.CreatedAt.After(*lastUpdated) {
			t := note.CreatedAt
			lastUpdated = &t
		}
	}

	var lastUpdatedText *string
	if lastUpdated != nil {
		text := lastUpdated.Format(time.RFC3339)
		lastUpdatedText = &text
	}

	return &models.RelativeDashboardResponse{
		ResidentID:    resident.ID,
		ResidentName:  strings.TrimSpace(resident.FirstName + " " + resident.LastName),
		Date:          selectedDate.Format("2006-01-02"),
		LastUpdatedAt: lastUpdatedText,
		Notes:         resultNotes,
	}, nil
}

func (uc *EmrUseCaseImpl) GetRelativePatientInfo(userID string) (*models.RelativePatientInfoResponse, error) {
	if err := uc.ensureRelative(userID); err != nil {
		return nil, err
	}

	relative, err := uc.emrrepo.GetRelativeByUserID(userID)
	if err != nil {
		return nil, errors.New("relative profile not found")
	}

	resident, err := uc.emrrepo.GetResidentByID(relative.ResidentID)
	if err != nil {
		return nil, errors.New("resident not found")
	}

	foodAllergyRows, err := uc.emrrepo.GetResidentAllergiesByResidentID(relative.ResidentID)
	if err != nil {
		return nil, errors.New("failed to get food allergies: " + err.Error())
	}
	drugAllergyRows, err := uc.emrrepo.GetResidentDrugAllergiesByResidentID(relative.ResidentID)
	if err != nil {
		return nil, errors.New("failed to get drug allergies: " + err.Error())
	}

	foodAllergies := make([]string, 0, len(foodAllergyRows))
	for _, item := range foodAllergyRows {
		if item == nil {
			continue
		}
		if item.Allergy.AllergyName != "" {
			foodAllergies = append(foodAllergies, item.Allergy.AllergyName)
		}
	}

	drugAllergies := make([]string, 0, len(drugAllergyRows))
	for _, item := range drugAllergyRows {
		if item == nil {
			continue
		}
		if item.DrugAllergy.AllergyName != "" {
			drugAllergies = append(drugAllergies, item.DrugAllergy.AllergyName)
		}
	}

	now := time.Now()
	age := now.Year() - resident.DateOfBirth.Year()
	if now.Month() < resident.DateOfBirth.Month() || (now.Month() == resident.DateOfBirth.Month() && now.Day() < resident.DateOfBirth.Day()) {
		age--
	}

	return &models.RelativePatientInfoResponse{
		ResidentID:                resident.ID,
		FirstName:                 resident.FirstName,
		LastName:                  resident.LastName,
		Nickname:                  resident.Nickname,
		Gender:                    resident.Gender,
		DateOfBirth:               resident.DateOfBirth.Format("2006-01-02"),
		Age:                       age,
		IdCardNumber:              resident.IdCardNumber,
		PurposeOfStay:             resident.PurposeOfStay,
		CheckInDate:               resident.CheckInDate.Format("2006-01-02"),
		Status:                    resident.Status,
		PreExistingConditions:     splitTextList(resident.PreExistingConditions),
		PreExistingConditionsNote: resident.PreExistingConditionsNotes,
		SurgicalHistory:           splitTextList(resident.SugicalHistory),
		FoodAllergies:             foodAllergies,
		DrugAllergies:             drugAllergies,
		EmergencyHospital:         resident.PreferredEmergencyHospital,
		EmergencyHospitalPhone:    resident.EmergencyHospitalPhone,
	}, nil
}

func (uc *EmrUseCaseImpl) CreateDoctorOrder(order *entities.DoctorOrder, userID string) (*entities.DoctorOrder, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	staff, err := uc.userrepo.GetStaffByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get staff ID: " + err.Error())
	}

	residentExists, err := uc.emrrepo.ResidentExists(order.ResidentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident does not exist")
	}

	order.Title = strings.TrimSpace(order.Title)
	if order.Title == "" {
		return nil, errors.New("title is required")
	}

	order.ID = uuid.New().String()
	order.CreatedByStaffID = staff.ID

	created, err := uc.emrrepo.CreateDoctorOrder(order)
	if err != nil {
		return nil, errors.New("failed to create doctor order: " + err.Error())
	}

	newData, _ := json.Marshal(created)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "doctor_orders",
		RecordID:  created.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for doctor order %s: %v", created.ID, err)
	}

	uc.populateDoctorOrderStaffName(created)

	return created, nil
}

func (uc *EmrUseCaseImpl) GetDoctorOrdersOverview(userID string) ([]*entities.DoctorOrder, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	orders, err := uc.emrrepo.GetDoctorOrdersOverview()
	if err != nil {
		return nil, err
	}

	uc.populateDoctorOrdersStaffNames(orders)
	return orders, nil
}

func (uc *EmrUseCaseImpl) GetDoctorOrdersByResidentID(residentID string, userID string) ([]*entities.DoctorOrder, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	orders, err := uc.emrrepo.GetDoctorOrdersByResidentID(residentID)
	if err != nil {
		return nil, err
	}

	uc.populateDoctorOrdersStaffNames(orders)
	return orders, nil
}

func (uc *EmrUseCaseImpl) UpdateDoctorOrderByID(orderID string, order *entities.DoctorOrder, userID string) (*entities.DoctorOrder, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	existing, err := uc.emrrepo.GetDoctorOrderByID(orderID)
	if err != nil {
		return nil, errors.New("doctor order not found: " + err.Error())
	}

	oldData, _ := json.Marshal(existing)

	if order.OrderDate != nil {
		existing.OrderDate = order.OrderDate
	}
	if order.OrderType != nil {
		existing.OrderType = order.OrderType
	}
	if strings.TrimSpace(order.Title) != "" {
		existing.Title = strings.TrimSpace(order.Title)
	}
	if order.Details != nil {
		existing.Details = order.Details
	}
	if order.StartDate != nil {
		existing.StartDate = order.StartDate
	}
	if order.EndDate != nil {
		existing.EndDate = order.EndDate
	}
	if order.Frequency != nil {
		existing.Frequency = order.Frequency
	}
	if order.OrderedBy != nil {
		existing.OrderedBy = order.OrderedBy
	}

	if strings.TrimSpace(existing.Title) == "" {
		return nil, errors.New("title is required")
	}

	updated, err := uc.emrrepo.UpdateDoctorOrderByID(existing)
	if err != nil {
		return nil, errors.New("failed to update doctor order: " + err.Error())
	}

	newData, _ := json.Marshal(updated)
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "doctor_orders",
		RecordID:  updated.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldData),
		NewValue:  string(newData),
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for doctor order %s: %v", updated.ID, err)
	}

	uc.populateDoctorOrderStaffName(updated)

	return updated, nil
}

func (uc *EmrUseCaseImpl) DeleteDoctorOrderByID(orderID string, userID string) error {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return err
	}

	existing, err := uc.emrrepo.GetDoctorOrderByID(orderID)
	if err != nil {
		return errors.New("doctor order not found: " + err.Error())
	}

	oldData, _ := json.Marshal(existing)
	if err := uc.emrrepo.DeleteDoctorOrderByID(orderID); err != nil {
		return errors.New("failed to delete doctor order: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "doctor_orders",
		RecordID:  orderID,
		UserID:    userID,
		Action:    audit_constants.AuditActionDelete,
		OldValue:  string(oldData),
		NewValue:  "",
	}
	_, err = uc.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for deleted doctor order %s: %v", orderID, err)
	}

	return nil
}
