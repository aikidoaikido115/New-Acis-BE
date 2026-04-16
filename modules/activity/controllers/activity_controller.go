package controllers

import (
	"errors"
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/modules/activity/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/activity/usecases"
	"github.com/gofiber/fiber/v2"
)

type ActivityController struct {
	activityUsecase usecases.ActivityUsecase
}

func NewActivityController(activityUsecase usecases.ActivityUsecase) *ActivityController {
	return &ActivityController{activityUsecase: activityUsecase}
}

// CreateActivityHandler godoc
// @Summary Create Activity
// @Description Create a new activity record
// @Tags Activity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateActivityRequest true "Activity payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.Activity}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activities [post]
func (c *ActivityController) CreateActivityHandler(ctx *fiber.Ctx) error {
	var req models.CreateActivityRequest
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

	createdActivity, err := c.activityUsecase.CreateActivity(req, userID)
	if err != nil {
		return ctx.Status(resolveStatusCode(err)).JSON(fiber.Map{
			"status":      resolveStatusText(err),
			"status_code": resolveStatusCode(err),
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "activity created successfully",
		"result":      createdActivity,
	})
}

// GetActivityByIDHandler godoc
// @Summary Get Activity by ID
// @Description Retrieve activity by ID
// @Tags Activity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.Activity}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activities/{id} [get]
func (c *ActivityController) GetActivityByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	id := ctx.Params("id")
	activity, err := c.activityUsecase.GetActivityByID(id)
	if err != nil {
		return ctx.Status(resolveStatusCode(err)).JSON(fiber.Map{
			"status":      resolveStatusText(err),
			"status_code": resolveStatusCode(err),
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "activity retrieved successfully",
		"result":      activity,
	})
}

// GetAllActivitiesHandler godoc
// @Summary Get All Activities
// @Description Retrieve all activity records
// @Tags Activity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]entities.Activity}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activities [get]
func (c *ActivityController) GetAllActivitiesHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	activities, err := c.activityUsecase.GetAllActivities()
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "activities retrieved successfully",
		"result":      activities,
	})
}

// UpdateActivityByIDHandler godoc
// @Summary Update Activity by ID
// @Description Update an activity record
// @Tags Activity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity ID"
// @Param request body models.UpdateActivityRequest true "Activity update payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.Activity}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activities/{id} [patch]
func (c *ActivityController) UpdateActivityByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateActivityRequest
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

	id := ctx.Params("id")
	updatedActivity, err := c.activityUsecase.UpdateActivityByID(id, req)
	if err != nil {
		return ctx.Status(resolveStatusCode(err)).JSON(fiber.Map{
			"status":      resolveStatusText(err),
			"status_code": resolveStatusCode(err),
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "activity updated successfully",
		"result":      updatedActivity,
	})
}

// DeleteActivityByIDHandler godoc
// @Summary Delete Activity by ID
// @Description Delete an activity record
// @Tags Activity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activities/{id} [delete]
func (c *ActivityController) DeleteActivityByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	id := ctx.Params("id")
	err := c.activityUsecase.DeleteActivityByID(id)
	if err != nil {
		return ctx.Status(resolveStatusCode(err)).JSON(fiber.Map{
			"status":      resolveStatusText(err),
			"status_code": resolveStatusCode(err),
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "activity deleted successfully",
		"result":      nil,
	})
}

func resolveStatusCode(err error) int {
	if errors.Is(err, usecases.ErrActivityNotFound) {
		return fiber.StatusNotFound
	}
	if errors.Is(err, usecases.ErrStaffProfileNotFound) {
		return fiber.StatusForbidden
	}

	message := strings.ToLower(err.Error())
	if strings.Contains(message, "required") || strings.Contains(message, "cannot be empty") || strings.Contains(message, "at least one field") {
		return fiber.StatusBadRequest
	}

	return fiber.StatusInternalServerError
}

func resolveStatusText(err error) string {
	statusCode := resolveStatusCode(err)
	switch statusCode {
	case fiber.StatusBadRequest:
		return fiber.ErrBadRequest.Message
	case fiber.StatusNotFound:
		return fiber.ErrNotFound.Message
	case fiber.StatusForbidden:
		return fiber.ErrForbidden.Message
	default:
		return fiber.ErrInternalServerError.Message
	}
}
