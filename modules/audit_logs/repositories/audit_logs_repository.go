package repositories

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"gorm.io/gorm"
)

type GormAuditLogRepository struct {
	db *gorm.DB
}

func NewGormAuditLogRepository(db *gorm.DB) *GormAuditLogRepository {
	return &GormAuditLogRepository{
		db: db,
	}
}

type AuditLogRepository interface {
	CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error)
	GetAuditLogByID(id string) (*entities.AuditLogs, error)
}

func (r *GormAuditLogRepository) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&auditLog).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetAuditLogByID(auditLog.ID)
}

func (r *GormAuditLogRepository) GetAuditLogByID(id string) (*entities.AuditLogs, error) {
	var auditLog entities.AuditLogs
	if err := r.db.First(&auditLog, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &auditLog, nil
}