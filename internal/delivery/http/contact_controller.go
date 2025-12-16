package http

import (
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ContactController struct {
	UseCase *usecase.ContactUseCase
	Log     *logrus.Logger
}

func NewContactController(useCase *usecase.ContactUseCase, log *logrus.Logger) *ContactController {
	return &ContactController{
		UseCase: useCase,
		Log:     log,
	}
}

// Create godoc
// @Summary      Create a new contact
// @Description  Create a new contact for the authenticated user
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.CreateContactRequest true "Contact creation details"
// @Success      200 {object} object{data=model.ContactResponse} "Successfully created contact"
// @Failure      400 {object} object{errors=string} "Invalid request body"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts [post]
func (c *ContactController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateContactRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}
	request.UserId = auth.ID

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

// List godoc
// @Summary      List contacts
// @Description  Search and list contacts for the authenticated user with pagination
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name query string false "Filter by name"
// @Param        email query string false "Filter by email"
// @Param        phone query string false "Filter by phone"
// @Param        page query int false "Page number" default(1)
// @Param        size query int false "Page size" default(10)
// @Success      200 {object} object{data=[]model.ContactResponse,paging=model.PageMetadata} "List of contacts with pagination"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts [get]
func (c *ContactController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.SearchContactRequest{
		UserId: auth.ID,
		Name:   ctx.Query("name", ""),
		Email:  ctx.Query("email", ""),
		Phone:  ctx.Query("phone", ""),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching contact")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.ContactResponse]{
		Data:   responses,
		Paging: paging,
	})
}

// Get godoc
// @Summary      Get a contact
// @Description  Get a specific contact by ID for the authenticated user
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        contactId path string true "Contact ID"
// @Success      200 {object} object{data=model.ContactResponse} "Contact details"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      404 {object} object{errors=string} "Contact not found"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts/{contactId} [get]
func (c *ContactController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetContactRequest{
		UserId: auth.ID,
		ID:     ctx.Params("contactId"),
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

// Update godoc
// @Summary      Update a contact
// @Description  Update a specific contact by ID for the authenticated user
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        contactId path string true "Contact ID"
// @Param        request body model.UpdateContactRequest true "Contact update details"
// @Success      200 {object} object{data=model.ContactResponse} "Successfully updated contact"
// @Failure      400 {object} object{errors=string} "Invalid request body"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      404 {object} object{errors=string} "Contact not found"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts/{contactId} [put]
func (c *ContactController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateContactRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ID = ctx.Params("contactId")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ContactResponse]{Data: response})
}

// Delete godoc
// @Summary      Delete a contact
// @Description  Delete a specific contact by ID for the authenticated user
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        contactId path string true "Contact ID"
// @Success      200 {object} object{data=bool} "Successfully deleted contact"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      404 {object} object{errors=string} "Contact not found"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts/{contactId} [delete]
func (c *ContactController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")

	request := &model.DeleteContactRequest{
		UserId: auth.ID,
		ID:     contactId,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
