# gen generates the client and server code
.PHONY: gen
gen:
	-rm -r gen
	mkdir -p gen/server && mkdir -p gen/client
	oto -template ./templates/server.go.plush -out ./gen/server/server.gen.go -ignore Ignorer -pkg server ./design
	oto -template ./templates/client.go.plush -out ./gen/client/client.gen.go -ignore Ignorer -pkg client ./design

build:
	go build -o person_server cmd/person/main.go
	go build -o person_server-cli cmd/person-cli/main.go
