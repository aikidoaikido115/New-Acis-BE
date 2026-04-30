package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/models"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	GetWarehouseItems(search string, category string) ([]*entities.WarehouseItem, error)
	GetWarehouseItemByID(id string) (*entities.WarehouseItem, error)
	GetWarehouseItemByCode(code string) (*entities.WarehouseItem, error)
	CreateWarehouseItem(item *entities.WarehouseItem) (*entities.WarehouseItem, error)
	UpdateWarehouseItem(item *entities.WarehouseItem) (*entities.WarehouseItem, error)
	DeleteWarehouseItem(id string) error
	NextItemCode(category string) (string, error)

	CreateTransaction(transaction *entities.WarehouseTransaction) (*entities.WarehouseTransaction, error)
	GetTransactions(query models.ListWarehouseTransactionsQuery) ([]*entities.WarehouseTransaction, error)
	GetTransactionByID(id string) (*entities.WarehouseTransaction, error)
	GetTransactionsByIDs(ids []string) ([]*entities.WarehouseTransaction, error)
	UpdateTransaction(transaction *entities.WarehouseTransaction) (*entities.WarehouseTransaction, error)
	NextTransactionCode(datePrefix string) (string, error)
}

type GormWarehouseRepository struct {
	db *gorm.DB
}

func NewGormWarehouseRepository(db *gorm.DB) *GormWarehouseRepository {
	return &GormWarehouseRepository{db: db}
}

func (r *GormWarehouseRepository) GetWarehouseItems(search string, category string) ([]*entities.WarehouseItem, error) {
	var items []*entities.WarehouseItem

	query := r.db.Model(&entities.WarehouseItem{})

	if strings.TrimSpace(search) != "" {
		searchLike := "%" + strings.TrimSpace(search) + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", searchLike, searchLike)
	}

	if strings.TrimSpace(category) != "" {
		query = query.Where("category = ?", strings.TrimSpace(category))
	}

	if err := query.Order("code ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *GormWarehouseRepository) GetWarehouseItemByID(id string) (*entities.WarehouseItem, error) {
	var item entities.WarehouseItem
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GormWarehouseRepository) GetWarehouseItemByCode(code string) (*entities.WarehouseItem, error) {
	var item entities.WarehouseItem
	if err := r.db.Where("LOWER(code) = LOWER(?)", code).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GormWarehouseRepository) CreateWarehouseItem(item *entities.WarehouseItem) (*entities.WarehouseItem, error) {
	if err := r.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return r.GetWarehouseItemByID(item.ID)
}

func (r *GormWarehouseRepository) UpdateWarehouseItem(item *entities.WarehouseItem) (*entities.WarehouseItem, error) {
	if err := r.db.Save(&item).Error; err != nil {
		return nil, err
	}
	return r.GetWarehouseItemByID(item.ID)
}

func (r *GormWarehouseRepository) DeleteWarehouseItem(id string) error {
	if err := r.db.Delete(&entities.WarehouseItem{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormWarehouseRepository) NextItemCode(category string) (string, error) {
	var item entities.WarehouseItem
	err := r.db.Where("category = ?", category).Order("code DESC").First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Sprintf("%s-%03d", category, 1), nil
		}
		return "", err
	}

	parts := strings.Split(item.Code, "-")
	if len(parts) != 2 {
		return fmt.Sprintf("%s-%03d", category, 1), nil
	}

	seq, convErr := strconv.Atoi(parts[1])
	if convErr != nil {
		return fmt.Sprintf("%s-%03d", category, 1), nil
	}

	return fmt.Sprintf("%s-%03d", category, seq+1), nil
}

func (r *GormWarehouseRepository) CreateTransaction(transaction *entities.WarehouseTransaction) (*entities.WarehouseTransaction, error) {
	if err := r.db.Create(&transaction).Error; err != nil {
		return nil, err
	}
	return r.GetTransactionByID(transaction.ID)
}

func (r *GormWarehouseRepository) GetTransactions(queryParams models.ListWarehouseTransactionsQuery) ([]*entities.WarehouseTransaction, error) {
	var transactions []*entities.WarehouseTransaction

	query := r.db.Model(&entities.WarehouseTransaction{})

	if strings.TrimSpace(queryParams.SearchItem) != "" {
		searchLike := "%" + strings.TrimSpace(queryParams.SearchItem) + "%"
		query = query.Where("item_name ILIKE ? OR item_code ILIKE ?", searchLike, searchLike)
	}

	if strings.TrimSpace(queryParams.SearchUser) != "" {
		searchLike := "%" + strings.TrimSpace(queryParams.SearchUser) + "%"
		query = query.Where("operator ILIKE ?", searchLike)
	}

	if strings.TrimSpace(queryParams.Status) != "" {
		query = query.Where("approval_status = ?", strings.TrimSpace(queryParams.Status))
	}

	if strings.TrimSpace(queryParams.Type) != "" {
		query = query.Where("type = ?", strings.TrimSpace(queryParams.Type))
	}

	if strings.TrimSpace(queryParams.StartDate) != "" {
		query = query.Where("DATE(created_at AT TIME ZONE 'Asia/Bangkok') >= ?", strings.TrimSpace(queryParams.StartDate))
	}

	if strings.TrimSpace(queryParams.EndDate) != "" {
		query = query.Where("DATE(created_at AT TIME ZONE 'Asia/Bangkok') <= ?", strings.TrimSpace(queryParams.EndDate))
	}

	if err := query.Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *GormWarehouseRepository) GetTransactionByID(id string) (*entities.WarehouseTransaction, error) {
	var transaction entities.WarehouseTransaction
	if err := r.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *GormWarehouseRepository) GetTransactionsByIDs(ids []string) ([]*entities.WarehouseTransaction, error) {
	var transactions []*entities.WarehouseTransaction
	if err := r.db.Where("id IN ?", ids).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *GormWarehouseRepository) UpdateTransaction(transaction *entities.WarehouseTransaction) (*entities.WarehouseTransaction, error) {
	if err := r.db.Save(&transaction).Error; err != nil {
		return nil, err
	}
	return r.GetTransactionByID(transaction.ID)
}

func (r *GormWarehouseRepository) NextTransactionCode(datePrefix string) (string, error) {
	var latest entities.WarehouseTransaction
	err := r.db.Where("code LIKE ?", datePrefix+"-%").Order("code DESC").First(&latest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Sprintf("%s-%03d", datePrefix, 1), nil
		}
		return "", err
	}

	parts := strings.Split(latest.Code, "-")
	if len(parts) != 2 {
		return fmt.Sprintf("%s-%03d", datePrefix, 1), nil
	}

	seq, convErr := strconv.Atoi(parts[1])
	if convErr != nil {
		return fmt.Sprintf("%s-%03d", datePrefix, 1), nil
	}

	return fmt.Sprintf("%s-%03d", datePrefix, seq+1), nil
}
