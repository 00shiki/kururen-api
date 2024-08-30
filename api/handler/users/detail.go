package users

import (
	"github.com/labstack/echo/v4"
	USERS_PRESENTATION "kururen/api/presentation/users"
	"kururen/pkg/response"
	"net/http"
)

// Detail godoc
// @Summary Detail User
// @Description Get User Detail
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=users.DetailResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /users/me [get]
func (handler *Controller) Detail(c echo.Context) error {
	userID, ok := Get(c, "user_id")
	if !ok {
		result := response.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid token",
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

	result := response.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data: USERS_PRESENTATION.DetailResponse{
			Username:      user.Username,
			Name:          user.Name,
			Email:         user.Email,
			DepositAmount: user.DepositAmount,
		},
	}
	return response.HandleResponse(c, &result)
}
