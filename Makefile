initdb:
	. scripts/env.sh; \
	cd backend/cmd/cli; \
	go run cli.go migrations drop; \
	go run cli.go migrations up; \
	go run cli.go fixtures;

golint:
	cd backend; \
	golangci-lint run ./... 

start-back-dev:
	. scripts/env.sh; \
	cd backend/cmd/panel; \
	go run panel.go;

start-front-dev:
	. scripts/env.sh; \
	cd frontend; \
	npm run dev;

build-front:
	cd frontend; \
	npm run build

build-image:
	docker build -t sslpanel .
