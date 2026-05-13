package usecases

import (
	"errors"
	"testing"
	"time"

	audit_repository "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repository "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/constants"
	"github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/models"
	warehouse_repository "github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fakeWarehouseP04UserRepo struct {
	userByID map[string]*entities.User
	roleByID map[string]*entities.Role
}

func (f *fakeWarehouseP04UserRepo) CreateUser(user *entities.User) (*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) CreateStaff(user *entities.User, staff *entities.Staff) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) CreateStaffFile(staffFile *entities.StaffsFiles) (*entities.StaffsFiles, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetUserByEmail(email string) (*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetUserByID(id string) (*entities.User, error) {
	if user, ok := f.userByID[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}
func (f *fakeWarehouseP04UserRepo) GetUsersByFirstAndLastName(firstName string, lastName string) ([]*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetStaffByID(id string) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetStaffByUserID(userID string) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetStaffFileByID(id string) (*entities.StaffsFiles, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetUserByUsername(username string) (*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetRoleByName(roleName string) (*entities.Role, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetRoleByID(roleID string) (*entities.Role, error) {
	if role, ok := f.roleByID[roleID]; ok {
		return role, nil
	}
	return nil, errors.New("role not found")
}
func (f *fakeWarehouseP04UserRepo) UsernameExists(username string) (bool, error) {
	return false, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) EmailExists(email string) (bool, error) {
	return false, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetAllUsers() ([]*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) UpdateUserByID(user *entities.User) error {
	return errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) UpdateUserApprovalByID(userID string, isApprove bool) error {
	return errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) DeleteStaffAndUserByStaffID(staffID string) error {
	return errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) DeleteRelativeAndUserByUserID(userID string) error {
	return errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) DeleteUserByID(userID string) error {
	return errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetRelativeUserByUserID(userID string) (*user_repository.AdminRelativeUser, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetRelativeUsersWithResident() ([]user_repository.AdminRelativeUser, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetStaffIDMapByUserIDs(userIDs []string) (map[string]string, error) {
	return map[string]string{}, nil
}
func (f *fakeWarehouseP04UserRepo) CreateOTP(otp *entities.OTP) error { return errors.New("not used") }
func (f *fakeWarehouseP04UserRepo) GetOTPByUserID(userID string) (*entities.OTP, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) DeleteOTP(userID string) error { return errors.New("not used") }
func (f *fakeWarehouseP04UserRepo) StoreResetToken(temptoken *entities.TempToken) error {
	return errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) GetResetToken(userID string) (string, error) {
	return "", errors.New("not used")
}
func (f *fakeWarehouseP04UserRepo) DeleteResetToken(userID string) error {
	return errors.New("not used")
}

type fakeWarehouseP04AuditRepo struct{ calls int }

func (f *fakeWarehouseP04AuditRepo) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	f.calls++
	return auditLog, nil
}
func (f *fakeWarehouseP04AuditRepo) GetAllAuditLogs() ([]*entities.AuditLogs, error) { return nil, nil }
func (f *fakeWarehouseP04AuditRepo) SearchAuditLogs(search string) ([]*entities.AuditLogs, error) {
	return nil, nil
}
func (f *fakeWarehouseP04AuditRepo) GetAuditLogByID(id string) (*entities.AuditLogs, error) {
	return nil, nil
}

type fakeWarehouseP04Repo struct {
	warehouse_repository.WarehouseRepository
	itemsByID        map[string]*entities.WarehouseItem
	transactionsByID map[string]*entities.WarehouseTransaction
	updatedItems     []*entities.WarehouseItem
	updatedTx        []*entities.WarehouseTransaction
	deleteItemIDs    []string
}

func (f *fakeWarehouseP04Repo) GetWarehouseItems(search string, category string) ([]*entities.WarehouseItem, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04Repo) GetWarehouseItemByID(id string) (*entities.WarehouseItem, error) {
	if item, ok := f.itemsByID[id]; ok {
		return item, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeWarehouseP04Repo) GetWarehouseItemByCode(code string) (*entities.WarehouseItem, error) {
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeWarehouseP04Repo) CreateWarehouseItem(item *entities.WarehouseItem) (*entities.WarehouseItem, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04Repo) UpdateWarehouseItem(item *entities.WarehouseItem) (*entities.WarehouseItem, error) {
	f.updatedItems = append(f.updatedItems, item)
	return item, nil
}
func (f *fakeWarehouseP04Repo) DeleteWarehouseItem(id string) error {
	f.deleteItemIDs = append(f.deleteItemIDs, id)
	return nil
}
func (f *fakeWarehouseP04Repo) NextItemCode(category string) (string, error) {
	return "", errors.New("not used")
}
func (f *fakeWarehouseP04Repo) CreateTransaction(transaction *entities.WarehouseTransaction) (*entities.WarehouseTransaction, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04Repo) GetTransactions(query models.ListWarehouseTransactionsQuery) ([]*entities.WarehouseTransaction, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04Repo) GetTransactionByID(id string) (*entities.WarehouseTransaction, error) {
	return nil, errors.New("not used")
}
func (f *fakeWarehouseP04Repo) GetTransactionsByIDs(ids []string) ([]*entities.WarehouseTransaction, error) {
	result := make([]*entities.WarehouseTransaction, 0, len(ids))
	for _, id := range ids {
		if tx, ok := f.transactionsByID[id]; ok {
			result = append(result, tx)
		}
	}
	return result, nil
}
func (f *fakeWarehouseP04Repo) UpdateTransaction(transaction *entities.WarehouseTransaction) (*entities.WarehouseTransaction, error) {
	f.updatedTx = append(f.updatedTx, transaction)
	return transaction, nil
}
func (f *fakeWarehouseP04Repo) NextTransactionCode(datePrefix string) (string, error) {
	return "", errors.New("not used")
}
func (f *fakeWarehouseP04Repo) GetWarehouseItemsCount() (int64, error) {
	return 0, errors.New("not used")
}
func (f *fakeWarehouseP04Repo) GetLowStockItemsCount(threshold int) (int64, error) {
	return 0, errors.New("not used")
}
func (f *fakeWarehouseP04Repo) GetPendingTransactionsCountByType(transactionType string) (int64, error) {
	return 0, errors.New("not used")
}

func TestApproveTransactions_RestockAndWithdraw_Success(t *testing.T) {
	userRepo := &fakeWarehouseP04UserRepo{
		userByID: map[string]*entities.User{
			"admin-1": {ID: "admin-1", RoleID: "role-admin", FirstName: "Ada", LastName: "Min"},
		},
		roleByID: map[string]*entities.Role{
			"role-admin": {ID: "role-admin", Name: user_constants.RoleAdmin},
		},
	}
	itemID := "item-1"
	baseTime := time.Date(2026, 5, 3, 8, 0, 0, 0, time.UTC)
	repo := &fakeWarehouseP04Repo{
		itemsByID: map[string]*entities.WarehouseItem{
			itemID: {ID: itemID, Code: "MED-001", Name: "Mask", Quantity: 10, MinimumQuantity: 1, Unit: "pcs", Category: constants.CategoryMedical},
		},
		transactionsByID: map[string]*entities.WarehouseTransaction{
			"tx-1": {ID: "tx-1", Type: constants.TransactionTypeRestock, ItemID: &itemID, ItemCode: "MED-001", ItemName: "Mask", Quantity: 5, ApprovalStatus: constants.ApprovalStatusPending, CreatedAt: baseTime},
			"tx-2": {ID: "tx-2", Type: constants.TransactionTypeWithdraw, ItemID: &itemID, ItemCode: "MED-001", ItemName: "Mask", Quantity: 3, ApprovalStatus: constants.ApprovalStatusPending, CreatedAt: baseTime.Add(time.Minute)},
		},
	}
	auditRepo := &fakeWarehouseP04AuditRepo{}
	uc := NewWarehouseUseCase(repo, auditRepo, userRepo)

	result, err := uc.ApproveTransactions(models.ApproveTransactionsRequest{TransactionIDs: []string{"tx-1", "tx-2"}}, "admin-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, len(repo.updatedTx))
	assert.Equal(t, 12, repo.itemsByID[itemID].Quantity)
	assert.Equal(t, 4, auditRepo.calls)
	assert.Equal(t, constants.ApprovalStatusApproved, repo.transactionsByID["tx-1"].ApprovalStatus)
	assert.Equal(t, constants.ApprovalStatusApproved, repo.transactionsByID["tx-2"].ApprovalStatus)
}

func TestApproveTransactions_ErrorWhenWithdrawExceedsStock(t *testing.T) {
	userRepo := &fakeWarehouseP04UserRepo{
		userByID: map[string]*entities.User{
			"admin-1": {ID: "admin-1", RoleID: "role-admin", FirstName: "Ada", LastName: "Min"},
		},
		roleByID: map[string]*entities.Role{
			"role-admin": {ID: "role-admin", Name: user_constants.RoleAdmin},
		},
	}
	itemID := "item-1"
	repo := &fakeWarehouseP04Repo{
		itemsByID: map[string]*entities.WarehouseItem{
			itemID: {ID: itemID, Code: "MED-001", Name: "Mask", Quantity: 2, MinimumQuantity: 1, Unit: "pcs", Category: constants.CategoryMedical},
		},
		transactionsByID: map[string]*entities.WarehouseTransaction{
			"tx-1": {ID: "tx-1", Type: constants.TransactionTypeWithdraw, ItemID: &itemID, ItemCode: "MED-001", ItemName: "Mask", Quantity: 5, ApprovalStatus: constants.ApprovalStatusPending, CreatedAt: time.Now()},
		},
	}
	uc := NewWarehouseUseCase(repo, &fakeWarehouseP04AuditRepo{}, userRepo)

	result, err := uc.ApproveTransactions(models.ApproveTransactionsRequest{TransactionIDs: []string{"tx-1"}}, "admin-1")

	assert.Nil(t, result)
	assert.EqualError(t, err, "failed to apply approved transaction effect: withdraw quantity exceeds available stock")
}

var _ user_repository.UserRepository = (*fakeWarehouseP04UserRepo)(nil)
var _ audit_repository.AuditLogRepository = (*fakeWarehouseP04AuditRepo)(nil)
