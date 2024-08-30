package users

import (
	"github.com/labstack/echo/v4"
	USERS_PRESENTATION "kururen/api/presentation/users"
	"kururen/pkg/response"
	"net/http"
)

// TopUp godoc
// @Summary TopUp User
// @Description Add user deposit amount
// @Tags Users
// @Accept json
// @Produce json
// @Param RequestBody body users.TopUpRequest true "Top Up Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /users/topup [post]
func (handler *Controller) TopUp(c echo.Context) error {
	userID, ok := Get(c, "user_id")
	if !ok {
		result := response.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid token",
		}
		return response.HandleResponse(c, &result)
	}

	payload := new(USERS_PRESENTATION.TopUpRequest)
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

	user, errUser := handler.ur.GetUserByID(userID)
	if errUser != nil {
		result := response.Response{
			Code:    http.StatusUnauthorized,
			Message: "User not found",
		}
		return response.HandleResponse(c, &result)
	}

	user.DepositAmount += payload.Amount
	errUpdate := handler.ur.UpdateUser(user)
	if errUpdate != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errUpdate.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	result := response.Response{
		Code:    http.StatusOK,
		Message: "Top Up Success",
	}
	return response.HandleResponse(c, &result)
}
