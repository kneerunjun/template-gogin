# template-gogin
Here is a quick start template for go-gin application behind a docker container. Docker-compose is the simple orchestrator for it. 

### Why do I need a template ?
---- 

Imagine working with a [microservices architecture]() where in there are multiple teams entrusted with the job of building an maintaining the service all through the lifecycle of the application. Each of the microservice will have their own business logic but in terms of the following it would be desirable to have a kind of uniformity across the microservices 

1. Container runtime used 
2. The way the services log, and the configuration behind it 
3. Error handling in the microservices 
4. Documentation of each of the functions 
5. Patterns of design used in all the services 

We are hence making up a template that can when used a __boiler-plate__, be re-used across the board. 

1. This saves a lot of setup time 
2. Teams are then fungible and members can be transient across the modules of the project.

### What is the template all about ?
-----

1. __docker-compose__ : A ready made docker-compose yml, that you can use to add your services. Basic api service is included as a sample
2. __Dockerfile__ for golang container that uses alpine linux. This can build a simple container with basic linux directory structure and go build and run commands
3. __Go GIN__ application with simple one `/api/ping` endpoint


### How to use the template?
-----

Clone the repository to your local machine. `cd` into directory and issue the following commands 

```sh
docker-compose --env-file dev.env build
```

```sh 
docker-compose --env-file dev.env up
```

### RabbitMQ communication :
---

We are planning to use [RabbitMQ](https://www.rabbitmq.com/) as a means of intra-services communication. Instead of the `gRPC` sphagetti a broker like RabbitMQ is far more convenient while also having plenty room for horizontal scaling. While socket communication can be an alternative, it has a reputation of becoming grossly ugly when the number of services goes up. Not to even mention having to publish messages to multiple listeners can be even more daunting. 

A message broker is ideally suited to shuttle messages between various microservices. 

#### Observations about the rabbit cotainer  
----

```
services:
  msgbroker:
    image: rabbitmq:3-management-alpine
    container_name: ctn_msgbrokr
    ports:
      - 5672:5672   # for sender and consumer connections
      - 15672:15672 # for serve RabbitMQ GUI
    volumes:
      - ${HOME}/dev-rabbitmq/data/:/var/lib/rabbitmq
      - ${HOME}/dev-rabbitmq/log/:/var/log/rabbitmq
    networks:
      - dev-network
    healthcheck:
      test: "exit 0"
```

port 5672 is used to communicate (dispatch / listen) on messages exchange. While 15672 is for starting a small observer web application. Observe `dev-network` is the VPN on which the message broker container actually operates. Your Go app should then be also on the same network if it has to talk. 
Sequence in which the containers are pulled up is also vital. Here since the broker has to come up before the microservice, `healthcheck` section defines how service can be tested for uptime. Other services (dependent) can have a litmus test to know if the broker is up and running. 

> Note: Incase you are on windows `$HOME` directive may not work and this can be replaced with any path that you choose to use for rabbit storage. 


connecting to the rabbit from within the code 

```go
// sample code on how to connect from the url 
```

#### Sending messages as json :
---

Upon hitting a post url - handler does the job 


#### Receiving the messages as json:
------

Background process that tracks in the incoming messages
can be interrupted when main process calls the channel 


At the center of all communications is basically the channel. Once the channel is setup communication is all about just sending the `[]byte` across the appropriate queue
Upon receiving the message since this is a sample code it just logs the message 

### Observing the Rabbit @ 15672 :
----
