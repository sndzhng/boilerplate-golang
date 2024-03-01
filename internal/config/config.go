package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	Datastore   DatastoreConfig
	Environment string
	JWT         JWTConfig
	Server      ServerConfig
)

type (
	DatastoreConfig struct {
		CloudStorage CloudStorageConfig
		Mongodb      MongodbConfig
		Postgresql   PostgresqlConfig
	}
	CloudStorageConfig struct {
		BucketName, ProjectID string
	}
	MongodbConfig struct {
		Database, Format, Host, Options, Password, User string
	}
	PostgresqlConfig struct {
		Database, Host, Password, Port, TimeZone, User string
	}
	JWTConfig struct {
		ExpireMinute, Key string
	}
	ServerConfig struct {
		Context, Port string
	}
)

func InitialConfig(args []string) {
	if len(args) > 1 {
		err := godotenv.Load(fmt.Sprintf("env/%s", args[1]))
		if err != nil {
			log.Fatal("Error environment file not found")
		}
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	Datastore = DatastoreConfig{
		CloudStorage: CloudStorageConfig{
			BucketName: getEnv("DATASTORE_CLOUD_STORAGE_BUCKET_NAME"),
			ProjectID:  getEnv("DATASTORE_CLOUD_STORAGE_PROJECT_ID"),
		},
		Mongodb: MongodbConfig{
			Database: getEnv("DATASTORE_MONGODB_DATABASE"),
			Format:   getEnv("DATASTORE_MONGODB_FORMAT"),
			Host:     getEnv("DATASTORE_MONGODB_HOST"),
			Options:  getEnv("DATASTORE_MONGODB_OPTIONS"),
			Password: getEnv("DATASTORE_MONGODB_PASSWORD"),
			User:     getEnv("DATASTORE_MONGODB_USER"),
		},
		Postgresql: PostgresqlConfig{
			Database: getEnv("DATASTORE_POSTGRESQL_DATABASE"),
			Host:     getEnv("DATASTORE_POSTGRESQL_HOST"),
			Password: getEnv("DATASTORE_POSTGRESQL_PASSWORD"),
			Port:     getEnv("DATASTORE_POSTGRESQL_PORT"),
			TimeZone: getEnv("DATASTORE_POSTGRESQL_TIME_ZONE"),
			User:     getEnv("DATASTORE_POSTGRESQL_USER"),
		},
	}
	Environment = getEnv("ENVIRONMENT")
	JWT = JWTConfig{
		ExpireMinute: getEnv("JWT_EXPIRE_MINUTE"),
		Key:          getEnv("JWT_KEY"),
	}
	Server = ServerConfig{
		Context: getEnv("SERVER_CONTEXT"),
		Port:    getEnv("SERVER_PORT"),
	}
}

func InitialTimeZone() {
	localTimeZone, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal("Error initial timezone")
	}
	time.Local = localTimeZone
}

func getEnv(key string) string {
	value, isExist := os.LookupEnv(key)
	if !isExist {
		log.Fatalf("%s environment variable not found", key)
	}
	return value
}
