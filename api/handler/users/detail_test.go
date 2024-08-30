package users

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"kururen/entity"
	MOCK_MAIL "kururen/pkg/mail/mocks"
	MOCK_USERS "kururen/repository/users/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDetail(t *testing.T) {
	type want struct {
		userID  uint
		ok      bool
		user    *entity.User
		errUser error
	}
	tests := []struct {
		name  string
		want  want
		token string
	}{
		{
			name: "200#Success",
			want: want{
				userID: 1,
				ok:     true,
				user:   &entity.User{},
			},
			token: "token",
		},
		{
			name: "401#Unauthorized",
			want: want{
				userID: 0,
				ok:     false,
			},
			token: "token",
		},
		{
			name: "401#Unauthorized",
			want: want{
				userID:  1,
				ok:      true,
				errUser: errors.New("unauthorized"),
			},
			token: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
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

			usersRepo.On("GetUserByID", mock.Anything).Return(tt.want.user, tt.want.errUser)
			user, errUser := usersRepo.GetUserByID(tt.want.userID)
			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("GetUserByID() = %v, want %v", user, tt.want.user)
			}
			if !reflect.DeepEqual(errUser, tt.want.errUser) {
				t.Errorf("GetUserByID() = %v, want %v", errUser, tt.want.errUser)
			}

			handler := Controller{
				ur: usersRepo,
				ms: mailService,
			}

			handler.Detail(c)
		})
	}
}
