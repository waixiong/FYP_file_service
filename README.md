<!-- # Where you should clone this project -->
<!-- At `$GOPATH/src/rj`, git clone this project. -->

# Pre
- Download Golang
- Download Docker
- Download protobuf compiler (http://google.github.io/proto-lens/installing-protoc.html)


### Folder Structure
- api
  - proto
    - ${services}.proto
- cmd
  - <server>
    - main.go (setup server)
  - key (dir that save key)
- pkg
  - api
    - ${services}
      - ${services}.pb.go
  - service
    - ${services}
      - ${services}.go (file, main part of the service)
      - src (dir for any extra go program)


For more of google services, please look at https://github.com/googleapis/google-api-go-client, currently google service is used for authentication.


# Run the services
The are three ways to run the services 
- use `go run cmd/server/main.go`
- compile docker image with `scripts/image.sh`, then run the image
- compile docker image with `scripts/image.sh`, then run docker-compose @ `deployments/docker-compose.yaml`