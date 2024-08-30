package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"kururen/api/presentation/users"
	"kururen/entity"
	MOCK_MAIL "kururen/pkg/mail/mocks"
	MOCK_USERS "kururen/repository/users/mocks"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestLogin(t *testing.T) {
	type want struct {
		errBind     error
		errValidate error
		user        *entity.User
		errUser     error
		errCompare  error
		tokenStr    string
		errToken    error
		errUpdate   error
	}
	tests := []struct {
		name        string
		want        want
		requestBody users.LoginRequest
	}{
		{
			name: "200#Success",
			want: want{
				user: &entity.User{},
			},
			requestBody: users.LoginRequest{
				Username: "test",
				Password: "test",
			},
		},
		{
			name: "400#BadRequest",
			want: want{
				errBind: errors.New("bad Request"),
				user:    &entity.User{},
			},
			requestBody: users.LoginRequest{
				Username: "test",
			},
		},
		{
			name: "400#BadRequest",
			want: want{
				errValidate: errors.New("bad Request"),
				user:        &entity.User{},
			},
			requestBody: users.LoginRequest{
				Username: "test",
				Password: "",
			},
		},
		{
			name: "401#Unauthorized",
			want: want{
				errUser: errors.New("unauthorized"),
			},
			requestBody: users.LoginRequest{
				Username: "test",
				Password: "test",
			},
		},
		{
			name: "401#Unauthorized",
			want: want{
				user:       &entity.User{},
				errCompare: errors.New("unauthorized"),
			},
			requestBody: users.LoginRequest{
				Username: "test",
				Password: "test",
			},
		},
		{
			name: "500#InternalServerError",
			want: want{
				user:      &entity.User{},
				errUpdate: errors.New("internal server error"),
			},
			requestBody: users.LoginRequest{
				Username: "test",
				Password: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			mockJSON, errJSON := json.Marshal(tt.requestBody)
			if errJSON != nil {
				log.Print(errJSON.Error())
			}
			req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(mockJSON))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			usersRepo := new(MOCK_USERS.Repository)
			mailService := new(MOCK_MAIL.Service)

			bind := Bind
			defer func() { Bind = bind }()
			Bind = func(echo.Context, interface{}) error {
				return tt.want.errBind
			}

			validate := Validate
			defer func() { Validate = validate }()
			Validate = func(echo.Context, interface{}) error {
				return tt.want.errValidate
			}

			usersRepo.On("GetUserByUsername", mock.Anything).Return(tt.want.user, tt.want.errUser)
			user, errUser := usersRepo.GetUserByUsername(tt.requestBody.Username)
			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("GetUserByUsername() = %v, want %v", user, tt.want.user)
			}
			if !reflect.DeepEqual(errUser, tt.want.errUser) {
				t.Errorf("GetUserByUsername() = %v, want %v", errUser, tt.want.errUser)
			}

			comparePassword := ComparePassword
			defer func() { ComparePassword = comparePassword }()
			ComparePassword = func([]byte, []byte) error {
				return tt.want.errCompare
			}

			usersRepo.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(tt.want.errUpdate)
			errUpdate := usersRepo.UpdateUser(user)
			if !reflect.DeepEqual(errUser, tt.want.errUser) {
				t.Errorf("UpdateUser() = %v, want %v", errUpdate, tt.want.errUpdate)
			}

			handler := Controller{
				ur: usersRepo,
				ms: mailService,
			}

			handler.Login(c)
		})
	}
}
