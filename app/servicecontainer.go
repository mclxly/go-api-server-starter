package app

import (	
	"sync"
	"os"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
    "github.com/go-redis/redis"

	// "github.com/irahardianto/service-pattern-go/controllers"
	// "github.com/irahardianto/service-pattern-go/repositories"
	// "github.com/irahardianto/service-pattern-go/infrastructures"
	// "github.com/irahardianto/service-pattern-go/services"
	"github.com/mclxly/go-api-server-starter/database"
)

type IServiceContainer interface {
	// InjectPlayerController() controllers.PlayerController
	// ConnectDB() *database.GlobalDB
	PreInit() *App
	Init() *App
}

type App struct{
    // public
    Port string
    Redis *redis.Client

    // private
	db *database.GlobalDB
	logFile *os.File    
}

// func (k *kernel) InjectPlayerController() controllers.PlayerController {

//   sqlConn, _ := sql.Open("sqlite3", "/var/tmp/tennis.db")
//   sqliteHandler := &infrastructures.SQLiteHandler{}
//   sqliteHandler.Conn = sqlConn

//   playerRepository := &repositories.PlayerRepository{sqliteHandler}
//   playerService := &services.PlayerService{&repositories.PlayerRepositoryWithCircuitBreaker{playerRepository}}
//   playerController := controllers.PlayerController{playerService}

//   return playerController
// }

func (k *App) PreInit() *App {
	// setting log
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(&log.TextFormatter{})

	// log file
    if k.logFile == nil {
        file, err := os.OpenFile("./storage/logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err == nil {
            // log.SetOutput(io.MultiWriter(os.Stdout, file))
            mio := io.MultiWriter(os.Stdout, file)
            log.SetOutput(mio)
            gin.DefaultWriter = mio
            k.logFile = file
        } else {
            log.Info("Failed to log to file, using default stderr")
        }
    }

    // redis client
    if k.Redis == nil {
        k.Redis = redis.NewClient(&redis.Options{
            Addr:     "localhost:6379",
            Password: "", // no password set
            DB:       0,  // use default DB
        })

        // test
        _, err := k.Redis.Ping().Result()
        if err != nil {
            log.Panic("Failed to connect Redis!")
        }
    }

    return k
}

func (k *App) Init() *App {
	// database
	if k.db == nil {
		db, err := database.NewConnect()
		if err != nil {
			log.Print("Database connect failed.")
			return k
		}
		log.Print("Database connected.")
		k.db = db
	}

	return k
}

func (k *App) CleanUp() {
    log.Info("App cleanUp.")
    
	if k != nil {
		k.logFile.Close()
	}
}

// func (k *App) GetPort() string {
//     if k != nil {
//         return k.appPort
//     }
// }

// func (k *App) connectDB() *database.GlobalDB {
// 	db, err := database.ConnectDB()
// 	if err != nil {
// 		return nil
// 	}	
// 	return db
// }

var (
	k             *App
	containerOnce sync.Once
)

func ServiceContainer() IServiceContainer {
	if k == nil {
		containerOnce.Do(func() {
			k = &App{Port: "8080"}
		})
	}
	return k
}