package usecases_test

import (
	"testing"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	auditRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	emrConstants "github.com/aikidoaikido115/New-Acis-BE/modules/emr/constants"
	emrModels "github.com/aikidoaikido115/New-Acis-BE/modules/emr/models"
	emrRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	emrUsecases "github.com/aikidoaikido115/New-Acis-BE/modules/emr/usecases"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	medicineUsecases "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/usecases"
	userConstants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	userRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/stretchr/testify/assert"
)

type fakeEmrCoreUserRepo struct {
	*userRepositories.GormUserRepository

	user  *entities.User
	role  *entities.Role
	staff *entities.Staff

	getUserCalls  int
	getRoleCalls  int
	getStaffCalls int

	getUserErr  error
	getRoleErr  error
	getStaffErr error
}

func newFakeEmrCoreUserRepo(roleName string) *fakeEmrCoreUserRepo {
	return &fakeEmrCoreUserRepo{
		GormUserRepository: userRepositories.NewGormUserRepository(nil),
		user: &entities.User{
			ID:     "user-1",
			RoleID: "role-1",
		},
		role: &entities.Role{
			ID:   "role-1",
			Name: roleName,
		},
		staff: &entities.Staff{
			ID:     "staff-1",
			UserID: "user-1",
		},
	}
}

func (f *fakeEmrCoreUserRepo) GetUserByID(id string) (*entities.User, error) {
	f.getUserCalls++
	if f.getUserErr != nil {
		return nil, f.getUserErr
	}
	return f.user, nil
}

func (f *fakeEmrCoreUserRepo) GetRoleByID(roleID string) (*entities.Role, error) {
	f.getRoleCalls++
	if f.getRoleErr != nil {
		return nil, f.getRoleErr
	}
	return f.role, nil
}

func (f *fakeEmrCoreUserRepo) GetStaffByUserID(userID string) (*entities.Staff, error) {
	f.getStaffCalls++
	if f.getStaffErr != nil {
		return nil, f.getStaffErr
	}
	return f.staff, nil
}

type fakeEmrCoreAuditRepo struct {
	*auditRepositories.GormAuditLogRepository
	createAuditLogCalls int
}

func newFakeEmrCoreAuditRepo() *fakeEmrCoreAuditRepo {
	return &fakeEmrCoreAuditRepo{GormAuditLogRepository: auditRepositories.NewGormAuditLogRepository(nil)}
}

func (f *fakeEmrCoreAuditRepo) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	f.createAuditLogCalls++
	return auditLog, nil
}

type fakeEmrCoreRepo struct {
	*emrRepositories.GormEmrRepository

	roomExistsResponse  bool
	idCardExists        bool
	residentExists      bool
	vitalSignSlotExist  bool
	laboratorySlotExist bool

	getResidentByIDFirst  *entities.Resident
	getResidentByIDSecond *entities.Resident
	createdResident       *entities.Resident
	updatedResident       *entities.Resident

	existingResidentLabels []*entities.ResidentLabels
	createdLabelByName     map[string]*entities.IntakeLabels
	createdResidentLabels  []*entities.ResidentLabels

	createdVitalSign *entities.VitalSign
	existingVital    *entities.VitalSign
	updatedVital     *entities.VitalSign

	createdLab  *entities.LaboratoryValue
	existingLab *entities.LaboratoryValue
	updatedLab  *entities.LaboratoryValue

	capturedCreatedResident *entities.Resident
	capturedUpdatedResident *entities.Resident
	capturedCreateVitalSign *entities.VitalSign
	capturedUpdateVitalSign *entities.VitalSign
	capturedCreateLab       *entities.LaboratoryValue
	capturedUpdateLab       *entities.LaboratoryValue

	getResidentByIDCalls           int
	createResidentCalls            int
	updateResidentCalls            int
	roomExistsCalls                int
	idCardNumberExistsCalls        int
	deleteResidentLabelsCalls      int
	createIntakeLabelCalls         int
	createResidentLabelCalls       int
	vitalSignSlotExistsCalls       int
	createVitalSignCalls           int
	getVitalSignByIDCalls          int
	updateVitalSignByIDCalls       int
	laboratorySlotExistsCalls      int
	createLaboratoryValueCalls     int
	getLaboratoryValueByIDCalls    int
	updateLaboratoryValueByIDCalls int
}

func newFakeEmrCoreRepo() *fakeEmrCoreRepo {
	return &fakeEmrCoreRepo{
		GormEmrRepository:  emrRepositories.NewGormEmrRepository(nil),
		roomExistsResponse: true,
		residentExists:     true,
		createdLabelByName: map[string]*entities.IntakeLabels{},
	}
}

func (f *fakeEmrCoreRepo) RoomExists(id string) (bool, error) {
	f.roomExistsCalls++
	return f.roomExistsResponse, nil
}

func (f *fakeEmrCoreRepo) IdCardNumberExists(idCardNumber string) (bool, error) {
	f.idCardNumberExistsCalls++
	return f.idCardExists, nil
}

func (f *fakeEmrCoreRepo) CreateResident(resident *entities.Resident) (*entities.Resident, error) {
	f.createResidentCalls++
	copied := *resident
	f.capturedCreatedResident = &copied
	if f.createdResident != nil {
		return f.createdResident, nil
	}
	return resident, nil
}

func (f *fakeEmrCoreRepo) GetResidentByID(id string) (*entities.Resident, error) {
	f.getResidentByIDCalls++
	if f.getResidentByIDCalls == 1 && f.getResidentByIDFirst != nil {
		return f.getResidentByIDFirst, nil
	}
	if f.getResidentByIDSecond != nil {
		return f.getResidentByIDSecond, nil
	}
	return f.getResidentByIDFirst, nil
}

func (f *fakeEmrCoreRepo) UpdateResident(resident *entities.Resident) (*entities.Resident, error) {
	f.updateResidentCalls++
	copied := *resident
	f.capturedUpdatedResident = &copied
	if f.updatedResident != nil {
		return f.updatedResident, nil
	}
	return resident, nil
}

func (f *fakeEmrCoreRepo) GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error) {
	return f.existingResidentLabels, nil
}

func (f *fakeEmrCoreRepo) DeleteResidentLabelsByResidentID(residentID string) error {
	f.deleteResidentLabelsCalls++
	return nil
}

func (f *fakeEmrCoreRepo) LabelExists(labelName string) (bool, error) {
	_, exists := f.createdLabelByName[labelName]
	return exists, nil
}

func (f *fakeEmrCoreRepo) GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error) {
	return f.createdLabelByName[labelName], nil
}

func (f *fakeEmrCoreRepo) CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error) {
	f.createIntakeLabelCalls++
	created := *label
	f.createdLabelByName[label.LabelName] = &created
	return &created, nil
}

func (f *fakeEmrCoreRepo) CreateIntakeLabelByResidentID(residentLabel *entities.ResidentLabels) (*entities.ResidentLabels, error) {
	f.createResidentLabelCalls++
	copied := *residentLabel
	f.createdResidentLabels = append(f.createdResidentLabels, &copied)
	return &copied, nil
}

func (f *fakeEmrCoreRepo) ResidentExists(id string) (bool, error) {
	return f.residentExists, nil
}

func (f *fakeEmrCoreRepo) VitalSignSlotExists(residentID string, measurementDate time.Time, timeOfDay string) (bool, error) {
	f.vitalSignSlotExistsCalls++
	return f.vitalSignSlotExist, nil
}

func (f *fakeEmrCoreRepo) CreateVitalSign(vitalSign *entities.VitalSign) (*entities.VitalSign, error) {
	f.createVitalSignCalls++
	copied := *vitalSign
	f.capturedCreateVitalSign = &copied
	if f.createdVitalSign != nil {
		return f.createdVitalSign, nil
	}
	return vitalSign, nil
}

func (f *fakeEmrCoreRepo) GetVitalSignByID(id string) (*entities.VitalSign, error) {
	f.getVitalSignByIDCalls++
	return f.existingVital, nil
}

func (f *fakeEmrCoreRepo) UpdateVitalSignByID(vitalSign *entities.VitalSign) (*entities.VitalSign, error) {
	f.updateVitalSignByIDCalls++
	copied := *vitalSign
	f.capturedUpdateVitalSign = &copied
	if f.updatedVital != nil {
		return f.updatedVital, nil
	}
	return vitalSign, nil
}

func (f *fakeEmrCoreRepo) LaboratoryValueSlotExists(residentID string, measurementDate time.Time, timeOfDay string) (bool, error) {
	f.laboratorySlotExistsCalls++
	return f.laboratorySlotExist, nil
}

func (f *fakeEmrCoreRepo) CreateLaboratoryValue(laboratoryValue *entities.LaboratoryValue) (*entities.LaboratoryValue, error) {
	f.createLaboratoryValueCalls++
	copied := *laboratoryValue
	f.capturedCreateLab = &copied
	if f.createdLab != nil {
		return f.createdLab, nil
	}
	return laboratoryValue, nil
}

func (f *fakeEmrCoreRepo) GetLaboratoryValueByID(id string) (*entities.LaboratoryValue, error) {
	f.getLaboratoryValueByIDCalls++
	return f.existingLab, nil
}

func (f *fakeEmrCoreRepo) UpdateLaboratoryValueByID(laboratoryValue *entities.LaboratoryValue) (*entities.LaboratoryValue, error) {
	f.updateLaboratoryValueByIDCalls++
	copied := *laboratoryValue
	f.capturedUpdateLab = &copied
	if f.updatedLab != nil {
		return f.updatedLab, nil
	}
	return laboratoryValue, nil
}

type fakeDrugUsecaseNoop struct {
	*medicineUsecases.DrugUseCaseImpl
}

func newEmrCoreUsecase(roleName string) (*emrUsecases.EmrUseCaseImpl, *fakeEmrCoreRepo, *fakeEmrCoreUserRepo, *fakeEmrCoreAuditRepo) {
	emrRepo := newFakeEmrCoreRepo()
	userRepo := newFakeEmrCoreUserRepo(roleName)
	auditRepo := newFakeEmrCoreAuditRepo()
	drugUsecase := &fakeDrugUsecaseNoop{DrugUseCaseImpl: &medicineUsecases.DrugUseCaseImpl{}}

	uc := emrUsecases.NewEmrUseCase(
		emrRepo,
		auditRepo,
		userRepo,
		drugUsecase,
		configs.Supabase{},
	).(*emrUsecases.EmrUseCaseImpl)

	return uc, emrRepo, userRepo, auditRepo
}

// func strPtr(v string) *string {
// 	return &v
// }

// func i16Ptr(v int16) *int16 {
// 	return &v
// }

func f64Ptr(v float64) *float64 {
	return &v
}

func TestCreateResident_Success(t *testing.T) {
	uc, emrRepo, userRepo, auditRepo := newEmrCoreUsecase(userConstants.RoleMedicalStaff)

	roomID := "room-1"
	idCard := "1234567890123"
	resident := &entities.Resident{
		RoomID:       &roomID,
		FirstName:    "Alice",
		LastName:     "Smith",
		DateOfBirth:  time.Date(1950, 1, 10, 0, 0, 0, 0, time.UTC),
		Gender:       "Male",
		IdCardNumber: &idCard,
		Status:       emrConstants.Active,
	}

	emrRepo.createdResident = &entities.Resident{ID: "resident-1", FirstName: "Alice", LastName: "Smith"}

	result, err := uc.CreateResident(resident, "user-1", nil)

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, "resident-1", result.ID)
	}
	assert.Equal(t, 1, userRepo.getUserCalls)
	assert.Equal(t, 1, userRepo.getRoleCalls)
	assert.Equal(t, 1, emrRepo.roomExistsCalls)
	assert.Equal(t, 1, emrRepo.idCardNumberExistsCalls)
	assert.Equal(t, 1, emrRepo.createResidentCalls)
	assert.Equal(t, 1, auditRepo.createAuditLogCalls)
	if assert.NotNil(t, emrRepo.capturedCreatedResident) {
		assert.Equal(t, "male", emrRepo.capturedCreatedResident.Gender)
		assert.Equal(t, roomID, *emrRepo.capturedCreatedResident.RoomID)
	}
}

func TestUpdateResidentByID_Success_WithLabelSyncAndDedup(t *testing.T) {
	uc, emrRepo, userRepo, auditRepo := newEmrCoreUsecase(userConstants.RoleMedicalStaff)

	emrRepo.getResidentByIDFirst = &entities.Resident{
		ID:        "resident-1",
		FirstName: "Old",
		LastName:  "Name",
		Gender:    "female",
		Status:    emrConstants.Active,
	}
	emrRepo.updatedResident = &entities.Resident{
		ID:        "resident-1",
		FirstName: "New",
		LastName:  "Name",
		Gender:    "female",
		Status:    emrConstants.Active,
	}
	emrRepo.getResidentByIDSecond = &entities.Resident{ID: "resident-1", FirstName: "New", LastName: "Name"}
	emrRepo.existingResidentLabels = []*entities.ResidentLabels{
		{ResidentID: "resident-1", LabelID: "old-label"},
	}

	newFirstName := "New"
	request := emrModels.UpdateResidentRequest{
		FirstName: &newFirstName,
		Labels: []emrModels.IntakeLabelRequest{
			{LabelName: "Diabetes"},
			{LabelName: "Diabetes"},
			{LabelName: "Low Sodium"},
		},
	}

	result, err := uc.UpdateResidentByID("resident-1", request, "user-1", nil)

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, "resident-1", result.ID)
	}
	assert.Equal(t, 1, userRepo.getUserCalls)
	assert.Equal(t, 1, userRepo.getRoleCalls)
	assert.Equal(t, 2, emrRepo.getResidentByIDCalls)
	assert.Equal(t, 1, emrRepo.updateResidentCalls)
	assert.Equal(t, 1, emrRepo.deleteResidentLabelsCalls)
	assert.Equal(t, 2, emrRepo.createIntakeLabelCalls)
	assert.Equal(t, 2, emrRepo.createResidentLabelCalls)
	assert.Equal(t, 4, auditRepo.createAuditLogCalls)
}

func TestCreateVitalSign_Success_NormalizeTimeOfDayAndAudit(t *testing.T) {
	uc, emrRepo, userRepo, auditRepo := newEmrCoreUsecase(userConstants.RoleMedicalStaff)

	emrRepo.residentExists = true
	emrRepo.vitalSignSlotExist = false
	temp := 36.8
	vital := &entities.VitalSign{
		ResidentID:    "resident-1",
		TimeOfDay:     "morning",
		Temperature:   &temp,
		BreathingRate: nil,
	}
	emrRepo.createdVitalSign = &entities.VitalSign{ID: "vs-1", ResidentID: "resident-1", CreatedByStaffID: "staff-1", TimeOfDay: "เช้า"}

	result, err := uc.CreateVitalSign(vital, "2026-05-01", "user-1")

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, "vs-1", result.ID)
	}
	assert.Equal(t, 1, userRepo.getUserCalls)
	assert.Equal(t, 1, userRepo.getRoleCalls)
	assert.Equal(t, 1, userRepo.getStaffCalls)
	assert.Equal(t, 1, emrRepo.vitalSignSlotExistsCalls)
	assert.Equal(t, 1, emrRepo.createVitalSignCalls)
	assert.Equal(t, 1, auditRepo.createAuditLogCalls)
	if assert.NotNil(t, emrRepo.capturedCreateVitalSign) {
		assert.Equal(t, "เช้า", emrRepo.capturedCreateVitalSign.TimeOfDay)
		assert.Equal(t, "staff-1", emrRepo.capturedCreateVitalSign.CreatedByStaffID)
		assert.Equal(t, "2026-05-01", emrRepo.capturedCreateVitalSign.MeasurementDate.Format("2006-01-02"))
	}
}

func TestCreateVitalSign_Error_WhenSlotExists(t *testing.T) {
	uc, emrRepo, _, _ := newEmrCoreUsecase(userConstants.RoleMedicalStaff)
	emrRepo.vitalSignSlotExist = true
	temp := 36.5

	_, err := uc.CreateVitalSign(&entities.VitalSign{
		ResidentID:  "resident-1",
		TimeOfDay:   "เช้า",
		Temperature: &temp,
	}, "2026-05-01", "user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vital sign already exists")
	assert.Equal(t, 0, emrRepo.createVitalSignCalls)
}

func TestUpdateVitalSignByID_Success(t *testing.T) {
	uc, emrRepo, _, auditRepo := newEmrCoreUsecase(userConstants.RoleMedicalStaff)

	oldTemp := 36.0
	newTemp := 38.1
	emrRepo.existingVital = &entities.VitalSign{
		ID:          "vs-1",
		ResidentID:  "resident-1",
		Temperature: &oldTemp,
		TimeOfDay:   "เช้า",
	}
	emrRepo.updatedVital = &entities.VitalSign{
		ID:          "vs-1",
		ResidentID:  "resident-1",
		Temperature: &newTemp,
		TimeOfDay:   "เช้า",
	}

	result, err := uc.UpdateVitalSignByID("vs-1", &entities.VitalSign{Temperature: &newTemp}, "user-1")

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, "vs-1", result.ID)
	}
	assert.Equal(t, 1, emrRepo.getVitalSignByIDCalls)
	assert.Equal(t, 1, emrRepo.updateVitalSignByIDCalls)
	assert.Equal(t, 1, auditRepo.createAuditLogCalls)
	if assert.NotNil(t, emrRepo.capturedUpdateVitalSign) {
		if assert.NotNil(t, emrRepo.capturedUpdateVitalSign.Temperature) {
			assert.InDelta(t, 38.1, *emrRepo.capturedUpdateVitalSign.Temperature, 0.001)
		}
	}
}

func TestCreateLaboratoryValue_Success_NormalizeTimeOfDayAndAudit(t *testing.T) {
	uc, emrRepo, userRepo, auditRepo := newEmrCoreUsecase(userConstants.RoleMedicalStaff)

	bloodGlucose := 130.0
	lab := &entities.LaboratoryValue{
		ResidentID:   "resident-1",
		BloodGlucose: &bloodGlucose,
	}
	emrRepo.createdLab = &entities.LaboratoryValue{ID: "lab-1", ResidentID: "resident-1", TimeOfDay: "กลางคืน", CreatedByStaffID: "staff-1"}

	result, err := uc.CreateLaboratoryValue(lab, "01/05/2026", "night", "user-1")

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, "lab-1", result.ID)
	}
	assert.Equal(t, 1, userRepo.getUserCalls)
	assert.Equal(t, 1, userRepo.getRoleCalls)
	assert.Equal(t, 1, userRepo.getStaffCalls)
	assert.Equal(t, 1, emrRepo.laboratorySlotExistsCalls)
	assert.Equal(t, 1, emrRepo.createLaboratoryValueCalls)
	assert.Equal(t, 1, auditRepo.createAuditLogCalls)
	if assert.NotNil(t, emrRepo.capturedCreateLab) {
		assert.Equal(t, "กลางคืน", emrRepo.capturedCreateLab.TimeOfDay)
		assert.Equal(t, "staff-1", emrRepo.capturedCreateLab.CreatedByStaffID)
		assert.Equal(t, "2026-05-01", emrRepo.capturedCreateLab.MeasurementDate.Format("2006-01-02"))
	}
}

func TestUpdateLaboratoryValueByID_Error_WhenUrinePairIncomplete(t *testing.T) {
	uc, emrRepo, _, _ := newEmrCoreUsecase(userConstants.RoleMedicalStaff)
	emrRepo.existingLab = &entities.LaboratoryValue{ID: "lab-1", ResidentID: "resident-1", BloodGlucose: f64Ptr(100)}

	_, err := uc.UpdateLaboratoryValueByID("lab-1", &entities.LaboratoryValue{
		UrineOutput: f64Ptr(200),
	}, "user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "urine_output and urine_type must be provided together")
	assert.Equal(t, 0, emrRepo.updateLaboratoryValueByIDCalls)
}

func TestUpdateLaboratoryValueByID_Success(t *testing.T) {
	uc, emrRepo, _, auditRepo := newEmrCoreUsecase(userConstants.RoleMedicalStaff)

	oldBG := 110.0
	newBG := 145.0
	emrRepo.existingLab = &entities.LaboratoryValue{
		ID:           "lab-1",
		ResidentID:   "resident-1",
		BloodGlucose: &oldBG,
	}
	emrRepo.updatedLab = &entities.LaboratoryValue{
		ID:           "lab-1",
		ResidentID:   "resident-1",
		BloodGlucose: &newBG,
	}

	result, err := uc.UpdateLaboratoryValueByID("lab-1", &entities.LaboratoryValue{BloodGlucose: &newBG}, "user-1")

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, "lab-1", result.ID)
	}
	assert.Equal(t, 1, emrRepo.getLaboratoryValueByIDCalls)
	assert.Equal(t, 1, emrRepo.updateLaboratoryValueByIDCalls)
	assert.Equal(t, 1, auditRepo.createAuditLogCalls)
	if assert.NotNil(t, emrRepo.capturedUpdateLab) {
		if assert.NotNil(t, emrRepo.capturedUpdateLab.BloodGlucose) {
			assert.InDelta(t, 145.0, *emrRepo.capturedUpdateLab.BloodGlucose, 0.001)
		}
	}
}
