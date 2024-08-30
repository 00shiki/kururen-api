package cars

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"kururen/entity"
	MOCK_CARS "kururen/repository/cars/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	type want struct {
		cars    []entity.Car
		errList error
	}
	tests := []struct {
		name  string
		want  want
		token string
	}{
		{
			name: "200#Success",
			want: want{
				cars: []entity.Car{
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
			},
			token: "token",
		},
		{
			name: "404#NotFound",
			want: want{
				errList: gorm.ErrRecordNotFound,
			},
			token: "token",
		},
		{
			name: "500#InternalServerError",
			want: want{
				errList: errors.New("internal server error"),
			},
			token: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/cars", nil)
			req.Header.Add("Authorization", "Bearer "+tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			carsRepo := new(MOCK_CARS.Repository)

			carsRepo.On("GetCars").Return(tt.want.cars, tt.want.errList)
			cars, errList := carsRepo.GetCars()
			if !reflect.DeepEqual(cars, tt.want.cars) {
				t.Errorf("GetCars() = %v, want %v", cars, tt.want.cars)
			}
			if !reflect.DeepEqual(errList, tt.want.errList) {
				t.Errorf("GetCars() = %v, want %v", errList, tt.want.errList)
			}

			handler := Controller{
				cr: carsRepo,
			}

			handler.List(c)
		})
	}
}
