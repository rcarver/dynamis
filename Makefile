
check: build test vet lint

test: start-dynamodb
	DYNAMODB_HOSTPORT=$$(docker-machine ip $${DOCKER_MACHINE_NAME}):${dynamodb-port} \
		go test ./...

vet:
	go vet ./...

lint:
	golint ./...

build:
	go install .

.PHONY: check test vet lint build


# Docker service commands
# =============================================================================

# A list of all docker services for startall/stopall.
docker-svcs=dynamodb

# Start all docker services.
startall: $(patsubst %,start-%,${docker-svcs})

# Stop all docker services.
stopall: $(patsubst %,stop-%,${docker-svcs})

# Start any docker service by name.
start-%:
	$(MAKE) tmp/docker/$*

# Stop any docker service by name.
stop-%: 
	test -f tmp/docker/$*
	docker stop $$(cat tmp/docker/$*)
	docker rm $$(cat tmp/docker/$*)
	rm -f tmp/docker/$* 


# =============================================================================

dynamodb-port=8000
tmp/docker/dynamodb:
	mkdir -p $(dir $@)
	docker run --cidfile=$@ --name dynamodb -d -p 8000:${dynamodb-port} dynamodb



