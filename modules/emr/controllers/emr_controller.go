package controllers

import (
	// "fmt"
	// "mime/multipart"
	// "strconv"
	// "mime/multipart"

	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/models"
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
// @Tags Resident
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{room_id=string,first_name=string,last_name=string,age=int,gender=string} true "Resident information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object} "Resident created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents [post]
func (c *EmrController) CreateResident(ctx *fiber.Ctx) error {
	var req models.CreateResidentRequest

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

	resident := &entities.Resident{
		RoomID:    req.RoomID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
		Gender:    req.Gender,
	}

	createdResident, err := c.emrUsecase.CreateResident(resident, userID)
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
		"message":     "resident created successfully",
		"result":      createdResident,
	})
}

// GetResidentByID godoc
// @Summary Get Resident by ID
// @Description Retrieve a single resident's information by their unique ID
// @Tags Resident
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Resident retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents/{id} [get]
func (c *EmrController) GetResidentByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	resident, err := c.emrUsecase.GetResidentByID(id)
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
		"message":     "resident retrieved successfully",
		"result":      resident,
	})
}

// GetResidentByRoomID godoc
// @Summary Get Residents by Room ID
// @Description Retrieve all residents assigned to a specific room using room_id query parameter
// @Tags Resident
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param room_id query string true "Room ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Residents retrieved successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing room_id query parameter"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents [get]
func (c *EmrController) GetResidentByRoomID(ctx *fiber.Ctx) error {
	roomID := ctx.Query("room_id")

	if roomID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "room_id query parameter is required",
			"result":      nil,
		})
	}

	residents, err := c.emrUsecase.GetResidentByRoomID(roomID)
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
		"message":     "residents retrieved successfully",
		"result":      residents,
	})
}

// GetAllResidents godoc
// @Summary Get All Residents
// @Description Retrieve a list of all residents in the system
// @Tags Resident
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "All residents retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents/all [get]
func (c *EmrController) GetAllResidents(ctx *fiber.Ctx) error {
	residents, err := c.emrUsecase.GetAllResidents()
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
		"message":     "residents retrieved successfully",
		"result":      residents,
	})
}

// UpdateResident godoc
// @Summary Update Resident
// @Description Partially update an existing resident's information by their unique ID. Only send fields that need to be updated.
// @Tags Resident
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resident ID"
// @Param request body object{room_id=string,first_name=string,last_name=string,age=int,gender=string,labels=[]object{label_name=string,note_text=string}} true "Fields to update (all optional)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Resident updated successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Invalid data"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents/{id} [patch]
func (c *EmrController) UpdateResidentByID(ctx *fiber.Ctx) error {
	var req models.UpdateResidentRequest

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

	residentID := ctx.Params("id")

	updatedResident, err := c.emrUsecase.UpdateResidentByID(residentID, req, userID)
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
		"message":     "resident updated successfully",
		"result":      updatedResident,
	})
}

// GetRoomByID godoc
// @Summary Get Room by ID
// @Description Retrieve room information by room ID
// @Tags Room
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Room retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/rooms/{id} [get]
func (c *EmrController) GetRoomByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	room, err := c.emrUsecase.GetRoomByID(id)
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
		"message":     "room retrieved successfully",
		"result":      room,
	})
}

// GetAllRooms godoc
// @Summary Get All Rooms
// @Description Retrieve a list of all rooms in the facility
// @Tags Room
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Rooms retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/rooms [get]
func (c *EmrController) GetAllRooms(ctx *fiber.Ctx) error {
	rooms, err := c.emrUsecase.GetAllRooms()
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
		"message":     "rooms retrieved successfully",
		"result":      rooms,
	})
}

// CreateRoom godoc
// @Summary Create Room
// @Description Create a new room with staff ID, floor number, and room number
// @Tags Room
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{staff_id=string,floor=int,room_number=string} true "Room information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object} "Room created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/rooms [post]
func (c *EmrController) CreateRoom(ctx *fiber.Ctx) error {
	var req models.CreateRoomRequest

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

	room := &entities.Room{
		StaffID:    req.StaffID,
		Floor:      req.Floor,
		RoomNumber: req.RoomNumber,
	}

	createdRoom, err := c.emrUsecase.CreateRoom(room, userID)
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
		"message":     "room created successfully",
		"result":      createdRoom,
	})
}

// UpdateRoomByID godoc
// @Summary Update Room
// @Description Partially update an existing room's information by its unique ID. Only send fields that need to be updated.
// @Tags Room
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Param request body object{staff_id=string} true "Fields to update (all optional)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Room updated successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Invalid data"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/rooms/{id} [patch]
func (c *EmrController) UpdateRoomByID(ctx *fiber.Ctx) error {
	var req models.UpdateRoomRequest

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

	roomID := ctx.Params("id")

	updatedRoom, err := c.emrUsecase.UpdateRoomByID(roomID, req, userID)
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
		"message":     "room updated successfully",
		"result":      updatedRoom,
	})
}

// GetNumberOfResidentsDashboard godoc
// @Summary Get Number of Residents for Dashboard
// @Description Retrieve the number of residents categorized by care levels for dashboard display
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Dashboard data retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/dashboard/residents [get]
func (c *EmrController) GetNumberOfResidentsDashboard(ctx *fiber.Ctx) error {
	response, err := c.emrUsecase.GetNumberOfResidentsDashboard()
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
		"message":     "dashboard data retrieved successfully",
		"result":      response,
	})
}

// GetResidentGenderStatsDashboard godoc
// @Summary Get Resident Gender Stats for Dashboard
// @Description Retrieve statistics on resident gender distribution for dashboard display
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Resident gender stats retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/dashboard/resident-gender-stats [get]
func (c *EmrController) GetResidentGenderStatsDashboard(ctx *fiber.Ctx) error {
	response, err := c.emrUsecase.GetResidentGenderStatsDashboard()
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
		"message":     "resident gender stats retrieved successfully",
		"result":      response,
	})
}

// GetAllIntakeLabels godoc
// @Summary Get All Intake Labels
// @Description Retrieve a list of all available intake label types (e.g., Blood Pressure, Temperature, Heart Rate)
// @Tags Intake
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Intake labels retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/intake-labels/all [get]
func (c *EmrController) GetAllIntakeLabels(ctx *fiber.Ctx) error {
	labels, err := c.emrUsecase.GetAllIntakeLabels()
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
		"message":     "intake labels retrieved successfully",
		"result":      labels,
	})
}

// GetResidentLabelsByResidentID godoc
// @Summary Get Resident Labels by Resident ID
// @Description Retrieve all intake labels associated with a specific resident using their resident ID
// @Tags Intake
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Resident labels retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/intake-labels [get]
func (c *EmrController) GetResidentLabelsByResidentID(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")

	if residentID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "resident_id query parameter is required",
			"result":      nil,
		})
	}
	labels, err := c.emrUsecase.GetResidentLabelsByResidentID(residentID)
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
		"message":     "resident labels retrieved successfully",
		"result":      labels,
	})
}

// CreateIntakeLabelByResidentID godoc
// @Summary Create Intake Labels for Resident
// @Description Create one or more intake labels (vital signs, notes, etc.) for a specific resident
// @Tags Intake
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{resident_id=string,labels=[]object{label_name=string,note_text=string}} true "Resident ID and array of labels"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=[]object} "Intake labels created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields or invalid data"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing or invalid authentication"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/intake-labels [post]
func (c *EmrController) CreateIntakeLabelByResidentID(ctx *fiber.Ctx) error {
	var req models.CreateIntakeLabelByResidentRequest

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

	result, err := c.emrUsecase.CreateIntakeLabelByResidentID(req.ResidentID, req.Labels, userID)
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
		"message":     "Intake labels created successfully",
		"result":      result,
	})
}

// CreateVitalSign godoc
// @Summary Create Vital Sign
// @Description Create a new vital sign entry for a resident
// @Tags VitalSign
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{resident_id=string,label_name=string,value=float64,unit=string,timestamp=string} true "Vital sign information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object} "Vital sign created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/vital-signs [post]
func (c *EmrController) CreateVitalSign(ctx *fiber.Ctx) error {
	var req models.CreateVitalSignRequest

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

	vitalSign := &entities.VitalSign{
		ResidentID:             req.ResidentID,
		Temperature:            req.Temperature,
		HeartRate:              req.HeartRate,
		BreathingRate:          req.BreathingRate,
		BloodPressureSystolic:  req.BloodPressureSystolic,
		BloodPressureDiastolic: req.BloodPressureDiastolic,
		OxygenSaturation:       req.OxygenSaturation,
	}
	createdVitalSign, err := c.emrUsecase.CreateVitalSign(vitalSign, userID)
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
		"message":     "vital sign created successfully",
		"result":      createdVitalSign,
	})

}
