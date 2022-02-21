package db

import (
	"context"
	"firebase-sso/helpers"
	"firebase-sso/helpers/env"
	"fmt"
	querybuilder "github.com/KirksFletcher/Go-SQL-Query-Builder-Golang"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type PostgresDBConnection struct {
	Db *pgxpool.Pool
}

var PostgresDB PostgresDBConnection

func (p *PostgresDBConnection) SetupPostgresDB() (*PostgresDBConnection, error) {

	connString, err := pgxpool.ParseConfig("postgresql://" + env.Get("DBHOST") + ":" + env.Get("DBPORT") + "/" + env.Get("DBNAME") + "?user=" + env.Get("DBUSER") + "&password=" + env.Get("DBPASSWORD") + "&statement_cache_mode=describe")
	if err != nil {
		log.Fatalln(err)
	}

	helpers.BasePath()
	connString.MaxConns = 50

	conn, err := pgxpool.ConnectConfig(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	PostgresDB.Db = conn

	return &PostgresDBConnection{
		Db: conn,
	}, nil
}

func (p *PostgresDBConnection) SetupPostgresDBX() {

	connString, err := pgxpool.ParseConfig("postgresql://" + env.Get("DBHOST") + ":" + env.Get("DBPORT") + "/" + env.Get("DBNAME") + "?user=" + env.Get("DBUSER") + "&password=" + env.Get("DBPASSWORD") + "&statement_cache_mode=describe")
	if err != nil {
		log.Fatalln(err)
	}

	connString.MaxConns = 50

	conn, err := pgxpool.ConnectConfig(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	PostgresDB.Db = conn
}

func (p *PostgresDBConnection) PingDB() {
	pg, err := PostgresDB.Db.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		pg.Release()
	}()

	fmt.Println("POSTGRESQL DB IS ALIVE : " + env.Get("DBPORT") + " @ " + env.Get("DBHOST"))
}

func CheckSchemaExists(schema string) bool {

	pg, err := PostgresDB.Db.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
		return false
	}

	defer func() {
		pg.Release()
	}()

	var dum string
	var sqlb querybuilder.Sqlbuilder

	sqlb.
		Select("schema_name").
		From("information_schema.schemata").
		Where(`schema_name`, `=`, schema)

	err = pg.QueryRow(context.Background(), sqlb.Build()).Scan(&dum)
	if err != nil {
		return false
	}

	return true
}

func CheckTableExists(schema string, table string) bool {

	pg, err := PostgresDB.Db.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
		return false
	}

	defer func() {
		pg.Release()
	}()

	var dum string
	var sqlb querybuilder.Sqlbuilder

	sqlb.
		From("pg_tables").
		Where(`schemaname`, `=`, schema).
		Where(`tablename`, `=`, table)

	err = pg.QueryRow(context.Background(), sqlb.Build()).Scan(&dum)
	if err != nil {
		fmt.Println(table + " table not found in " + schema + " schema.")
		fmt.Println(err.Error())
		return false
	}

	return true
}
