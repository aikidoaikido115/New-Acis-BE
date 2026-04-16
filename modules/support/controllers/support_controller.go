package controllers

import (
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/support/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/support/usecases"
	"github.com/gofiber/fiber/v2"
)

type SupportController struct {
	supportUsecase usecases.SupportUsecase
}

func NewSupportController(supportUsecase usecases.SupportUsecase) *SupportController {
	return &SupportController{supportUsecase: supportUsecase}
}

func mapSupportError(err error) (int, string) {
	errMessage := err.Error()

	switch {
	case strings.Contains(errMessage, "only users with"):
		return fiber.ErrForbidden.Code, fiber.ErrForbidden.Message
	case strings.Contains(errMessage, "required"), strings.Contains(errMessage, "invalid"), strings.Contains(errMessage, "must be one of"):
		return fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message
	case strings.Contains(errMessage, "not found"):
		return fiber.ErrNotFound.Code, fiber.ErrNotFound.Message
	default:
		return fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message
	}
}

// CreateSupportTicketHandler godoc
// @Summary Create support ticket
// @Description Create a support ticket from medical or kitchen staff
// @Tags Support
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateSupportTicketRequest true "Support ticket payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.SupportTicket}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 403 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/support/tickets [post]
func (c *SupportController) CreateSupportTicketHandler(ctx *fiber.Ctx) error {
	var req models.CreateSupportTicketRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	ticket := &entities.SupportTicket{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
	}

	createdTicket, err := c.supportUsecase.CreateSupportTicket(ticket, userID)
	if err != nil {
		statusCode, statusText := mapSupportError(err)

		return ctx.Status(statusCode).JSON(fiber.Map{
			"status":      statusText,
			"status_code": statusCode,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "support ticket created successfully",
		"result":      createdTicket,
	})
}

// GetSupportTicketsHandler godoc
// @Summary Get support tickets
// @Description Get support tickets with optional status and search filters (Medical Staff only)
// @Tags Support
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by reporter name/email/subject/message"
// @Param status query string false "Ticket status (open, in_progress, resolved)"
// @Param reporterRole query string false "Reporter role (Medical Staff, Kitchen Staff)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]entities.SupportTicket}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 403 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/support/tickets [get]
func (c *SupportController) GetSupportTicketsHandler(ctx *fiber.Ctx) error {
	var query models.ListSupportTicketsQuery
	if err := ctx.QueryParser(&query); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	tickets, err := c.supportUsecase.GetSupportTickets(query, userID)
	if err != nil {
		statusCode, statusText := mapSupportError(err)

		return ctx.Status(statusCode).JSON(fiber.Map{
			"status":      statusText,
			"status_code": statusCode,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "support tickets retrieved successfully",
		"result":      tickets,
	})
}

// GetSupportTicketByIDHandler godoc
// @Summary Get support ticket by ID
// @Description Get support ticket detail by ID (Medical Staff only)
// @Tags Support
// @Produce json
// @Security BearerAuth
// @Param id path string true "Support ticket ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.SupportTicket}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 403 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/support/tickets/{id} [get]
func (c *SupportController) GetSupportTicketByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	ticket, err := c.supportUsecase.GetSupportTicketByID(id, userID)
	if err != nil {
		statusCode, statusText := mapSupportError(err)

		return ctx.Status(statusCode).JSON(fiber.Map{
			"status":      statusText,
			"status_code": statusCode,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "support ticket retrieved successfully",
		"result":      ticket,
	})
}

// UpdateSupportTicketStatusHandler godoc
// @Summary Update support ticket status
// @Description Update support ticket status (Medical Staff only)
// @Tags Support
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Support ticket ID"
// @Param request body models.UpdateSupportTicketStatusRequest true "Support ticket status payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.SupportTicket}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 403 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/support/tickets/{id}/status [patch]
func (c *SupportController) UpdateSupportTicketStatusHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req models.UpdateSupportTicketStatusRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	updatedTicket, err := c.supportUsecase.UpdateSupportTicketStatus(id, req, userID)
	if err != nil {
		statusCode, statusText := mapSupportError(err)

		return ctx.Status(statusCode).JSON(fiber.Map{
			"status":      statusText,
			"status_code": statusCode,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "support ticket status updated successfully",
		"result":      updatedTicket,
	})
}

// DeleteSupportTicketByIDHandler godoc
// @Summary Delete support ticket by ID
// @Description Delete support ticket by ID when status is resolved (Medical Staff only)
// @Tags Support
// @Produce json
// @Security BearerAuth
// @Param id path string true "Support ticket ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 403 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/support/tickets/{id} [delete]
func (c *SupportController) DeleteSupportTicketByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	err := c.supportUsecase.DeleteSupportTicketByID(id, userID)
	if err != nil {
		statusCode, statusText := mapSupportError(err)

		return ctx.Status(statusCode).JSON(fiber.Map{
			"status":      statusText,
			"status_code": statusCode,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "support ticket deleted successfully",
		"result":      nil,
	})
}
