package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"kururen/api"
	"kururen/config"
	"kururen/docs"
	"os"
)

//	@title						Kururen API
//	@version					1.0
//	@description				Kuruma (Car) Renting API
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@BasePath					/api/v1
// 	@securityDefinitions.apiKey BearerAuth
// 	@in 						header
// 	@name 						Authorization

func main() {
	_ = godotenv.Load()

	var url = os.Getenv("URL")

	docs.SwaggerInfo.Host = url

	db, err := gorm.Open(postgres.Open(config.DatabaseConfig()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	api.Init(e, db)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
