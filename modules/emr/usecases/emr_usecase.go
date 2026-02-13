package usecases

import (
	"encoding/json"
	"errors"

	// "io"
	"log"
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
	UpdateResidentByID(residentID string, data models.UpdateResidentRequest, userID string) (*entities.Resident, error)

	// Dashboard operations
	GetNumberOfResidentsDashboard() (models.NumberOfResidentsDashboardResponse, error)
	GetResidentGenderStatsDashboard() (models.ResidentGenderStatsDashboardResponse, error)

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

	// ResidentLabel operations (many-to-many)
	CreateIntakeLabelByResidentID(residentID string, labels []models.IntakeLabelRequest, userID string) ([]*entities.ResidentLabels, error)
	GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error)

	// VitalSign operations
	CreateVitalSign(vitalSign *entities.VitalSign, userID string) (*entities.VitalSign, error)
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

	if resident.Age != nil && *resident.Age < 0 {
		return nil, errors.New("age cannot be negative")
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

	resident.ID = uuid.New().String()

	createdResident, err := uc.emrrepo.CreateResident(resident)
	if err != nil {
		return nil, errors.New("failed to create resident: " + err.Error())
	}

	newResidentData, _ := json.Marshal(map[string]interface{}{
		"first_name": createdResident.FirstName,
		"last_name":  createdResident.LastName,
		"age":        createdResident.Age,
		"gender":     createdResident.Gender,
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

func (uc *EmrUseCaseImpl) UpdateResidentByID(residentID string, data models.UpdateResidentRequest, userID string) (*entities.Resident, error) {
	resident, err := uc.emrrepo.GetResidentByID(residentID)
	if err != nil {
		return nil, errors.New("resident not found: " + err.Error())
	}

	oldResidentData, _ := json.Marshal(map[string]interface{}{
		"room_id":    resident.RoomID,
		"first_name": resident.FirstName,
		"last_name":  resident.LastName,
		"age":        resident.Age,
		"gender":     resident.Gender,
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

	if data.Age != nil {
		if *data.Age < 0 {
			return nil, errors.New("age cannot be negative")
		}
		resident.Age = data.Age
	}

	if data.Gender != nil {
		gender := strings.ToLower(strings.TrimSpace(*data.Gender))
		if gender != "male" && gender != "female" && gender != "other" {
			return nil, errors.New("gender must be either 'male', 'female', or 'other'")
		}
		resident.Gender = gender
	}

	updatedResident, err := uc.emrrepo.UpdateResident(resident)
	if err != nil {
		return nil, errors.New("failed to update resident: " + err.Error())
	}

	newResidentData, _ := json.Marshal(map[string]interface{}{
		"room_id":    updatedResident.RoomID,
		"first_name": updatedResident.FirstName,
		"last_name":  updatedResident.LastName,
		"age":        updatedResident.Age,
		"gender":     updatedResident.Gender,
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
		return nil, errors.New("resident does not exist")
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
	return createdVitalSign, nil
}
