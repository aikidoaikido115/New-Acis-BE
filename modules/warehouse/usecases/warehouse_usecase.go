package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	audit_constants "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/constants"
	audit_repository "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repository "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/constants"
	"github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/models"
	warehouse_repository "github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WarehouseUsecase interface {
	GetWarehouseItems(query models.ListWarehouseItemsQuery, userID string) ([]*entities.WarehouseItem, error)
	CreateWarehouseItem(req models.CreateWarehouseItemRequest, userID string) (*entities.WarehouseItem, error)
	UpdateWarehouseItemByID(id string, req models.UpdateWarehouseItemRequest, userID string) (*entities.WarehouseItem, error)
	DeleteWarehouseItemByID(id string, userID string) error
	AdjustWarehouseItemByID(id string, req models.AdjustWarehouseItemRequest, userID string) (*entities.WarehouseItem, error)

	GetWarehouseTransactions(query models.ListWarehouseTransactionsQuery, userID string) ([]*models.WarehouseTransactionResponse, error)
	GetWarehouseTransactionByID(id string, userID string) (*models.WarehouseTransactionResponse, error)
	ApproveTransactions(req models.ApproveTransactionsRequest, userID string) ([]*models.WarehouseTransactionResponse, error)
	RejectTransactions(req models.RejectTransactionsRequest, userID string) ([]*models.WarehouseTransactionResponse, error)
}

type WarehouseUseCaseImpl struct {
	warehouseRepo warehouse_repository.WarehouseRepository
	auditLogRepo  audit_repository.AuditLogRepository
	userRepo      user_repository.UserRepository
}

func NewWarehouseUseCase(
	warehouseRepo warehouse_repository.WarehouseRepository,
	auditLogRepo audit_repository.AuditLogRepository,
	userRepo user_repository.UserRepository,
) *WarehouseUseCaseImpl {
	return &WarehouseUseCaseImpl{
		warehouseRepo: warehouseRepo,
		auditLogRepo:  auditLogRepo,
		userRepo:      userRepo,
	}
}

func (uc *WarehouseUseCaseImpl) GetWarehouseItems(query models.ListWarehouseItemsQuery, userID string) ([]*entities.WarehouseItem, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	search := strings.TrimSpace(query.Search)
	category := strings.TrimSpace(query.Category)
	if strings.EqualFold(category, "all") {
		category = ""
	}
	if category != "" && !isValidCategory(category) {
		return nil, errors.New("category must be one of MED, EQU, CON")
	}

	items, err := uc.warehouseRepo.GetWarehouseItems(search, category)
	if err != nil {
		return nil, errors.New("failed to get warehouse items: " + err.Error())
	}

	return items, nil
}

func (uc *WarehouseUseCaseImpl) CreateWarehouseItem(req models.CreateWarehouseItemRequest, userID string) (*entities.WarehouseItem, error) {
	if err := uc.ensureWarehouseItemCreator(userID); err != nil {
		return nil, err
	}

	operatorUser, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get operator user: " + err.Error())
	}

	category := strings.ToUpper(strings.TrimSpace(req.Category))
	if !isValidCategory(category) {
		return nil, errors.New("category must be one of MED, EQU, CON")
	}

	name := strings.TrimSpace(req.Name)
	unit := strings.TrimSpace(req.Unit)
	description := strings.TrimSpace(req.Description)

	if name == "" {
		return nil, errors.New("name is required")
	}
	if unit == "" {
		return nil, errors.New("unit is required")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if req.MinimumQuantity < 0 {
		return nil, errors.New("minimumQuantity cannot be negative")
	}

	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if code == "" {
		nextCode, nextErr := uc.warehouseRepo.NextItemCode(category)
		if nextErr != nil {
			return nil, errors.New("failed to generate item code: " + nextErr.Error())
		}
		code = nextCode
	}

	if err := validateItemCode(code, category); err != nil {
		return nil, err
	}

	if _, findErr := uc.warehouseRepo.GetWarehouseItemByCode(code); findErr == nil {
		return nil, errors.New("item code already exists")
	} else if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to verify item code uniqueness: " + findErr.Error())
	}

	item := &entities.WarehouseItem{
		ID:              uuid.New().String(),
		Code:            code,
		Name:            name,
		Description:     description,
		Quantity:        0,
		MinimumQuantity: req.MinimumQuantity,
		Unit:            unit,
		Category:        category,
	}

	createdItem, err := uc.warehouseRepo.CreateWarehouseItem(item)
	if err != nil {
		return nil, errors.New("failed to create warehouse item: " + err.Error())
	}

	if _, err := uc.createTransactionRecord(createdItem, constants.TransactionTypeAddItem, req.Quantity, operatorUser); err != nil {
		if rollbackErr := uc.warehouseRepo.DeleteWarehouseItem(createdItem.ID); rollbackErr != nil {
			log.Printf("[WARN] failed to rollback warehouse item after transaction creation failure: %v", rollbackErr)
		}
		return nil, errors.New("failed to create initial stock transaction: " + err.Error())
	}

	newValue, _ := json.Marshal(createdItem)
	uc.createAuditLog(userID, audit_constants.AuditActionInsert, "warehouse_items", createdItem.ID, "", string(newValue))

	return createdItem, nil
}

func (uc *WarehouseUseCaseImpl) UpdateWarehouseItemByID(id string, req models.UpdateWarehouseItemRequest, userID string) (*entities.WarehouseItem, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	operatorUser, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get operator user: " + err.Error())
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("item id is required")
	}

	existing, err := uc.warehouseRepo.GetWarehouseItemByID(id)
	if err != nil {
		return nil, errors.New("warehouse item not found")
	}

	oldValue, _ := json.Marshal(existing)
	itemChanged := false
	quantityAdjustmentType := ""
	quantityAdjustmentAmount := 0

	if req.Category != nil {
		category := strings.ToUpper(strings.TrimSpace(*req.Category))
		if !isValidCategory(category) {
			return nil, errors.New("category must be one of MED, EQU, CON")
		}
		if existing.Category != category {
			existing.Category = category
			itemChanged = true
		}
	}

	if req.Code != nil {
		code := strings.ToUpper(strings.TrimSpace(*req.Code))
		if code == "" {
			return nil, errors.New("code cannot be empty")
		}

		if err := validateItemCode(code, existing.Category); err != nil {
			return nil, err
		}

		if !strings.EqualFold(code, existing.Code) {
			if _, findErr := uc.warehouseRepo.GetWarehouseItemByCode(code); findErr == nil {
				return nil, errors.New("item code already exists")
			} else if !errors.Is(findErr, gorm.ErrRecordNotFound) {
				return nil, errors.New("failed to verify item code uniqueness: " + findErr.Error())
			}
		}

		if existing.Code != code {
			existing.Code = code
			itemChanged = true
		}
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}
		if existing.Name != name {
			existing.Name = name
			itemChanged = true
		}
	}

	if req.Description != nil {
		description := strings.TrimSpace(*req.Description)
		if existing.Description != description {
			existing.Description = description
			itemChanged = true
		}
	}

	if req.Unit != nil {
		unit := strings.TrimSpace(*req.Unit)
		if unit == "" {
			return nil, errors.New("unit cannot be empty")
		}
		if existing.Unit != unit {
			existing.Unit = unit
			itemChanged = true
		}
	}

	if req.MinimumQuantity != nil {
		if *req.MinimumQuantity < 0 {
			return nil, errors.New("minimumQuantity cannot be negative")
		}
		if existing.MinimumQuantity != *req.MinimumQuantity {
			existing.MinimumQuantity = *req.MinimumQuantity
			itemChanged = true
		}
	}

	if req.Quantity != nil {
		if *req.Quantity < 0 {
			return nil, errors.New("quantity cannot be negative")
		}

		if *req.Quantity != existing.Quantity {
			if *req.Quantity > existing.Quantity {
				quantityAdjustmentType = constants.TransactionTypeRestock
				quantityAdjustmentAmount = *req.Quantity - existing.Quantity
			} else {
				quantityAdjustmentType = constants.TransactionTypeWithdraw
				quantityAdjustmentAmount = existing.Quantity - *req.Quantity
			}
		}
	}

	updatedItem := existing
	if itemChanged {
		updatedItem, err = uc.warehouseRepo.UpdateWarehouseItem(existing)
		if err != nil {
			return nil, errors.New("failed to update warehouse item: " + err.Error())
		}

		newValue, _ := json.Marshal(updatedItem)
		uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "warehouse_items", updatedItem.ID, string(oldValue), string(newValue))
	}

	if quantityAdjustmentAmount > 0 {
		if _, txErr := uc.createTransactionRecord(updatedItem, quantityAdjustmentType, quantityAdjustmentAmount, operatorUser); txErr != nil {
			if itemChanged {
				log.Printf("warning: warehouse item %s updated successfully but failed to create quantity adjustment transaction: %v", updatedItem.ID, txErr)
				return updatedItem, nil
			}
			return nil, errors.New("failed to create warehouse quantity adjustment transaction: " + txErr.Error())
		}
	}

	return updatedItem, nil
}

func (uc *WarehouseUseCaseImpl) DeleteWarehouseItemByID(id string, userID string) error {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return err
	}

	operatorUser, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get operator user: " + err.Error())
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("item id is required")
	}

	existing, err := uc.warehouseRepo.GetWarehouseItemByID(id)
	if err != nil {
		return errors.New("warehouse item not found")
	}

	if existing.Quantity <= 0 {
		return errors.New("cannot submit remove request for item with zero stock")
	}

	if _, txErr := uc.createTransactionRecord(existing, constants.TransactionTypeRemove, existing.Quantity, operatorUser); txErr != nil {
		return errors.New("failed to create warehouse remove transaction: " + txErr.Error())
	}

	return nil
}

func (uc *WarehouseUseCaseImpl) AdjustWarehouseItemByID(id string, req models.AdjustWarehouseItemRequest, userID string) (*entities.WarehouseItem, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	operatorUser, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get operator user: " + err.Error())
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("item id is required")
	}

	mode := strings.ToLower(strings.TrimSpace(req.Mode))
	if mode != constants.AdjustModeRestock && mode != constants.AdjustModeWithdraw {
		return nil, errors.New("mode must be either restock or withdraw")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	item, err := uc.warehouseRepo.GetWarehouseItemByID(id)
	if err != nil {
		return nil, errors.New("warehouse item not found")
	}

	txType := constants.TransactionTypeRestock
	if mode == constants.AdjustModeWithdraw {
		if req.Quantity > item.Quantity {
			return nil, errors.New("withdraw quantity exceeds available stock")
		}
		txType = constants.TransactionTypeWithdraw
	}

	if _, txErr := uc.createTransactionRecord(item, txType, req.Quantity, operatorUser); txErr != nil {
		return nil, errors.New("failed to create warehouse adjustment transaction: " + txErr.Error())
	}

	return item, nil
}

func (uc *WarehouseUseCaseImpl) GetWarehouseTransactions(query models.ListWarehouseTransactionsQuery, userID string) ([]*models.WarehouseTransactionResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	query.StartDate = strings.TrimSpace(query.StartDate)
	query.EndDate = strings.TrimSpace(query.EndDate)
	query.SearchItem = strings.TrimSpace(query.SearchItem)
	query.SearchUser = strings.TrimSpace(query.SearchUser)
	query.Status = strings.TrimSpace(query.Status)
	query.Type = strings.TrimSpace(query.Type)

	if query.StartDate != "" {
		if _, err := time.Parse("2006-01-02", query.StartDate); err != nil {
			return nil, errors.New("startDate must be in YYYY-MM-DD format")
		}
	}
	if query.EndDate != "" {
		if _, err := time.Parse("2006-01-02", query.EndDate); err != nil {
			return nil, errors.New("endDate must be in YYYY-MM-DD format")
		}
	}
	if query.Status != "" && !isValidApprovalStatus(query.Status) {
		return nil, errors.New("status must be one of รออนุมัติ, อนุมัติ, ไม่อนุมัติ")
	}
	if query.Type != "" && !isValidTransactionType(query.Type) {
		return nil, errors.New("type must be one of เพิ่มสินค้าใหม่, เติมสินค้า, เบิกสินค้า, นำออก")
	}

	transactions, err := uc.warehouseRepo.GetTransactions(query)
	if err != nil {
		return nil, errors.New("failed to get warehouse transactions: " + err.Error())
	}

	result := make([]*models.WarehouseTransactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		result = append(result, mapTransactionResponse(transaction))
	}

	return result, nil
}

func (uc *WarehouseUseCaseImpl) GetWarehouseTransactionByID(id string, userID string) (*models.WarehouseTransactionResponse, error) {
	if err := uc.ensureMedicalStaff(userID); err != nil {
		return nil, err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("transaction id is required")
	}

	transaction, err := uc.warehouseRepo.GetTransactionByID(id)
	if err != nil {
		return nil, errors.New("warehouse transaction not found")
	}

	return mapTransactionResponse(transaction), nil
}

func (uc *WarehouseUseCaseImpl) ApproveTransactions(req models.ApproveTransactionsRequest, userID string) ([]*models.WarehouseTransactionResponse, error) {
	if err := uc.ensureWarehouseTransactionManager(userID); err != nil {
		return nil, err
	}

	if len(req.TransactionIDs) == 0 {
		return nil, errors.New("transactionIds is required")
	}

	approverUser, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get approver user: " + err.Error())
	}

	transactions, err := uc.warehouseRepo.GetTransactionsByIDs(req.TransactionIDs)
	if err != nil {
		return nil, errors.New("failed to retrieve transactions: " + err.Error())
	}

	if len(transactions) != len(req.TransactionIDs) {
		return nil, errors.New("some transactionIds are invalid")
	}

	now := bangkokNow()
	approverName := buildDisplayName(approverUser)
	updatedTransactions := make([]*entities.WarehouseTransaction, 0, len(transactions))

	for _, transaction := range transactions {
		if transaction.ApprovalStatus != constants.ApprovalStatusPending {
			continue
		}

		if applyErr := uc.applyApprovedTransactionEffect(transaction, userID); applyErr != nil {
			return nil, errors.New("failed to apply approved transaction effect: " + applyErr.Error())
		}

		oldValue, _ := json.Marshal(transaction)

		transaction.ApprovalStatus = constants.ApprovalStatusApproved
		transaction.ApprovedByUserID = &userID
		transaction.ApprovedBy = &approverName
		approvedAt := now
		transaction.ApprovedAt = &approvedAt
		transaction.RejectedByUserID = nil
		transaction.RejectedBy = nil
		transaction.RejectedAt = nil
		transaction.RejectionReason = nil

		updatedTransaction, updateErr := uc.warehouseRepo.UpdateTransaction(transaction)
		if updateErr != nil {
			return nil, errors.New("failed to approve transaction: " + updateErr.Error())
		}

		newValue, _ := json.Marshal(updatedTransaction)
		uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "warehouse_transactions", updatedTransaction.ID, string(oldValue), string(newValue))
		updatedTransactions = append(updatedTransactions, updatedTransaction)
	}

	if len(updatedTransactions) == 0 {
		return nil, errors.New("no pending transactions were selected")
	}

	sort.Slice(updatedTransactions, func(i, j int) bool {
		return updatedTransactions[i].CreatedAt.After(updatedTransactions[j].CreatedAt)
	})

	response := make([]*models.WarehouseTransactionResponse, 0, len(updatedTransactions))
	for _, transaction := range updatedTransactions {
		response = append(response, mapTransactionResponse(transaction))
	}

	return response, nil
}

func (uc *WarehouseUseCaseImpl) applyApprovedTransactionEffect(transaction *entities.WarehouseTransaction, userID string) error {
	if transaction == nil {
		return errors.New("transaction is required")
	}

	if transaction.Type == constants.TransactionTypeRemove {
		if transaction.ItemID == nil || strings.TrimSpace(*transaction.ItemID) == "" {
			return nil
		}

		item, err := uc.warehouseRepo.GetWarehouseItemByID(strings.TrimSpace(*transaction.ItemID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return errors.New("failed to get warehouse item: " + err.Error())
		}

		oldItemValue, _ := json.Marshal(item)
		if err := uc.warehouseRepo.DeleteWarehouseItem(item.ID); err != nil {
			return errors.New("failed to remove warehouse item: " + err.Error())
		}

		uc.createAuditLog(userID, audit_constants.AuditActionDelete, "warehouse_items", item.ID, string(oldItemValue), "")
		return nil
	}

	if transaction.ItemID == nil || strings.TrimSpace(*transaction.ItemID) == "" {
		return errors.New("transaction item reference is missing")
	}

	item, err := uc.warehouseRepo.GetWarehouseItemByID(strings.TrimSpace(*transaction.ItemID))
	if err != nil {
		return errors.New("failed to get warehouse item: " + err.Error())
	}

	oldItemValue, _ := json.Marshal(item)

	switch transaction.Type {
	case constants.TransactionTypeAddItem, constants.TransactionTypeRestock:
		item.Quantity += transaction.Quantity
	case constants.TransactionTypeWithdraw:
		if transaction.Quantity > item.Quantity {
			return errors.New("withdraw quantity exceeds available stock")
		}
		item.Quantity -= transaction.Quantity
	default:
		return errors.New("unsupported transaction type")
	}

	updatedItem, err := uc.warehouseRepo.UpdateWarehouseItem(item)
	if err != nil {
		return errors.New("failed to update warehouse item: " + err.Error())
	}

	newItemValue, _ := json.Marshal(updatedItem)
	uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "warehouse_items", updatedItem.ID, string(oldItemValue), string(newItemValue))

	return nil
}

func (uc *WarehouseUseCaseImpl) RejectTransactions(req models.RejectTransactionsRequest, userID string) ([]*models.WarehouseTransactionResponse, error) {
	if err := uc.ensureWarehouseTransactionManager(userID); err != nil {
		return nil, err
	}

	if len(req.TransactionIDs) == 0 {
		return nil, errors.New("transactionIds is required")
	}

	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return nil, errors.New("reason is required")
	}

	rejectorUser, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get rejector user: " + err.Error())
	}

	transactions, err := uc.warehouseRepo.GetTransactionsByIDs(req.TransactionIDs)
	if err != nil {
		return nil, errors.New("failed to retrieve transactions: " + err.Error())
	}

	if len(transactions) != len(req.TransactionIDs) {
		return nil, errors.New("some transactionIds are invalid")
	}

	now := bangkokNow()
	rejectorName := buildDisplayName(rejectorUser)
	updatedTransactions := make([]*entities.WarehouseTransaction, 0, len(transactions))

	for _, transaction := range transactions {
		if transaction.ApprovalStatus != constants.ApprovalStatusPending {
			continue
		}

		oldValue, _ := json.Marshal(transaction)

		transaction.ApprovalStatus = constants.ApprovalStatusRejected
		transaction.ApprovedByUserID = nil
		transaction.ApprovedBy = nil
		transaction.ApprovedAt = nil
		transaction.RejectedByUserID = &userID
		transaction.RejectedBy = &rejectorName
		rejectedAt := now
		transaction.RejectedAt = &rejectedAt
		transaction.RejectionReason = &reason

		updatedTransaction, updateErr := uc.warehouseRepo.UpdateTransaction(transaction)
		if updateErr != nil {
			return nil, errors.New("failed to reject transaction: " + updateErr.Error())
		}

		newValue, _ := json.Marshal(updatedTransaction)
		uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "warehouse_transactions", updatedTransaction.ID, string(oldValue), string(newValue))
		updatedTransactions = append(updatedTransactions, updatedTransaction)
	}

	if len(updatedTransactions) == 0 {
		return nil, errors.New("no pending transactions were selected")
	}

	sort.Slice(updatedTransactions, func(i, j int) bool {
		return updatedTransactions[i].CreatedAt.After(updatedTransactions[j].CreatedAt)
	})

	response := make([]*models.WarehouseTransactionResponse, 0, len(updatedTransactions))
	for _, transaction := range updatedTransactions {
		response = append(response, mapTransactionResponse(transaction))
	}

	return response, nil
}

func (uc *WarehouseUseCaseImpl) createTransactionRecord(item *entities.WarehouseItem, txType string, quantity int, operator *entities.User) (*entities.WarehouseTransaction, error) {
	if item == nil {
		return nil, errors.New("item is required")
	}

	if operator == nil {
		return nil, errors.New("operator user is required")
	}

	now := bangkokNow()
	code, err := uc.warehouseRepo.NextTransactionCode(now.Format("20060102"))
	if err != nil {
		return nil, err
	}

	transaction := &entities.WarehouseTransaction{
		ID:             uuid.New().String(),
		Code:           code,
		Type:           txType,
		ItemID:         &item.ID,
		ItemCode:       item.Code,
		ItemName:       item.Name,
		Quantity:       quantity,
		OperatorUserID: operator.ID,
		Operator:       buildDisplayName(operator),
		ApprovalStatus: constants.ApprovalStatusPending,
	}

	return uc.warehouseRepo.CreateTransaction(transaction)
}

func (uc *WarehouseUseCaseImpl) ensureMedicalStaff(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	role, err := uc.userRepo.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}

	if role.Name != user_constants.RoleMedicalStaff && role.Name != user_constants.RoleSuperUser && role.Name != user_constants.RoleAdmin {
		return errors.New("only users with 'Medical Staff', 'Super User', or 'Admin' role can manage warehouse data")
	}

	return nil
}

func (uc *WarehouseUseCaseImpl) ensureWarehouseItemCreator(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	role, err := uc.userRepo.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}

	if role.Name != user_constants.RoleSuperUser && role.Name != user_constants.RoleAdmin {
		return errors.New("only users with 'Super User' or 'Admin' role can create warehouse items")
	}

	return nil
}

func (uc *WarehouseUseCaseImpl) ensureWarehouseTransactionManager(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	role, err := uc.userRepo.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}

	if role.Name != user_constants.RoleSuperUser && role.Name != user_constants.RoleAdmin {
		return errors.New("only users with 'Super User' or 'Admin' role can approve or reject warehouse transactions")
	}

	return nil
}

func (uc *WarehouseUseCaseImpl) createAuditLog(userID string, action string, tableName string, recordID string, oldValue string, newValue string) {
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: tableName,
		RecordID:  recordID,
		UserID:    userID,
		Action:    action,
		OldValue:  oldValue,
		NewValue:  newValue,
	}

	if _, err := uc.auditLogRepo.CreateAuditLog(auditLog); err != nil {
		log.Printf("[ERROR] Failed to create audit log for %s %s: %v", tableName, recordID, err)
	}
}

func isValidCategory(category string) bool {
	return category == constants.CategoryMedical ||
		category == constants.CategoryEquipment ||
		category == constants.CategoryConsumer
}

func isValidTransactionType(txType string) bool {
	return txType == constants.TransactionTypeAddItem ||
		txType == constants.TransactionTypeRestock ||
		txType == constants.TransactionTypeWithdraw ||
		txType == constants.TransactionTypeRemove
}

func isValidApprovalStatus(status string) bool {
	return status == constants.ApprovalStatusPending ||
		status == constants.ApprovalStatusApproved ||
		status == constants.ApprovalStatusRejected
}

func validateItemCode(code string, category string) error {
	pattern := regexp.MustCompile(`^[A-Z]{3}-\d{3}$`)
	if !pattern.MatchString(code) {
		return errors.New("code must be in format XXX-000")
	}

	prefix := strings.Split(code, "-")[0]
	if prefix != category {
		return fmt.Errorf("code prefix must match category (%s)", category)
	}

	return nil
}

func mapTransactionResponse(tx *entities.WarehouseTransaction) *models.WarehouseTransactionResponse {
	response := &models.WarehouseTransactionResponse{
		ID:              tx.ID,
		Code:            tx.Code,
		Type:            tx.Type,
		ItemCode:        tx.ItemCode,
		ItemName:        tx.ItemName,
		Quantity:        tx.Quantity,
		Operator:        tx.Operator,
		Date:            formatDateTime(tx.CreatedAt),
		ApprovalStatus:  tx.ApprovalStatus,
		ApprovedBy:      tx.ApprovedBy,
		RejectedBy:      tx.RejectedBy,
		RejectionReason: tx.RejectionReason,
	}

	if tx.ApprovedAt != nil {
		value := formatDateTime(*tx.ApprovedAt)
		response.ApprovedAt = &value
	}

	if tx.RejectedAt != nil {
		value := formatDateTime(*tx.RejectedAt)
		response.RejectedAt = &value
	}

	return response
}

func buildDisplayName(user *entities.User) string {
	if user == nil {
		return ""
	}

	fullName := strings.TrimSpace(strings.TrimSpace(user.FirstName) + " " + strings.TrimSpace(user.LastName))
	if fullName != "" {
		return fullName
	}

	nickName := strings.TrimSpace(user.Nickname)
	if nickName != "" {
		return nickName
	}

	return strings.TrimSpace(user.Username)
}

func formatDateTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.In(getBangkokLocation()).Format("02/01/2006 15:04")
}

func bangkokNow() time.Time {
	return time.Now().In(getBangkokLocation())
}

func getBangkokLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Local
	}
	return location
}
