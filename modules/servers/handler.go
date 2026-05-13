package servers

import (
	"fmt"
	"log"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	activityController "github.com/aikidoaikido115/New-Acis-BE/modules/activity/controllers"
	activityRepository "github.com/aikidoaikido115/New-Acis-BE/modules/activity/repositories"
	activityUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/activity/usecases"

	auditLogController "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/controllers"
	auditLogUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/usecases"
	mealController "github.com/aikidoaikido115/New-Acis-BE/modules/meal/controllers"
	mealRepository "github.com/aikidoaikido115/New-Acis-BE/modules/meal/repositories"
	mealUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/meal/usecases"
	medicineController "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/controllers"
	medicineRepository "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/repositories"
	medicineUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/usecases"
	supportController "github.com/aikidoaikido115/New-Acis-BE/modules/support/controllers"
	supportRepository "github.com/aikidoaikido115/New-Acis-BE/modules/support/repositories"
	supportUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/support/usecases"
	userController "github.com/aikidoaikido115/New-Acis-BE/modules/user/controllers"
	userRepository "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	userUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/user/usecases"
	warehouseController "github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/controllers"
	warehouseRepository "github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/repositories"
	warehouseUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/usecases"

	auditLogRepository "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"

	emrController "github.com/aikidoaikido115/New-Acis-BE/modules/emr/controllers"
	emrRepository "github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	emrUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/emr/usecases"

	aiinfra "github.com/aikidoaikido115/New-Acis-BE/pkg/ai"
	"github.com/aikidoaikido115/New-Acis-BE/pkg/database"
	"github.com/aikidoaikido115/New-Acis-BE/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"gorm.io/gorm"
)

func SetupServer(server configs.Server, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               server.AppName,
		BodyLimit:             1024 * 1024 * 1024, // 1GB limit
		DisableStartupMessage: true,
		ReduceMemoryUsage:     true,    // ลดการจองหน่วยความจำ
		Concurrency:           1000000, // ปรับจำนวน concurrent requests สูงสุด

	})

	setupMiddlewares(app, server.CORS)
	setupRoutes(app, server, jwt, supa, mail)

	return app
}

func setupMiddlewares(app *fiber.App, cor configs.CORS) {
	// Recovery middleware - จับ panic และแปลงเป็น 500 error
	app.Use(recover.New())

	// Logger middleware - บันทึก request/response
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cor.AllowOrigins,
		AllowMethods:     cor.AllowMethods,
		AllowHeaders:     cor.AllowHeaders,
		AllowCredentials: cor.AllowCredentials,
	}))
}

func setupRoutes(app *fiber.App, server configs.Server, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) {
	// Initialize database connection
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	SetupUserRoutes(app, db, jwt, supa, mail)
	SetupEmrRoutes(app, db, jwt, supa)
	SetupMedicineRoutes(app, db, jwt)
	SetupMealRoutes(app, db, jwt)
	SetupActivityRoutes(app, db, jwt, supa)
	SetupWarehouseRoutes(app, db, jwt)
	SetupSupportRoutes(app, db, jwt)
	SetupAuditLogRoutes(app, db, jwt)

	// API group
	api := app.Group("/api")

	api.Get("/hello", func(c *fiber.Ctx) error {
		message := fmt.Sprintf("Misha Necron มา สำรวจ ภาษา Go แล้ว %s จะถูกสร้างในไม่ช้า", server.AppName)

		return c.JSON(fiber.Map{
			"message": message,
		})
	})
}

func SetupUserRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) {

	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	userRepository := userRepository.NewGormUserRepository(db)
	userUsecase := userUsecase.NewUserUseCase(userRepository, auditLogRepository, jwt, supa, mail)
	userController := userController.NewUserController(userUsecase)

	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", userController.RegisterHandler)
	authGroup.Post("/login", userController.LoginHandler)
	authGroup.Post("/forgotpassword", userController.ForgotPasswordHandler)
	authGroup.Post("/forgotpassword/otp", userController.VerifyOTPHandler)
	authGroup.Patch("/forgotpassword/changepassword", userController.ChangePasswordHandler)
	authGroup.Patch("/resetpassword", middlewares.JWTMiddleware(jwt), userController.ResetPasswordHandler)
	authGroup.Post("/logout", middlewares.JWTMiddleware(jwt), userController.LogoutHandler)

	userGroup := app.Group("/api/user")
	userGroup.Get("/", middlewares.JWTMiddleware(jwt), userController.GetUserByIDHandler)
	userGroup.Get("/search", middlewares.JWTMiddleware(jwt), userController.GetUsersByFirstAndLastNameHandler)
	userGroup.Patch("/", middlewares.JWTMiddleware(jwt), userController.UpdateUserByIDHandler)
	userGroup.Post("/staff/files", middlewares.JWTMiddleware(jwt), userController.CreateStaffFileHandler)

	adminGroup := app.Group("/api/admin")
	adminGroup.Get("/users", middlewares.JWTMiddleware(jwt), userController.GetAllUsersHandler)
	adminGroup.Get("/users/relatives", middlewares.JWTMiddleware(jwt), userController.GetRelativeUsersHandler)
	adminGroup.Patch("/users/:user_id/approval", middlewares.JWTMiddleware(jwt), userController.UpdateUserApprovalHandler)
	adminGroup.Patch("/users/staffs/:staff_id/role", middlewares.JWTMiddleware(jwt), userController.UpdateStaffRoleByIDHandler)
	adminGroup.Delete("/users/staffs/:staff_id", middlewares.JWTMiddleware(jwt), userController.DeleteStaffByIDHandler)
	adminGroup.Delete("/users/relatives/:user_id", middlewares.JWTMiddleware(jwt), userController.DeleteRelativeByUserIDHandler)
	adminGroup.Delete("/users/:user_id", middlewares.JWTMiddleware(jwt), userController.DeleteUserByIDHandler)
}

func SetupEmrRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	userRepository := userRepository.NewGormUserRepository(db)

	emrRepository := emrRepository.NewGormEmrRepository(db)
	drugRepository := medicineRepository.NewGormDrugRepository(db)
	drugUsecase := medicineUsecase.NewDrugUseCase(drugRepository, auditLogRepository, userRepository, emrRepository)
	emrUsecase := emrUsecase.NewEmrUseCase(emrRepository, auditLogRepository, userRepository, drugUsecase, supa, jwt)
	emrController := emrController.NewEmrController(emrUsecase)

	residentGroup := app.Group("/api/emr/residents")
	residentGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateResidentHandler)
	residentGroup.Get("/all", middlewares.JWTMiddleware(jwt), emrController.GetAllResidentsHandler)
	residentGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetResidentOverviewHandler)
	residentGroup.Get("/:id", middlewares.JWTMiddleware(jwt), emrController.GetResidentByIDHandler)
	residentGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetResidentByRoomIDHandler)
	residentGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateResidentByIDHandler)

	roomGroup := app.Group("/api/emr/rooms")
	roomGroup.Get("/:id", middlewares.JWTMiddleware(jwt), emrController.GetRoomByIDHandler)
	roomGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetAllRoomsHandler)
	roomGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateRoomHandler)
	roomGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateRoomByIDHandler)

	dashboardGroup := app.Group("/api/emr/dashboard")
	dashboardGroup.Get("/residents", middlewares.JWTMiddleware(jwt), emrController.GetNumberOfResidentsDashboardHandler)
	dashboardGroup.Get("/resident-gender-stats", middlewares.JWTMiddleware(jwt), emrController.GetResidentGenderStatsDashboardHandler)
	dashboardGroup.Get("/vital-sign-stats", middlewares.JWTMiddleware(jwt), emrController.GetVitalSignStatsDashboardHandler)
	dashboardGroup.Get("/drug-plan-time-of-day-stats", middlewares.JWTMiddleware(jwt), emrController.GetDrugPlanTimeOfDayStatsDashboardHandler)
	dashboardGroup.Get("/resident-allergy-stats", middlewares.JWTMiddleware(jwt), emrController.GetResidentAllergyStatsDashboardHandler)
	dashboardGroup.Get("/resident-drug-allergy-stats", middlewares.JWTMiddleware(jwt), emrController.GetResidentDrugAllergyStatsDashboardHandler)

	labelGroup := app.Group("/api/emr/intake-labels")
	labelGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetResidentLabelsByResidentIDHandler)
	labelGroup.Get("/all", middlewares.JWTMiddleware(jwt), emrController.GetAllIntakeLabelsHandler)
	labelGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateIntakeLabelByResidentIDHandler)
	labelGroup.Post("/master", middlewares.JWTMiddleware(jwt), emrController.CreateIntakeLabelMasterHandler)

	allergyGroup := app.Group("/api/emr/allergies")
	allergyGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetResidentAllergiesByResidentIDHandler)
	allergyGroup.Get("/all", middlewares.JWTMiddleware(jwt), emrController.GetAllAllergiesHandler)
	allergyGroup.Get("/residents/all", middlewares.JWTMiddleware(jwt), emrController.GetAllResidentAllergiesHandler)
	allergyGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateAllergyByResidentIDHandler)

	drugAllergyGroup := app.Group("/api/emr/drug-allergies")
	drugAllergyGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetResidentDrugAllergiesByResidentIDHandler)
	drugAllergyGroup.Get("/all", middlewares.JWTMiddleware(jwt), emrController.GetAllDrugAllergiesHandler)
	drugAllergyGroup.Get("/residents/all", middlewares.JWTMiddleware(jwt), emrController.GetAllResidentDrugAllergiesHandler)
	drugAllergyGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateDrugAllergyByResidentIDHandler)

	vitalSignGroup := app.Group("/api/emr/vital-signs")
	vitalSignGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateVitalSignHandler)
	vitalSignGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetVitalSignsOverviewHandler)
	vitalSignGroup.Get("/resident", middlewares.JWTMiddleware(jwt), emrController.GetVitalSignsByResidentHandler)
	vitalSignGroup.Get("/room", middlewares.JWTMiddleware(jwt), emrController.GetRoomVitalSignsHandler)
	vitalSignGroup.Get("/history/:resident_id", middlewares.JWTMiddleware(jwt), emrController.GetVitalSignsHistoryHandler)
	vitalSignGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateVitalSignByIDHandler)

	laboratoryValueGroup := app.Group("/api/emr/laboratory-values")
	laboratoryValueGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateLaboratoryValueHandler)
	laboratoryValueGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetLaboratoryValuesOverviewHandler)
	laboratoryValueGroup.Get("/resident", middlewares.JWTMiddleware(jwt), emrController.GetLaboratoryValuesByResidentHandler)
	laboratoryValueGroup.Get("/room", middlewares.JWTMiddleware(jwt), emrController.GetRoomLaboratoryValuesHandler)
	laboratoryValueGroup.Get("/history/:resident_id", middlewares.JWTMiddleware(jwt), emrController.GetLaboratoryValuesHistoryHandler)
	laboratoryValueGroup.Get("/urine-output-sum/:resident_id", middlewares.JWTMiddleware(jwt), emrController.GetUrineOutputSumByResidentIDHandler)
	laboratoryValueGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateLaboratoryValueByIDHandler)

	nurseNoteGroup := app.Group("/api/emr/nurse-notes")
	nurseNoteGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateNurseNoteHandler)
	nurseNoteGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetNurseNotesOverviewHandler)
	nurseNoteGroup.Get("/resident/all", middlewares.JWTMiddleware(jwt), emrController.GetNurseNotesByResidentHandler)
	nurseNoteGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateNurseNoteByIDHandler)
	nurseNoteGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), emrController.DeleteNurseNoteByIDHandler)

	woundCareNoteGroup := app.Group("/api/emr/wound-care-notes")
	woundCareNoteGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateWoundCareNoteHandler)
	woundCareNoteGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetWoundCareNotesOverviewHandler)
	woundCareNoteGroup.Get("/resident/all", middlewares.JWTMiddleware(jwt), emrController.GetWoundCareNotesByResidentHandler)
	woundCareNoteGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateWoundCareNoteByIDHandler)
	woundCareNoteGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), emrController.DeleteWoundCareNoteByIDHandler)

	relativeNoteGroup := app.Group("/api/emr/relative-notes")
	relativeNoteGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateRelativeNoteHandler)
	relativeNoteGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetRelativeNotesOverviewHandler)
	relativeNoteGroup.Get("/resident/all", middlewares.JWTMiddleware(jwt), emrController.GetRelativeNotesByResidentHandler)
	relativeNoteGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateRelativeNoteByIDHandler)
	relativeNoteGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), emrController.DeleteRelativeNoteByIDHandler)

	relativePortalStaffGroup := app.Group("/api/emr/relatives")
	relativePortalStaffGroup.Get("/magic-link", middlewares.JWTMiddleware(jwt), emrController.GetRelativeMagicLinkHandler)
	relativePortalStaffGroup.Post("/magic-link/issue", middlewares.JWTMiddleware(jwt), emrController.IssueRelativeMagicLinkHandler)

	relativePortalAuthGroup := app.Group("/api/relative/auth")
	relativePortalAuthGroup.Post("/login", emrController.RelativePortalLoginHandler)

	relativePortalGroup := app.Group("/api/relative")
	relativePortalGroup.Get("/dashboard", middlewares.JWTMiddleware(jwt), emrController.GetRelativeDashboardHandler)
	relativePortalGroup.Get("/patient-info", middlewares.JWTMiddleware(jwt), emrController.GetRelativePatientInfoHandler)

	doctorOrderGroup := app.Group("/api/emr/doctor-orders")
	doctorOrderGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateDoctorOrderHandler)
	doctorOrderGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetDoctorOrdersOverviewHandler)
	doctorOrderGroup.Get("/resident/all", middlewares.JWTMiddleware(jwt), emrController.GetDoctorOrdersByResidentHandler)
	doctorOrderGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateDoctorOrderByIDHandler)
	doctorOrderGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), emrController.DeleteDoctorOrderByIDHandler)
}

func SetupMealRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	userRepository := userRepository.NewGormUserRepository(db)
	emrRepo := emrRepository.NewGormEmrRepository(db)
	allergyAIClient := aiinfra.NewClientFromEnv(aiinfra.WithTimeout(10 * time.Second))

	mealRepository := mealRepository.NewGormMealRepository(db)
	mealUsecase := mealUsecase.NewMealUseCase(mealRepository, emrRepo, auditLogRepository, userRepository, allergyAIClient)
	mealController := mealController.NewMealController(mealUsecase)

	menuGroup := app.Group("/api/meals/menus")
	menuGroup.Post("/", middlewares.JWTMiddleware(jwt), mealController.CreateMenuHandler)
	menuGroup.Get("/", middlewares.JWTMiddleware(jwt), mealController.GetAllMenusHandler)
	menuGroup.Get("/:id", middlewares.JWTMiddleware(jwt), mealController.GetMenuByIDHandler)
	menuGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), mealController.UpdateMenuHandler)

	mealPlanGroup := app.Group("/api/meals/meal-plans")
	mealPlanGroup.Post("/", middlewares.JWTMiddleware(jwt), mealController.CreateMealPlanHandler)
	mealPlanGroup.Post("/manual", middlewares.JWTMiddleware(jwt), mealController.CreateMealPlanManualHandler)
	mealPlanGroup.Get("/", middlewares.JWTMiddleware(jwt), mealController.GetAllMealPlansHandler)
	mealPlanGroup.Get("/history", middlewares.JWTMiddleware(jwt), mealController.GetMealHistoryHandler)
	mealPlanGroup.Get("/date", middlewares.JWTMiddleware(jwt), mealController.GetMealPlansTodayHandler)
	mealPlanGroup.Get("/:id", middlewares.JWTMiddleware(jwt), mealController.GetMealPlanByIDHandler)
	// mealPlanGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), mealController.UpdateMealPlanHandler)
}

func SetupMedicineRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	userRepository := userRepository.NewGormUserRepository(db)
	drugRepository := medicineRepository.NewGormDrugRepository(db)
	emrRepository := emrRepository.NewGormEmrRepository(db)
	drugUsecase := medicineUsecase.NewDrugUseCase(drugRepository, auditLogRepository, userRepository, emrRepository)
	drugController := medicineController.NewDrugController(drugUsecase)

	personalDrugGroup := app.Group("/api/emr/personal-drugs")
	personalDrugGroup.Post("/", middlewares.JWTMiddleware(jwt), drugController.CreatePersonalDrugHandler)
	personalDrugGroup.Get("/overview", middlewares.JWTMiddleware(jwt), drugController.GetPersonalDrugsOverviewHandler)
	personalDrugGroup.Get("/resident/all", middlewares.JWTMiddleware(jwt), drugController.GetPersonalDrugsByResidentHandler)
	personalDrugGroup.Get("/resident", middlewares.JWTMiddleware(jwt), drugController.GetPersonalDrugsByResidentTodayHandler)
	personalDrugGroup.Get("/:id", middlewares.JWTMiddleware(jwt), drugController.GetPersonalDrugByIDHandler)
	personalDrugGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), drugController.UpdatePersonalDrugByIDHandler)
	personalDrugGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), drugController.DeletePersonalDrugByIDHandler)

	drugMasterGroup := app.Group("/api/emr/drug-masters")
	drugMasterGroup.Post("/", middlewares.JWTMiddleware(jwt), drugController.CreateDrugMasterHandler)
	drugMasterGroup.Get("/", middlewares.JWTMiddleware(jwt), drugController.GetDrugMastersHandler)
	drugMasterGroup.Get("/:id", middlewares.JWTMiddleware(jwt), drugController.GetDrugMasterByIDHandler)
	drugMasterGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), drugController.UpdateDrugMasterByIDHandler)
	drugMasterGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), drugController.DeleteDrugMasterByIDHandler)

	drugPlanGroup := app.Group("/api/emr/drug-plans")
	drugPlanGroup.Post("/", middlewares.JWTMiddleware(jwt), drugController.CreateDrugPlanHandler)
	drugPlanGroup.Post("/generate-today", middlewares.JWTMiddleware(jwt), drugController.ForceGenerateTodayDrugPlansHandler)
	drugPlanGroup.Post("/generate-today/resident/:resident_id", middlewares.JWTMiddleware(jwt), drugController.ForceGenerateTodayDrugPlansByResidentHandler)
	drugPlanGroup.Get("/istaken-summary", middlewares.JWTMiddleware(jwt), drugController.GetDrugPlansTodayResidentSummaryHandler)
	drugPlanGroup.Get("/today", middlewares.JWTMiddleware(jwt), drugController.GetDrugPlansTodayHandler)
	drugPlanGroup.Get("/overview", middlewares.JWTMiddleware(jwt), drugController.GetDrugPlansOverviewHandler)
	drugPlanGroup.Get("/history", middlewares.JWTMiddleware(jwt), drugController.GetDrugAdministrationHistoryHandler)
	drugPlanGroup.Get("/resident/all", middlewares.JWTMiddleware(jwt), drugController.GetDrugPlansByResidentHandler)
	drugPlanGroup.Get("/resident", middlewares.JWTMiddleware(jwt), drugController.GetDrugPlansByResidentTodayHandler)
	drugPlanGroup.Patch("/resident/:resident_id/take", middlewares.JWTMiddleware(jwt), drugController.TakeDrugPlansByResidentTodayHandler)
	drugPlanGroup.Patch("/resident/:resident_id/omit", middlewares.JWTMiddleware(jwt), drugController.OmitDrugPlansByResidentTodayHandler)
	drugPlanGroup.Get("/", middlewares.JWTMiddleware(jwt), drugController.GetDrugPlansHandler)
	drugPlanGroup.Get("/:id", middlewares.JWTMiddleware(jwt), drugController.GetDrugPlanByIDHandler)
	drugPlanGroup.Patch("/:id/take", middlewares.JWTMiddleware(jwt), drugController.TakeDrugPlanByIDHandler)
	drugPlanGroup.Patch("/:id/omit", middlewares.JWTMiddleware(jwt), drugController.OmitDrugPlanByIDHandler)
	drugPlanGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), drugController.UpdateDrugPlanByIDHandler)
	drugPlanGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), drugController.DeleteDrugPlanByIDHandler)
}

func SetupActivityRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	activityRepository := activityRepository.NewGormActivityRepository(db)
	userRepository := userRepository.NewGormUserRepository(db)
	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	activityUsecase := activityUsecase.NewActivityUseCase(activityRepository, userRepository, auditLogRepository, supa)
	activityController := activityController.NewActivityController(activityUsecase)

	activityGroup := app.Group("/api/activities")
	activityGroup.Post("/", middlewares.JWTMiddleware(jwt), activityController.CreateActivityHandler)
	activityGroup.Get("/", middlewares.JWTMiddleware(jwt), activityController.GetAllActivitiesHandler)
	activityGroup.Get("/:id", middlewares.JWTMiddleware(jwt), activityController.GetActivityByIDHandler)
	activityGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), activityController.UpdateActivityByIDHandler)
	activityGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), activityController.DeleteActivityByIDHandler)

	activityScheduleGroup := app.Group("/api/activity-schedules")
	activityScheduleGroup.Post("/sync", middlewares.JWTMiddleware(jwt), activityController.CreateActivityScheduleWithActivitySyncHandler)
	activityScheduleGroup.Get("/sync", middlewares.JWTMiddleware(jwt), activityController.GetAllActivitySchedulesWithActivitySyncHandler)
	activityScheduleGroup.Get("/sync/:id", middlewares.JWTMiddleware(jwt), activityController.GetActivityScheduleWithActivitySyncByIDHandler)
	activityScheduleGroup.Patch("/sync/:id", middlewares.JWTMiddleware(jwt), activityController.UpdateActivityScheduleWithActivitySyncByIDHandler)
	activityScheduleGroup.Post("/", middlewares.JWTMiddleware(jwt), activityController.CreateActivityScheduleHandler)
	activityScheduleGroup.Get("/", middlewares.JWTMiddleware(jwt), activityController.GetAllActivitySchedulesHandler)
	activityScheduleGroup.Get("/:id/residents", middlewares.JWTMiddleware(jwt), activityController.GetResidentsByScheduleIDCustomHandler)
	activityScheduleGroup.Get("/:id", middlewares.JWTMiddleware(jwt), activityController.GetActivityScheduleByIDHandler)
	activityScheduleGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), activityController.UpdateActivityScheduleByIDHandler)
	activityScheduleGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), activityController.DeleteActivityScheduleByIDHandler)

	participationGroup := app.Group("/api/activity-participations")
	participationGroup.Post("/", middlewares.JWTMiddleware(jwt), activityController.CreateParticipationHandler)
	participationGroup.Get("/", middlewares.JWTMiddleware(jwt), activityController.GetAllParticipationsHandler)
	participationGroup.Patch("/is-participating/bulk", middlewares.JWTMiddleware(jwt), activityController.BulkUpdateParticipationIsParticipatingByResidentIDsHandler)
	participationGroup.Get("/:resident_id/:as_id", middlewares.JWTMiddleware(jwt), activityController.GetParticipationByResidentIDAndASIDHandler)
	participationGroup.Patch("/:resident_id/:as_id", middlewares.JWTMiddleware(jwt), activityController.UpdateParticipationByResidentIDAndASIDHandler)
	participationGroup.Delete("/:resident_id/:as_id", middlewares.JWTMiddleware(jwt), activityController.DeleteParticipationByResidentIDAndASIDHandler)
}
func SetupWarehouseRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	userRepository := userRepository.NewGormUserRepository(db)
	warehouseRepository := warehouseRepository.NewGormWarehouseRepository(db)
	warehouseUsecase := warehouseUsecase.NewWarehouseUseCase(warehouseRepository, auditLogRepository, userRepository)
	warehouseController := warehouseController.NewWarehouseController(warehouseUsecase)

	warehouseItemGroup := app.Group("/api/warehouse/items")
	warehouseItemGroup.Get("/", middlewares.JWTMiddleware(jwt), warehouseController.GetWarehouseItemsHandler)
	warehouseItemGroup.Post("/", middlewares.JWTMiddleware(jwt), warehouseController.CreateWarehouseItemHandler)
	warehouseItemGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), warehouseController.UpdateWarehouseItemByIDHandler)
	warehouseItemGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), warehouseController.DeleteWarehouseItemByIDHandler)
	warehouseItemGroup.Post("/:id/adjust", middlewares.JWTMiddleware(jwt), warehouseController.AdjustWarehouseItemByIDHandler)

	warehouseTransactionGroup := app.Group("/api/warehouse/transactions")
	warehouseTransactionGroup.Get("/", middlewares.JWTMiddleware(jwt), warehouseController.GetWarehouseTransactionsHandler)
	warehouseTransactionGroup.Get("/:id", middlewares.JWTMiddleware(jwt), warehouseController.GetWarehouseTransactionByIDHandler)
	warehouseTransactionGroup.Patch("/approve", middlewares.JWTMiddleware(jwt), warehouseController.ApproveTransactionsHandler)
	warehouseTransactionGroup.Patch("/reject", middlewares.JWTMiddleware(jwt), warehouseController.RejectTransactionsHandler)

	warehouseDashboardGroup := app.Group("/api/warehouse/dashboard")
	warehouseDashboardGroup.Get("/summary", middlewares.JWTMiddleware(jwt), warehouseController.GetWarehouseDashboardSummaryHandler)
}

func SetupSupportRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	userRepository := userRepository.NewGormUserRepository(db)
	supportRepository := supportRepository.NewGormSupportRepository(db)
	supportUsecase := supportUsecase.NewSupportUseCase(supportRepository, userRepository)
	supportController := supportController.NewSupportController(supportUsecase)

	supportGroup := app.Group("/api/support/tickets")
	supportGroup.Get("/", middlewares.JWTMiddleware(jwt), supportController.GetSupportTicketsHandler)
	supportGroup.Get("/:id", middlewares.JWTMiddleware(jwt), supportController.GetSupportTicketByIDHandler)
	supportGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), supportController.DeleteSupportTicketByIDHandler)
	supportGroup.Patch("/:id/status", middlewares.JWTMiddleware(jwt), supportController.UpdateSupportTicketStatusHandler)
	supportGroup.Post("/", middlewares.JWTMiddleware(jwt), supportController.CreateSupportTicketHandler)
}

func SetupAuditLogRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	userRepository := userRepository.NewGormUserRepository(db)
	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	auditLogUsecaseInstance := auditLogUsecase.NewAuditLogUseCase(auditLogRepository, userRepository)
	auditLogControllerInstance := auditLogController.NewAuditLogController(auditLogUsecaseInstance)

	auditLogGroup := app.Group("/api/admin/audit-logs")
	auditLogGroup.Get("/", middlewares.JWTMiddleware(jwt), auditLogControllerInstance.GetAuditLogsHandler)
	auditLogGroup.Get("/search", middlewares.JWTMiddleware(jwt), auditLogControllerInstance.SearchAuditLogsHandler)
	auditLogGroup.Get("/:id", middlewares.JWTMiddleware(jwt), auditLogControllerInstance.GetAuditLogByIDHandler)
}
