package controllers

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

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
	updatedActivity, err := c.activityUsecase.UpdateActivityByID(id, req, userID)
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

// CreateActivityScheduleHandler godoc
// @Summary Create Activity Schedule
// @Description Create a new activity schedule record
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateActivityScheduleRequest true "Activity schedule payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.ActivitySchedule}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules [post]
func (c *ActivityController) CreateActivityScheduleHandler(ctx *fiber.Ctx) error {
	var req models.CreateActivityScheduleRequest
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

	createdActivitySchedule, err := c.activityUsecase.CreateActivitySchedule(req, userID)
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
		"message":     "activity schedule created successfully",
		"result":      createdActivitySchedule,
	})
}

// CreateActivityScheduleWithActivitySyncHandler godoc
// @Summary Create Activity Schedule with Activity Sync
// @Description Create a new schedule by activity name. Reuse existing activity when attributes match, update activity when attributes changed, or create new activity when name does not exist.
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateActivityScheduleWithActivitySyncRequest true "Activity and schedule payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.ActivitySchedule}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules/sync [post]
func (c *ActivityController) CreateActivityScheduleWithActivitySyncHandler(ctx *fiber.Ctx) error {
	var req models.CreateActivityScheduleWithActivitySyncRequest
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

	createdActivitySchedule, err := c.activityUsecase.CreateActivityScheduleWithActivitySync(req, userID)
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
		"message":     "activity schedule synced and created successfully",
		"result":      createdActivitySchedule,
	})
}

// UpdateActivityScheduleWithActivitySyncByIDHandler godoc
// @Summary Update Activity Schedule with Activity Sync by ID
// @Description Partially update activity and activity schedule in one request. If activity_name is changed to an existing name, an error is returned.
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity Schedule ID"
// @Param request body models.UpdateActivityScheduleWithActivitySyncRequest true "Activity and schedule partial update payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.ActivitySchedule}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 409 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules/sync/{id} [patch]
func (c *ActivityController) UpdateActivityScheduleWithActivitySyncByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateActivityScheduleWithActivitySyncRequest
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
	updatedActivitySchedule, err := c.activityUsecase.UpdateActivityScheduleWithActivitySyncByID(id, req, userID)
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
		"message":     "activity and activity schedule updated successfully",
		"result":      updatedActivitySchedule,
	})
}

// GetActivityScheduleWithActivitySyncByIDHandler godoc
// @Summary Get Activity Schedule with Activity Sync by ID
// @Description Retrieve a clean activity-schedule view by schedule ID with only activity_name, activity_type, date, start_time, end_time, location, and description.
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity Schedule ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.ActivityScheduleWithActivitySyncResponse}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules/sync/{id} [get]
func (c *ActivityController) GetActivityScheduleWithActivitySyncByIDHandler(ctx *fiber.Ctx) error {
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
	result, err := c.activityUsecase.GetActivityScheduleWithActivitySyncByID(id)
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
		"message":     "activity schedule sync data retrieved successfully",
		"result":      result,
	})
}

// GetAllActivitySchedulesWithActivitySyncHandler godoc
// @Summary Get All Activity Schedules with Activity Sync
// @Description Retrieve joined activity and activity schedule records. Optional query param date (YYYY-MM-DD) to filter by date.
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string false "Filter by date (YYYY-MM-DD)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]models.ActivityScheduleWithActivitySyncResponse}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules/sync [get]
func (c *ActivityController) GetAllActivitySchedulesWithActivitySyncHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	dateQuery := strings.TrimSpace(ctx.Query("date"))
	var filterDate *time.Time
	if dateQuery != "" {
		loc := time.FixedZone("ICT", 7*60*60)
		parsedDate, err := time.ParseInLocation("2006-01-02", dateQuery, loc)
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     "date must be in YYYY-MM-DD format",
				"result":      nil,
			})
		}
		filterDate = &parsedDate
	}

	activitySchedules, err := c.activityUsecase.GetAllActivitySchedulesWithActivitySync(filterDate)
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
		"message":     "activity schedules with activity data retrieved successfully",
		"result":      activitySchedules,
	})
}

// GetActivityScheduleByIDHandler godoc
// @Summary Get Activity Schedule by ID
// @Description Retrieve activity schedule by ID
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity Schedule ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.ActivitySchedule}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules/{id} [get]
func (c *ActivityController) GetActivityScheduleByIDHandler(ctx *fiber.Ctx) error {
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
	activitySchedule, err := c.activityUsecase.GetActivityScheduleByID(id)
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
		"message":     "activity schedule retrieved successfully",
		"result":      activitySchedule,
	})
}

// GetAllActivitySchedulesHandler godoc
// @Summary Get All Activity Schedules
// @Description Retrieve all activity schedule records
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]entities.ActivitySchedule}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules [get]
func (c *ActivityController) GetAllActivitySchedulesHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	activitySchedules, err := c.activityUsecase.GetAllActivitySchedules()
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
		"message":     "activity schedules retrieved successfully",
		"result":      activitySchedules,
	})
}

// UpdateActivityScheduleByIDHandler godoc
// @Summary Update Activity Schedule by ID
// @Description Update an activity schedule record
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity Schedule ID"
// @Param request body models.UpdateActivityScheduleRequest true "Activity schedule update payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.ActivitySchedule}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules/{id} [patch]
func (c *ActivityController) UpdateActivityScheduleByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateActivityScheduleRequest
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
	updatedActivitySchedule, err := c.activityUsecase.UpdateActivityScheduleByID(id, req, userID)
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
		"message":     "activity schedule updated successfully",
		"result":      updatedActivitySchedule,
	})
}

// DeleteActivityScheduleByIDHandler godoc
// @Summary Delete Activity Schedule by ID
// @Description Delete an activity schedule record
// @Tags ActivitySchedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity Schedule ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules/{id} [delete]
func (c *ActivityController) DeleteActivityScheduleByIDHandler(ctx *fiber.Ctx) error {
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
	err := c.activityUsecase.DeleteActivityScheduleByID(id)
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
		"message":     "activity schedule deleted successfully",
		"result":      nil,
	})
}

// CreateParticipationHandler godoc
// @Summary Create Participation
// @Description Create a new participation record
// @Tags Participation
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param resident_id formData string true "Resident ID"
// @Param as_id formData string true "Activity Schedule ID"
// @Param is_participating formData bool false "Is participating"
// @Param file formData file false "Image files to upload (can select multiple)"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.Participation}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 409 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-participations [post]
func (c *ActivityController) CreateParticipationHandler(ctx *fiber.Ctx) error {
	var req models.CreateParticipationRequest
	var files []*multipart.FileHeader

	if strings.HasPrefix(strings.ToLower(ctx.Get("Content-Type")), "multipart/form-data") {
		form, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     err.Error(),
				"result":      nil,
			})
		}

		req.ResidentID = strings.TrimSpace(ctx.FormValue("resident_id"))
		req.ASID = strings.TrimSpace(ctx.FormValue("as_id"))

		if isParticipatingRaw := strings.TrimSpace(ctx.FormValue("is_participating")); isParticipatingRaw != "" {
			parsed, parseErr := strconv.ParseBool(isParticipatingRaw)
			if parseErr != nil {
				return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
					"status":      fiber.ErrBadRequest.Message,
					"status_code": fiber.ErrBadRequest.Code,
					"message":     "is_participating must be boolean",
					"result":      nil,
				})
			}
			req.IsParticipating = parsed
		}

		files = getMultipartImageFiles(form)
	} else {
		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     err.Error(),
				"result":      nil,
			})
		}
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

	createdParticipation, err := c.activityUsecase.CreateParticipation(req, userID, files)
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
		"message":     "participation created successfully",
		"result":      createdParticipation,
	})
}

// GetParticipationByResidentIDAndASIDHandler godoc
// @Summary Get Participation by Composite Key
// @Description Retrieve participation by resident_id and as_id
// @Tags Participation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID"
// @Param as_id path string true "Activity Schedule ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.Participation}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-participations/{resident_id}/{as_id} [get]
func (c *ActivityController) GetParticipationByResidentIDAndASIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	residentID := ctx.Params("resident_id")
	asID := ctx.Params("as_id")
	participation, err := c.activityUsecase.GetParticipationByResidentIDAndASID(residentID, asID)
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
		"message":     "participation retrieved successfully",
		"result":      participation,
	})
}

// GetAllParticipationsHandler godoc
// @Summary Get All Participations
// @Description Retrieve all participation records
// @Tags Participation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]entities.Participation}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-participations [get]
func (c *ActivityController) GetAllParticipationsHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	participations, err := c.activityUsecase.GetAllParticipations()
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
		"message":     "participations retrieved successfully",
		"result":      participations,
	})
}

// GetResidentsByScheduleIDCustomHandler godoc
// @Summary Get Residents by Activity Schedule ID (Custom)
// @Description Retrieve residents in a schedule with room and intake labels. Supports filters by search, floor, and intake label IDs with pagination.
// @Tags Participation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity Schedule ID"
// @Param search query string false "Search by resident first_name, last_name, nickname"
// @Param floor query int false "Filter by room floor"
// @Param label_ids query []string false "Filter by intake label IDs (match all provided labels)"
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 20, max 100)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.ResidentsByScheduleListResponse}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-schedules/{id}/residents [get]
func (c *ActivityController) GetResidentsByScheduleIDCustomHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	asID := strings.TrimSpace(ctx.Params("id"))
	if asID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "activity schedule id is required",
			"result":      nil,
		})
	}

	var req models.ResidentsByScheduleQueryParams
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	result, err := c.activityUsecase.GetResidentsByScheduleIDCustom(asID, req)
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
		"message":     "residents by activity schedule retrieved successfully",
		"result":      result,
	})
}

// UpdateParticipationByResidentIDAndASIDHandler godoc
// @Summary Update Participation by Composite Key
// @Description Update participation by resident_id and as_id. resident_id and as_id are immutable and accepted from path only.
// @Tags Participation
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID (immutable key)"
// @Param as_id path string true "Activity Schedule ID (immutable key)"
// @Param is_participating formData bool false "Is participating"
// @Param file formData file false "Image files to upload (can select multiple)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.Participation}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-participations/{resident_id}/{as_id} [patch]
func (c *ActivityController) UpdateParticipationByResidentIDAndASIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateParticipationRequest
	var files []*multipart.FileHeader

	if strings.HasPrefix(strings.ToLower(ctx.Get("Content-Type")), "multipart/form-data") {
		form, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     err.Error(),
				"result":      nil,
			})
		}

		if hasImmutableCompositeKeyInMultipart(form, "resident_id", "as_id") {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     "resident_id and as_id are immutable and must not be sent in request body",
				"result":      nil,
			})
		}

		isParticipating, parseErr := parseOptionalBoolFromMultipart(form, "is_participating")
		if parseErr != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     parseErr.Error(),
				"result":      nil,
			})
		}
		req.IsParticipating = isParticipating

		files = getMultipartImageFiles(form)
	} else {
		if hasImmutableCompositeKeyInJSON(ctx.Body(), "resident_id", "as_id") {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     "resident_id and as_id are immutable and must not be sent in request body",
				"result":      nil,
			})
		}

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     err.Error(),
				"result":      nil,
			})
		}
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

	residentID := ctx.Params("resident_id")
	asID := ctx.Params("as_id")
	updatedParticipation, err := c.activityUsecase.UpdateParticipationByResidentIDAndASID(residentID, asID, req, userID, files)
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
		"message":     "participation updated successfully",
		"result":      updatedParticipation,
	})
}

// DeleteParticipationByResidentIDAndASIDHandler godoc
// @Summary Delete Participation by Composite Key
// @Description Delete participation by resident_id and as_id
// @Tags Participation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID"
// @Param as_id path string true "Activity Schedule ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-participations/{resident_id}/{as_id} [delete]
func (c *ActivityController) DeleteParticipationByResidentIDAndASIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	residentID := ctx.Params("resident_id")
	asID := ctx.Params("as_id")
	err := c.activityUsecase.DeleteParticipationByResidentIDAndASID(residentID, asID)
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
		"message":     "participation deleted successfully",
		"result":      nil,
	})
}

// BulkUpdateParticipationIsParticipatingByResidentIDsHandler godoc
// @Summary Bulk Update Participation IsParticipating by Resident IDs
// @Description Update is_participating for one or more resident IDs in a single request (scoped by as_id)
// @Tags Participation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BulkUpdateParticipationIsParticipatingByResidentIDsRequest true "Bulk participation update payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]entities.Participation}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/activity-participations/is-participating/bulk [patch]
func (c *ActivityController) BulkUpdateParticipationIsParticipatingByResidentIDsHandler(ctx *fiber.Ctx) error {
	var req models.BulkUpdateParticipationIsParticipatingByResidentIDsRequest
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

	updatedParticipations, err := c.activityUsecase.BulkUpdateParticipationIsParticipatingByResidentIDs(req, userID)
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
		"message":     "participations updated successfully",
		"result":      updatedParticipations,
	})
}

func resolveStatusCode(err error) int {
	if errors.Is(err, usecases.ErrActivityNotFound) {
		return fiber.StatusNotFound
	}
	if errors.Is(err, usecases.ErrActivityScheduleNotFound) {
		return fiber.StatusNotFound
	}
	if errors.Is(err, usecases.ErrActivityAlreadyExists) {
		return fiber.StatusConflict
	}
	if errors.Is(err, usecases.ErrParticipationNotFound) {
		return fiber.StatusNotFound
	}
	if errors.Is(err, usecases.ErrParticipationAlreadyExists) {
		return fiber.StatusConflict
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
	case fiber.StatusConflict:
		return fiber.ErrConflict.Message
	case fiber.StatusNotFound:
		return fiber.ErrNotFound.Message
	case fiber.StatusForbidden:
		return fiber.ErrForbidden.Message
	default:
		return fiber.ErrInternalServerError.Message
	}
}

func parseOptionalBoolFromMultipart(form *multipart.Form, key string) (*bool, error) {
	values, ok := form.Value[key]
	if !ok || len(values) == 0 {
		return nil, nil
	}

	raw := strings.TrimSpace(values[0])
	if raw == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, errors.New(key + " must be boolean")
	}

	return &parsed, nil
}

func getMultipartImageFiles(form *multipart.Form) []*multipart.FileHeader {
	files := make([]*multipart.FileHeader, 0, len(form.File["file"])+len(form.File["images"])+len(form.File["image"]))

	files = append(files, form.File["file"]...)

	files = append(files, form.File["images"]...)
	files = append(files, form.File["image"]...)
	return files
}

func hasImmutableCompositeKeyInMultipart(form *multipart.Form, keys ...string) bool {
	for _, key := range keys {
		values := form.Value[key]
		for _, value := range values {
			if strings.TrimSpace(value) != "" {
				return true
			}
		}
	}

	return false
}

func hasImmutableCompositeKeyInJSON(body []byte, keys ...string) bool {
	if strings.TrimSpace(string(body)) == "" {
		return false
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}

	for _, key := range keys {
		if _, exists := payload[key]; exists {
			return true
		}
	}

	return false
}
