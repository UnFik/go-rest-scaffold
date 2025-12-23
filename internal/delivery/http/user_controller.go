package http

import (
	"go-rest-scaffold/internal/delivery/http/middleware"
	"go-rest-scaffold/internal/model"
	"go-rest-scaffold/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with ID, name, and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body model.RegisterUserRequest true "User registration details"
// @Success      200 {object} object{data=model.UserResponse} "Successfully registered user"
// @Failure      400 {object} object{errors=string} "Invalid request body"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /users [post]
func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and receive access token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body model.LoginUserRequest true "User login credentials"
// @Success      200 {object} object{data=model.UserResponse} "Successfully logged in with token"
// @Failure      400 {object} object{errors=string} "Invalid request body"
// @Failure      401 {object} object{errors=string} "Invalid credentials"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /users/_login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Current godoc
// @Summary      Get current user
// @Description  Get the currently authenticated user's information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} object{data=model.UserResponse} "Current user information"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /users/_current [get]
func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Logout godoc
// @Summary      User logout
// @Description  Logout the currently authenticated user and invalidate token
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} object{data=bool} "Successfully logged out"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /users [delete]
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	response, err := c.UseCase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout user")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}

// Update godoc
// @Summary      Update current user
// @Description  Update the currently authenticated user's information (name and/or password)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.UpdateUserRequest true "User update details"
// @Success      200 {object} object{data=model.UserResponse} "Successfully updated user"
// @Failure      400 {object} object{errors=string} "Invalid request body"
// @Failure      401 {object} object{errors=string} "Unauthorized"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /users/_current [patch]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = auth.ID
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to update user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a valid refresh token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body model.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} object{data=model.UserResponse} "New access token"
// @Failure      400 {object} object{errors=string} "Invalid request body"
// @Failure      401 {object} object{errors=string} "Invalid refresh token"
// @Failure      500 {object} object{errors=string} "Internal server error"
// @Router       /users/refresh-token [post]
func (c *UserController) RefreshToken(ctx *fiber.Ctx) error {
	request := new(model.RefreshTokenRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.RefreshToken(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to refresh token : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
