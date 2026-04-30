package controllers

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/warehouse/usecases"
	"github.com/gofiber/fiber/v2"
)

type WarehouseController struct {
	warehouseUsecase usecases.WarehouseUsecase
}

func NewWarehouseController(warehouseUsecase usecases.WarehouseUsecase) *WarehouseController {
	return &WarehouseController{warehouseUsecase: warehouseUsecase}
}

// @Summary Get Warehouse Items
// @Description Get warehouse items with optional search and category filters
// @Tags Warehouse
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by item name or code"
// @Param category query string false "Category (MED, EQU, CON)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/items [get]
func (c *WarehouseController) GetWarehouseItemsHandler(ctx *fiber.Ctx) error {
	var query models.ListWarehouseItemsQuery
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

	result, err := c.warehouseUsecase.GetWarehouseItems(query, userID)
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
		"message":     "warehouse items retrieved successfully",
		"result":      result,
	})
}

// @Summary Create Warehouse Item
// @Description Create a new warehouse item
// @Tags Warehouse
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateWarehouseItemRequest true "Warehouse item payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/items [post]
func (c *WarehouseController) CreateWarehouseItemHandler(ctx *fiber.Ctx) error {
	var req models.CreateWarehouseItemRequest
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

	result, err := c.warehouseUsecase.CreateWarehouseItem(req, userID)
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
		"message":     "warehouse item created and waiting for approval",
		"result":      result,
	})
}

// @Summary Update Warehouse Item By ID
// @Description Update warehouse item fields by item ID
// @Tags Warehouse
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Warehouse item ID"
// @Param request body models.UpdateWarehouseItemRequest true "Warehouse item update payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/items/{id} [patch]
func (c *WarehouseController) UpdateWarehouseItemByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.UpdateWarehouseItemRequest
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

	result, err := c.warehouseUsecase.UpdateWarehouseItemByID(id, req, userID)
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
		"message":     "warehouse item updated successfully",
		"result":      result,
	})
}

// @Summary Delete Warehouse Item By ID
// @Description Delete a warehouse item by item ID
// @Tags Warehouse
// @Produce json
// @Security BearerAuth
// @Param id path string true "Warehouse item ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/items/{id} [delete]
func (c *WarehouseController) DeleteWarehouseItemByIDHandler(ctx *fiber.Ctx) error {
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

	if err := c.warehouseUsecase.DeleteWarehouseItemByID(id, userID); err != nil {
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
		"message":     "warehouse remove request submitted for approval",
		"result":      nil,
	})
}

// @Summary Adjust Warehouse Item Quantity
// @Description Adjust item quantity by restock or withdraw mode
// @Tags Warehouse
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Warehouse item ID"
// @Param request body models.AdjustWarehouseItemRequest true "Adjust warehouse item payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/items/{id}/adjust [post]
func (c *WarehouseController) AdjustWarehouseItemByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req models.AdjustWarehouseItemRequest
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

	result, err := c.warehouseUsecase.AdjustWarehouseItemByID(id, req, userID)
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
		"message":     "warehouse adjustment request submitted for approval",
		"result":      result,
	})
}

// @Summary Get Warehouse Transactions
// @Description Get warehouse transaction history with optional filters
// @Tags Warehouse
// @Produce json
// @Security BearerAuth
// @Param startDate query string false "Start date in YYYY-MM-DD"
// @Param endDate query string false "End date in YYYY-MM-DD"
// @Param searchItem query string false "Search by item name or item code"
// @Param searchUser query string false "Search by operator name"
// @Param status query string false "Approval status (รออนุมัติ, อนุมัติ, ไม่อนุมัติ)"
// @Param type query string false "Transaction type (เพิ่มสินค้าใหม่, เติมสินค้า, เบิกสินค้า, นำออก)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]models.WarehouseTransactionResponse}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/transactions [get]
func (c *WarehouseController) GetWarehouseTransactionsHandler(ctx *fiber.Ctx) error {
	var query models.ListWarehouseTransactionsQuery
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

	result, err := c.warehouseUsecase.GetWarehouseTransactions(query, userID)
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
		"message":     "warehouse transactions retrieved successfully",
		"result":      result,
	})
}

// @Summary Get Warehouse Transaction By ID
// @Description Get warehouse transaction detail by transaction ID
// @Tags Warehouse
// @Produce json
// @Security BearerAuth
// @Param id path string true "Warehouse transaction ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.WarehouseTransactionResponse}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/transactions/{id} [get]
func (c *WarehouseController) GetWarehouseTransactionByIDHandler(ctx *fiber.Ctx) error {
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

	result, err := c.warehouseUsecase.GetWarehouseTransactionByID(id, userID)
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
		"message":     "warehouse transaction retrieved successfully",
		"result":      result,
	})
}

// @Summary Approve Warehouse Transactions
// @Description Approve selected pending warehouse transactions
// @Tags Warehouse
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ApproveTransactionsRequest true "Approve transactions payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]models.WarehouseTransactionResponse}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/transactions/approve [patch]
func (c *WarehouseController) ApproveTransactionsHandler(ctx *fiber.Ctx) error {
	var req models.ApproveTransactionsRequest
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

	result, err := c.warehouseUsecase.ApproveTransactions(req, userID)
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
		"message":     "warehouse transactions approved successfully",
		"result":      result,
	})
}

// @Summary Reject Warehouse Transactions
// @Description Reject selected pending warehouse transactions with reason
// @Tags Warehouse
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RejectTransactionsRequest true "Reject transactions payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]models.WarehouseTransactionResponse}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/warehouse/transactions/reject [patch]
func (c *WarehouseController) RejectTransactionsHandler(ctx *fiber.Ctx) error {
	var req models.RejectTransactionsRequest
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

	result, err := c.warehouseUsecase.RejectTransactions(req, userID)
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
		"message":     "warehouse transactions rejected successfully",
		"result":      result,
	})
}
