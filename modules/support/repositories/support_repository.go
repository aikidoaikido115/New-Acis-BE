package repositories

import (
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/support/models"

	"gorm.io/gorm"
)

type GormSupportRepository struct {
	db *gorm.DB
}

func NewGormSupportRepository(db *gorm.DB) *GormSupportRepository {
	return &GormSupportRepository{db: db}
}

type SupportRepository interface {
	CreateSupportTicket(ticket *entities.SupportTicket) (*entities.SupportTicket, error)
	GetSupportTickets(query models.ListSupportTicketsQuery) ([]*entities.SupportTicket, error)
	GetSupportTicketByID(id string) (*entities.SupportTicket, error)
	UpdateSupportTicket(ticket *entities.SupportTicket) (*entities.SupportTicket, error)
	DeleteSupportTicketByID(id string) error
}

func (r *GormSupportRepository) CreateSupportTicket(ticket *entities.SupportTicket) (*entities.SupportTicket, error) {
	if err := r.db.Create(&ticket).Error; err != nil {
		return nil, err
	}

	return r.GetSupportTicketByID(ticket.ID)
}

func (r *GormSupportRepository) GetSupportTickets(queryParams models.ListSupportTicketsQuery) ([]*entities.SupportTicket, error) {
	var tickets []*entities.SupportTicket

	query := r.db.Model(&entities.SupportTicket{})

	if strings.TrimSpace(queryParams.Search) != "" {
		searchLike := "%" + strings.TrimSpace(queryParams.Search) + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR subject ILIKE ? OR message ILIKE ?", searchLike, searchLike, searchLike, searchLike)
	}

	if strings.TrimSpace(queryParams.Status) != "" {
		query = query.Where("status = ?", strings.TrimSpace(queryParams.Status))
	}

	if strings.TrimSpace(queryParams.ReporterRole) != "" {
		query = query.Where("reporter_role = ?", strings.TrimSpace(queryParams.ReporterRole))
	}

	if err := query.Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *GormSupportRepository) GetSupportTicketByID(id string) (*entities.SupportTicket, error) {
	var ticket entities.SupportTicket
	if err := r.db.Where("id = ?", id).First(&ticket).Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *GormSupportRepository) UpdateSupportTicket(ticket *entities.SupportTicket) (*entities.SupportTicket, error) {
	if err := r.db.Save(&ticket).Error; err != nil {
		return nil, err
	}

	return r.GetSupportTicketByID(ticket.ID)
}

func (r *GormSupportRepository) DeleteSupportTicketByID(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&entities.SupportTicket{}).Error; err != nil {
		return err
	}

	return nil
}
