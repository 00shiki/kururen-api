package cars

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"kururen/entity"
	MOCK_CARS "kururen/repository/cars/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDetail(t *testing.T) {
	type want struct {
		id      int
		errConv error
		car     *entity.Car
		errCar  error
	}
	tests := []struct {
		name  string
		want  want
		token string
	}{
		{
			name: "200#Success",
			want: want{
				id:  1,
				car: &entity.Car{},
			},
			token: "token",
		},
		{
			name: "400#BadRequest",
			want: want{
				id:      1,
				errConv: errors.New("bad request"),
			},
			token: "token",
		},
		{
			name: "401#Unauthorized",
			want: want{
				id:  1,
				car: &entity.Car{},
			},
			token: "",
		},
		{
			name: "404#NotFound",
			want: want{
				id:     1,
				errCar: gorm.ErrRecordNotFound,
			},
			token: "token",
		},
		{
			name: "500#InternalServerError",
			want: want{
				id:     1,
				errCar: errors.New("internal server error"),
			},
			token: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/cars/%d", tt.want.id),
				nil,
			)
			req.Header.Add("Authorization", "Bearer "+tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			carsRepo := new(MOCK_CARS.Repository)

			atoi := Atoi
			defer func() { Atoi = atoi }()
			Atoi = func(string) (int, error) {
				return tt.want.id, tt.want.errConv
			}

			carsRepo.On("GetCarByID", mock.AnythingOfType("uint")).Return(tt.want.car, tt.want.errCar)
			car, errCar := carsRepo.GetCarByID(uint(tt.want.id))
			if !reflect.DeepEqual(car, tt.want.car) {
				t.Errorf("GetCarByID() = %v, want %v", car, tt.want.car)
			}
			if !reflect.DeepEqual(errCar, tt.want.errCar) {
				t.Errorf("GetCarByID() = %v, want %v", errCar, tt.want.errCar)
			}

			handler := Controller{
				cr: carsRepo,
			}

			handler.Detail(c)
		})
	}
}
