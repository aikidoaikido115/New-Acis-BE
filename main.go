package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	"github.com/aikidoaikido115/New-Acis-BE/docs"
	"github.com/aikidoaikido115/New-Acis-BE/modules/servers"
	"github.com/aikidoaikido115/New-Acis-BE/pkg/database"
)

// @title New-Acis API
// @description This is a sample for New-Acis API.
// @BasePath /
// @schemes https http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	
	cfg := configs.LoadConfigs()

	if host := os.Getenv("SWAGGER_HOST"); host != "" {
		docs.SwaggerInfo.Host = host
	}
	if schemes := os.Getenv("SWAGGER_SCHEMES"); schemes != "" {
		docs.SwaggerInfo.Schemes = strings.Split(schemes, ",")
	}

	database.InitDB(cfg.PostgreSQL)
	app := servers.SetupServer(cfg.Server, cfg.JWT, cfg.Supabase, cfg.Mail)

	serverAddress := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server is starting on %s", serverAddress)

	if err := app.Listen(serverAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
