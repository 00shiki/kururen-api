package rentals

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/mock"
	RENTALS_PRESENTATION "kururen/api/presentation/rentals"
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

func TestCreate(t *testing.T) {
	type want struct {
		userID        uint
		ok            bool
		user          *entity.User
		errUser       error
		errBind       error
		errValidate   error
		startDate     time.Time
		errParseStart error
		endDate       time.Time
		errParseEnd   error
		car           *entity.Car
		errCar        error
		errInvoice    error
		errUpdateUser error
		errSend       error
		errCreate     error
		errUpdateCar  error
	}
	tests := []struct {
		name          string
		want          want
		requestBody   RENTALS_PRESENTATION.CreateRequest
		rentalHistory *entity.RentalHistory
		token         string
	}{
		{
			name: "200#Success",
			want: want{
				userID: 1,
				ok:     true,
				user: &entity.User{
					Email:         "test@test.com",
					DepositAmount: 200000,
				},
				startDate: time.Date(2024, time.August, 30, 0, 0, 0, 0, time.Local),
				endDate:   time.Date(2024, time.August, 31, 0, 0, 0, 0, time.Local),
				car: &entity.Car{
					ID:           1,
					Availability: "available",
					RentalCost:   100000,
				},
			},
			requestBody: RENTALS_PRESENTATION.CreateRequest{
				Cars: []RENTALS_PRESENTATION.CarRequest{
					{
						CarID: 1,
					},
				},
				StartDate:   "2024-08-30",
				EndDate:     "2024-08-31",
				PaymentType: "MANDIRI",
			},
			rentalHistory: &entity.RentalHistory{},
			token:         "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			mockJSON, errJSON := json.Marshal(tt.requestBody)
			if errJSON != nil {
				log.Print(errJSON.Error())
			}
			req := httptest.NewRequest(http.MethodPost, "/rentals", bytes.NewReader(mockJSON))
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

			usersRepo.On("GetUserByID", mock.AnythingOfType("uint")).Return(tt.want.user, tt.want.errUser)
			user, errUser := usersRepo.GetUserByID(tt.want.userID)
			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("GetUserByID() = %v, want %v", user, tt.want.user)
			}
			if !reflect.DeepEqual(errUser, tt.want.errUser) {
				t.Errorf("GetUserByID() = %v, want %v", errUser, tt.want.errUser)
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

			parseTime := ParseTime
			defer func() { ParseTime = parseTime }()
			ParseTime = func(string, string) (time.Time, error) {
				return tt.want.startDate, tt.want.errParseStart
			}

			carsRepo.On("GetCarByID", mock.AnythingOfType("uint")).Return(tt.want.car, tt.want.errCar)
			car, errCar := carsRepo.GetCarByID(tt.want.car.ID)
			if !reflect.DeepEqual(car, tt.want.car) {
				t.Errorf("GetCarByID() = %v, want %v", car, tt.want.car)
			}
			if !reflect.DeepEqual(errCar, tt.want.errCar) {
				t.Errorf("GetCarByID() = %v, want %v", errCar, tt.want.errCar)
			}

			xenditService.On("CreateInvoice", mock.AnythingOfType("*entity.RentalHistory")).Return(tt.want.errInvoice)
			errInvoice := xenditService.CreateInvoice(tt.rentalHistory)
			if !reflect.DeepEqual(errInvoice, tt.want.errInvoice) {
				t.Errorf("CreateInvoice() = %v, want %v", errInvoice, tt.want.errInvoice)
			}

			usersRepo.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(tt.want.errUpdateUser)
			errUpdateUser := usersRepo.UpdateUser(tt.want.user)
			if !reflect.DeepEqual(errUpdateUser, tt.want.errUpdateUser) {
				t.Errorf("UpdateUser() = %v, want %v", errUpdateUser, tt.want.errUpdateUser)
			}

			mailService.On("SendMail", mock.Anything, mock.Anything, mock.Anything).Return(tt.want.errSend)
			errSend := mailService.SendMail(tt.want.user.Email, "", "")
			if !reflect.DeepEqual(errSend, tt.want.errSend) {
				t.Errorf("SendMail() = %v, want %v", errSend, tt.want.errSend)
			}

			rentalsRepo.On("CreateRentalHistory", mock.AnythingOfType("*entity.RentalHistory")).Return(tt.want.errCreate)
			errCreate := rentalsRepo.CreateRentalHistory(tt.rentalHistory)
			if !reflect.DeepEqual(errCreate, tt.want.errCreate) {
				t.Errorf("CreateRentalHistory() = %v, want %v", errCreate, tt.want.errCreate)
			}

			carsRepo.On("UpdateCarStatus", mock.AnythingOfType("uint"), mock.Anything).Return(tt.want.errUpdateCar)
			errUpdateCar := carsRepo.UpdateCarStatus(tt.requestBody.Cars[0].CarID, "booked")
			if !reflect.DeepEqual(errUpdateCar, tt.want.errUpdateCar) {
				t.Errorf("UpdateCarStatus() = %v, want %v", errUpdateCar, tt.want.errUpdateCar)
			}

			handler := Controller{
				rr: rentalsRepo,
				ur: usersRepo,
				cr: carsRepo,
				xs: xenditService,
				ms: mailService,
			}

			handler.Create(c)
		})
	}
}
