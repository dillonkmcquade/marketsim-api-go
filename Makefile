swagger:
	GO111MODULE=off swagger generate spec -o ./api/swagger.yml --scan-models
build:
	go build -pgo=auto -o ./bin/server ./cmd/server/
install: build
	cp ./bin/server ~/.local/bin/
uninstall:
	rm ~/.local/bin/server
