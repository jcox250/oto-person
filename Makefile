# gen generates the client and server code 
.PHONY: gen
gen:
	-rm -r gen
	mkdir -p gen/server && mkdir -p gen/client
	oto -template ./templates/server.go.plush -out ./gen/server/server.gen.go -ignore Ignorer -pkg server ./design
	oto -template ./templates/client.go.plush -out ./gen/client/client.gen.go -ignore Ignorer -pkg client ./design

.PHONY: certs
certs: 
	go build -o ./tools/certgen/certgen ./tools/certgen/main.go 
	./tools/certgen/certgen
	mv cert.pem clientcert.pem
	mv key.pem clientkey.pem
	./tools/certgen/certgen


.PHONY: build
build:
	go build -o person-service cmd/person-service/main.go
	go build -o person-service-cli cmd/person-service-cli/main.go


.PHONY: redis
redis:
	docker run -d -p 6379:6379 redis


.PHONY: run
run: redis certs build 
	./person-service -debug

.PHONY: clean
clean:
	-docker kill $$(docker ps | grep redis | awk '{ print $$1}')
	-rm *.pem
	-rm personservice
	-rm personservice-cli
