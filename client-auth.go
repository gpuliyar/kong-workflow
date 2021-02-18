package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Kong/go-pdk"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
)

const (
	dbHost     = "db"
	dbPort     = 5432
	dbUser     = "kong"
	dbPassword = "kong"
	dbName     = "kong"
	dbDriver   = "postgres"

	redisHost    = "redis"
	redisPort    = 6379
	redisNetwork = "tcp"
)

type Config struct {
	ClientID    string
	ClientToken string
	Jwt         string
}

func New() interface{} {
	return &Config{}
}

func (conf *Config) Access(kong *pdk.PDK) {
	responseHeaders := make(map[string][]string)
	responseHeaders["Content-Type"] = append(responseHeaders["Content-Type"], "application/json")

	dbConn, err := connectDb(kong)
	defer dbConn.Close()
	if err != nil {
		kong.Log.Err("unable to connect db", err.Error())
		kong.Response.SetStatus(401)
		kong.Response.Exit(500, `{"message": "db connection failed", "status": "failed"}`, responseHeaders)
		return
	}

	redisConn, err := connectRedis(kong)
	defer redisConn.Close()
	if err != nil {
		kong.Log.Err("unable to connect redis", err.Error())
		kong.Response.SetStatus(401)
		kong.Response.Exit(500, `{"message": "redis connection failed", "status": "failed"}`, responseHeaders)
		return
	}

	xClientID, err := kong.Request.GetHeader("X-Client-ID")
	if err != nil {
		kong.Log.Err("auth header - x-client-id - not found", err.Error())
		kong.Response.SetStatus(401)
		kong.Response.Exit(500, `{"message": "x-client-id header not found", "status": "failed"}`, responseHeaders)
		return
	}

	xClientToken, err := kong.Request.GetHeader("X-Client-Token")
	if err != nil {
		kong.Log.Err("auth header - x-client-token - not found", err.Error())
		kong.Response.SetStatus(401)
		kong.Response.Exit(500, `{"message": "x-client-token header not found", "status": "failed"}`, responseHeaders)
		return
	}

	conf.ClientID = xClientID
	conf.ClientToken = xClientToken
	if err = conf.validate(dbConn, kong); err != nil {
		kong.Log.Err("x-client-id and x-client-token validation failed", err.Error())
		kong.Response.SetStatus(401)
		kong.Response.Exit(500, `{"message": "client validation failed", "status": "failed"}`, responseHeaders)
		return
	}

	conf.fetchJwt(redisConn, kong)

	bearerText := fmt.Sprintf("Bearer %s", conf.Jwt)
	kong.Response.SetHeader("Authorization", bearerText)
}

func (conf *Config) fetchJwt(conn redis.Conn, kong *pdk.PDK) {
	if jwt, err := redis.String(conn.Do("HGET", conf.ClientID, "jwt")); err != nil {
		kong.Log.Err("jwt cache fetch failed", err)
		conf.Jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImV4cCI6MTkwMDAwMDAwMCwiaXNzIjoiZW52b3kiLCJzdWIiOiJlbnZveSIsImF1ZCI6InZlbmRvcnMifQ.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.lPpSEDd6aDLJPO6UpL4C-cunqWkZxOTJR30etdyHLD0"
		if _, err := conn.Do("HMSET", conf.ClientID, "jwt", jwt); err != nil {
			kong.Log.Err("unable to cache jwt", err)
		}
	} else {
		conf.Jwt = jwt
	}
}

func (conf *Config) validate(conn *sql.DB, kong *pdk.PDK) error {
	sqlSelect := "select client_token from vendors where client_id = $1"
	row := conn.QueryRow(sqlSelect, conf.ClientID)

	var dbClientToken string

	switch err := row.Scan(&dbClientToken); err {
	case sql.ErrNoRows:
		kong.Log.Err("unable to find client id", err.Error())
		return err
	case nil:
		kong.Log.Info("client id exists in db")
	default:
		kong.Log.Err("error : %x", err.Error())
		return err
	}

	if dbClientToken != conf.ClientToken {
		kong.Log.Err("client token doesn't match db record")
		return errors.New("client token doesn't match db record")
	} else {
		kong.Log.Info("client id and token matched successfully")
	}

	return nil
}

func connectDb(kong *pdk.PDK) (*sql.DB, error) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open(dbDriver, connInfo)
	if err != nil {
		kong.Log.Err("unable to connect db", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		kong.Log.Err("unable to ping db", err)
		return nil, err
	}

	kong.Log.Info("db connected!")

	return db, nil
}

func connectRedis(kong *pdk.PDK) (redis.Conn, error) {
	connInfo := fmt.Sprintf("%s:%d", redisHost, redisPort)

	conn, err := redis.Dial(redisNetwork, connInfo)
	if err != nil {
		kong.Log.Err("unable to dial redis", err)
		return nil, err
	}

	kong.Log.Info("redis connected!")

	return conn, nil
}
