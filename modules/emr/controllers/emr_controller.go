package controllers

import (
	// "fmt"
	// "mime/multipart"
	// "strconv"
	// "mime/multipart"
	"encoding/json"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/usecases"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"

	"github.com/google/uuid"
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
// @Description Create a new resident. Required: room_id, first_name, last_name, gender, id_card_number (13 digits), date_of_birth, check_in_date, status (active/inactive). Optional: nickname, purpose_of_stay, expected_check_out_date, pre_existing_conditions, pre_existing_conditions_notes, resuscitation_status (CPR/DNR), surgical_history, preferred_emergency_hospital, emergency_hospital_phone (10 digits)
// @Tags Resident
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{room_id=string,first_name=string,last_name=string,gender=string,nickname=string,id_card_number=string,date_of_birth=string,purpose_of_stay=string,check_in_date=string,expected_check_out_date=string,status=string,pre_existing_conditions=string,pre_existing_conditions_notes=string,resuscitation_status=string,surgical_history=string,preferred_emergency_hospital=string,emergency_hospital_phone=string} true "Resident information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object} "Resident created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing or invalid fields"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents [post]
func (c *EmrController) CreateResidentHandler(ctx *fiber.Ctx) error {
	var req models.CreateResidentRequest
	var file multipart.File
	if form, err := ctx.MultipartForm(); err == nil && (len(form.Value) > 0 || len(form.File) > 0) {
		var parseErr error
		req, file, parseErr = parseResidentCreateForm(form)
		if parseErr != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     parseErr.Error(),
				"result":      nil,
			})
		}
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

	emergencyContactsJSON := datatypes.JSON(nil)
	if len(req.EmergencyContacts) > 0 {
		cleanedContacts := make([]models.EmergencyContact, 0, len(req.EmergencyContacts))
		for _, contact := range req.EmergencyContacts {
			name := strings.TrimSpace(contact.Name)
			relation := strings.TrimSpace(contact.Relation)
			phone := strings.TrimSpace(contact.Phone)
			if name == "" && relation == "" && phone == "" {
				continue
			}
			cleanedContacts = append(cleanedContacts, models.EmergencyContact{
				Name:     name,
				Relation: relation,
				Phone:    phone,
			})
		}
		if len(cleanedContacts) > 0 {
			if raw, err := json.Marshal(cleanedContacts); err == nil {
				emergencyContactsJSON = datatypes.JSON(raw)
			}
		}
	}

	resident := &entities.Resident{
		RoomID:                     req.RoomID,
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		Gender:                     req.Gender,
		Nickname:                   req.Nickname,
		IdCardNumber:               req.IdCardNumber,
		DateOfBirth:                req.DateOfBirth,
		PurposeOfStay:              req.PurposeOfStay,
		CheckInDate:                req.CheckInDate,
		ExpectedCheckOutDate:       req.ExpectedCheckOutDate,
		Status:                     req.Status,
		PreExistingConditions:      req.PreExistingConditions,
		PreExistingConditionsNotes: req.PreExistingConditionsNotes,
		ResucitationStatus:         req.ResucitationStatus,
		SugicalHistory:             req.SugicalHistory,
		PreferredEmergencyHospital: req.PreferredEmergencyHospital,
		EmergencyHospitalPhone:     req.EmergencyHospitalPhone,
		ProfileImage:               req.ProfileImage,
		EmergencyContacts:          emergencyContactsJSON,
	}

	if file != nil {
		defer file.Close()
	}

	createdResident, err := c.emrUsecase.CreateResident(resident, userID, file)
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
func (c *EmrController) GetResidentByIDHandler(ctx *fiber.Ctx) error {
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

	resident, err := c.emrUsecase.GetResidentByID(id, userID)
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
func (c *EmrController) GetResidentByRoomIDHandler(ctx *fiber.Ctx) error {
	roomID := ctx.Query("room_id")

	if roomID == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "room_id query parameter is required",
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

	residents, err := c.emrUsecase.GetResidentByRoomID(roomID, userID)
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
func (c *EmrController) GetAllResidentsHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	residents, err := c.emrUsecase.GetAllResidents(userID)
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

// GetResidentOverview godoc
// @Summary Get Resident Overview
// @Description Retrieve an overview of residents with optional filtering by floor, labels, and status for dashboard display
// @Tags Resident
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param floor query int false "Floor number (optional)"
// @Param label_ids query []string false "List of label IDs (optional)"
// @Param status query string false "Resident status (optional)"
// @Param search query string false "Search by first name, last name, or nickname (optional)"
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 20, max 100)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.ResidentOverviewListResponse} "Resident overview retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents/overview [get]
func (c *EmrController) GetResidentOverviewHandler(ctx *fiber.Ctx) error {
	var req models.ResidentQueryParams
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

	response, err := c.emrUsecase.GetResidentOverview(req, userID)
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
		"message":     "resident overview retrieved successfully",
		"result":      response,
	})
}

// UpdateResident godoc
// @Summary Update Resident
// @Description Partially update an existing resident's information by their unique ID. All fields are optional — only send fields that need to be updated.
// @Tags Resident
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Resident ID"
// @Param request body object{room_id=string,first_name=string,last_name=string,gender=string,nickname=string,id_card_number=string,date_of_birth=string,purpose_of_stay=string,check_in_date=string,expected_check_out_date=string,status=string,pre_existing_conditions=string,pre_existing_conditions_notes=string,resuscitation_status=string,surgical_history=string,preferred_emergency_hospital=string,emergency_hospital_phone=string,labels=[]object{label_name=string,note_text=string}} false "Fields to update (all optional)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Resident updated successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Invalid data"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/residents/{id} [patch]
func (c *EmrController) UpdateResidentByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateResidentRequest
	var file multipart.File

	if form, err := ctx.MultipartForm(); err == nil && (len(form.Value) > 0 || len(form.File) > 0) {
		var parseErr error
		req, file, parseErr = parseResidentUpdateForm(form)
		if parseErr != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     parseErr.Error(),
				"result":      nil,
			})
		}
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

	residentID := ctx.Params("id")

	if file != nil {
		defer file.Close()
	}

	updatedResident, err := c.emrUsecase.UpdateResidentByID(residentID, req, userID, file)
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

func parseResidentCreateForm(form *multipart.Form) (models.CreateResidentRequest, multipart.File, error) {
	var req models.CreateResidentRequest
	var file multipart.File

	if value, ok := getFormValue(form, "first_name"); ok {
		req.FirstName = value
	}
	if value, ok := getFormValue(form, "last_name"); ok {
		req.LastName = value
	}
	if value, ok := getFormValue(form, "gender"); ok {
		req.Gender = value
	}
	if value, ok := getFormValue(form, "status"); ok {
		req.Status = value
	}
	if value, ok := getFormValue(form, "date_of_birth"); ok {
		parsed, err := parseOptionalTime(value)
		if err != nil || parsed == nil {
			return req, nil, err
		}
		req.DateOfBirth = *parsed
	}

	if value, ok := getFormValue(form, "nickname"); ok {
		req.Nickname = &value
	}
	if value, ok := getFormValue(form, "id_card_number"); ok {
		req.IdCardNumber = &value
	}
	if value, ok := getFormValue(form, "purpose_of_stay"); ok {
		req.PurposeOfStay = &value
	}
	if value, ok := getFormValue(form, "check_in_date"); ok {
		req.CheckInDate, _ = parseOptionalTime(value)
	}
	if value, ok := getFormValue(form, "expected_check_out_date"); ok {
		req.ExpectedCheckOutDate, _ = parseOptionalTime(value)
	}
	if value, ok := getFormValue(form, "room_id"); ok {
		req.RoomID = &value
	}
	if value, ok := getFormValue(form, "pre_existing_conditions"); ok {
		req.PreExistingConditions = &value
	}
	if value, ok := getFormValue(form, "pre_existing_conditions_notes"); ok {
		req.PreExistingConditionsNotes = &value
	}
	if value, ok := getFormValue(form, "resuscitation_status"); ok {
		req.ResucitationStatus = &value
	}
	if value, ok := getFormValue(form, "surgical_history"); ok {
		req.SugicalHistory = &value
	}
	if value, ok := getFormValue(form, "preferred_emergency_hospital"); ok {
		req.PreferredEmergencyHospital = &value
	}
	if value, ok := getFormValue(form, "emergency_hospital_phone"); ok {
		req.EmergencyHospitalPhone = &value
	}
	if value, ok := getFormValue(form, "profile_image"); ok {
		req.ProfileImage = &value
	}
	if value, ok := getFormValue(form, "emergency_contacts"); ok {
		var contacts []models.EmergencyContact
		if err := json.Unmarshal([]byte(value), &contacts); err != nil {
			return req, nil, err
		}
		req.EmergencyContacts = contacts
	}

	if files := form.File["profile_image"]; len(files) > 0 {
		fileHeader := files[0]
		opened, err := fileHeader.Open()
		if err != nil {
			return req, nil, err
		}
		file = opened
	}

	if req.FirstName == "" || req.LastName == "" || req.Gender == "" || req.Status == "" || req.DateOfBirth.IsZero() {
		return req, file, fiber.ErrBadRequest
	}

	return req, file, nil
}

func parseResidentUpdateForm(form *multipart.Form) (models.UpdateResidentRequest, multipart.File, error) {
	var req models.UpdateResidentRequest
	var file multipart.File

	if value, ok := getFormValue(form, "room_id"); ok {
		req.RoomID = &value
	}
	if value, ok := getFormValue(form, "first_name"); ok {
		req.FirstName = &value
	}
	if value, ok := getFormValue(form, "last_name"); ok {
		req.LastName = &value
	}
	if value, ok := getFormValue(form, "gender"); ok {
		req.Gender = &value
	}
	if value, ok := getFormValue(form, "nickname"); ok {
		req.Nickname = &value
	}
	if value, ok := getFormValue(form, "id_card_number"); ok {
		req.IdCardNumber = &value
	}
	if value, ok := getFormValue(form, "date_of_birth"); ok {
		req.DateOfBirth, _ = parseOptionalTime(value)
	}
	if value, ok := getFormValue(form, "purpose_of_stay"); ok {
		req.PurposeOfStay = &value
	}
	if value, ok := getFormValue(form, "check_in_date"); ok {
		req.CheckInDate, _ = parseOptionalTime(value)
	}
	if value, ok := getFormValue(form, "expected_check_out_date"); ok {
		req.ExpectedCheckOutDate, _ = parseOptionalTime(value)
	}
	if value, ok := getFormValue(form, "status"); ok {
		req.Status = &value
	}
	if value, ok := getFormValue(form, "pre_existing_conditions"); ok {
		req.PreExistingConditions = &value
	}
	if value, ok := getFormValue(form, "pre_existing_conditions_notes"); ok {
		req.PreExistingConditionsNotes = &value
	}
	if value, ok := getFormValue(form, "resuscitation_status"); ok {
		req.ResucitationStatus = &value
	}
	if value, ok := getFormValue(form, "surgical_history"); ok {
		req.SugicalHistory = &value
	}
	if value, ok := getFormValue(form, "preferred_emergency_hospital"); ok {
		req.PreferredEmergencyHospital = &value
	}
	if value, ok := getFormValue(form, "emergency_hospital_phone"); ok {
		req.EmergencyHospitalPhone = &value
	}
	if value, ok := getFormValue(form, "profile_image"); ok {
		req.ProfileImage = &value
	}
	if value, ok := getFormValue(form, "emergency_contacts"); ok {
		var contacts []models.EmergencyContact
		if err := json.Unmarshal([]byte(value), &contacts); err != nil {
			return req, nil, err
		}
		req.EmergencyContacts = &contacts
	}
	if value, ok := getFormValue(form, "labels"); ok {
		var labels []models.IntakeLabelRequest
		if err := json.Unmarshal([]byte(value), &labels); err != nil {
			return req, nil, err
		}
		req.Labels = labels
	}

	if files := form.File["profile_image"]; len(files) > 0 {
		fileHeader := files[0]
		opened, err := fileHeader.Open()
		if err != nil {
			return req, nil, err
		}
		file = opened
	}

	return req, file, nil
}

func getFormValue(form *multipart.Form, key string) (string, bool) {
	values := form.Value[key]
	if len(values) == 0 {
		return "", false
	}
	return strings.TrimSpace(values[0]), true
}

func parseOptionalTime(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	if parsed, err := time.Parse(time.RFC3339, trimmed); err == nil {
		return &parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02", trimmed); err == nil {
		return &parsed, nil
	}

	return nil, fiber.ErrBadRequest
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
func (c *EmrController) GetRoomByIDHandler(ctx *fiber.Ctx) error {
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

	room, err := c.emrUsecase.GetRoomByID(id, userID)
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
func (c *EmrController) GetAllRoomsHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	rooms, err := c.emrUsecase.GetAllRooms(userID)
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
func (c *EmrController) CreateRoomHandler(ctx *fiber.Ctx) error {
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
func (c *EmrController) UpdateRoomByIDHandler(ctx *fiber.Ctx) error {
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
func (c *EmrController) GetNumberOfResidentsDashboardHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	response, err := c.emrUsecase.GetNumberOfResidentsDashboard(userID)
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
func (c *EmrController) GetResidentGenderStatsDashboardHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	response, err := c.emrUsecase.GetResidentGenderStatsDashboard(userID)
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

// GetVitalSignStatsDashboard godoc
// @Summary Get Vital Sign Stats for Dashboard
// @Description Retrieve today's vital sign summary (normal/abnormal) for dashboard display
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.VitalSignDashboardSummary} "Vital sign stats retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/dashboard/vital-sign-stats [get]
func (c *EmrController) GetVitalSignStatsDashboardHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	response, err := c.emrUsecase.GetVitalSignStatsDashboard(userID)
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
		"message":     "vital sign stats retrieved successfully",
		"result":      response,
	})
}

// GetDrugPlanTimeOfDayStatsDashboard godoc
// @Summary Get Drug Plan Time-of-Day Stats for Dashboard
// @Description Retrieve today's drug plan summary by meal/time-of-day for dashboard display
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]models.DrugPlanTimeOfDayDashboardSummary} "Drug plan time-of-day stats retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/dashboard/drug-plan-time-of-day-stats [get]
func (c *EmrController) GetDrugPlanTimeOfDayStatsDashboardHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	response, err := c.emrUsecase.GetDrugPlanTimeOfDayStatsDashboard(userID)
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
		"message":     "drug plan time-of-day stats retrieved successfully",
		"result":      response,
	})
}

// GetResidentAllergyStatsDashboard godoc
// @Summary Get Resident Allergy Stats for Dashboard
// @Description Retrieve allergy summary for dashboard including allergic count, non-allergic count, and grouped allergy details
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.ResidentAllergyStatsDashboardResponse} "Resident allergy stats retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/dashboard/resident-allergy-stats [get]
func (c *EmrController) GetResidentAllergyStatsDashboardHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	response, err := c.emrUsecase.GetResidentAllergyStatsDashboard(userID)
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
		"message":     "resident allergy stats retrieved successfully",
		"result":      response,
	})
}

// GetResidentDrugAllergyStatsDashboard godoc
// @Summary Get Resident Drug Allergy Stats for Dashboard
// @Description Retrieve drug allergy summary for dashboard including drug-allergic count, non-drug-allergic count, and grouped drug allergy details
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.ResidentDrugAllergyStatsDashboardResponse} "Resident drug allergy stats retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/dashboard/resident-drug-allergy-stats [get]
func (c *EmrController) GetResidentDrugAllergyStatsDashboardHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	response, err := c.emrUsecase.GetResidentDrugAllergyStatsDashboard(userID)
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
		"message":     "resident drug allergy stats retrieved successfully",
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
func (c *EmrController) GetAllIntakeLabelsHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	labels, err := c.emrUsecase.GetAllIntakeLabels(userID)
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

// CreateIntakeLabelMaster godoc
// @Summary Create Intake Label Master
// @Description Create a new intake label master entry
// @Tags Intake
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{label_name=string} true "Label name"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object} "Intake label created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/intake-labels/master [post]
func (c *EmrController) CreateIntakeLabelMasterHandler(ctx *fiber.Ctx) error {
	var req models.IntakeLabelRequest
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

	newLabel, err := c.emrUsecase.CreateIntakeLabel(&entities.IntakeLabels{
		ID:        uuid.New().String(),
		LabelName: req.LabelName,
	})
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
		"message":     "intake label created successfully",
		"result":      newLabel,
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
func (c *EmrController) GetResidentLabelsByResidentIDHandler(ctx *fiber.Ctx) error {
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

	labels, err := c.emrUsecase.GetResidentLabelsByResidentID(residentID, userID)
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
func (c *EmrController) CreateIntakeLabelByResidentIDHandler(ctx *fiber.Ctx) error {
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

// GetAllAllergies godoc
// @Summary Get All Allergies
// @Description Retrieve all allergy master records
// @Tags Allergy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Allergies retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/allergies/all [get]
func (c *EmrController) GetAllAllergiesHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	allergies, err := c.emrUsecase.GetAllAllergies(userID)
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
		"message":     "allergies retrieved successfully",
		"result":      allergies,
	})
}

// GetResidentAllergiesByResidentID godoc
// @Summary Get Resident Allergies by Resident ID
// @Description Retrieve all allergies associated with a specific resident using resident_id query parameter
// @Tags Allergy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Resident allergies retrieved successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing resident_id query parameter"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/allergies [get]
func (c *EmrController) GetResidentAllergiesByResidentIDHandler(ctx *fiber.Ctx) error {
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

	allergies, err := c.emrUsecase.GetResidentAllergiesByResidentID(residentID, userID)
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
		"message":     "resident allergies retrieved successfully",
		"result":      allergies,
	})
}

// GetAllResidentAllergies godoc
// @Summary Get All Resident Allergies
// @Description Retrieve all residents with first name, last name, resident ID, and their allergies list
// @Tags Allergy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]models.ResidentAllergyListResponse} "All resident allergies retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/allergies/residents/all [get]
func (c *EmrController) GetAllResidentAllergiesHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetAllResidentAllergies(userID)
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
		"message":     "all resident allergies retrieved successfully",
		"result":      result,
	})
}

// CreateAllergyByResidentID godoc
// @Summary Create Allergies for Resident
// @Description Create one or more allergy records for a specific resident. New allergy master values are auto-created when they do not exist.
// @Tags Allergy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateAllergyByResidentRequest true "Resident ID and array of allergies"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=[]object} "Allergies created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields or invalid data"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing or invalid authentication"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/allergies [post]
func (c *EmrController) CreateAllergyByResidentIDHandler(ctx *fiber.Ctx) error {
	var req models.CreateAllergyByResidentRequest

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

	result, err := c.emrUsecase.CreateAllergyByResidentID(req.ResidentID, req.Allergies, userID)
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
		"message":     "allergies created successfully",
		"result":      result,
	})
}

// GetAllDrugAllergies godoc
// @Summary Get All Drug Allergies
// @Description Retrieve all drug allergy master records
// @Tags Drug Allergy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Drug allergies retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/drug-allergies/all [get]
func (c *EmrController) GetAllDrugAllergiesHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	drugAllergies, err := c.emrUsecase.GetAllDrugAllergies(userID)
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
		"message":     "drug allergies retrieved successfully",
		"result":      drugAllergies,
	})
}

// GetResidentDrugAllergiesByResidentID godoc
// @Summary Get Resident Drug Allergies by Resident ID
// @Description Retrieve all drug allergies associated with a specific resident using resident_id query parameter
// @Tags Drug Allergy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Resident drug allergies retrieved successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing resident_id query parameter"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/drug-allergies [get]
func (c *EmrController) GetResidentDrugAllergiesByResidentIDHandler(ctx *fiber.Ctx) error {
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

	drugAllergies, err := c.emrUsecase.GetResidentDrugAllergiesByResidentID(residentID, userID)
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
		"message":     "resident drug allergies retrieved successfully",
		"result":      drugAllergies,
	})
}

// GetAllResidentDrugAllergies godoc
// @Summary Get All Resident Drug Allergies
// @Description Retrieve all residents with first name, last name, resident ID, and their drug allergies list
// @Tags Drug Allergy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]models.ResidentDrugAllergyListResponse} "All resident drug allergies retrieved successfully"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/drug-allergies/residents/all [get]
func (c *EmrController) GetAllResidentDrugAllergiesHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetAllResidentDrugAllergies(userID)
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
		"message":     "all resident drug allergies retrieved successfully",
		"result":      result,
	})
}

// CreateDrugAllergyByResidentID godoc
// @Summary Create Drug Allergies for Resident
// @Description Create one or more drug allergy records for a specific resident. New drug allergy master values are auto-created when they do not exist.
// @Tags Drug Allergy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateDrugAllergyByResidentRequest true "Resident ID and array of drug allergies"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=[]object} "Drug allergies created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields or invalid data"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing or invalid authentication"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/drug-allergies [post]
func (c *EmrController) CreateDrugAllergyByResidentIDHandler(ctx *fiber.Ctx) error {
	var req models.CreateDrugAllergyByResidentRequest

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

	result, err := c.emrUsecase.CreateDrugAllergyByResidentID(req.ResidentID, req.DrugAllergies, userID)
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
		"message":     "drug allergies created successfully",
		"result":      result,
	})
}

// CreateVitalSign godoc
// @Summary Create Vital Sign
// @Description Create a new vital sign entry for a resident on a selected date and time slot
// @Tags VitalSign
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateVitalSignRequest true "Vital sign information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.VitalSign} "Vital sign created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/vital-signs [post]
func (c *EmrController) CreateVitalSignHandler(ctx *fiber.Ctx) error {
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
		TimeOfDay:              req.TimeOfDay,
		Temperature:            req.Temperature,
		HeartRate:              req.HeartRate,
		BreathingRate:          req.BreathingRate,
		BloodPressureSystolic:  req.BloodPressureSystolic,
		BloodPressureDiastolic: req.BloodPressureDiastolic,
		OxygenSaturation:       req.OxygenSaturation,
	}
	createdVitalSign, err := c.emrUsecase.CreateVitalSign(vitalSign, req.Date, userID)
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

// @Summary Get Vital Signs Overview
// @Description Get vital signs by selected date with optional filters. vitalsign_status: 'all' (default), 'normal', or 'abnormal'
// @Tags VitalSign
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string true "Selected date (YYYY-MM-DD)"
// @Param time_of_day query string false "Time of day filter" Enums(เช้า,สาย,บ่าย,เย็น,กลางคืน,morning,late_morning,afternoon,evening,night)
// @Param floor query int false "Filter by floor"
// @Param label_ids query []string false "Filter by label IDs"
// @Param vitalsign_status query string false "Filter by vital sign status" Enums(all, normal, abnormal)
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 20, max 100)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.VitalSignsOverviewResponse}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/vital-signs/overview [get]
func (c *EmrController) GetVitalSignsOverviewHandler(ctx *fiber.Ctx) error {
	var req models.VitalSignQueryParams
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

	response, err := c.emrUsecase.GetVitalSignsOverview(req, userID)
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
		"message":     "vital signs overview retrieved successfully",
		"result":      response,
	})
}

// @Summary Get Vital Signs by Resident ID
// @Description Retrieve vital signs on a selected date for a specific resident, with an option to get only the latest entry. is_latest must be 'true' or 'false'
// @Tags VitalSign
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Param date query string true "Selected date (YYYY-MM-DD)"
// @Param is_latest query string false "Retrieve only the latest vital sign entry ('true' or 'false')" Enums(true, false)
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/vital-signs/resident [get]
func (c *EmrController) GetVitalSignsByResidentHandler(ctx *fiber.Ctx) error {
	var req models.VitalSignQueryParams
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	residentID := ""
	if req.ResidentID != nil {
		residentID = *req.ResidentID
	}
	isLatest := ctx.Query("is_latest", "false")
	selectedDate := ""
	if req.Date != nil {
		selectedDate = *req.Date
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

	vitalSigns, err := c.emrUsecase.GetVitalSignsByResident(residentID, selectedDate, isLatest, userID)
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
		"message":     "vital signs retrieved successfully",
		"result":      vitalSigns,
	})
}

// @Summary Get Vital Signs by Room ID
// @Description Retrieve vital signs today for all residents in a specific room, with an option to get only the latest entry per resident. is_latest must be 'true' or 'false'
// @Tags VitalSign
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param room_id query string true "Room ID"
// @Param is_latest query string true "Retrieve only the latest vital sign entry per resident ('true' or 'false')" Enums(true, false)
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/vital-signs/room [get]
func (c *EmrController) GetRoomVitalSignsHandler(ctx *fiber.Ctx) error {
	roomID := ctx.Query("room_id")
	isLatest := ctx.Query("is_latest")

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	vitalSigns, err := c.emrUsecase.GetRoomVitalSigns(roomID, isLatest, userID)
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
		"message":     "vital signs retrieved successfully",
		"result":      vitalSigns,
	})
}

// @Summary Get Vital Signs History by Resident ID
// @Description Retrieve the full history of vital signs for a specific resident
// @Tags VitalSign
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/vital-signs/history/{resident_id} [get]
func (c *EmrController) GetVitalSignsHistoryHandler(ctx *fiber.Ctx) error {
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

	vitalSigns, err := c.emrUsecase.GetVitalSignsHistory(residentID, userID)
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
		"message":     "vital signs history retrieved successfully",
		"result":      vitalSigns,
	})
}

// @Summary Update Vital Sign by ID
// @Description Update an existing vital sign entry by its unique ID. Only send fields that need to be updated.
// @Tags VitalSign
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vital Sign ID"
// @Param request body object{temperature=number,heart_rate=int,breathing_rate=int,blood_pressure_systolic=int,blood_pressure_diastolic=int,oxygen_saturation=int} true "Fields to update (all optional)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Vital sign updated successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Invalid data"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/vital-signs/{id} [patch]
func (c *EmrController) UpdateVitalSignByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateVitalSignRequest

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

	vitalSignID := ctx.Params("id")
	vitalSign := &entities.VitalSign{
		Temperature:            req.Temperature,
		HeartRate:              req.HeartRate,
		BreathingRate:          req.BreathingRate,
		BloodPressureSystolic:  req.BloodPressureSystolic,
		BloodPressureDiastolic: req.BloodPressureDiastolic,
		OxygenSaturation:       req.OxygenSaturation,
	}

	updatedVitalSign, err := c.emrUsecase.UpdateVitalSignByID(vitalSignID, vitalSign, userID)
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
		"message":     "vital sign updated successfully",
		"result":      updatedVitalSign,
	})
}

// CreateLaboratoryValue godoc
// @Summary Create Laboratory Value
// @Description Create a new laboratory value entry for a resident. urine_type must be 'ml' or 'times' and must be provided together with urine_output.
// @Tags LaboratoryValue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateLaboratoryValueRequest true "Laboratory value information"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=entities.LaboratoryValue} "Laboratory value created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/laboratory-values [post]
func (c *EmrController) CreateLaboratoryValueHandler(ctx *fiber.Ctx) error {
	var req models.CreateLaboratoryValueRequest

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

	laboratoryValue := &entities.LaboratoryValue{
		ResidentID:   req.ResidentID,
		BloodGlucose: req.BloodGlucose,
		FluidIn:      req.FluidIn,
		FluidOut:     req.FluidOut,
		UrineOutput:  req.UrineOutput,
		UrineType:    req.UrineType,
		Stool:        req.Stool,
		DiaperChange: req.DiaperChange,
	}

	createdLaboratoryValue, err := c.emrUsecase.CreateLaboratoryValue(laboratoryValue, req.Date, req.TimeOfDay, userID)
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
		"message":     "laboratory value created successfully",
		"result":      createdLaboratoryValue,
	})
}

// @Summary Get Laboratory Values Overview
// @Description Get laboratory values for a specific date with optional filters. laboratory_value_status: 'all' (default), 'normal', or 'abnormal'. urine_type must be 'ml' or 'times' and must be provided together with urine_output.
// @Tags LaboratoryValue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string true "Filter by date (format: YYYY-MM-DD)"
// @Param time_of_day query string false "Filter by time of day (เช้า/สาย/บ่าย/เย็น/กลางคืน)"
// @Param floor query int false "Filter by floor"
// @Param label_ids query []string false "Filter by label IDs"
// @Param urine_type query string false "Filter by urine type (must be 'ml' or 'times' if urine_output is provided)" Enums(ml, times)
// @Param laboratory_value_status query string false "Filter by laboratory value status" Enums(all, normal, abnormal)
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 20, max 100)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=models.LaboratoryValuesOverviewResponse}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/laboratory-values/overview [get]
func (c *EmrController) GetLaboratoryValuesOverviewHandler(ctx *fiber.Ctx) error {
	var req models.LaboratoryValueQueryParams
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

	response, err := c.emrUsecase.GetLaboratoryValuesOverview(req, userID)
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
		"message":     "laboratory values overview retrieved successfully",
		"result":      response,
	})
}

// @Summary Get Laboratory Values by Resident ID
// @Description Retrieve laboratory values for a specific resident on a specific date, with an option to get only the latest entry. is_latest must be 'true' or 'false'
// @Tags LaboratoryValue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Param date query string true "Filter by date (format: YYYY-MM-DD)"
// @Param is_latest query string true "Retrieve only the latest laboratory value entry ('true' or 'false')" Enums(true, false)
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/laboratory-values/resident [get]
func (c *EmrController) GetLaboratoryValuesByResidentHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	dateStr := ctx.Query("date")
	isLatest := ctx.Query("is_latest")

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	laboratoryValues, err := c.emrUsecase.GetLaboratoryValuesByResident(residentID, dateStr, isLatest, userID)
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
		"message":     "laboratory values retrieved successfully",
		"result":      laboratoryValues,
	})
}

// @Summary Get Laboratory Values by Room ID
// @Description Retrieve laboratory values today for all residents in a specific room, with an option to get only the latest entry per resident. is_latest must be 'true' or 'false'
// @Tags LaboratoryValue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param room_id query string true "Room ID"
// @Param is_latest query string true "Retrieve only the latest laboratory value entry per resident ('true' or 'false')" Enums(true, false)
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/laboratory-values/room [get]
func (c *EmrController) GetRoomLaboratoryValuesHandler(ctx *fiber.Ctx) error {
	roomID := ctx.Query("room_id")
	isLatest := ctx.Query("is_latest")

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	laboratoryValues, err := c.emrUsecase.GetRoomLaboratoryValues(roomID, isLatest, userID)
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
		"message":     "laboratory values retrieved successfully",
		"result":      laboratoryValues,
	})
}

// @Summary Get Laboratory Values History by Resident ID
// @Description Retrieve the full history of laboratory values for a specific resident
// @Tags LaboratoryValue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object}
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any}
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any}
// @Router /api/emr/laboratory-values/history/{resident_id} [get]
func (c *EmrController) GetLaboratoryValuesHistoryHandler(ctx *fiber.Ctx) error {
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

	laboratoryValues, err := c.emrUsecase.GetLaboratoryValuesHistory(residentID, userID)
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
		"message":     "laboratory values history retrieved successfully",
		"result":      laboratoryValues,
	})
}

// @Summary Get Urine Output Sum by Resident ID
// @Description Calculate the total urine output (ml and times) for a specific resident. Always returns both total_ml and total_times regardless of how data was recorded. Optional filters: start_date, end_date.
// @Tags LaboratoryValue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id path string true "Resident ID"
// @Param start_date query string false "Filter from date (RFC3339)"
// @Param end_date query string false "Filter to date (RFC3339)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object{resident_id=string,total_ml=number,total_times=number}} "Urine output sum retrieved successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Invalid parameters"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/laboratory-values/urine-output-sum/{resident_id} [get]
func (c *EmrController) GetUrineOutputSumByResidentIDHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Params("resident_id")

	var req models.LaboratoryValueQueryParams
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
	urineOutputSum, err := c.emrUsecase.GetUrineOutputSumByResidentID(residentID, req, userID)
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
		"message":     "urine output sum retrieved successfully",
		"result":      urineOutputSum,
	})
}

// @Summary Update Laboratory Value by ID
// @Description Update an existing laboratory value entry by its unique ID. Only send fields that need to be updated. urine_type must be 'ml' or 'times' and must be provided together with urine_output.
// @Tags LaboratoryValue
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Laboratory Value ID"
// @Param request body object{blood_glucose=number,fluid_in=number,fluid_out=number,urine_output=number,urine_type=string,stool=string,diaper_change=boolean} true "Fields to update (all optional)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Laboratory value updated successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Invalid data"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/laboratory-values/{id} [patch]
func (c *EmrController) UpdateLaboratoryValueByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateLaboratoryValueRequest

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

	laboratoryValueID := ctx.Params("id")
	laboratoryValue := &entities.LaboratoryValue{
		BloodGlucose: req.BloodGlucose,
		FluidIn:      req.FluidIn,
		FluidOut:     req.FluidOut,
		UrineOutput:  req.UrineOutput,
		UrineType:    req.UrineType,
		Stool:        req.Stool,
		DiaperChange: req.DiaperChange,
	}

	updatedLaboratoryValue, err := c.emrUsecase.UpdateLaboratoryValueByID(laboratoryValueID, laboratoryValue, userID)
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
		"message":     "laboratory value updated successfully",
		"result":      updatedLaboratoryValue,
	})
}

func (c *EmrController) CreateNurseNoteHandler(ctx *fiber.Ctx) error {
	var req models.CreateNurseNoteRequest
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

	note := &entities.NurseNote{
		ResidentID: req.ResidentID,
		Category:   req.Category,
		Content:    req.Content,
		Priority:   req.Priority,
		SendNote:   req.SendNote,
	}

	created, err := c.emrUsecase.CreateNurseNote(note, userID)
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
		"message":     "nurse note created successfully",
		"result":      created,
	})
}

func (c *EmrController) GetNurseNotesOverviewHandler(ctx *fiber.Ctx) error {
	selectedDate := ctx.Query("date")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetNurseNotesOverview(selectedDate, userID)
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
		"message":     "nurse notes overview retrieved successfully",
		"result":      result,
	})
}

func (c *EmrController) GetNurseNotesByResidentHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	selectedDate := ctx.Query("date")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetNurseNotesByResidentID(residentID, selectedDate, userID)
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
		"message":     "nurse notes retrieved successfully",
		"result":      result,
	})
}

func (c *EmrController) UpdateNurseNoteByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateNurseNoteRequest
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

	note := &entities.NurseNote{}
	if req.Category != nil {
		note.Category = *req.Category
	}
	if req.Content != nil {
		note.Content = *req.Content
	}
	if req.Priority != nil {
		note.Priority = *req.Priority
	}
	if req.SendNote != nil {
		note.SendNote = *req.SendNote
	}

	updated, err := c.emrUsecase.UpdateNurseNoteByID(ctx.Params("id"), note, userID)
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
		"message":     "nurse note updated successfully",
		"result":      updated,
	})
}

func (c *EmrController) DeleteNurseNoteByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	if err := c.emrUsecase.DeleteNurseNoteByID(ctx.Params("id"), userID); err != nil {
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
		"message":     "nurse note deleted successfully",
		"result":      nil,
	})
}

func (c *EmrController) CreateWoundCareNoteHandler(ctx *fiber.Ctx) error {
	var req models.CreateWoundCareNoteRequest

	var imageFile multipart.File
	if strings.HasPrefix(strings.ToLower(ctx.Get("Content-Type")), "multipart/form-data") {
		req = models.CreateWoundCareNoteRequest{
			ResidentID: ctx.FormValue("resident_id"),
			Location:   ctx.FormValue("location"),
			WoundType:  ctx.FormValue("wound_type"),
			Size:       optionalFormString(ctx.FormValue("size")),
			Treatment:  optionalFormString(ctx.FormValue("treatment")),
			Supplies:   optionalFormString(ctx.FormValue("supplies")),
			Status:     optionalFormString(ctx.FormValue("status")),
			ImageURL:   optionalFormString(ctx.FormValue("image_url")),
			Note:       optionalFormString(ctx.FormValue("note")),
		}

		fileHeader, err := ctx.FormFile("image")
		if err == nil && fileHeader != nil {
			file, openErr := fileHeader.Open()
			if openErr != nil {
				return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
					"status":      fiber.ErrBadRequest.Message,
					"status_code": fiber.ErrBadRequest.Code,
					"message":     "failed to open uploaded image: " + openErr.Error(),
					"result":      nil,
				})
			}
			imageFile = file
			defer imageFile.Close()
		}
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

	note := &entities.WoundCareNote{
		ResidentID: req.ResidentID,
		Location:   req.Location,
		WoundType:  req.WoundType,
		Size:       req.Size,
		Treatment:  req.Treatment,
		Supplies:   req.Supplies,
		Status:     req.Status,
		ImageURL:   req.ImageURL,
		Note:       req.Note,
	}

	created, err := c.emrUsecase.CreateWoundCareNote(note, userID, imageFile)
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
		"message":     "wound care note created successfully",
		"result":      created,
	})
}

func (c *EmrController) GetWoundCareNotesOverviewHandler(ctx *fiber.Ctx) error {
	selectedDate := ctx.Query("date")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetWoundCareNotesOverview(selectedDate, userID)
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
		"message":     "wound care notes overview retrieved successfully",
		"result":      result,
	})
}

func (c *EmrController) GetWoundCareNotesByResidentHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	selectedDate := ctx.Query("date")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetWoundCareNotesByResidentID(residentID, selectedDate, userID)
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
		"message":     "wound care notes retrieved successfully",
		"result":      result,
	})
}

func (c *EmrController) UpdateWoundCareNoteByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateWoundCareNoteRequest

	var imageFile multipart.File
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

		if v, ok := getMultipartField(form, "location"); ok {
			req.Location = &v
		}
		if v, ok := getMultipartField(form, "wound_type"); ok {
			req.WoundType = &v
		}
		if v, ok := getMultipartField(form, "size"); ok {
			req.Size = &v
		}
		if v, ok := getMultipartField(form, "treatment"); ok {
			req.Treatment = &v
		}
		if v, ok := getMultipartField(form, "supplies"); ok {
			req.Supplies = &v
		}
		if v, ok := getMultipartField(form, "status"); ok {
			req.Status = &v
		}
		if v, ok := getMultipartField(form, "image_url"); ok {
			req.ImageURL = &v
		}
		if v, ok := getMultipartField(form, "note"); ok {
			req.Note = &v
		}

		fileHeader, err := ctx.FormFile("image")
		if err == nil && fileHeader != nil {
			file, openErr := fileHeader.Open()
			if openErr != nil {
				return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
					"status":      fiber.ErrBadRequest.Message,
					"status_code": fiber.ErrBadRequest.Code,
					"message":     "failed to open uploaded image: " + openErr.Error(),
					"result":      nil,
				})
			}
			imageFile = file
			defer imageFile.Close()
		}
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

	note := &entities.WoundCareNote{
		Location:  derefString(req.Location),
		WoundType: derefString(req.WoundType),
		Size:      req.Size,
		Treatment: req.Treatment,
		Supplies:  req.Supplies,
		Status:    req.Status,
		ImageURL:  req.ImageURL,
		Note:      req.Note,
	}

	updated, err := c.emrUsecase.UpdateWoundCareNoteByID(ctx.Params("id"), note, userID, imageFile)
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
		"message":     "wound care note updated successfully",
		"result":      updated,
	})
}

func (c *EmrController) DeleteWoundCareNoteByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	if err := c.emrUsecase.DeleteWoundCareNoteByID(ctx.Params("id"), userID); err != nil {
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
		"message":     "wound care note deleted successfully",
		"result":      nil,
	})
}

func (c *EmrController) CreateRelativeNoteHandler(ctx *fiber.Ctx) error {
	var req models.CreateRelativeNoteRequest
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

	note := &entities.RelativeNote{
		ResidentID: req.ResidentID,
		Relation:   req.Relation,
		Content:    req.Content,
		SendNote:   req.SendNote,
	}

	created, err := c.emrUsecase.CreateRelativeNote(note, userID)
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
		"message":     "relative note created successfully",
		"result":      created,
	})
}

func (c *EmrController) GetRelativeNotesOverviewHandler(ctx *fiber.Ctx) error {
	selectedDate := ctx.Query("date")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetRelativeNotesOverview(selectedDate, userID)
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
		"message":     "relative notes overview retrieved successfully",
		"result":      result,
	})
}

func (c *EmrController) GetRelativeNotesByResidentHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	selectedDate := ctx.Query("date")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetRelativeNotesByResidentID(residentID, selectedDate, userID)
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
		"message":     "relative notes retrieved successfully",
		"result":      result,
	})
}

func (c *EmrController) UpdateRelativeNoteByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateRelativeNoteRequest
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

	note := &entities.RelativeNote{}
	if req.Relation != nil {
		note.Relation = *req.Relation
	}
	if req.Content != nil {
		note.Content = *req.Content
	}
	if req.SendNote != nil {
		note.SendNote = *req.SendNote
	}

	updated, err := c.emrUsecase.UpdateRelativeNoteByID(ctx.Params("id"), note, userID)
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
		"message":     "relative note updated successfully",
		"result":      updated,
	})
}

func (c *EmrController) DeleteRelativeNoteByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	if err := c.emrUsecase.DeleteRelativeNoteByID(ctx.Params("id"), userID); err != nil {
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
		"message":     "relative note deleted successfully",
		"result":      nil,
	})
}

func (c *EmrController) IssueRelativeMagicLinkHandler(ctx *fiber.Ctx) error {
	var req models.IssueRelativeMagicLinkRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if strings.TrimSpace(req.ResidentID) == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "resident_id is required",
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

	result, err := c.emrUsecase.IssueRelativeMagicLink(req.ResidentID, userID)
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
		"message":     "relative magic link issued successfully",
		"result":      result,
	})
}

func (c *EmrController) GetRelativeMagicLinkHandler(ctx *fiber.Ctx) error {
	residentID := strings.TrimSpace(ctx.Query("resident_id"))
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

	result, err := c.emrUsecase.GetRelativeMagicLink(residentID, userID)
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
		"message":     "relative magic link retrieved successfully",
		"result":      result,
	})
}

func (c *EmrController) RelativePortalLoginHandler(ctx *fiber.Ctx) error {
	var req models.RelativePortalLoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.RelativePortalLogin(req)
	if err != nil {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "relative login successful",
		"result":      result,
	})
}

func (c *EmrController) GetRelativeDashboardHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	dateInput := strings.TrimSpace(ctx.Query("date"))
	result, err := c.emrUsecase.GetRelativeDashboard(userID, dateInput)
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
		"message":     "relative dashboard retrieved successfully",
		"result":      result,
	})
}

func (c *EmrController) GetRelativePatientInfoHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetRelativePatientInfo(userID)
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
		"message":     "relative patient info retrieved successfully",
		"result":      result,
	})
}

// CreateDoctorOrderHandler godoc
// @Summary Create Doctor Order
// @Description Create a new doctor order for a resident
// @Tags DoctorOrder
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateDoctorOrderRequest true "Doctor order payload"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object} "Doctor order created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/doctor-orders [post]
func (c *EmrController) CreateDoctorOrderHandler(ctx *fiber.Ctx) error {
	var req models.CreateDoctorOrderRequest
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

	order := &entities.DoctorOrder{
		ResidentID: req.ResidentID,
		OrderDate:  req.OrderDate,
		OrderType:  req.OrderType,
		Title:      req.Title,
		Details:    req.Details,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Frequency:  req.Frequency,
		OrderedBy:  req.OrderedBy,
	}

	created, err := c.emrUsecase.CreateDoctorOrder(order, userID)
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
		"message":     "doctor order created successfully",
		"result":      created,
	})
}

// GetDoctorOrdersOverviewHandler godoc
// @Summary Get Doctor Orders Overview
// @Description Retrieve doctor orders for overview screen on a selected date
// @Tags DoctorOrder
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string true "Selected date (YYYY-MM-DD)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Doctor orders overview retrieved successfully"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/doctor-orders/overview [get]
func (c *EmrController) GetDoctorOrdersOverviewHandler(ctx *fiber.Ctx) error {
	selectedDate := ctx.Query("date")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetDoctorOrdersOverview(selectedDate, userID)
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
		"message":     "doctor orders overview retrieved successfully",
		"result":      result,
	})
}

// GetDoctorOrdersByResidentHandler godoc
// @Summary Get Doctor Orders By Resident
// @Description Retrieve doctor orders for a specific resident on a selected date
// @Tags DoctorOrder
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resident_id query string true "Resident ID"
// @Param date query string true "Selected date (YYYY-MM-DD)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=[]object} "Doctor orders retrieved successfully"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/doctor-orders/resident/all [get]
func (c *EmrController) GetDoctorOrdersByResidentHandler(ctx *fiber.Ctx) error {
	residentID := ctx.Query("resident_id")
	selectedDate := ctx.Query("date")
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	result, err := c.emrUsecase.GetDoctorOrdersByResidentID(residentID, selectedDate, userID)
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
		"message":     "doctor orders retrieved successfully",
		"result":      result,
	})
}

// UpdateDoctorOrderByIDHandler godoc
// @Summary Update Doctor Order By ID
// @Description Update doctor order fields by ID
// @Tags DoctorOrder
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Doctor Order ID"
// @Param request body models.UpdateDoctorOrderRequest true "Doctor order payload"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "Doctor order updated successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/doctor-orders/{id} [patch]
func (c *EmrController) UpdateDoctorOrderByIDHandler(ctx *fiber.Ctx) error {
	var req models.UpdateDoctorOrderRequest
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

	order := &entities.DoctorOrder{
		OrderDate: req.OrderDate,
		OrderType: req.OrderType,
		Details:   req.Details,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Frequency: req.Frequency,
		OrderedBy: req.OrderedBy,
	}
	if req.Title != nil {
		order.Title = *req.Title
	}

	updated, err := c.emrUsecase.UpdateDoctorOrderByID(ctx.Params("id"), order, userID)
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
		"message":     "doctor order updated successfully",
		"result":      updated,
	})
}

// DeleteDoctorOrderByIDHandler godoc
// @Summary Delete Doctor Order By ID
// @Description Delete a doctor order by ID
// @Tags DoctorOrder
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Doctor Order ID"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any} "Doctor order deleted successfully"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/emr/doctor-orders/{id} [delete]
func (c *EmrController) DeleteDoctorOrderByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	if err := c.emrUsecase.DeleteDoctorOrderByID(ctx.Params("id"), userID); err != nil {
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
		"message":     "doctor order deleted successfully",
		"result":      nil,
	})
}

func derefString(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

func optionalFormString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func getMultipartField(form *multipart.Form, key string) (string, bool) {
	if form == nil || form.Value == nil {
		return "", false
	}
	values, ok := form.Value[key]
	if !ok || len(values) == 0 {
		return "", false
	}
	return values[0], true
}
