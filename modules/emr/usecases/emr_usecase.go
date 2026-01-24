package usecases

import (
	"encoding/json"
	"errors"

	// "io"
	"log"
	// "mime/multipart"
	// "os"
	// "sync"
	// "time"

	audit_repo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	// "github.com/aikidoaikido115/New-Acis-BE/pkg/utils"
	"github.com/google/uuid"
	// "golang.org/x/text/unicode/norm"
)

type EmrUsecase interface {
	CreateResident(resident *entities.Resident) (*entities.Resident, error)
	GetResidentByID(id string) (*entities.Resident, error)
	GetResidentByRoomID(roomID string) ([]*entities.Resident, error)
	GetAllResidents() ([]*entities.Resident, error)
	GetRoomByID(id string) (*entities.Room, error)
	GetAllRooms() ([]*entities.Room, error)
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

func (uc *EmrUseCaseImpl) CreateResident(resident *entities.Resident) (*entities.Resident, error) {

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
		UserID:    "system", // ในกรณีนี้ใช้ "system" เป็นผู้ใช้ที่สร้าง
		Action:    "CREATE",
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
