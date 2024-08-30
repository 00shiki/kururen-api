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

func TestTopUp(t *testing.T) {
	type want struct {
		userID      uint
		ok          bool
		errBind     error
		errValidate error
		user        *entity.User
		errUser     error
		errUpdate   error
	}
	tests := []struct {
		name        string
		want        want
		requestBody USERS_PRESENTATION.TopUpRequest
		token       string
	}{
		{
			name: "200#Success",
			want: want{
				userID: 1,
				ok:     true,
				user:   &entity.User{},
			},
			requestBody: USERS_PRESENTATION.TopUpRequest{
				Amount: 500000,
			},
		},
		{
			name: "400#BadRequest",
			want: want{
				userID:  1,
				ok:      true,
				errBind: errors.New("bad request"),
			},
			requestBody: USERS_PRESENTATION.TopUpRequest{
				Amount: 500000,
			},
		},
		{
			name: "400#BadRequest",
			want: want{
				userID:      1,
				ok:          true,
				errValidate: errors.New("bad request"),
			},
			requestBody: USERS_PRESENTATION.TopUpRequest{
				Amount: 500000,
			},
		},
		{
			name: "401#Unauthorized",
			want: want{
				userID: 0,
				ok:     false,
			},
			requestBody: USERS_PRESENTATION.TopUpRequest{
				Amount: 500000,
			},
		},
		{
			name: "401#Unauthorized",
			want: want{
				userID:  1,
				ok:      true,
				errUser: errors.New("unauthorized"),
			},
			requestBody: USERS_PRESENTATION.TopUpRequest{
				Amount: 500000,
			},
		},
		{
			name: "500#InternalServerError",
			want: want{
				userID:    1,
				ok:        true,
				user:      &entity.User{},
				errUpdate: errors.New("internal server error"),
			},
			requestBody: USERS_PRESENTATION.TopUpRequest{
				Amount: 500000,
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
			req := httptest.NewRequest(http.MethodPost, "/users/topup", bytes.NewReader(mockJSON))
			req.Header.Add("Authorization", "Bearer "+tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			usersRepo := new(MOCK_USERS.Repository)
			mailService := new(MOCK_MAIL.Service)

			get := Get
			defer func() { Get = get }()
			Get = func(echo.Context, string) (uint, bool) {
				return tt.want.userID, tt.want.ok
			}

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

			usersRepo.On("GetUserByID", mock.AnythingOfType("uint")).Return(tt.want.user, tt.want.errUser)
			user, errUser := usersRepo.GetUserByID(tt.want.userID)
			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("GetUserByID() = %v, want %v", user, tt.want.user)
			}
			if !reflect.DeepEqual(errUser, tt.want.errUser) {
				t.Errorf("GetUserByID() = %v, want %v", errUser, tt.want.errUser)
			}

			usersRepo.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(tt.want.errUpdate)
			errUpdate := usersRepo.UpdateUser(tt.want.user)
			if !reflect.DeepEqual(errUpdate, tt.want.errUpdate) {
				t.Errorf("UpdateUser() = %v, want %v", errUpdate, tt.want.errUpdate)
			}

			handler := Controller{
				ur: usersRepo,
				ms: mailService,
			}

			handler.TopUp(c)
		})
	}
}
