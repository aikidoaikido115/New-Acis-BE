package controllers

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/medicine/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/medicine/usecases"
	"github.com/gofiber/fiber/v2"
)

type DrugController struct {
	drugUsecase usecases.DrugUsecase
}

func NewDrugController(drugUsecase usecases.DrugUsecase) *DrugController {
	return &DrugController{drugUsecase: drugUsecase}
}

// CreateDrugMaster godoc
// @Summary Create Drug Master
// @Description Create a new drug master record
// @Tags DrugMaster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateDrugMasterRequest true "Drug master payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-masters [post]
func (c *DrugController) CreateDrugMasterHandler(ctx *fiber.Ctx) error {
	var req models.CreateDrugMasterRequest
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

	result, err := c.drugUsecase.CreateDrugMaster(req, userID)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "drug master created successfully",
		"result":      result,
	})
}

// GetDrugMasters godoc
// @Summary Get All Drug Masters
// @Description Retrieve all drug master records
// @Tags DrugMaster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-masters [get]
func (c *DrugController) GetDrugMastersHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.drugUsecase.GetDrugMasters(userID)
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
		"message":     "drug masters retrieved successfully",
		"result":      result,
	})
}

// GetDrugMasterByID godoc
// @Summary Get Drug Master By ID
// @Description Retrieve a drug master by ID
// @Tags DrugMaster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Drug Master ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-masters/{id} [get]
func (c *DrugController) GetDrugMasterByIDHandler(ctx *fiber.Ctx) error {
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

	result, err := c.drugUsecase.GetDrugMasterByID(id, userID)
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
		"message":     "drug master retrieved successfully",
		"result":      result,
	})
}

// UpdateDrugMasterByID godoc
// @Summary Update Drug Master By ID
// @Description Update a drug master record
// @Tags DrugMaster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Drug Master ID"
// @Param request body models.UpdateDrugMasterRequest true "Drug master update payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-masters/{id} [patch]
func (c *DrugController) UpdateDrugMasterByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.UpdateDrugMasterRequest
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

	result, err := c.drugUsecase.UpdateDrugMasterByID(id, req, userID)
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
		"message":     "drug master updated successfully",
		"result":      result,
	})
}

// DeleteDrugMasterByID godoc
// @Summary Delete Drug Master By ID
// @Description Delete a drug master record
// @Tags DrugMaster
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Drug Master ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-masters/{id} [delete]
func (c *DrugController) DeleteDrugMasterByIDHandler(ctx *fiber.Ctx) error {
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

	if err := c.drugUsecase.DeleteDrugMasterByID(id, userID); err != nil {
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
		"message":     "drug master deleted successfully",
		"result":      nil,
	})
}

// CreatePersonalDrug godoc
// @Summary Create Personal Drug
// @Description Create a new personal drug record for resident
// @Tags PersonalDrug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreatePersonalDrugRequest true "Personal drug payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/personal-drugs [post]
func (c *DrugController) CreatePersonalDrugHandler(ctx *fiber.Ctx) error {
	var req models.CreatePersonalDrugRequest
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

	result, err := c.drugUsecase.CreatePersonalDrug(req, userID)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "personal drug created successfully",
		"result":      result,
	})
}

// GetPersonalDrugsOverview godoc
// @Summary Get Personal Drugs Overview (Today)
// @Description Retrieve today's personal drugs with optional filters by time_of_day, take_type, and resident name search
// @Tags PersonalDrug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param time_of_day query string false "Time of day filter"
// @Param take_type query string false "Take type filter" Enums(regular, as_needed)
// @Param search query string false "Search by resident first_name, last_name, nickname"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/personal-drugs/overview [get]
func (c *DrugController) GetPersonalDrugsOverviewHandler(ctx *fiber.Ctx) error {
	var req models.PersonalDrugOverviewQueryParams
	if err := ctx.QueryParser(&req); err != nil {
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

	result, err := c.drugUsecase.GetPersonalDrugsOverview(req, userID)
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
		"message":     "personal drugs overview retrieved successfully",
		"result":      result,
	})
}

// GetPersonalDrugsByResidentToday godoc
// @Summary Get Resident Personal Drugs Today
// @Description Retrieve today's personal drugs by resident ID
// @Tags PersonalDrug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/personal-drugs/resident [get]
func (c *DrugController) GetPersonalDrugsByResidentTodayHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	if residentID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "resident_id query parameter is required",
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

	result, err := c.drugUsecase.GetPersonalDrugsByResidentIDToday(residentID, userID)
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
		"message":     "personal drugs by resident retrieved successfully",
		"result":      result,
	})
}

// GetPersonalDrugsByResident godoc
// @Summary Get Resident Personal Drugs (All)
// @Description Retrieve all personal drugs by resident ID
// @Tags PersonalDrug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/personal-drugs/resident/all [get]
func (c *DrugController) GetPersonalDrugsByResidentHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	if residentID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "resident_id query parameter is required",
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

	result, err := c.drugUsecase.GetPersonalDrugsByResidentID(residentID, userID)
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
		"message":     "all personal drugs by resident retrieved successfully",
		"result":      result,
	})
}

// GetPersonalDrugByID godoc
// @Summary Get Personal Drug By ID
// @Description Retrieve a personal drug record by ID
// @Tags PersonalDrug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Personal Drug ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/personal-drugs/{id} [get]
func (c *DrugController) GetPersonalDrugByIDHandler(ctx *fiber.Ctx) error {
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

	result, err := c.drugUsecase.GetPersonalDrugByID(id, userID)
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
		"message":     "personal drug retrieved successfully",
		"result":      result,
	})
}

// UpdatePersonalDrugByID godoc
// @Summary Update Personal Drug By ID
// @Description Update a personal drug record
// @Tags PersonalDrug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Personal Drug ID"
// @Param request body models.UpdatePersonalDrugRequest true "Personal drug update payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/personal-drugs/{id} [patch]
func (c *DrugController) UpdatePersonalDrugByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.UpdatePersonalDrugRequest
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

	result, err := c.drugUsecase.UpdatePersonalDrugByID(id, req, userID)
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
		"message":     "personal drug updated successfully",
		"result":      result,
	})
}

// DeletePersonalDrugByID godoc
// @Summary Delete Personal Drug By ID
// @Description Delete a personal drug record
// @Tags PersonalDrug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Personal Drug ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/personal-drugs/{id} [delete]
func (c *DrugController) DeletePersonalDrugByIDHandler(ctx *fiber.Ctx) error {
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

	if err := c.drugUsecase.DeletePersonalDrugByID(id, userID); err != nil {
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
		"message":     "personal drug deleted successfully",
		"result":      nil,
	})
}

// CreateDrugPlan godoc
// @Summary Create Drug Plan
// @Description Create a new drug plan record
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateDrugPlanRequest true "Drug plan payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans [post]
func (c *DrugController) CreateDrugPlanHandler(ctx *fiber.Ctx) error {
	var req models.CreateDrugPlanRequest
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

	result, err := c.drugUsecase.CreateDrugPlan(req, userID)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "drug plan created successfully",
		"result":      result,
	})
}

// ForceGenerateTodayDrugPlans godoc
// @Summary Force Generate Today's Drug Plans (All)
// @Description Manually trigger lazy generation for today's drug plans across all residents
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.DrugPlanGenerationResponse}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/generate-today [post]
func (c *DrugController) ForceGenerateTodayDrugPlansHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.drugUsecase.ForceGenerateTodayDrugPlans(userID)
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
		"message":     "today's drug plans generated successfully",
		"result":      result,
	})
}

// ForceGenerateTodayDrugPlansByResident godoc
// @Summary Force Generate Today's Drug Plans (Resident)
// @Description Manually trigger lazy generation for today's drug plans of a specific resident
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.DrugPlanGenerationResponse}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/generate-today/resident/{resident_id} [post]
func (c *DrugController) ForceGenerateTodayDrugPlansByResidentHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Params("resident_id")

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.drugUsecase.ForceGenerateTodayDrugPlansByResidentID(residentID, userID)
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
		"message":     "today's resident drug plans generated successfully",
		"result":      result,
	})
}

// GetDrugPlansTodayResidentSummary godoc
// @Summary Get Drug Plans Resident Summary (Today)
// @Description Retrieve resident-level summary for today's drug plans using is_taken status
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.DrugPlanResidentSummaryResponse}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/istaken-summary [get]
func (c *DrugController) GetDrugPlansTodayResidentSummaryHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.drugUsecase.GetDrugPlansTodayResidentSummary(userID)
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
		"message":     "drug plans summary retrieved successfully",
		"result":      result,
	})
}

// GetDrugPlansToday godoc
// @Summary Get Drug Plans Today
// @Description Retrieve today's drug plan records
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/today [get]
func (c *DrugController) GetDrugPlansTodayHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.drugUsecase.GetDrugPlansToday(userID)
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
		"message":     "today's drug plans retrieved successfully",
		"result":      result,
	})
}

// GetDrugPlansOverview godoc
// @Summary Get Drug Plans Overview (Today)
// @Description Retrieve today's drug plans with optional filters by time_of_day, take_type, and resident name search
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param time_of_day query string false "Time of day filter"
// @Param take_type query string false "Take type filter" Enums(regular, as_needed)
// @Param search query string false "Search by resident first_name, last_name, nickname"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/overview [get]
func (c *DrugController) GetDrugPlansOverviewHandler(ctx *fiber.Ctx) error {
	var req models.DrugPlanOverviewQueryParams
	if err := ctx.QueryParser(&req); err != nil {
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

	result, err := c.drugUsecase.GetDrugPlansOverview(req, userID)
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
		"message":     "drug plans overview retrieved successfully",
		"result":      result,
	})
}

// GetDrugAdministrationHistory godoc
// @Summary Get Drug Administration History
// @Description Retrieve drug administration history with filters and pagination
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string false "Date filter (YYYY-MM-DD). Defaults to today"
// @Param search query string false "Search by resident first_name, last_name, nickname"
// @Param time_of_day query string false "Time of day filter"
// @Param status query string false "Status filter" Enums(taken, omitted, pending)
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 20, max 100)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.DrugAdministrationHistoryResponse}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/history [get]
func (c *DrugController) GetDrugAdministrationHistoryHandler(ctx *fiber.Ctx) error {
	var req models.DrugAdministrationHistoryQueryParams
	if err := ctx.QueryParser(&req); err != nil {
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

	result, err := c.drugUsecase.GetDrugAdministrationHistory(req, userID)
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
		"message":     "drug administration history retrieved successfully",
		"result":      result,
	})
}

// GetDrugPlansByResident godoc
// @Summary Get Resident Drug Plans (All)
// @Description Retrieve all drug plans by resident ID
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/resident/all [get]
func (c *DrugController) GetDrugPlansByResidentHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	if residentID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "resident_id query parameter is required",
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

	result, err := c.drugUsecase.GetDrugPlansByResidentID(residentID, userID)
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
		"message":     "all drug plans by resident retrieved successfully",
		"result":      result,
	})
}

// GetDrugPlansByResidentToday godoc
// @Summary Get Resident Drug Plans Today
// @Description Retrieve today's drug plans by resident ID
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/resident [get]
func (c *DrugController) GetDrugPlansByResidentTodayHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	if residentID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "resident_id query parameter is required",
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

	result, err := c.drugUsecase.GetDrugPlansByResidentIDToday(residentID, userID)
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
		"message":     "today's drug plans by resident retrieved successfully",
		"result":      result,
	})
}

// GetDrugPlans godoc
// @Summary Get All Drug Plans
// @Description Retrieve all drug plan records
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans [get]
func (c *DrugController) GetDrugPlansHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.drugUsecase.GetDrugPlans(userID)
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
		"message":     "drug plans retrieved successfully",
		"result":      result,
	})
}

// GetDrugPlanByID godoc
// @Summary Get Drug Plan By ID
// @Description Retrieve a drug plan record by ID
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Drug Plan ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/{id} [get]
func (c *DrugController) GetDrugPlanByIDHandler(ctx *fiber.Ctx) error {
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

	result, err := c.drugUsecase.GetDrugPlanByID(id, userID)
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
		"message":     "drug plan retrieved successfully",
		"result":      result,
	})
}

// UpdateDrugPlanByID godoc
// @Summary Update Drug Plan By ID
// @Description Update a drug plan record
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Drug Plan ID"
// @Param request body models.UpdateDrugPlanRequest true "Drug plan update payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/{id} [patch]
func (c *DrugController) UpdateDrugPlanByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.UpdateDrugPlanRequest
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

	result, err := c.drugUsecase.UpdateDrugPlanByID(id, req, userID)
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
		"message":     "drug plan updated successfully",
		"result":      result,
	})
}

// TakeDrugPlanByID godoc
// @Summary Take Drug Plan By ID (Today)
// @Description Mark a specific today's drug plan as taken and update given_by_staff_id from provided staff name
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Drug Plan ID"
// @Param request body models.TakeDrugPlanByIDRequest true "Take drug plan payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/{id}/take [patch]
func (c *DrugController) TakeDrugPlanByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.TakeDrugPlanByIDRequest
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

	result, err := c.drugUsecase.TakeDrugPlanByID(id, req, userID)
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
		"message":     "drug plan marked as taken successfully",
		"result":      result,
	})
}

// OmitDrugPlanByID godoc
// @Summary Omit Drug Plan By ID (Today)
// @Description Mark a specific today's drug plan as omitted and update given_by_staff_id from provided staff name
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Drug Plan ID"
// @Param request body models.OmitDrugPlanByIDRequest true "Omit drug plan payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/{id}/omit [patch]
func (c *DrugController) OmitDrugPlanByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.OmitDrugPlanByIDRequest
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

	result, err := c.drugUsecase.OmitDrugPlanByID(id, req, userID)
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
		"message":     "drug plan marked as omitted successfully",
		"result":      result,
	})
}

// TakeDrugPlansByResidentToday godoc
// @Summary Take All Resident Drug Plans (Today)
// @Description Mark all today's drug plans of a resident as taken by applying the single-item action to each record
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID"
// @Param request body models.TakeDrugPlansByResidentRequest true "Bulk take payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/resident/{resident_id}/take [patch]
func (c *DrugController) TakeDrugPlansByResidentTodayHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Params("resident_id")
	var req models.TakeDrugPlansByResidentRequest
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

	result, err := c.drugUsecase.TakeDrugPlansByResidentIDToday(residentID, req, userID)
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
		"message":     "today's resident drug plans marked as taken successfully",
		"result":      result,
	})
}

// OmitDrugPlansByResidentToday godoc
// @Summary Omit All Resident Drug Plans (Today)
// @Description Mark all today's drug plans of a resident as omitted by applying the single-item action to each record
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID"
// @Param request body models.OmitDrugPlansByResidentRequest true "Bulk omit payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/resident/{resident_id}/omit [patch]
func (c *DrugController) OmitDrugPlansByResidentTodayHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Params("resident_id")
	var req models.OmitDrugPlansByResidentRequest
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

	result, err := c.drugUsecase.OmitDrugPlansByResidentIDToday(residentID, req, userID)
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
		"message":     "today's resident drug plans marked as omitted successfully",
		"result":      result,
	})
}

// DeleteDrugPlanByID godoc
// @Summary Delete Drug Plan By ID
// @Description Delete a drug plan record
// @Tags DrugPlan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Drug Plan ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/drug-plans/{id} [delete]
func (c *DrugController) DeleteDrugPlanByIDHandler(ctx *fiber.Ctx) error {
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

	if err := c.drugUsecase.DeleteDrugPlanByID(id, userID); err != nil {
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
		"message":     "drug plan deleted successfully",
		"result":      nil,
	})
}
