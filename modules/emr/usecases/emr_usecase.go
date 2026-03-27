package usecases

import (
	"encoding/json"
	"errors"

	// "io"
	"log"
	"strconv"

	// "mime/multipart"
	// "os"
	// "sync"
	"strings"
	"time"

	audit_constants "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/constants"
	audit_repo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	emr_constants "github.com/aikidoaikido115/New-Acis-BE/modules/emr/constants"
	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"

	"github.com/google/uuid"
	// "golang.org/x/text/unicode/norm"
)

type EmrUsecase interface {

	// Resident operations
	CreateResident(resident *entities.Resident, userID string) (*entities.Resident, error)
	GetResidentByID(id string) (*entities.Resident, error)
	GetResidentByRoomID(roomID string) ([]*entities.Resident, error)
	GetAllResidents() ([]*entities.Resident, error)
	GetResidentOverview(req models.ResidentQueryParams) ([]*models.ResidentOverviewResponse, error)
	UpdateResidentByID(residentID string, data models.UpdateResidentRequest, userID string) (*entities.Resident, error)

	// Dashboard operations
	GetNumberOfResidentsDashboard() (models.NumberOfResidentsDashboardResponse, error)
	GetResidentGenderStatsDashboard() (models.ResidentGenderStatsDashboardResponse, error)
	GetResidentAllergyStatsDashboard() (models.ResidentAllergyStatsDashboardResponse, error)

	// Room operations
	GetRoomByID(id string) (*entities.Room, error)
	GetAllRooms() ([]*entities.Room, error)
	CreateRoom(room *entities.Room, userID string) (*entities.Room, error)
	UpdateRoomByID(roomID string, data models.UpdateRoomRequest, userID string) (*entities.Room, error)

	// IntakeLabel operations
	CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error)
	// GetIntakeLabelByID(id string) (*entities.IntakeLabels, error)
	GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error)
	GetAllIntakeLabels() ([]*entities.IntakeLabels, error)
	// Allergy operations
	CreateAllergy(allergy *entities.Allergy) (*entities.Allergy, error)
	GetAllergyByName(allergyName string) (*entities.Allergy, error)
	GetAllAllergies() ([]*entities.Allergy, error)

	// ResidentLabel operations (many-to-many)
	CreateIntakeLabelByResidentID(residentID string, labels []models.IntakeLabelRequest, userID string) ([]*entities.ResidentLabels, error)
	GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error)
	// ResidentAllergy operations (many-to-many)
	CreateAllergyByResidentID(residentID string, allergies []models.AllergyRequest, userID string) ([]*entities.ResidentAllergies, error)
	GetResidentAllergiesByResidentID(residentID string) ([]*entities.ResidentAllergies, error)
	GetAllResidentAllergies() ([]*models.ResidentAllergyListResponse, error)

	// VitalSign operations
	CreateVitalSign(vitalSign *entities.VitalSign, userID string) (*entities.VitalSign, error)

	GetVitalSignsOverview(req models.VitalSignQueryParams, userID string) ([]*entities.VitalSign, error)
	GetVitalSignsByResident(residentID string, isLatest string, userID string) ([]*entities.VitalSign, error)
	GetRoomVitalSigns(roomID string, isLatest string, userID string) ([]*entities.VitalSign, error)
	GetVitalSignsHistory(residentID string, userID string) ([]*entities.VitalSign, error)
	GetAbnormalVitalSigns(floor string, isLatest string, userID string) ([]*entities.VitalSign, error)

	UpdateVitalSignByID(vitalSignID string, vitalSign *entities.VitalSign, userID string) (*entities.VitalSign, error)

	// LaboratoryValue operations
	CreateLaboratoryValue(laboratoryValue *entities.LaboratoryValue, userID string) (*entities.LaboratoryValue, error)
	GetLaboratoryValuesOverview(req models.LaboratoryValueQueryParams, userID string) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesByResident(residentID string, isLatest string, userID string) ([]*entities.LaboratoryValue, error)
	GetRoomLaboratoryValues(roomID string, isLatest string, userID string) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesHistory(residentID string, userID string) ([]*entities.LaboratoryValue, error)
	GetAbnormalLaboratoryValues(floor string, isLatest string, userID string) ([]*entities.LaboratoryValue, error)
	GetUrineOutputSumByResidentID(residentID string, req models.LaboratoryValueQueryParams, userID string) (*models.UrineOutputSummaryByResidentResponse, error)

	UpdateLaboratoryValueByID(laboratoryValueID string, laboratoryValue *entities.LaboratoryValue, userID string) (*entities.LaboratoryValue, error)
	//todo search resident by like sql
	//todo overview resident
}

type EmrUseCaseImpl struct {
	emrrepo      repositories.EmrRepository
	auditlogrepo audit_repo.AuditLogRepository
	userrepo     user_repo.UserRepository
}

func NewEmrUseCase(
	emrrepo repositories.EmrRepository,
	auditlogrepo audit_repo.AuditLogRepository,
	userrepo user_repo.UserRepository) EmrUsecase {
	return &EmrUseCaseImpl{
		emrrepo:      emrrepo,
		auditlogrepo: auditlogrepo,
		userrepo:     userrepo,
	}
}

func (uc *EmrUseCaseImpl) CreateResident(resident *entities.Resident, userID string) (*entities.Resident, error) {

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

func (uc *EmrUseCaseImpl) GetResidentByID(id string) (*entities.Resident, error) {
	resident, err := uc.emrrepo.GetResidentByID(id)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}
	return resident, nil
}

func (uc *EmrUseCaseImpl) GetResidentByRoomID(roomID string) ([]*entities.Resident, error) {
	residents, err := uc.emrrepo.GetResidentByRoomID(roomID)
	if err != nil {
		return nil, errors.New("failed to get residents by room ID: " + err.Error())
	}
	return residents, nil
}

func (uc *EmrUseCaseImpl) GetAllResidents() ([]*entities.Resident, error) {
	residents, err := uc.emrrepo.GetAllResidents()
	if err != nil {
		return nil, errors.New("failed to get all residents: " + err.Error())
	}
	return residents, nil
}

func (uc *EmrUseCaseImpl) GetResidentOverview(req models.ResidentQueryParams) ([]*models.ResidentOverviewResponse, error) {
	var (
		residents []*entities.Resident
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

	if req.Status != nil && *req.Status != "" && *req.Status != emr_constants.Active && *req.Status != emr_constants.InActive {
		return nil, errors.New("status must be 'active' or 'inactive'")
	}

	log.Printf("เข้ามาใน overview")
	if hasFilter {
		residents, err = uc.emrrepo.GetResidentsCustom(req)
		log.Printf("มีการใช้ filter ใน GetResidentOverview: floor=%v, label_ids=%v, search=%v, status=%v | จำนวน residents ที่ได้จาก custom query: %d", req.Floor, req.LabelIDs, req.Search, req.Status, len(residents))
	} else {
		residents, err = uc.emrrepo.GetAllResidents()
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
	return response, nil
}

func (uc *EmrUseCaseImpl) UpdateResidentByID(residentID string, data models.UpdateResidentRequest, userID string) (*entities.Resident, error) {
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

func (uc *EmrUseCaseImpl) GetNumberOfResidentsDashboard() (models.NumberOfResidentsDashboardResponse, error) {
	response, err := uc.emrrepo.GetNumberOfResidentsDashboard()
	if err != nil {
		return models.NumberOfResidentsDashboardResponse{}, errors.New("failed to get dashboard data: " + err.Error())
	}
	return response, nil
}

func (uc *EmrUseCaseImpl) GetResidentGenderStatsDashboard() (models.ResidentGenderStatsDashboardResponse, error) {
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

func (uc *EmrUseCaseImpl) GetResidentAllergyStatsDashboard() (models.ResidentAllergyStatsDashboardResponse, error) {
	response, err := uc.emrrepo.GetResidentAllergyStatsDashboard()
	if err != nil {
		return models.ResidentAllergyStatsDashboardResponse{}, errors.New("failed to get resident allergy stats: " + err.Error())
	}

	return response, nil
}

func (uc *EmrUseCaseImpl) GetRoomByID(id string) (*entities.Room, error) {
	room, err := uc.emrrepo.GetRoomByID(id)
	if err != nil {
		return nil, errors.New("room not found: " + err.Error())
	}
	return room, nil
}

func (uc *EmrUseCaseImpl) GetAllRooms() ([]*entities.Room, error) {
	rooms, err := uc.emrrepo.GetAllRooms()
	if err != nil {
		return nil, errors.New("failed to get all rooms: " + err.Error())
	}
	return rooms, nil
}

func (uc *EmrUseCaseImpl) CreateRoom(room *entities.Room, userID string) (*entities.Room, error) {
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

func (uc *EmrUseCaseImpl) GetAllIntakeLabels() ([]*entities.IntakeLabels, error) {
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

func (uc *EmrUseCaseImpl) GetAllAllergies() ([]*entities.Allergy, error) {
	allergies, err := uc.emrrepo.GetAllAllergies()
	if err != nil {
		return nil, errors.New("failed to get all allergies: " + err.Error())
	}
	return allergies, nil
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

func (uc *EmrUseCaseImpl) CreateIntakeLabelByResidentID(residentID string, labels []models.IntakeLabelRequest, userID string) ([]*entities.ResidentLabels, error) {

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

func (uc *EmrUseCaseImpl) GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error) {
	residentLabels, err := uc.emrrepo.GetResidentLabelsByResidentID(residentID)
	if err != nil {
		return nil, errors.New("failed to get resident labels: " + err.Error())
	}
	return residentLabels, nil
}

func (uc *EmrUseCaseImpl) CreateAllergyByResidentID(residentID string, allergies []models.AllergyRequest, userID string) ([]*entities.ResidentAllergies, error) {
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

func (uc *EmrUseCaseImpl) GetResidentAllergiesByResidentID(residentID string) ([]*entities.ResidentAllergies, error) {
	residentAllergies, err := uc.emrrepo.GetResidentAllergiesByResidentID(residentID)
	if err != nil {
		return nil, errors.New("failed to get resident allergies: " + err.Error())
	}
	return residentAllergies, nil
}

func (uc *EmrUseCaseImpl) GetAllResidentAllergies() ([]*models.ResidentAllergyListResponse, error) {
	residentAllergies, err := uc.emrrepo.GetAllResidentAllergies()
	if err != nil {
		return nil, errors.New("failed to get all resident allergies: " + err.Error())
	}
	return residentAllergies, nil
}

func (uc *EmrUseCaseImpl) CreateVitalSign(vitalSign *entities.VitalSign, userID string) (*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can create vital signs")
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
func filterVitalSignsByStatus(vitalSigns []*entities.VitalSign, status string) []*entities.VitalSign {
	if status == "" || status == "all" {
		return vitalSigns
	}

	wantAbnormal := status == "abnormal"
	// log.Printf("Filtering vital signs for status '%s' (wantAbnormal=%t)", status, wantAbnormal)

	result := make([]*entities.VitalSign, 0)
	for _, vs := range vitalSigns {
		isAbnormal := false

		if vs.Temperature != nil && (*vs.Temperature < emr_constants.NormalTempLow || *vs.Temperature > emr_constants.NormalTempHigh) {
			isAbnormal = true
		}
		if vs.HeartRate != nil && (*vs.HeartRate < emr_constants.NormalHeartRateLow || *vs.HeartRate > emr_constants.NormalHeartRateHigh) {
			isAbnormal = true
		}
		if vs.BreathingRate != nil && (*vs.BreathingRate < emr_constants.NormalBreathingRateLow || *vs.BreathingRate > emr_constants.NormalBreathingRateHigh) {
			isAbnormal = true
		}
		if vs.BloodPressureSystolic != nil && (*vs.BloodPressureSystolic < emr_constants.NormalSystolicLow || *vs.BloodPressureSystolic > emr_constants.NormalSystolicHigh) {
			isAbnormal = true
		}
		if vs.BloodPressureDiastolic != nil && (*vs.BloodPressureDiastolic < emr_constants.NormalDiastolicLow || *vs.BloodPressureDiastolic > emr_constants.NormalDiastolicHigh) {
			isAbnormal = true
		}
		if vs.OxygenSaturation != nil && *vs.OxygenSaturation < emr_constants.NormalOxygenSaturationLow {
			isAbnormal = true
		}

		if isAbnormal == wantAbnormal {
			result = append(result, vs)
		}
	}
	return result
}

func (uc *EmrUseCaseImpl) GetVitalSignsOverview(req models.VitalSignQueryParams, userID string) ([]*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view vital signs overview")
	}

	if req.VitalSignStatus != "" && req.VitalSignStatus != "all" && req.VitalSignStatus != "normal" && req.VitalSignStatus != "abnormal" {
		return nil, errors.New("vitalsign_status must be 'all', 'normal', or 'abnormal'")
	}

	// กรณีธรรมดา: ทั้งหมด วันนี้
	var vitalSigns []*entities.VitalSign

	if req.Floor == nil && len(req.LabelIDs) == 0 {
		vitalSigns, err = uc.emrrepo.GetVitalSignsToday(false)
	} else {
		// กรณีมี filter: ใช้ Custom
		params := models.VitalSignQueryParams{
			Floor:    req.Floor,
			LabelIDs: req.LabelIDs,
			IsLatest: false,
			Limit:    100,
		}
		vitalSigns, err = uc.emrrepo.GetVitalSignsCustom(params)
	}
	if err != nil {
		return nil, err
	}
	// log.Printf("ก่อนจะกรอง vital signs ทั้งหมด: %d รายการ | vitalsign_status='%s'", len(vitalSigns), req.VitalSignStatus)
	return filterVitalSignsByStatus(vitalSigns, req.VitalSignStatus), nil
}

func (uc *EmrUseCaseImpl) GetVitalSignsByResident(residentID string, isLatest string, userID string) ([]*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view vital signs by resident")
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	isLatestBool, err := strconv.ParseBool(isLatest)
	if err != nil {
		return nil, errors.New("invalid isLatest parameter you must provide a boolean value: " + err.Error())
	}

	vitalSigns, err := uc.emrrepo.GetVitalSignsByResidentIDToday(residentID, isLatestBool)
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

	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view vital signs by room")
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

	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view vital signs history")
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

func (uc *EmrUseCaseImpl) GetAbnormalVitalSigns(floor string, isLatest string, userID string) ([]*entities.VitalSign, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view abnormal vital signs")
	}

	var vitalSigns []*entities.VitalSign
	var isLatestBool bool
	// var err error

	var floor_pointer_int16 *int16
	if floor != "" {
		floor64, err := strconv.ParseInt(floor, 10, 16)
		if err != nil {
			return nil, errors.New("invalid floor parameter")
		}
		f := int16(floor64)
		floor_pointer_int16 = &f
	}

	isLatestBool, err = strconv.ParseBool(isLatest)
	if err != nil {
		return nil, errors.New("invalid isLatest parameter you must provide a boolean value: " + err.Error())
	}

	// ถ้าไม่ระบุ floor → ดึงทั้งหมด, ถ้าระบุ → ดึงเฉพาะชั้นนั้น
	if floor_pointer_int16 == nil {
		vitalSigns, err = uc.emrrepo.GetVitalSignsToday(isLatestBool)
	} else {
		vitalSigns, err = uc.emrrepo.GetVitalSignsByFloorToday(*floor_pointer_int16, isLatestBool)
	}
	if err != nil {
		return nil, errors.New("failed to get vital signs: " + err.Error())
	}

	abnormalVitalSigns := filterVitalSignsByStatus(vitalSigns, "abnormal")
	return abnormalVitalSigns, nil
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

	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can update vital signs")
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

func (uc *EmrUseCaseImpl) CreateLaboratoryValue(laboratoryValue *entities.LaboratoryValue, userID string) (*entities.LaboratoryValue, error) {

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}

	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}

	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can create laboratory values")
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

// filterLaboratoryValuesByStatus filters laboratory values by status: "all", "normal", or "abnormal".
func filterLaboratoryValuesByStatus(labs []*entities.LaboratoryValue, status string) []*entities.LaboratoryValue {
	if status == "" || status == "all" {
		return labs
	}

	wantAbnormal := status == "abnormal"
	result := make([]*entities.LaboratoryValue, 0)
	for _, lab := range labs {
		isAbnormal := false

		if lab.BloodGlucose != nil && (*lab.BloodGlucose < emr_constants.NormalBloodGlucoseLow || *lab.BloodGlucose > emr_constants.NormalBloodGlucoseHigh) {
			isAbnormal = true
		}
		if lab.FluidIn != nil && (*lab.FluidIn < emr_constants.NormalFluidInLow || *lab.FluidIn > emr_constants.NormalFluidInHigh) {
			isAbnormal = true
		}
		if lab.FluidOut != nil && (*lab.FluidOut < emr_constants.NormalFluidOutLow || *lab.FluidOut > emr_constants.NormalFluidOutHigh) {
			isAbnormal = true
		}
		if lab.UrineOutput != nil && lab.UrineType != nil {
			if *lab.UrineType == emr_constants.UrineTypeML {
				if *lab.UrineOutput < emr_constants.NormalUrineOutputMLLow || *lab.UrineOutput > emr_constants.NormalUrineOutputMLHigh {
					isAbnormal = true
				}
			} else {
				if *lab.UrineOutput < emr_constants.NormalUrineOutputTimesLow || *lab.UrineOutput > emr_constants.NormalUrineOutputTimesHigh {
					isAbnormal = true
				}
			}
		}
		if lab.Stool != nil && *lab.Stool > emr_constants.NormalStoolHigh {
			isAbnormal = true
		}
		if lab.DiaperChange != nil && *lab.DiaperChange > emr_constants.NormalDiaperChangeHigh {
			isAbnormal = true
		}

		if isAbnormal == wantAbnormal {
			result = append(result, lab)
		}
	}
	return result
}

func (uc *EmrUseCaseImpl) GetLaboratoryValuesOverview(req models.LaboratoryValueQueryParams, userID string) ([]*entities.LaboratoryValue, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view laboratory values overview")
	}

	if req.LaboratoryValueStatus != "" && req.LaboratoryValueStatus != "all" && req.LaboratoryValueStatus != "normal" && req.LaboratoryValueStatus != "abnormal" {
		return nil, errors.New("laboratory_value_status must be 'all', 'normal', or 'abnormal'")
	}

	var labs []*entities.LaboratoryValue
	if req.Floor == nil && len(req.LabelIDs) == 0 {
		labs, err = uc.emrrepo.GetLaboratoryValuesToday(false)
	} else {
		params := models.LaboratoryValueQueryParams{
			Floor:    req.Floor,
			LabelIDs: req.LabelIDs,
			IsLatest: false,
			Limit:    100,
		}
		labs, err = uc.emrrepo.GetLaboratoryValuesCustom(params)
	}
	if err != nil {
		return nil, err
	}
	return filterLaboratoryValuesByStatus(labs, req.LaboratoryValueStatus), nil
}

func (uc *EmrUseCaseImpl) GetLaboratoryValuesByResident(residentID string, isLatest string, userID string) ([]*entities.LaboratoryValue, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view laboratory values by resident")
	}

	residentExists, err := uc.emrrepo.ResidentExists(residentID)
	if err != nil {
		return nil, errors.New("failed to verify resident existence: " + err.Error())
	}
	if !residentExists {
		return nil, errors.New("resident not found")
	}

	isLatestBool, err := strconv.ParseBool(isLatest)
	if err != nil {
		return nil, errors.New("invalid isLatest parameter you must provide a boolean value: " + err.Error())
	}

	labs, err := uc.emrrepo.GetLaboratoryValuesByResidentIDToday(residentID, isLatestBool)
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
	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view laboratory values by room")
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
	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view laboratory values history")
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

func (uc *EmrUseCaseImpl) GetAbnormalLaboratoryValues(floor string, isLatest string, userID string) ([]*entities.LaboratoryValue, error) {
	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user: " + err.Error())
	}
	userRole, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("failed to get user role: " + err.Error())
	}
	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view abnormal laboratory values")
	}

	var floorPtr *int16
	if floor != "" {
		floor64, err := strconv.ParseInt(floor, 10, 16)
		if err != nil {
			return nil, errors.New("invalid floor parameter")
		}
		f := int16(floor64)
		floorPtr = &f
	}

	isLatestBool, err := strconv.ParseBool(isLatest)
	if err != nil {
		return nil, errors.New("invalid isLatest parameter you must provide a boolean value: " + err.Error())
	}

	var labs []*entities.LaboratoryValue
	if floorPtr == nil {
		labs, err = uc.emrrepo.GetLaboratoryValuesToday(isLatestBool)
	} else {
		labs, err = uc.emrrepo.GetLaboratoryValuesByFloorToday(*floorPtr, isLatestBool)
	}
	if err != nil {
		return nil, errors.New("failed to get laboratory values: " + err.Error())
	}

	return filterLaboratoryValuesByStatus(labs, "abnormal"), nil
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
	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can view urine output summary")
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
	if userRole.Name != user_constants.RoleMedicalStaff {
		return nil, errors.New("only users with 'Medical Staff' role can update laboratory values")
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
