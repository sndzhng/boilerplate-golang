package datastore

import (
	"fmt"
	"log"
	"strings"

	"github.com/sndzhng/gin-template/internal/config"
	"github.com/sndzhng/gin-template/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Postgresql *gorm.DB
)

func ConnectPostgresql() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s sslmode=disable",
		config.Datastore.Postgresql.Host,
		config.Datastore.Postgresql.User,
		config.Datastore.Postgresql.Password,
		config.Datastore.Postgresql.Database,
		config.Datastore.Postgresql.Port,
		config.Datastore.Postgresql.TimeZone,
	)

	err := error(nil)
	Postgresql, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}

	// description: validate connection timezone
	timezone := ""
	err = Postgresql.Raw("SHOW TIME ZONE").Scan(&timezone).Error
	if err != nil {
		log.Fatal(err)
	} else if timezone != config.Datastore.Postgresql.TimeZone {
		log.Fatal("postgresql timezone is not Asia/Bangkok")
	}

	// description: migrate enum with values map type name
	migrateEnums(
		map[string][]interface{}{
			// "name": {"value", "value"},
		},
	)

	// description: migrate table
	err = Postgresql.AutoMigrate(
		&entity.Admin{},
		&entity.Role{},
		&entity.User{},
	)
	if err != nil {
		log.Fatal(err)
	}

	// description: migrate constraints with table, constraint name and condition for migrate
	migrateConstraints(
		[][3]string{
			// {"users", "password_is_not_null", "password IS NOT NULL"},
		},
	)
}

func migrateConstraints(constraints [][3]string) {
	query := "DO $$	BEGIN "
	for _, constraint := range constraints {
		query += fmt.Sprintf("IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = '%s' AND conrelid = '%s'::regclass) THEN ALTER TABLE %s ADD CONSTRAINT %s CHECK (%s); END IF; ", constraint[1], constraint[0], constraint[0], constraint[1], constraint[2])
	}
	query += "END $$;"

	err := Postgresql.Exec(query).Error
	if err != nil {
		log.Fatal(err)
	}
}

func migrateEnums(enumMapDataType map[string][]interface{}) {
	existedDataTypes := []string{}
	err := Postgresql.Raw("SELECT enumtypid::regtype AS enum_name FROM pg_enum	GROUP BY enum_name;").Scan(&existedDataTypes).Error
	if err != nil {
		log.Fatal(err)
	}

	for dataType := range enumMapDataType {
		isExist := false
		for _, existedDataType := range existedDataTypes {
			if existedDataType == dataType {
				isExist = true
				break
			}
		}

		if !isExist {
			stringEnums := []string{}
			for _, enum := range enumMapDataType[dataType] {
				stringEnums = append(stringEnums, fmt.Sprintf("'%v'", enum))
			}

			err = Postgresql.Exec(fmt.Sprintf("CREATE TYPE %s AS ENUM (%v)", dataType, strings.Join(stringEnums, ", "))).Error
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
