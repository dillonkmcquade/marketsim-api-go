swagger:
	GO111MODULE=off swagger generate spec -o ./api/swagger.yml --scan-models
build:
	 go build -o ./bin/server

