// Package migrate to migrate the schema
package migrate

import (
	"fmt"

	gconfig "github.com/pilinux/gorest/config"
	gdatabase "github.com/pilinux/gorest/database"
	gmodel "github.com/pilinux/gorest/database/model"

	"apidev/database/model"
)

// Load all the models
type auth gmodel.Auth
type twoFA gmodel.TwoFA
type user model.User
type note model.Note

// DropAllTables - careful! It will drop all the tables!
func DropAllTables() error {
	db := gdatabase.GetDB()

	if err := db.Migrator().DropTable(
		&note{},
		&user{},
		&twoFA{},
		&auth{},
	); err != nil {
		return err
	}

	fmt.Println("old tables are deleted!")
	return nil
}

// StartMigration - automatically migrate all the tables
// - Only create tables with missing columns and missing indexes
// - Will not change/delete any existing columns and their types
func StartMigration(configure gconfig.Configuration) error {
	db := gdatabase.GetDB()
	configureDB := configure.Database.RDBMS
	driver := configureDB.Env.Driver

	if driver == "mysql" {
		// db.Set() --> add table suffix during auto migration
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
			&auth{},
			&twoFA{},
			&user{},
			&note{},
		); err != nil {
			return err
		}

		fmt.Println("new tables are  migrated successfully!")
		return nil
	}

	if err := db.AutoMigrate(
		&auth{},
		&twoFA{},
		&user{},
		&note{},
	); err != nil {
		return err
	}

	fmt.Println("new tables are  migrated successfully!")
	return nil
}

// SetPkFk - manually set foreign key for MySQL and PostgreSQL
func SetPkFk() error {
	db := gdatabase.GetDB()

	if !db.Migrator().HasConstraint(&user{}, "Notes") {
		err := db.Migrator().CreateConstraint(&user{}, "Notes")
		if err != nil {
			return err
		}
	}

	return nil
}
