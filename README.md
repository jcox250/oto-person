# Oto-person


## Build

`make certs` - will create the server and client certs that you'll need to run the service and hit it using the client
`make build` to build the `person_server` binary and `person_server-cli`

## Run

`make run` bring up the service and everything it needs to function which includes the client & server certs and a redis container.

Add a person 
`./person-service-cli -method Add -payload '{"id": "1", "name": "Peter Rabbit", "age": 5}'`

Get a person
```
./person-service-cli -method Show -payload '{"id": "1"}'
&{Peter Rabbit 5}
```

When you're finished you can run `make clean` which will remove all build artifacts and stop the running redis container
