package models

import (
	"context"
	"database/sql"
	"errors"
	"firebase-sso/caching"
	"firebase-sso/db"
	"fmt"
	querybuilder "github.com/KirksFletcher/Go-SQL-Query-Builder-Golang"
	"github.com/georgysavva/scany/pgxscan"
	"strconv"
	"strings"
	"time"
)

// This model contains some core functionality for interacting with Users

// User
// Primary user definition
type User struct {
	ID              int           `db:"id" json:"id"`
	Role            int           `db:"role" json:"role"`
	UUID            string        `db:"uuid" json:"uuid"`
	Title           string        `db:"title" jsosn:"title"`
	FirstName       string        `db:"first_name" json:"first_name"`
	LastName        string        `db:"last_name" json:"last_name"`
	Email           string        `db:"email" json:"email"`
	Picture         string        `db:"-" json:"picture"`
	CreatedAt       sql.NullInt64 `db:"created_at" json:"-"`
	UpdatedAt       sql.NullInt64 `db:"updated_at" json:"-"`
	DeletedAt       sql.NullInt64 `db:"deleted_at" json:"-"`
	PublicCreatedAt int64         `db:"-" json:"created_at"`
	PublicUpdatedAt int64         `db:"-" json:"updated_at"`
	PublicDeletedAt int64         `db:"-" json:"deleted_at"`
}

// UserEditData
// This exists as struct for binding availability on API controller
type UserEditData struct {
	ID        int    `db:"-" json:"user_id"`
	UUID      string `db:"-" json:"-"`
	Title     string `db:"title" json:"title"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Role      int    `db:"role" json:"role"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
}

// UserCreateData
// This exists as struct for binding availability on API controller
type UserCreateData struct {
	UUID      string `db:"uuid" json:"uuid"`
	Title     string `db:"title" json:"title"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Role      int    `db:"role" json:"role"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
}

func (u User) FormatDates() User {

	// Method for cleaning up struct sql.Null types for json return

	if u.CreatedAt.Valid {
		u.PublicCreatedAt = u.CreatedAt.Int64
	}
	if u.UpdatedAt.Valid {
		u.PublicUpdatedAt = u.UpdatedAt.Int64
	}
	if u.DeletedAt.Valid {
		u.PublicDeletedAt = u.DeletedAt.Int64
	}

	return u
}

func GetUserById(userId int) (User, error) {
	var userData User

	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return userData, err
	}

	var sqlb querybuilder.Sqlbuilder

	defer func() {
		pg.Release()
	}()

	query := sqlb.Reset().From(`system.users`).Where(`id`, `=`, strconv.Itoa(userId)).Build()

	err = pgxscan.Get(ctx, pg, &userData, query)
	if err != nil {
		return userData, err
	}

	return userData, nil
}

func GetUserByEmail(userEmail string) (User, error) {
	var userData User

	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return userData, err
	}

	var sqlb querybuilder.Sqlbuilder

	defer func() {
		pg.Release()
	}()

	query := sqlb.Reset().From(`system.users`).Where(`email`, `=`, userEmail).Build()

	err = pgxscan.Get(ctx, pg, &userData, query)
	if err != nil {
		return userData, err
	}

	return userData, nil
}

func GetUsers() ([]User, error) {
	var userList []User

	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return userList, err
	}

	var sqlb querybuilder.Sqlbuilder

	defer func() {
		pg.Release()
	}()

	query := sqlb.Reset().From(`system.users`).Build()

	err = pgxscan.Select(ctx, pg, &userList, query)
	if err != nil {
		return userList, err
	}

	return userList, nil
}

func (u UserCreateData) RegisterUser() (int, error) {

	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return -1, err
	}

	var sqlb querybuilder.Sqlbuilder

	defer func() {
		pg.Release()
	}()

	var insertID interface{}

	insert, err := sqlb.BuildInsert(`system.users`, u, `RETURNING id`)
	if err != nil {
		return -1, err
	}
	err = pg.QueryRow(ctx, insert).Scan(&insertID)
	if err != nil {
		if strings.Contains(err.Error(), `constraint "users_email_uindex"`) {
			return -1, errors.New("email already in use")
		} else {
			return -1, err
		}
	}
	insertIDInt, ok := insertID.(int64)
	if ok {
		return int(insertIDInt), nil
	}
	return -1, errors.New("could not decode user ID")

}

func (u UserEditData) UpdateUser() error {
	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return err
	}

	var sqlb querybuilder.Sqlbuilder

	defer func() {
		pg.Release()
	}()

	updateStruct := struct {
		Title     string `db:"title" json:"title"`
		FirstName string `db:"first_name" json:"first_name"`
		LastName  string `db:"last_name" json:"last_name"`
		Role      int    `db:"role" json:"role"`
		UpdatedAt int64  `db:"updated_at" json:"updated_at"`
	}{
		Title:     strings.TrimSpace(u.Title),
		FirstName: strings.TrimSpace(u.FirstName),
		LastName:  strings.TrimSpace(u.LastName),
		Role:      u.Role,
		UpdatedAt: time.Now().Unix(),
	}
	//var updateUUID interface{}
	userUpdateSql, err := sqlb.Where(`id`, `=`, strconv.Itoa(u.ID)).BuildUpdate(`system.users`, updateStruct, ``)
	if err != nil {
		return err
	}

	_, err = pg.Exec(ctx, userUpdateSql)
	//err = pg.QueryRow(ctx, userUpdateSql).Scan(&updateUUID)
	if err != nil {
		return err
	}
	caching.SystemCache.Delete("token:" + u.UUID)
	return nil
}

func (u User) SetUUID() error {

	if u.UUID == "" {
		return errors.New("uuid required")
	}

	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return err
	}

	var sqlb querybuilder.Sqlbuilder

	defer func() {
		pg.Release()
	}()

	updateStruct := struct {
		UUID string
	}{
		UUID: u.UUID,
	}

	userUpdateSql, err := sqlb.Where(`email`, `=`, u.Email).BuildUpdate(`system.users`, updateStruct, ``)
	if err != nil {
		return err
	}

	_, err = pg.Exec(ctx, userUpdateSql)

	if err != nil {
		return err
	}

	return nil
}

func (u User) DisableAccount() error {
	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return err
	}

	defer func() {
		pg.Release()
	}()

	dbQuery := fmt.Sprintf(`UPDATE "system"."users" SET "updated_at" = %v, "deleted_at" = %v WHERE "id" = '%v' `,
		u.PublicUpdatedAt, u.PublicDeletedAt, u.ID)

	_, err = pg.Exec(ctx, dbQuery)
	if err != nil {
		return err
	}
	caching.SystemCache.Delete("token:" + u.UUID)
	return nil
}

func (u User) EnableAccount() error {
	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return err
	}

	defer func() {
		pg.Release()
	}()

	dbQuery := fmt.Sprintf(`UPDATE "system"."users" SET "updated_at" = %v, "deleted_at" = null WHERE "id" = '%v' `,
		u.PublicUpdatedAt, u.ID)

	_, err = pg.Exec(ctx, dbQuery)
	if err != nil {
		return err
	}
	caching.SystemCache.Delete("token:" + u.UUID)
	return nil
}
