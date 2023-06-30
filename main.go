package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	FVerbose, FLogF bool
	logFile         string
)

/*
==================
- CORS enabling all cross origin requests for all verbs except OPTIONS
- this will be applied to all api across the board during the develpment stages
- do not apply this middleware though for routes that deliver web static content
====================
*/
func CORS(c *gin.Context) {
	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		// Everytime we receive an OPTIONS request,
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real
		// request using any other method than OPTIONS
		c.AbortWithStatus(http.StatusOK)
	}
}

func init() {
	/* ======================
	-verbose=true would mean log.Debug can work
	-verbose=false would mean log.Debug will be hidden
	-flog=true: all the log output shall be onto a file
	-flog=false: all the log output shall be on stdout
	- We are setting the default log level to be Info level
	======================= */
	flag.BoolVar(&FVerbose, "verbose", false, "Level of logging messages are set here")
	flag.BoolVar(&FLogF, "flog", false, "Direction in which the log should output")
	// Setting up log configuration for the api
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)
	// By default the log output is stdout and the level is info
	log.SetOutput(os.Stdout)     // FLogF will set it main, but dfault is stdout
	log.SetLevel(log.DebugLevel) // default level info debug but FVerbose will set it main
	logFile = os.Getenv("LOGF")
}

func main() {

	/* ++++++++++++++++++++++
	command line arguments as configuration for
	- logging verbosity
	- direction of logs - stdout , file
	++++++++++++++++++++++*/

	flag.Parse() // command line flags are parsed
	log.WithFields(log.Fields{
		"verbose": FVerbose,
		"flog":    FLogF,
	}).Info("Log configuration..")
	if FVerbose {
		log.SetLevel(log.DebugLevel)
	}
	if FLogF {
		lf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to connect to log file, kindly check the privileges")
		} else {
			log.Infof("Check log file for entries @ %s", logFile)
			log.SetOutput(lf)
		}
	}
	log.Info("Starting api server..")
	defer log.Warn("Exiting api server..")
	/* +++++++++++++++++++++
	Setting up the gin server, add your routes here

	+++++++++++++++++++++*/
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	api := r.Group("api", CORS)
	// just to test if the api server is running
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app": "wicwug",
			// "logs":      logFile,
			// "verblog":   FVerbose,
			// "logtofile": FLogF,
		})
	})
	log.Fatal(r.Run(":8080"))
}
