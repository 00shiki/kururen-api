package cars

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	CARS_PRESENTATION "kururen/api/presentation/cars"
	MOCK_CARS "kururen/repository/cars/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAvailability(t *testing.T) {
	type want struct {
		errBind     error
		errValidate error
		id          int
		errConv     error
		errUpdate   error
	}
	tests := []struct {
		name        string
		want        want
		requestBody CARS_PRESENTATION.AvailabilityRequest
		token       string
	}{
		{
			name: "200#Success",
			want: want{
				id: 1,
			},
			requestBody: CARS_PRESENTATION.AvailabilityRequest{
				Status: "available",
			},
			token: "token",
		},
		{
			name: "400#BadRequest",
			want: want{
				id:      1,
				errBind: errors.New("bad request"),
			},
			requestBody: CARS_PRESENTATION.AvailabilityRequest{},
			token:       "token",
		},
		{
			name: "400#BadRequest",
			want: want{
				id:          1,
				errValidate: errors.New("bad request"),
			},
			requestBody: CARS_PRESENTATION.AvailabilityRequest{},
			token:       "token",
		},
		{
			name: "400#BadRequest",
			want: want{
				id:      0,
				errConv: errors.New("bad request"),
			},
			requestBody: CARS_PRESENTATION.AvailabilityRequest{
				Status: "available",
			},
			token: "token",
		},
		{
			name: "404#NotFound",
			want: want{
				id:        1,
				errUpdate: gorm.ErrRecordNotFound,
			},
			requestBody: CARS_PRESENTATION.AvailabilityRequest{
				Status: "available",
			},
			token: "token",
		},
		{
			name: "500#InternalServerError",
			want: want{
				id:        1,
				errUpdate: errors.New("internal server error"),
			},
			requestBody: CARS_PRESENTATION.AvailabilityRequest{
				Status: "available",
			},
			token: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			mockJSON, errJSON := json.Marshal(tt.requestBody)
			if errJSON != nil {
				log.Print(errJSON.Error())
			}
			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/cars/%d/availability", tt.want.id),
				bytes.NewReader(mockJSON),
			)
			req.Header.Add("Authorization", "Bearer "+tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			carsRepo := new(MOCK_CARS.Repository)

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

			atoi := Atoi
			defer func() { Atoi = atoi }()
			Atoi = func(string) (int, error) {
				return tt.want.id, tt.want.errConv
			}

			carsRepo.On("UpdateCarStatus", mock.AnythingOfType("uint"), mock.Anything).Return(tt.want.errUpdate)
			errUpdate := carsRepo.UpdateCarStatus(uint(tt.want.id), tt.requestBody.Status)
			if !reflect.DeepEqual(errUpdate, tt.want.errUpdate) {
				t.Errorf("UpdateCarStatus() = %v, want %v", errUpdate, tt.want.errUpdate)
			}

			handler := Controller{
				cr: carsRepo,
			}

			handler.Availability(c)
		})
	}
}
