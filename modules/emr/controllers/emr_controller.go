package controllers

import (
	// "fmt"
	// "mime/multipart"
	// "strconv"
	// "mime/multipart"

	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/usecases"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"github.com/gofiber/fiber/v2"
)

type EmrController struct {
	emrUsecase usecases.EmrUsecase
}

func NewEmrController(emrUsecase usecases.EmrUsecase) *EmrController {
	return &EmrController{
		emrUsecase: emrUsecase,
	}
}

// CreateResidentHandler godoc
// @Summary Create Resident
// @Description Create a new resident with room ID, first name, last name, age, and gender
// @Tags EMR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{room_id=string,first_name=string,last_name=string,age=int,gender=string} true "Resident information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object} "Resident created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents [post]
func (c *EmrController) CreateResident(ctx *fiber.Ctx) error {
	var req struct {
		Name      string `json:"name"`
		RoomID    string `json:"room_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       *int16 `json:"age"`
		Gender    string `json:"gender"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.FirstName == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "First name is missing",
			"result":      nil,
		})
	}

	if req.LastName == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Last name is missing",
			"result":      nil,
		})
	}

	if req.Age == nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Age is missing",
			"result":      nil,
		})
	}

	if req.Gender == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Gender is missing",
			"result":      nil,
		})
	}

	resident := &entities.Resident{
		RoomID:    req.RoomID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
		Gender:    req.Gender,
	}

	createdResident, err := c.emrUsecase.CreateResident(resident)
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
		"status_code": fiber.StatusOK,
		"message":     "resident created successfully",
		"result":      createdResident,
	})
}
