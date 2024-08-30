package users

import (
	"github.com/labstack/echo/v4"
	USERS_PRESENTATION "kururen/api/presentation/users"
	"kururen/entity"
	"kururen/pkg/response"
	"net/http"
)

// Register godoc
// @Summary Register User
// @Description Registering new user
// @Tags Users
// @Accept json
// @Produce json
// @Param RequestBody body users.RegisterRequest true "Register Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/register [post]
func (handler *Controller) Register(c echo.Context) error {
	payload := new(USERS_PRESENTATION.RegisterRequest)
	if errBind := Bind(c, payload); errBind != nil {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: errBind.Error(),
		}
		return response.HandleResponse(c, &result)
	}
	if errValidate := Validate(c, payload); errValidate != nil {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: errValidate.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	password, errPassword := GeneratePassword([]byte(payload.Password), 14)
	if errPassword != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errPassword.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	user := entity.User{
		Username: payload.Username,
		Password: string(password),
		Name:     payload.Name,
		Email:    payload.Email,
	}
	errCreate := handler.ur.CreateUser(&user)
	if errCreate != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errCreate.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	errSend := handler.ms.SendMail(
		user.Email,
		user.Name,
		"Registration Successful",
		"Thank you for registering!",
	)
	if errSend != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errSend.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	result := response.Response{
		Code:    http.StatusCreated,
		Message: "User created",
	}
	return response.HandleResponse(c, &result)
}
