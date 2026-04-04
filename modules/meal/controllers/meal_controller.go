package controllers

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/meal/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/meal/usecases"
	"github.com/gofiber/fiber/v2"
)

type MealController struct {
	mealUsecase usecases.MealUsecase
}

func NewMealController(mealUsecase usecases.MealUsecase) *MealController {
	return &MealController{mealUsecase: mealUsecase}
}

// CreateMenuHandler godoc
// @Summary Create Menu
// @Description Create a new menu. Only users with Kitchen Staff role can manage meals.
// @Tags Meal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateMenuRequest true "Menu information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.Menu}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/meals/menus [post]
func (c *MealController) CreateMenuHandler(ctx *fiber.Ctx) error {
	var req models.CreateMenuRequest
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

	menu := &entities.Menu{
		MenuName:    req.MenuName,
		Description: req.Description,
	}

	createdMenu, err := c.mealUsecase.CreateMenu(menu, userID)
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
		"message":     "menu created successfully",
		"result":      createdMenu,
	})
}

// GetMenuByIDHandler godoc
// @Summary Get Menu by ID
// @Description Get a menu by ID. Only users with Kitchen Staff role can manage meals.
// @Tags Meal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Menu ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.Menu}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/meals/menus/{id} [get]
func (c *MealController) GetMenuByIDHandler(ctx *fiber.Ctx) error {
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
	menu, err := c.mealUsecase.GetMenuByID(id, userID)
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
		"message":     "menu retrieved successfully",
		"result":      menu,
	})
}

// GetAllMenusHandler godoc
// @Summary Get All Menus
// @Description Get all menus. Only users with Kitchen Staff role can manage meals.
// @Tags Meal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]entities.Menu}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/meals/menus [get]
func (c *MealController) GetAllMenusHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	menus, err := c.mealUsecase.GetAllMenus(userID)
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
		"message":     "menus retrieved successfully",
		"result":      menus,
	})
}

// UpdateMenuHandler godoc
// @Summary Update Menu
// @Description Update menu by ID. Only users with Kitchen Staff role can manage meals.
// @Tags Meal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Menu ID"
// @Param request body models.UpdateMenuRequest true "Menu fields to update"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.Menu}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/meals/menus/{id} [patch]
func (c *MealController) UpdateMenuHandler(ctx *fiber.Ctx) error {
	var req models.UpdateMenuRequest
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

	menuID := ctx.Params("id")
	updatedMenu, err := c.mealUsecase.UpdateMenu(menuID, req, userID)
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
		"message":     "menu updated successfully",
		"result":      updatedMenu,
	})
}

// CreateMealPlanHandler godoc
// @Summary Create Meal Plan
// @Description Create a new meal plan. Only users with Kitchen Staff role can manage meals.
// @Tags Meal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateMealPlanRequest true "Meal plan information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.MealPlan}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/meals/meal-plans [post]
func (c *MealController) CreateMealPlanHandler(ctx *fiber.Ctx) error {
	var req models.CreateMealPlanRequest
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

	mealPlan := &entities.MealPlan{
		MenuID:       req.MenuID,
		BackUpMenuID: req.BackUpMenuID,
		MainAmount:   req.MainAmount,
		BackUpAmount: req.BackUpAmount,
		MealType:     req.MealType,
	}

	humanInTheLoop := false
	if req.HumanInTheLoop != nil {
		humanInTheLoop = *req.HumanInTheLoop
	}

	createdMealPlan, warningSummary, err := c.mealUsecase.CreateMealPlan(mealPlan, userID, humanInTheLoop)
	if err != nil {
		// Check if this is an allergy check error with AI response data
		if allergyErr, ok := err.(*models.AllergyCheckError); ok {
			// Return the AI response data along with the error
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      allergyErr.Status,
				"status_code": fiber.StatusBadRequest,
				"message":     allergyErr.Error(),
				"result":      allergyErr.Response,
			})
		}

		// For other errors, return standard error response
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	message := "meal plan created successfully"
	if warningSummary != nil {
		message = "meal plan created with warning: main menu may not be safe for all residents"
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     message,
		"result": fiber.Map{
			"meal_plan":      createdMealPlan,
			"allergy_check":  warningSummary,
			"has_ai_warning": warningSummary != nil,
		},
	})
}

// GetMealPlanByIDHandler godoc
// @Summary Get Meal Plan by ID
// @Description Get a meal plan by ID. Only users with Kitchen Staff role can manage meals.
// @Tags Meal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Meal Plan ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.MealPlan}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/meals/meal-plans/{id} [get]
func (c *MealController) GetMealPlanByIDHandler(ctx *fiber.Ctx) error {
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
	mealPlan, err := c.mealUsecase.GetMealPlanByID(id, userID)
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
		"message":     "meal plan retrieved successfully",
		"result":      mealPlan,
	})
}

// GetAllMealPlansHandler godoc
// @Summary Get All Meal Plans
// @Description Get all meal plans. Only users with Kitchen Staff role can manage meals.
// @Tags Meal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]entities.MealPlan}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/meals/meal-plans [get]
func (c *MealController) GetAllMealPlansHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	mealPlans, err := c.mealUsecase.GetAllMealPlans(userID)
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
		"message":     "meal plans retrieved successfully",
		"result":      mealPlans,
	})
}

// GetMealPlansTodayHandler godoc
// @Summary Get Today's Meal Plans
// @Description Get all meal plans created today. Only users with Kitchen Staff role can manage meals.
// @Tags Meal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]entities.MealPlan}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/meals/meal-plans/today [get]
func (c *MealController) GetMealPlansTodayHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	mealPlans, err := c.mealUsecase.GetMealPlansToday(userID)
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
		"message":     "today's meal plans retrieved successfully",
		"result":      mealPlans,
	})
}

// UpdateMealPlanHandler godoc
// DEPRECATED: This endpoint is temporarily disabled
// // @Summary Update Meal Plan
// // @Description Update meal plan by ID. Only users with Kitchen Staff role can manage meals.
// // @Tags Meal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "Meal Plan ID"
// // @Param request body models.UpdateMealPlanRequest true "Meal plan fields to update"
// // @Success 200 {object} object{status=string,status_code=int,message=string,result=entities.MealPlan}
// // @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// // @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// // @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// // @Router /api/meals/meal-plans/{id} [patch]
func (c *MealController) UpdateMealPlanHandler(ctx *fiber.Ctx) error {
	var req models.UpdateMealPlanRequest
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

	mealPlanID := ctx.Params("id")
	mealPlan := &entities.MealPlan{}
	if req.MenuID != nil {
		mealPlan.MenuID = *req.MenuID
	}
	mealPlan.BackUpMenuID = req.BackUpMenuID
	if req.MainAmount != nil {
		mealPlan.MainAmount = *req.MainAmount
	}
	mealPlan.BackUpAmount = req.BackUpAmount
	if req.MealType != nil {
		mealPlan.MealType = *req.MealType
	}

	updatedMealPlan, err := c.mealUsecase.UpdateMealPlan(mealPlanID, mealPlan, userID)
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
		"message":     "meal plan updated successfully",
		"result":      updatedMealPlan,
	})
}
