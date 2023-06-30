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

Clone the repository to your local machine. `cd` into directory 

```sh
docker-compose --env-file dev.env build
```

```sh 
docker-compose --env-file dev.env up
```