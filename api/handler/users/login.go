package users

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	USERS_PRESENTATION "kururen/api/presentation/users"
	PKG_JWT "kururen/pkg/jwt"
	"kururen/pkg/response"
	"net/http"
	"os"
	"time"
)

// Login godoc
// @Summary Login User
// @Description Logged in to user account
// @Tags Users
// @Accept json
// @Produce json
// @Param RequestBody body users.LoginRequest true "Login Request"
// @Success 201 {object} response.Response{data=users.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/login [post]
func (handler *Controller) Login(c echo.Context) error {
	payload := new(USERS_PRESENTATION.LoginRequest)
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

	user, errUser := handler.ur.GetUserByUsername(payload.Username)
	if errUser != nil {
		result := response.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid username or password",
		}
		return response.HandleResponse(c, &result)
	}
	errCompare := ComparePassword([]byte(user.Password), []byte(payload.Password))
	if errCompare != nil {
		result := response.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid username or password",
		}
		return response.HandleResponse(c, &result)
	}

	claims := PKG_JWT.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, errToken := SignedString(token, []byte(os.Getenv("JWT_SECRET")))
	if errToken != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: "Error generating token",
		}
		return response.HandleResponse(c, &result)
	}

	user.JwtToken = tokenStr
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
		Message: "Success",
		Data: USERS_PRESENTATION.LoginResponse{
			Token: tokenStr,
		},
	}
	return response.HandleResponse(c, &result)
}
