package main

import (
	"flag"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
	"os"
)

var (
	Host                  = flag.String("host", "localhost:8080", "The host to bind to")
	DatabaseName          = flag.String("dbname", "osf20130903", "The name of your OSF database")
	DatabaseAddress       = flag.String("dbaddress", "localhost:27017", "The address of your mongodb. ie: localhost:27017")
	DatabaseSession       mgo.Session
	UserCollection        *mgo.Collection
	AccessTokenCollection *mgo.Collection
)

func main() {
	flag.Parse()
	e := echo.New()
	e.Use(mw.LoggerFromConfig(mw.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${response_time} ${response_size}\n",
		Output: os.Stdout,
	}))
	e.Use(mw.Recover())

	e.Use(standard.WrapMiddleware(cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders:   []string{"Range", "Content-Type", "Authorization", "X-Requested-With"},
		ExposedHeaders:   []string{"Range", "Content-Type", "Authorization", "X-Requested-With"},
	}).Handler))

	e.Get("/login", Login)
	e.Post("/login", Login)
	e.Get("/logout", Logout)
	e.Get("/oauth2/profile", OAuth)
	e.Get("/p3/serviceValidate", ServiceValidate)

	fmt.Println("Expecting database", *DatabaseName, "to be running at", *DatabaseAddress)

	DatabaseSession, err := mgo.Dial(*DatabaseAddress)
	if err != nil {
		panic(err)
	}
	defer DatabaseSession.Close()

	UserCollection = DatabaseSession.DB(*DatabaseName).C("user")
	AccessTokenCollection = DatabaseSession.DB(*DatabaseName).C("apioauth2personaltoken")

	fmt.Println("Listening on", *Host)
	e.Run(standard.New(*Host))
}
