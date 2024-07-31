DEST:="root@81.200.152.157"

lint:
	golangci-lint run

serve:
	go run cmd/app/main.go

build: clean $(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) ./app/cmd/app/main.go

bin/snake-tournament: clean
	env GOOS=linux go  build -o bin/snake-tournament cmd/app/main.go

deploy: bin/snake-tournament
	scp bin/snake-tournament $(DEST):/opt/challenge/snake-tournament/
	# scp .env $(DEST):/opt/challenge/users-service.env

clean:
	rm -rf bin/* || true
	rm -rf ./app/build || true

swagger:
	swag init -g ./cmd/app/main.go --exclude go -o ./docs

migrate:
	$(APP_BIN) migrate -version $(version)

migrate.down:
	$(APP_BIN) migrate -seq down

migrate.up:
	$(APP_BIN) migrate -seq up