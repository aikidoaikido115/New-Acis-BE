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

	audit_repo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"github.com/aikidoaikido115/New-Acis-BE/pkg/utils"
	"github.com/google/uuid"
	// "golang.org/x/text/unicode/norm"
)

type EmrUsecase interface {

	// Resident operations
	CreateResident(resident *entities.Resident, userID string) (*entities.Resident, error)
	GetResidentByID(id string) (*entities.Resident, error)
	GetResidentByRoomID(roomID string) ([]*entities.Resident, error)
	GetAllResidents() ([]*entities.Resident, error)

	// Room operations
	GetRoomByID(id string) (*entities.Room, error)
	GetAllRooms() ([]*entities.Room, error)

	// IntakeLabel operations
	CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error)
	// GetIntakeLabelByID(id string) (*entities.IntakeLabels, error)
	GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error)
	GetAllIntakeLabels() ([]*entities.IntakeLabels, error)

	// ResidentLabel operations (many-to-many)
	CreateIntakeLabelByResidentID(residentID string, labels []IntakeLabelInput, userID string) ([]*entities.ResidentLabels, error)
	GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error)
}

type EmrUseCaseImpl struct {
	emrrepo      repositories.EmrRepository
	auditlogrepo audit_repo.AuditLogRepository
}

func NewEmrUseCase(
	emrrepo repositories.EmrRepository,
	auditlogrepo audit_repo.AuditLogRepository) EmrUsecase {

	return &EmrUseCaseImpl{
		emrrepo:      emrrepo,
		auditlogrepo: auditlogrepo,
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

	resident.ID = uuid.New().String()

	createdResident, err := uc.emrrepo.CreateResident(resident)
	if err != nil {
		return nil, errors.New("failed to create resident: " + err.Error())
	}

	newResidentData, _ := json.Marshal(map[string]interface{}{
		"first_name": createdResident.FirstName,
		"last_name":  createdResident.LastName,
		"age":        createdResident.Age,
	})
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "residents",
		RecordID:  createdResident.ID,
		UserID:    userID,
		Action:    utils.AuditActionInsert,
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

func (uc *EmrUseCaseImpl) CreateIntakeLabelByResidentID(residentID string, labels []IntakeLabelInput, userID string) ([]*entities.ResidentLabels, error) {

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
			Action:    utils.AuditActionInsert,
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
