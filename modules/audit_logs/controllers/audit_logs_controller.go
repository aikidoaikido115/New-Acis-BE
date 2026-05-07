package controllers

import (
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/usecases"
	"github.com/gofiber/fiber/v2"
)

type AuditLogController struct {
	auditLogUsecase usecases.AuditLogUsecase
}

func NewAuditLogController(auditLogUsecase usecases.AuditLogUsecase) *AuditLogController {
	return &AuditLogController{auditLogUsecase: auditLogUsecase}
}

// @Summary Get all audit logs
// @Description Get all audit logs. Admin only.
// @Tags Audit Logs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=array} "Audit logs retrieved successfully"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 403 {object} object{status=string,status_code=int,message=string,result=any} "Forbidden - Admin only"
// @Router /api/admin/audit-logs [get]
func (c *AuditLogController) GetAuditLogsHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": fiber.ErrUnauthorized.Message, "status_code": fiber.StatusUnauthorized, "message": "Unauthorized: Missing user ID", "result": nil})
	}

	auditLogs, err := c.auditLogUsecase.GetAllAuditLogs(userID)
	if err != nil {
		if strings.Contains(err.Error(), "Admin") {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": fiber.ErrForbidden.Message, "status_code": fiber.StatusForbidden, "message": err.Error(), "result": nil})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": fiber.ErrInternalServerError.Message, "status_code": fiber.StatusInternalServerError, "message": err.Error(), "result": nil})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "status_code": fiber.StatusOK, "message": "audit logs retrieved successfully", "result": auditLogs})
}

// @Summary Search audit logs
// @Description Search audit logs by table name, record id, user id, action, old value, or new value. Admin only.
// @Tags Audit Logs
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search text"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=array} "Audit logs retrieved successfully"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 403 {object} object{status=string,status_code=int,message=string,result=any} "Forbidden - Admin only"
// @Router /api/admin/audit-logs/search [get]
func (c *AuditLogController) SearchAuditLogsHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": fiber.ErrUnauthorized.Message, "status_code": fiber.StatusUnauthorized, "message": "Unauthorized: Missing user ID", "result": nil})
	}

	search := ctx.Query("search")
	auditLogs, err := c.auditLogUsecase.SearchAuditLogs(search, userID)
	if err != nil {
		if strings.Contains(err.Error(), "Admin") {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": fiber.ErrForbidden.Message, "status_code": fiber.StatusForbidden, "message": err.Error(), "result": nil})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": fiber.ErrInternalServerError.Message, "status_code": fiber.StatusInternalServerError, "message": err.Error(), "result": nil})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "status_code": fiber.StatusOK, "message": "audit logs retrieved successfully", "result": auditLogs})
}

// @Summary Get audit log by ID
// @Description Get a single audit log by ID. Admin only.
// @Tags Audit Logs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audit Log ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Audit log retrieved successfully"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 403 {object} object{status=string,status_code=int,message=string,result=any} "Forbidden - Admin only"
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any} "Audit log not found"
// @Router /api/admin/audit-logs/{id} [get]
func (c *AuditLogController) GetAuditLogByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": fiber.ErrUnauthorized.Message, "status_code": fiber.StatusUnauthorized, "message": "Unauthorized: Missing user ID", "result": nil})
	}

	auditLogID := ctx.Params("id")
	auditLog, err := c.auditLogUsecase.GetAuditLogByID(auditLogID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "Admin") {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": fiber.ErrForbidden.Message, "status_code": fiber.StatusForbidden, "message": err.Error(), "result": nil})
		}
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": fiber.ErrNotFound.Message, "status_code": fiber.StatusNotFound, "message": err.Error(), "result": nil})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": fiber.ErrInternalServerError.Message, "status_code": fiber.StatusInternalServerError, "message": err.Error(), "result": nil})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Success", "status_code": fiber.StatusOK, "message": "audit log retrieved successfully", "result": auditLog})
}
