package main

/* ==================
 - This microservice will host an api endpoint implementation
 - this serves as the bolier plate code for implementation of api endpoint behind a u-service
 - also since we plan to use rabbitmq as the broker of messages between the u-services,
 - we get to see how we can integrate amqp code in endpoint handler
================== */

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

var (
	FVerbose, FLogF bool
	logFile         string
)

const (
	AMQP_SERVER_URL = "amqp://guest:guest@msgbroker:5672" //broker url
	QUEUE_KEY       = "TestQueue1"                        // key of the queue to be used when sending receiving
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

/*
	=============

Handler to test rabbitmq server and dispatching messages to the same
Declares a simple channel + queue and sends a json hello world message across
Incases of error a 502 error is reported back
================
*/
func TestRabbit(rabbitCh *amqp.Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		body := struct {
			Msg string `json:"msg"`
		}{Msg: "Hello world"}
		jsonMsg, _ := json.Marshal(body)
		err := rabbitCh.PublishWithContext(context.TODO(), "", QUEUE_KEY, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonMsg),
		})
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("failed to send message to broker")
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		c.AbortWithStatus(http.StatusOK)
	}
}

// listenOnRabbit : sets up a background process
func listenOnRabbit(rchn, rurl string) (chan bool, *amqp.Channel, error) {
	cancel := make(chan bool) // zero channel only for interrupt flagging
	conn, err := amqp.Dial(AMQP_SERVER_URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed: listenOnRabbit")
	}
	rabbitCh, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to rabbit channel")
	}
	// Setup the queue before even starting to listening to it
	_, err = rabbitCh.QueueDeclare(QUEUE_KEY, true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to declare queue with broker")
		return nil, nil, fmt.Errorf("failed: listenOnRabbit")
	}
	msgs, err := rabbitCh.Consume(QUEUE_KEY, "", true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed: listenOnRabbit")
	}
	log.Debug("success! : listenOnRabbit")
	go func() {
		defer conn.Close()
		defer rabbitCh.Close()
		for {
			select {
			case <-cancel:
				log.Warn("interruption, listenOnRabbit")
				return
			case m := <-msgs:
				log.WithFields(log.Fields{
					"message": m.Body,
				}).Debug("received message on rabbit..")
			}
		}
	}()
	return cancel, rabbitCh, nil
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

	// start listening on the rabbit channel
	// this will start a gorutine that would run in the background listening to all the incoming messages
	// cancel channel is to interrupt the listening when program exits
	cancel, rChn, err := listenOnRabbit(QUEUE_KEY, AMQP_SERVER_URL)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("failed to start the listening channel")
	}
	defer close(cancel)
	// Hit this endpoint to see this u-service send a message to the messaging queue
	api.POST("/rabbit/test", TestRabbit(rChn))
	log.Fatal(r.Run(":8080"))
}
