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
	GetAllAuditLogs() ([]*entities.AuditLogs, error)
	SearchAuditLogs(search string) ([]*entities.AuditLogs, error)
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

func (r *GormAuditLogRepository) GetAllAuditLogs() ([]*entities.AuditLogs, error) {
	var auditLogs []*entities.AuditLogs
	if err := r.db.Order("created_at desc").Find(&auditLogs).Error; err != nil {
		return nil, err
	}
	return auditLogs, nil
}

func (r *GormAuditLogRepository) SearchAuditLogs(search string) ([]*entities.AuditLogs, error) {
	var auditLogs []*entities.AuditLogs
	like := "%" + search + "%"
	if err := r.db.
		Where("table_name ILIKE ? OR record_id ILIKE ? OR user_id ILIKE ? OR action ILIKE ? OR old_value ILIKE ? OR new_value ILIKE ?", like, like, like, like, like, like).
		Order("created_at desc").
		Find(&auditLogs).Error; err != nil {
		return nil, err
	}
	return auditLogs, nil
}

func (r *GormAuditLogRepository) GetAuditLogByID(id string) (*entities.AuditLogs, error) {
	var auditLog entities.AuditLogs
	if err := r.db.First(&auditLog, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &auditLog, nil
}
