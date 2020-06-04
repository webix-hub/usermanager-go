package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"webix/usermanager/api"
	"webix/usermanager/data"
	"webix/usermanager/demodata"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	Config.LoadFromFile("./config.yml")

	var dbType, connectString string
	// DB
	if Config.DB.Path != "" {
		dbType = "sqlite3"
		connectString = Config.DB.Path
	} else {
		dbType = "mysql"
		connectString = fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true",
			Config.DB.User, Config.DB.Password, Config.DB.Host, Config.DB.Database)
	}

	log.Println(dbType, connectString)
	conn, err := gorm.Open(dbType, connectString)
	if err != nil {
		log.Fatal("Can't connect to the database", err.Error())
	}
	defer conn.Close()
	dao := data.NewDAO(conn)

	demodata.ImportDemoData(conn)

	// File storage
	err = os.MkdirAll(filepath.Join(Config.Server.Data, "avatars"), 0770)
	if err != nil {
		log.Fatal("Can't create data folder", err)
	}

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(crs.Handler)

	api.Routes(r, dao, Config.Server.Data, Config.Server.Public)

	fmt.Println("Listen at port ", Config.Server.Port)
	err = http.ListenAndServe(Config.Server.Port, r)
	log.Println(err.Error())
}
