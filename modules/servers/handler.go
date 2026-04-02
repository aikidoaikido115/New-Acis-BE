package servers

import (
	"fmt"
	"log"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"

	mealController "github.com/aikidoaikido115/New-Acis-BE/modules/meal/controllers"
	mealRepository "github.com/aikidoaikido115/New-Acis-BE/modules/meal/repositories"
	mealUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/meal/usecases"
	userController "github.com/aikidoaikido115/New-Acis-BE/modules/user/controllers"
	userRepository "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	userUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/user/usecases"

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
	SetupEmrRoutes(app, db, jwt)
	SetupMealRoutes(app, db, jwt)

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
	userGroup.Patch("/", middlewares.JWTMiddleware(jwt), userController.UpdateUserByIDHandler)
	userGroup.Post("/staff/files", middlewares.JWTMiddleware(jwt), userController.CreateStaffFileHandler)
}

func SetupEmrRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	userRepository := userRepository.NewGormUserRepository(db)

	emrRepository := emrRepository.NewGormEmrRepository(db)
	emrUsecase := emrUsecase.NewEmrUseCase(emrRepository, auditLogRepository, userRepository)
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
	dashboardGroup.Get("/resident-allergy-stats", middlewares.JWTMiddleware(jwt), emrController.GetResidentAllergyStatsDashboardHandler)

	labelGroup := app.Group("/api/emr/intake-labels")
	labelGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetResidentLabelsByResidentIDHandler)
	labelGroup.Get("/all", middlewares.JWTMiddleware(jwt), emrController.GetAllIntakeLabelsHandler)
	labelGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateIntakeLabelByResidentIDHandler)

	allergyGroup := app.Group("/api/emr/allergies")
	allergyGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetResidentAllergiesByResidentIDHandler)
	allergyGroup.Get("/all", middlewares.JWTMiddleware(jwt), emrController.GetAllAllergiesHandler)
	allergyGroup.Get("/residents/all", middlewares.JWTMiddleware(jwt), emrController.GetAllResidentAllergiesHandler)
	allergyGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateAllergyByResidentIDHandler)

	vitalSignGroup := app.Group("/api/emr/vital-signs")
	vitalSignGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateVitalSignHandler)
	vitalSignGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetVitalSignsOverviewHandler)
	vitalSignGroup.Get("/resident", middlewares.JWTMiddleware(jwt), emrController.GetVitalSignsByResidentHandler)
	vitalSignGroup.Get("/room", middlewares.JWTMiddleware(jwt), emrController.GetRoomVitalSignsHandler)
	vitalSignGroup.Get("/history/:resident_id", middlewares.JWTMiddleware(jwt), emrController.GetVitalSignsHistoryHandler)
	vitalSignGroup.Get("/abnormal", middlewares.JWTMiddleware(jwt), emrController.GetAbnormalVitalSignsHandler)
	vitalSignGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateVitalSignByIDHandler)

	laboratoryValueGroup := app.Group("/api/emr/laboratory-values")
	laboratoryValueGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateLaboratoryValueHandler)
	laboratoryValueGroup.Get("/overview", middlewares.JWTMiddleware(jwt), emrController.GetLaboratoryValuesOverviewHandler)
	laboratoryValueGroup.Get("/resident", middlewares.JWTMiddleware(jwt), emrController.GetLaboratoryValuesByResidentHandler)
	laboratoryValueGroup.Get("/room", middlewares.JWTMiddleware(jwt), emrController.GetRoomLaboratoryValuesHandler)
	laboratoryValueGroup.Get("/history/:resident_id", middlewares.JWTMiddleware(jwt), emrController.GetLaboratoryValuesHistoryHandler)
	laboratoryValueGroup.Get("/urine-output-sum/:resident_id", middlewares.JWTMiddleware(jwt), emrController.GetUrineOutputSumByResidentIDHandler)
	laboratoryValueGroup.Get("/abnormal", middlewares.JWTMiddleware(jwt), emrController.GetAbnormalLaboratoryValuesHandler)
	laboratoryValueGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), emrController.UpdateLaboratoryValueByIDHandler)
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
	mealPlanGroup.Get("/", middlewares.JWTMiddleware(jwt), mealController.GetAllMealPlansHandler)
	mealPlanGroup.Get("/today", middlewares.JWTMiddleware(jwt), mealController.GetMealPlansTodayHandler)
	mealPlanGroup.Get("/:id", middlewares.JWTMiddleware(jwt), mealController.GetMealPlanByIDHandler)
	// mealPlanGroup.Patch("/:id", middlewares.JWTMiddleware(jwt), mealController.UpdateMealPlanHandler)
}
