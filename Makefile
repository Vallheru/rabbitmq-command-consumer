PROJECT_NAME = rabbitmq-command-consumer

build : dependencies lint compile
	
compile :
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)
	
run :
	go run ./cmd/$(PROJECT_NAME)

test : lint compile
	go test -v ./cmd/$(PROJECT_NAME)/

lint :
	golint -set_exit_status ./cmd/$(PROJECT_NAME)

dependencies :
	go get -u gopkg.in/yaml.v2
	go get -u golang.org/x/lint/golint
	go get -u github.com/stretchr/testify/assert
	go get -u github.com/google/logger
	go get -u go.uber.org/zap
	go get -u github.com/google/uuid