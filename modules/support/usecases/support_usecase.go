package usecases

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/support/models"
	supportRepo "github.com/aikidoaikido115/New-Acis-BE/modules/support/repositories"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	userRepo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SupportUsecase interface {
	CreateSupportTicket(ticket *entities.SupportTicket, userID string) (*entities.SupportTicket, error)
	GetSupportTickets(query models.ListSupportTicketsQuery, userID string) ([]*entities.SupportTicket, error)
	GetSupportTicketByID(id string, userID string) (*entities.SupportTicket, error)
	UpdateSupportTicketStatus(id string, req models.UpdateSupportTicketStatusRequest, userID string) (*entities.SupportTicket, error)
	DeleteSupportTicketByID(id string, userID string) error
}

type SupportUseCaseImpl struct {
	repo     supportRepo.SupportRepository
	userrepo userRepo.UserRepository
}

func NewSupportUseCase(repo supportRepo.SupportRepository, userrepo userRepo.UserRepository) SupportUsecase {
	return &SupportUseCaseImpl{
		repo:     repo,
		userrepo: userrepo,
	}
}

func (uc *SupportUseCaseImpl) CreateSupportTicket(ticket *entities.SupportTicket, userID string) (*entities.SupportTicket, error) {
	reporterRole, err := uc.ensureAllowedReporter(userID)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, errors.New("support ticket payload is required")
	}

	ticket.Name = strings.TrimSpace(ticket.Name)
	ticket.Email = strings.TrimSpace(ticket.Email)
	ticket.Subject = strings.TrimSpace(ticket.Subject)
	ticket.Message = strings.TrimSpace(ticket.Message)

	if ticket.Name == "" {
		return nil, errors.New("name is required")
	}

	if ticket.Email == "" {
		return nil, errors.New("email is required")
	}

	if _, err := mail.ParseAddress(ticket.Email); err != nil {
		return nil, errors.New("email format is invalid")
	}

	if ticket.Subject == "" {
		return nil, errors.New("subject is required")
	}

	if ticket.Message == "" {
		return nil, errors.New("message is required")
	}

	ticket.ID = uuid.New().String()
	ticket.Status = "open"
	ticket.ReporterRole = reporterRole
	ticket.CreatedByUserID = userID
	now := time.Now().UTC()
	ticket.CreatedAt = now
	ticket.UpdatedAt = now

	return uc.repo.CreateSupportTicket(ticket)
}

func (uc *SupportUseCaseImpl) GetSupportTickets(query models.ListSupportTicketsQuery, userID string) ([]*entities.SupportTicket, error) {
	if err := uc.ensureMedicalStaffAdmin(userID); err != nil {
		return nil, err
	}

	query.Search = strings.TrimSpace(query.Search)
	query.Status = strings.TrimSpace(strings.ToLower(query.Status))
	query.ReporterRole = strings.TrimSpace(query.ReporterRole)

	if query.Status != "" && !isValidSupportStatus(query.Status) {
		return nil, errors.New("status must be one of open, in_progress, resolved")
	}

	return uc.repo.GetSupportTickets(query)
}

func (uc *SupportUseCaseImpl) GetSupportTicketByID(id string, userID string) (*entities.SupportTicket, error) {
	if err := uc.ensureMedicalStaffAdmin(userID); err != nil {
		return nil, err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("support ticket id is required")
	}

	ticket, err := uc.repo.GetSupportTicketByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("support ticket not found")
		}

		return nil, err
	}

	return ticket, nil
}

func (uc *SupportUseCaseImpl) UpdateSupportTicketStatus(id string, req models.UpdateSupportTicketStatusRequest, userID string) (*entities.SupportTicket, error) {
	if err := uc.ensureMedicalStaffAdmin(userID); err != nil {
		return nil, err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("support ticket id is required")
	}

	status := strings.TrimSpace(strings.ToLower(req.Status))
	if status == "" {
		return nil, errors.New("status is required")
	}

	if !isValidSupportStatus(status) {
		return nil, errors.New("status must be one of open, in_progress, resolved")
	}

	ticket, err := uc.repo.GetSupportTicketByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("support ticket not found")
		}

		return nil, err
	}

	ticket.Status = status
	ticket.UpdatedAt = time.Now().UTC()

	return uc.repo.UpdateSupportTicket(ticket)
}

func (uc *SupportUseCaseImpl) DeleteSupportTicketByID(id string, userID string) error {
	if err := uc.ensureMedicalStaffAdmin(userID); err != nil {
		return err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("support ticket id is required")
	}

	ticket, err := uc.repo.GetSupportTicketByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("support ticket not found")
		}

		return err
	}

	if ticket.Status != "resolved" {
		return errors.New("only resolved support tickets can be deleted")
	}

	return uc.repo.DeleteSupportTicketByID(id)
}

func (uc *SupportUseCaseImpl) ensureAllowedReporter(userID string) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", errors.New("user id is required")
	}

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return "", errors.New("failed to get user: " + err.Error())
	}

	role, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return "", errors.New("failed to get user role: " + err.Error())
	}

	if role.Name != user_constants.RoleMedicalStaff && role.Name != user_constants.RoleKitchenStaff {
		return "", errors.New("only users with 'Medical Staff' or 'Kitchen Staff' role can create support tickets")
	}

	return role.Name, nil
}

func (uc *SupportUseCaseImpl) ensureMedicalStaffAdmin(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := uc.userrepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	role, err := uc.userrepo.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}

	if role.Name != user_constants.RoleMedicalStaff {
		return errors.New("only users with 'Medical Staff' role can manage support tickets")
	}

	return nil
}

func isValidSupportStatus(status string) bool {
	switch status {
	case "open", "in_progress", "resolved":
		return true
	default:
		return false
	}
}
