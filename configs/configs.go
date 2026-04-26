package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	PostgreSQL PostgreSQL
	Server     Server
	JWT        JWT
	Supabase   Supabase
	Mail       Mail
	SeedAdmin  SeedAdmin
}

type Server struct {
	Host    string
	Port    string
	AppName string
	CORS    CORS
}

type CORS struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

type JWT struct {
	Secret string
}

type Supabase struct {
	URL        string
	ServiceKey string
	Bucket     string
}

type Mail struct {
	Host   string
	Port   string
	Sender string
	Key    string
}

type SeedAdmin struct {
	Username  string
	Email     string
	Password  string
	RoleName  string
	FirstName string
	LastName  string
	Nickname  string
}

func LoadConfigs() *Configs {
	// โหลด .env ตาม environment
	env := getEnv("GO_ENV", "development")

	var envFile string
	switch env {
	case "production":
		envFile = ".env.prod"
	case "development":
		envFile = ".env.dev"
	default:
		envFile = ".env"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("No %s file found, trying .env...\n", envFile)
		err = godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using environment variables")
		}
	} else {
		log.Printf("Loaded configuration from %s\n", envFile)
	}

	return &Configs{
		PostgreSQL: PostgreSQL{
			Host:     getEnv("DB_HOST", "db"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_NAME", "test"),
			SSLMode:  getEnv("SSL_Mode", "disable"),
		},
		Server: Server{
			Host:    getEnv("SERVER_HOST", "0.0.0.0"),
			Port:    getEnv("PORT", getEnv("SERVER_PORT", "8080")), // Heroku ใช้ PORT env
			AppName: "New Acis",
			CORS: CORS{
				AllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000"),
				AllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
				AllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization,X-Requested-With"),
				AllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
			},
		},
		JWT: JWT{
			Secret: getEnv("JWT_SECRET", "secret-key-change-in-production"),
		},
		Supabase: Supabase{
			URL:        os.Getenv("SUPABASE_URL"),
			ServiceKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
			Bucket:     os.Getenv("SUPABASE_BUCKET"),
		},
		Mail: Mail{
			Host:   os.Getenv("EMAIL_HOST"),
			Port:   os.Getenv("EMAIL_PORT"),
			Sender: os.Getenv("SENDER_EMAIL"),
			Key:    os.Getenv("APP_PASSWORD"),
		},
		SeedAdmin: SeedAdmin{
			Username:  getEnv("ADMIN_USERNAME", "admin"),
			Email:     getEnv("ADMIN_EMAIL", "admin@example.com"),
			Password:  getEnv("ADMIN_PASSWORD", "password"),
			RoleName:  getEnv("ROLE_NAME", "Admin"),
			FirstName: getEnv("FIRST_NAME", "Admin"),
			LastName:  getEnv("LAST_NAME", "User"),
			Nickname:  getEnv("NICKNAME", "Admin"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}
