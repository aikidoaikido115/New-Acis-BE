package servers

import (
	"fmt"
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/configs"

	userController "github.com/aikidoaikido115/New-Acis-BE/modules/user/controllers"
	userRepository "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	userUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/user/usecases"

	auditLogRepository "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"

	emrController "github.com/aikidoaikido115/New-Acis-BE/modules/emr/controllers"
	emrRepository "github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	emrUsecase "github.com/aikidoaikido115/New-Acis-BE/modules/emr/usecases"

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
	authGroup.Put("/forgotpassword/changepassword", userController.ChangePasswordHandler)
	authGroup.Put("/resetpassword", middlewares.JWTMiddleware(jwt), userController.ResetPasswordHandler)
	authGroup.Post("/logout", middlewares.JWTMiddleware(jwt), userController.LogoutHandler)

	userGroup := app.Group("/api/user")
	userGroup.Get("/", middlewares.JWTMiddleware(jwt), userController.GetUserByIDHandler)
	userGroup.Put("/", middlewares.JWTMiddleware(jwt), userController.UpdateUserByIDHandler)
	userGroup.Post("/staff/files", middlewares.JWTMiddleware(jwt), userController.CreateStaffFileHandler)
}

func SetupEmrRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	auditLogRepository := auditLogRepository.NewGormAuditLogRepository(db)
	emrRepository := emrRepository.NewGormEmrRepository(db)
	emrUsecase := emrUsecase.NewEmrUseCase(emrRepository, auditLogRepository)
	emrController := emrController.NewEmrController(emrUsecase)

	residentGroup := app.Group("/api/emr/residents")
	residentGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateResident)
	residentGroup.Get("/all", middlewares.JWTMiddleware(jwt), emrController.GetAllResidents)
	residentGroup.Get("/:id", middlewares.JWTMiddleware(jwt), emrController.GetResidentByID)
	residentGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetResidentByRoomID)

	roomGroup := app.Group("/api/emr/rooms")
	roomGroup.Get("/:id", middlewares.JWTMiddleware(jwt), emrController.GetRoomByID)
	roomGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetAllRooms)

	labelGroup := app.Group("/api/emr/intake-labels")
	labelGroup.Get("/", middlewares.JWTMiddleware(jwt), emrController.GetResidentLabelsByResidentID)
	labelGroup.Get("/all", middlewares.JWTMiddleware(jwt), emrController.GetAllIntakeLabels)
	labelGroup.Post("/", middlewares.JWTMiddleware(jwt), emrController.CreateIntakeLabelByResidentID)

}
