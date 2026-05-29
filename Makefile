all: build deploy

lint:
	golangci-lint run --fix

serve:
	go run cmd/app/main.go

build: clean bin/snake-tournament

bin/snake-tournament: clean
	env GOOS=linux GOARCH=amd64 go  build -o bin/snake-tournament cmd/app/main.go

deploy: bin/snake-tournament
	mkdir -p /opt/challenge/snake-tournament
	cd /opt/challenge && docker compose stop
	cp bin/snake-tournament /opt/challenge/snake-tournament/
	#cp .env /opt/challenge/users-service.env
	cd /opt/challenge && docker compose start

clean:
	rm -rf bin/* || true
	rm -rf ./app/build || true

swagger:
	swag init --parseDependency --parseInternal -g ./cmd/app/main.go --exclude go -o ./docs

migrate:
	$(APP_BIN) migrate -version $(version)

migrate.down:
	$(APP_BIN) migrate -seq down

migrate.up:
	$(APP_BIN) migrate -seq up

mockery:
	go tool mockery