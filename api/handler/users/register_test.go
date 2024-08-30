package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/mock"
	USERS_PRESENTATION "kururen/api/presentation/users"
	"kururen/entity"
	MOCK_MAIL "kururen/pkg/mail/mocks"
	MOCK_USERS "kururen/repository/users/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	type want struct {
		errBind     error
		errValidate error
		password    []byte
		errPassword error
		errCreate   error
		errSend     error
	}
	tests := []struct {
		name        string
		want        want
		user        *entity.User
		requestBody USERS_PRESENTATION.RegisterRequest
	}{
		{
			name: "201#Success",
			want: want{
				password: []byte("hashedpassword"),
			},
			user: &entity.User{
				Username: "test",
				Password: "hashedpassword",
				Name:     "test",
				Email:    "test@test.com",
			},
			requestBody: USERS_PRESENTATION.RegisterRequest{
				Username: "test",
				Password: "password",
				Name:     "test",
				Email:    "test@test.com",
			},
		},
		{
			name: "400#BadRequest",
			want: want{
				errBind:  errors.New("bad request"),
				password: []byte("hashedpassword"),
			},
			user: &entity.User{},
			requestBody: USERS_PRESENTATION.RegisterRequest{
				Username: "test",
			},
		},
		{
			name: "400#BadRequest",
			want: want{
				errValidate: errors.New("bad request"),
				password:    []byte("hashedpassword"),
			},
			user: &entity.User{},
			requestBody: USERS_PRESENTATION.RegisterRequest{
				Username: "test",
			},
		},
		{
			name: "500#InternalServerError",
			want: want{
				password:    []byte("hashedpassword"),
				errPassword: errors.New("internal server error"),
			},
			user: &entity.User{
				Username: "test",
				Password: "hashedpassword",
				Name:     "test",
				Email:    "test@test.com",
			},
			requestBody: USERS_PRESENTATION.RegisterRequest{
				Username: "test",
				Password: "password",
			},
		},
		{
			name: "500#InternalServerError",
			want: want{
				password:  []byte("hashedpassword"),
				errCreate: errors.New("internal server error"),
			},
			user: &entity.User{
				Username: "test",
				Password: "hashedpassword",
				Name:     "test",
				Email:    "test@test.com",
			},
			requestBody: USERS_PRESENTATION.RegisterRequest{
				Username: "test",
				Password: "password",
			},
		},
		{
			name: "500#InternalServerError",
			want: want{
				password: []byte("hashedpassword"),
				errSend:  errors.New("internal server error"),
			},
			user: &entity.User{
				Username: "test",
				Password: "hashedpassword",
				Name:     "test",
				Email:    "test@test.com",
			},
			requestBody: USERS_PRESENTATION.RegisterRequest{
				Username: "test",
				Password: "password",
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
			req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader(mockJSON))
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

			generatePassword := GeneratePassword
			defer func() { GeneratePassword = generatePassword }()
			GeneratePassword = func([]byte, int) ([]byte, error) {
				return tt.want.password, tt.want.errPassword
			}

			usersRepo.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(tt.want.errCreate)
			errCreate := usersRepo.CreateUser(tt.user)
			if !reflect.DeepEqual(errCreate, tt.want.errCreate) {
				t.Errorf("CreateUser() = %v, want %v", errCreate, tt.want.errCreate)
			}

			mailService.On("SendMail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.want.errSend)
			errSend := mailService.SendMail(tt.user.Email, tt.user.Name, "", "")
			if !reflect.DeepEqual(errSend, tt.want.errSend) {
				t.Errorf("SendMail() = %v, want %v", errSend, tt.want.errSend)
			}

			handler := Controller{
				ur: usersRepo,
				ms: mailService,
			}

			handler.Register(c)
		})
	}
}
