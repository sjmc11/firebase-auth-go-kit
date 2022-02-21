package models

import (
	"context"
	"errors"
	"firebase-sso/db"
	"fmt"
	querybuilder "github.com/KirksFletcher/Go-SQL-Query-Builder-Golang"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gosimple/slug"
	"strconv"
)

type Project struct {
	ID        int    `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	Slug      string `db:"slug" json:"slug"`
	Created   int64  `db:"created" json:"created"`
	CreatedBy int    `db:"created_by" json:"created_by"`
}

func FetchProjectByID(projectID int) (Project, error) {
	var projectData Project

	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return projectData, err
	}

	var sqlb querybuilder.Sqlbuilder

	defer func() {
		pg.Release()
	}()

	query := sqlb.Reset().From(`system.user_projects`).Where(`id`, `=`, strconv.Itoa(projectID)).Build()

	err = pgxscan.Get(ctx, pg, &projectData, query)
	if err != nil {
		return projectData, err
	}

	return projectData, nil
}

func (proj Project) Create() (int, error) {
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
	projectInsert := struct {
		Title     string `db:"title" json:"title"`
		Slug      string `db:"slug" json:"slug"`
		Created   int64  `db:"created" json:"created"`
		CreatedBy int    `db:"created_by" json:"created_by"`
	}{
		Title:     proj.Title,
		Slug:      slug.Make(proj.Title),
		Created:   proj.Created,
		CreatedBy: proj.CreatedBy,
	}

	insert, err := sqlb.BuildInsert(`system.user_projects`, projectInsert, `RETURNING id`)
	if err != nil {
		return -1, err
	}

	err = pg.QueryRow(ctx, insert).Scan(&insertID)
	if err != nil {
		fmt.Println("err")
		return -1, err
	}
	insertIDInt, ok := insertID.(int64)
	if ok {
		return int(insertIDInt), nil
	}
	return -1, errors.New("could not decode project ID")

}

func (proj Project) Delete() error {
	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return err
	}

	defer func() {
		pg.Release()
	}()

	commandTag, err := pg.Exec(ctx, "delete from system.user_projects where id=$1", proj.ID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no project found to delete")
	}

	return nil
}

func (u User) GetProjects() ([]Project, error) {
	var projectList []Project

	ctx := context.Background()
	pg, err := db.PostgresDB.Db.Acquire(ctx)

	if err != nil {
		return projectList, err
	}

	var sqlb querybuilder.Sqlbuilder

	defer func() {
		pg.Release()
	}()

	query := sqlb.Reset().From(`system.user_projects`).Where(`created_by`, `=`, strconv.Itoa(u.ID)).Build()

	err = pgxscan.Select(ctx, pg, &projectList, query)
	if err != nil {
		return projectList, err
	}

	return projectList, nil
}
