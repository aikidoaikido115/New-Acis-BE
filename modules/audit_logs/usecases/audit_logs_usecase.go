package usecases

import (
	"errors"
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"gorm.io/gorm"
)

type AuditLogUsecase interface {
	GetAllAuditLogs(userID string) ([]*entities.AuditLogs, error)
	SearchAuditLogs(search string, userID string) ([]*entities.AuditLogs, error)
	GetAuditLogByID(id string, userID string) (*entities.AuditLogs, error)
}

type AuditLogUseCaseImpl struct {
	repo     repositories.AuditLogRepository
	userrepo user_repo.UserRepository
}

func NewAuditLogUseCase(repo repositories.AuditLogRepository, userrepo user_repo.UserRepository) AuditLogUsecase {
	return &AuditLogUseCaseImpl{repo: repo, userrepo: userrepo}
}

func (uc *AuditLogUseCaseImpl) GetAllAuditLogs(userID string) ([]*entities.AuditLogs, error) {
	if err := uc.ensureAdmin(userID); err != nil {
		return nil, err
	}

	return uc.repo.GetAllAuditLogs()
}

func (uc *AuditLogUseCaseImpl) SearchAuditLogs(search string, userID string) ([]*entities.AuditLogs, error) {
	if err := uc.ensureAdmin(userID); err != nil {
		return nil, err
	}

	search = strings.TrimSpace(search)
	if search == "" {
		return uc.repo.GetAllAuditLogs()
	}

	return uc.repo.SearchAuditLogs(search)
}

func (uc *AuditLogUseCaseImpl) GetAuditLogByID(id string, userID string) (*entities.AuditLogs, error) {
	if err := uc.ensureAdmin(userID); err != nil {
		return nil, err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("audit log id is required")
	}

	auditLog, err := uc.repo.GetAuditLogByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audit log not found")
		}
		return nil, err
	}

	return auditLog, nil
}

func (uc *AuditLogUseCaseImpl) ensureAdmin(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	if user.Role.Name != user_constants.RoleAdmin {
		return errors.New("only users with 'Admin' role can manage audit logs")
	}

	return nil
}
