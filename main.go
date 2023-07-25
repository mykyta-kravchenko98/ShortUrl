package main

import (
	"database/sql"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/cache"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/config"
	repositories "github.com/mykyta-kravchenko98/ShortUrl/internal/db/postgres"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/handler"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/router"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/service"
	"github.com/mykyta-kravchenko98/ShortUrl/pkg/generator"
)

var (
	logger echo.Logger
)

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Could not create CPU profile:", err)
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	r := router.New()

	if logger == nil {
		logger = r.Logger
	}

	env := os.Getenv("environment")
	if env == "" {
		env = "dev"
	}

	//Load configuration
	conf, confErr := config.LoadConfig(env)
	if confErr != nil {
		r.Logger.Fatal("Config load failed")
	}

	// connection string
	psqlConf := conf.PostgresDB
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		psqlConf.Host, psqlConf.Port, psqlConf.User, psqlConf.Password, psqlConf.DBName)

	//open postgres connection
	db, err := sql.Open("postgres", psqlconn)
	checkError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	checkError(err)

	//Init Repository
	urlRepo := repositories.NewCurrencySnapshotDataService(db)

	//Init server group
	v1 := r.Group("/api/v1")

	//Init cache
	c := cache.InitLRUCache(100)

	//Init Id Generator
	idGen, err := generator.NewSnowflake(int64(conf.Server.DataCenterID), int64(conf.Server.MashineID))

	checkError(err)

	urlService := service.NewURLService(idGen, c, urlRepo)

	h := handler.NewHandler(urlService)
	h.Register(v1)

	r.Logger.Fatal(r.Start(fmt.Sprintf("127.0.0.1:%s", conf.Server.RESTPort)))
}

func checkError(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
