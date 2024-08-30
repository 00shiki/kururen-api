package rentals

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"kururen/entity"
	MOCK_MAIL "kururen/pkg/mail/mocks"
	MOCK_XENDIT "kururen/pkg/xendit/mocks"
	MOCK_CARS "kururen/repository/cars/mocks"
	MOCK_RENTALS "kururen/repository/rentals/mocks"
	MOCK_USERS "kururen/repository/users/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	type want struct {
		userID          uint
		ok              bool
		rentalHistories []entity.RentalHistory
		errRentals      error
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
				rentalHistories: []entity.RentalHistory{
					{
						ID:        1,
						StartDate: time.Now(),
						EndDate:   time.Now().Add(time.Hour * 24),
						Cars: []*entity.Car{
							{
								ID:           1,
								Model:        "Model",
								Brand:        "Brand",
								Color:        "Color",
								Year:         "Year",
								RentalCost:   100000,
								Availability: "available",
							},
						},
						Payment: entity.Payment{
							ID:         1,
							Type:       "Mandiri",
							Amount:     100000,
							InvoiceURL: "invoice url",
						},
					},
				},
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
			name: "500#InternalServerError",
			want: want{
				userID:     1,
				ok:         true,
				errRentals: errors.New("internal server error"),
			},
			token: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/rentals", nil)
			req.Header.Add("Authorization", "Bearer "+tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			rentalsRepo := new(MOCK_RENTALS.Repository)
			usersRepo := new(MOCK_USERS.Repository)
			carsRepo := new(MOCK_CARS.Repository)

			xenditService := new(MOCK_XENDIT.Service)
			mailService := new(MOCK_MAIL.Service)

			get := Get
			defer func() { Get = get }()
			Get = func(echo.Context, string) (uint, bool) {
				return tt.want.userID, tt.want.ok
			}

			rentalsRepo.On("GetUserRentalHistories", mock.AnythingOfType("uint")).Return(tt.want.rentalHistories, tt.want.errRentals)
			rentalHistories, errRentals := rentalsRepo.GetUserRentalHistories(tt.want.userID)
			if !reflect.DeepEqual(rentalHistories, tt.want.rentalHistories) {
				t.Errorf("GetUserRentalHistory() = %v, want %v", rentalHistories, tt.want.rentalHistories)
			}
			if !reflect.DeepEqual(errRentals, tt.want.errRentals) {
				t.Errorf("GetUserRentalHistory() = %v, want %v", errRentals, tt.want.errRentals)
			}

			handler := Controller{
				rr: rentalsRepo,
				ur: usersRepo,
				cr: carsRepo,
				xs: xenditService,
				ms: mailService,
			}

			handler.List(c)
		})
	}
}
