// Package config provides a singleton instance of the configuration.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// EnvType is the type of the environment.
type EnvType string

const (
	test EnvType = "test"
	prod EnvType = "prod"
	dev  EnvType = "dev"
)

type (
	// Config is the configuration for the application.
	Config struct {
		Environment    EnvType `env:"ENVIRONMENT" default:"dev"` // required:"true"`
		HTTP           HTTP
		Postgres       Postgres
		Logger         Logger
		SigexEndpoints SigexEndpoints
		// Token       Token
		CORS  CORS
		Niger string `env:"NIEGR"`
	}

	// HTTP is the configuration for the HTTP server.
	HTTP struct {
		Host           string        `envconfig:"HTTP_HOST" default:"0.0.0.0"` //               required:"true"`
		Port           string        `envconfig:"HTTP_PORT" default:"8080"`    //               required:"true"`
		MaxHeaderBytes int           `envconfig:"HTTP_MAX_HEADER_BYTES"                 default:"1"`
		ReadTimeout    time.Duration `envconfig:"HTTP_READ_TIMEOUT"                     default:"10s"`
		WriteTimeout   time.Duration `envconfig:"HTTP_WRITE_TIMEOUT"                    default:"10s"`
	}

	// Postgres is the configuration for the Postgres database.
	Postgres struct {
		Host     string `env:"POSTGRES_HOST" default:"db"`                 // required:"true"`
		Port     string `env:"POSTGRES_PORT" default:"5432"`               //    required:"true"`
		DBName   string `env:"POSTGRES_DBNAME" default:"petition_service"` //     required:"true"`
		User     string `env:"POSTGRES_USER" default:"postgres"`           //  required:"true"`
		Password string `env:"POSTGRES_PASSWORD" default:"LiftKZ2023"`     //   required:"true" json:"-"`
		SSLMode  string `env:"POSTGRES_SSLMODE"                               default:"disable"`
	}

	// Logger is the configuration for the logger.
	Logger struct {
		Level string `env:"LOGGER_LEVEL" default:"info"`
	}

	SigexEndpoints struct {
		BaseUrl string `env:"BASE_URL"  default:"https://sigex.kz"`
	}

	// Token is the configuration for the token.
	// Token struct {
	// 	SecretKey string        `env:"TOKEN_SECRET_KEY" required:"true" json:"-"`
	// 	Expired   time.Duration `env:"TOKEN_EXPIRED"                             default:"15m"`
	// }

	// CORS is the configuration for the CORS.
	CORS struct {
		AllowOrigins []string `env:"CORS_ALLOW_ORIGINS" default:"http://localhost:3000"`
		// required:"true"`
	}
)

// IsDev check that environment is dev.
func (c *Config) IsDev() bool {
	return c.Environment == dev
}

// IsTest check that environment is test.
func (c *Config) IsTest() bool {
	return c.Environment == test
}

// IsProd check that environment is prod.
func (c *Config) IsProd() bool {
	return c.Environment == prod
}

var (
	instance Config
	once     sync.Once
)

// Get returns the singleton instance of the configuration.
func Get() *Config {
	once.Do(func() {

		// Store the current working directory.
		// This is important because we'll be changing the working directory temporarily to locate and load the .env file,
		// and we want to ensure we can return to the original directory later.
		// originalDir, err := os.Getwd()
		// if err != nil {

		// 	log.Fatalf("Error getting current working directory: %v", err)
		// }

		// Change the working directory to where the .env file is located.
		// Adjust the path as needed.
		// err = os.Chdir("../../")
		// if err != nil {
		// 	log.Fatalf("Error changing working directory: %v", err)
		// }

		// Restore the original working directory when the test finishes.
		// defer func() {
		// 	err := os.Chdir(originalDir)
		// 	if err != nil {
		// 		log.Fatalf("Error restoring working directory: %v", err)
		// 	}
		// }()

		// Load the .env file and read its contents.
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		myEnv, err := godotenv.Read()
		if err != nil {
			log.Fatalf("Error read .env file: %v", err)
		}
		fmt.Println("myEnv: ", myEnv)

		instance.Environment = EnvType(myEnv["ENVIRONMENT"])
		instance.HTTP.Host = myEnv["HTTP_HOST"]
		instance.HTTP.Port = myEnv["HTTP_PORT"]
		maxHeaderBytes, err := strconv.Atoi(myEnv["HTTP_MAX_HEADER_BYTES"])
		if err != nil {
			log.Fatalf("Error converting HTTP_MAX_HEADER_BYTES to int: %v", err)
		}
		instance.HTTP.MaxHeaderBytes = maxHeaderBytes

		readTimeoutStr := myEnv["HTTP_READ_TIMEOUT"]
		readTimeout, err := time.ParseDuration(readTimeoutStr)
		if err != nil {
			log.Fatalf("Error converting HTTP_READ_TIMEOUT to time.Duration: %v", err)
		}
		instance.HTTP.ReadTimeout = readTimeout

		writeTimeoutStr := myEnv["HTTP_WRITE_TIMEOUT"]
		writeTimeout, err := time.ParseDuration(writeTimeoutStr)
		if err != nil {
			log.Fatalf("Error converting HTTP_WRITE_TIMEOUT to time.Duration: %v", err)
		}
		instance.HTTP.WriteTimeout = writeTimeout

		instance.Postgres.Host = myEnv["POSTGRES_HOST"]
		instance.Postgres.Port = myEnv["POSTGRES_PORT"]
		instance.Postgres.DBName = myEnv["POSTGRES_DBNAME"]
		instance.Postgres.User = myEnv["POSTGRES_USER"]
		instance.Postgres.Password = myEnv["POSTGRES_PASSWORD"]
		instance.Postgres.SSLMode = myEnv["POSTGRES_SSLMODE"]
		instance.Logger.Level = myEnv["LOGGER_LEVEL"]
		instance.SigexEndpoints.BaseUrl = myEnv["BASE_URL"]
		instance.CORS.AllowOrigins = strings.Split(myEnv["CORS_ALLOW_ORIGINS"], ",")
		instance.Niger = myEnv["NIEGR"]

		fmt.Println("instance: ", instance)

		fmt.Println("instance of the application: ", instance)

		switch instance.Environment {
		case test, prod, dev:
		default:
			log.Fatal("config environment should be test, prod or dev")
		}
		if instance.IsDev() {
			configBytes, err := json.MarshalIndent(instance, "", " ")
			if err != nil {
				log.Fatal(fmt.Errorf("error marshaling indent config: %w", err))
			}
			fmt.Println("Configuration:", string(configBytes))
		}
	})

	return &instance
}
