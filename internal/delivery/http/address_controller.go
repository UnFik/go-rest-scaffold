package http

import (
	"go-rest-scaffold/internal/delivery/http/middleware"
	"go-rest-scaffold/internal/model"
	"go-rest-scaffold/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	UseCase *usecase.AddressUseCase
	Log     *logrus.Logger
}

func NewAddressController(useCase *usecase.AddressUseCase, log *logrus.Logger) *AddressController {
	return &AddressController{
		Log:     log,
		UseCase: useCase,
	}
}

// Create godoc
// @Summary      Create a new address
// @Description  Create a new address for a specific contact
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        contactId path string true "Contact ID"
// @Param        request body model.CreateAddressRequest true "Address creation details"
// @Success      200 {object} object{data=model.AddressResponse} "Successfully created address"
// @Failure      400 {object} object{errors=string} "Invalid request body"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      404 {object} object{errors=string} "Contact not found"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts/{contactId}/addresses [post]
func (c *AddressController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateAddressRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ContactId = ctx.Params("contactId")

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create address")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

// List godoc
// @Summary      List addresses
// @Description  Get all addresses for a specific contact
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        contactId path string true "Contact ID"
// @Success      200 {object} object{data=[]model.AddressResponse} "List of addresses"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      404 {object} object{errors=string} "Contact not found"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts/{contactId}/addresses [get]
func (c *AddressController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")

	request := &model.ListAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
	}

	responses, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list addresses")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.AddressResponse]{Data: responses})
}

// Get godoc
// @Summary      Get an address
// @Description  Get a specific address by ID for a contact
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        contactId path string true "Contact ID"
// @Param        addressId path string true "Address ID"
// @Success      200 {object} object{data=model.AddressResponse} "Address details"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      404 {object} object{errors=string} "Address not found"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts/{contactId}/addresses/{addressId} [get]
func (c *AddressController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")
	addressId := ctx.Params("addressId")

	request := &model.GetAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to get address")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

// Update godoc
// @Summary      Update an address
// @Description  Update a specific address by ID for a contact
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        contactId path string true "Contact ID"
// @Param        addressId path string true "Address ID"
// @Param        request body model.UpdateAddressRequest true "Address update details"
// @Success      200 {object} object{data=model.AddressResponse} "Successfully updated address"
// @Failure      400 {object} object{errors=string} "Invalid request body"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      404 {object} object{errors=string} "Address not found"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts/{contactId}/addresses/{addressId} [put]
func (c *AddressController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateAddressRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ContactId = ctx.Params("contactId")
	request.ID = ctx.Params("addressId")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update address")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{Data: response})
}

// Delete godoc
// @Summary      Delete an address
// @Description  Delete a specific address by ID for a contact
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        contactId path string true "Contact ID"
// @Param        addressId path string true "Address ID"
// @Success      200 {object} object{data=bool} "Successfully deleted address"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      404 {object} object{errors=string} "Address not found"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /contacts/{contactId}/addresses/{addressId} [delete]
func (c *AddressController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")
	addressId := ctx.Params("addressId")

	request := &model.DeleteAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("failed to delete address")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
