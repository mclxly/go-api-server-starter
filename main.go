/*
Desc:
1. Init Log
2. Init App
*/
package main

import (
	"fmt"
    "net/http"
    "os"
    "os/signal"    
    "syscall"
    "time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper" // Application Config
    "github.com/gin-gonic/gin" // http
    "github.com/appleboy/gin-jwt" // jwt
    // "github.com/go-redis/redis" // redis client
    "github.com/ulule/limiter" // rate limiter
    mgin "github.com/ulule/limiter/drivers/middleware/gin"
    sredis "github.com/ulule/limiter/drivers/store/redis"

	"github.com/mclxly/go-api-server-starter/app"
)

func readConfig() {
	log.Printf("Reading config...")

	viper.SetConfigName("env")  // name of config file (without extension)
	viper.AddConfigPath("./")   // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func myInit() (*app.App) {    
    appInst := app.ServiceContainer().PreInit()

	readConfig()

    appInst.Init()

	app_name := viper.GetString("app_name")    
	// fmt.Printf("\n%s\n\n", app_name)
	// log.Printf("Starting %s...", app_name)
	// Log as JSON instead of the default ASCII formatter.
	
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// log.SetOutput(os.Stdout)

	// log.WithFields(log.Fields{"request_id": request_id, "user_ip": user_ip})

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Printf("Init %s...", app_name)

    appInst.Port = viper.GetString("app_port")
    log.Printf("%v", appInst)
    return appInst
}

func startHttp(appInst *app.App) {
    // Disable Console Color, you don't need console color when writing the logs to file.
    // gin.DisableConsoleColor()
    
    // r := gin.Default()

    r := gin.New()

    // Global middleware
    // Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
    // By default gin.DefaultWriter = os.Stdout
    r.Use(gin.Logger())

    // Recovery middleware recovers from any panics and writes a 500 if there was one.
    r.Use(gin.Recovery())

    // ----------------------------------Middleware
    // rate limiter
    rate, err := limiter.NewRateFromFormatted("3-M")
    if err != nil {
        panic(err)
    }

    // Create a store with the redis client.
    store, err := sredis.NewStoreWithOptions(appInst.Redis, limiter.StoreOptions{
        Prefix:   viper.GetString("app_name"),
        MaxRetry: 3,
    })
    if err != nil {
        log.Fatal(err)
        return
    }

    // Create a new middleware with the limiter instance.
    rateLimiterMiddleware := mgin.NewMiddleware(limiter.New(store, rate))
    _ = rateLimiterMiddleware

    // setup
    r.Use(rateLimiterMiddleware)

    // the jwt middleware
    authMiddleware := &jwt.GinJWTMiddleware{
        Realm:      "test zone",
        Key:        []byte("secret key"),
        Timeout:    time.Hour,
        MaxRefresh: time.Hour,
        Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
            if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
                return userId, true
            }

            return userId, false
        },
        Authorizator: func(userId string, c *gin.Context) bool {
            if userId == "admin" {
                return true
            }

            return false
        },
        Unauthorized: func(c *gin.Context, code int, message string) {
            c.JSON(code, gin.H{
                "code":    code,
                "message": message,
            })
        },
        // TokenLookup is a string in the form of "<source>:<name>" that is used
        // to extract token from the request.
        // Optional. Default value "header:Authorization".
        // Possible values:
        // - "header:<name>"
        // - "query:<name>"
        // - "cookie:<name>"
        TokenLookup: "header:Authorization",
        // TokenLookup: "query:token",
        // TokenLookup: "cookie:token",

        // TokenHeadName is a string in the header. Default value is "Bearer"
        TokenHeadName: "Bearer",

        // TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
        TimeFunc: time.Now,
    }
    
    // ************************************************
    // router list
    // ************************************************
    // ---------------------------------test purpose
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "goooooooooo.")
    })

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
            "version": "20171202",
        })
    })

    // ----------------------------------public route
    r.POST("/login", authMiddleware.LoginHandler)

    // ----------------------------------secure route
    auth := r.Group("/v1")
    auth.Use(authMiddleware.MiddlewareFunc())
    {
        auth.GET("/hello", helloHandler)
        auth.GET("/refresh_token", authMiddleware.RefreshHandler)
    }


    log.Info("listen and serve on 0.0.0.0:8080")
    // r.Run() // listen and serve on 0.0.0.0:8080    
    r.Run(appInst.Port)
}

// -------------------------------------------
// handle request
// -------------------------------------------
func helloHandler(c *gin.Context) {
    claims := jwt.ExtractClaims(c)
    c.JSON(200, gin.H{
        "userID": claims["id"],
        "text":   "Hello World.",
    })
}

func main() {
    // Prepare for graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
    signal.Notify(stop, os.Interrupt)

    appInst := myInit()
	// log.Print(connStr)
	// log.Printf("%+v\n", appInst)

    go startHttp(appInst)

    // Block until the OS signal
    <-stop

    log.Info("Caught graceful shutdown signal")
    
	// cleanUp()

    // must be last line
    appInst.CleanUp()
}

// func cleanUp() {	
// }
