package models

import (
	"context"
	"errors"
	"firebase-sso/db"
	"fmt"
	"github.com/logrusorgru/aurora"
)

func CreateSchema(schema string) error {
	pg, err := db.PostgresDB.Db.Acquire(context.Background())
	if err != nil {
		return errors.New("could not acquire PostGre connection")
	}

	defer func() {
		pg.Release()
	}()

	if !db.CheckSchemaExists(schema) {
		createInstallSchemaStmnt := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %v", schema)

		_, err = pg.Exec(context.Background(), createInstallSchemaStmnt)
		if err != nil {
			return err
		}

		fmt.Println(aurora.BgGreen("Created schema: " + schema))
		return nil
	}

	return errors.New(schema + " schema already created")
}

func CreateUsersTable() error {
	pg, err := db.PostgresDB.Db.Acquire(context.Background())
	if err != nil {
		return errors.New("could not acquire PostGre connection")
	}

	defer func() {
		pg.Release()
	}()

	if db.CheckSchemaExists("system") {

		createTableStmnt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v.users(
				id bigserial not null,
				uuid varchar(255),
				email varchar(255),
				title varchar(255),
				first_name varchar(255),
				last_name varchar(255),
				role int default 2,
				created_at bigint,
				updated_at bigint,
				deleted_at bigint,
				PRIMARY KEY (id)
			);
			create unique index users_email_uindex on system.users (email);`, "system")

		_, err := pg.Exec(context.Background(), createTableStmnt)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("system schema does not exist")
}

func CreateUserProjectsTable() error {
	pg, err := db.PostgresDB.Db.Acquire(context.Background())
	if err != nil {
		return errors.New("could not acquire PostGre connection")
	}

	defer func() {
		pg.Release()
	}()

	if db.CheckSchemaExists("system") {

		createTableStmnt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v.user_projects(
				id bigserial not null,
				title varchar(255),
				slug varchar(255),
				created bigint,
				created_by int not null,
				PRIMARY KEY (id)
			);`, "system")

		_, err := pg.Exec(context.Background(), createTableStmnt)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("system schema does not exist")
}
