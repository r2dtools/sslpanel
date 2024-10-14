initdb:
	. scripts/env.sh; \
	cd backend/cmd/cli; \
	go run main.go migrations drop; \
	go run main.go migrations up; \
	go run main.go fixtures;

golint:
	cd backend; \
	golangci-lint run ./... 

start-back-dev:
	. scripts/env.sh; \
	cd backend/cmd/app; \
	go run main.go;

start-front-dev:
	. scripts/env.sh; \
	cd frontend; \
	npm run dev;

start-front-prod:
	. scripts/env.sh; \
	cd frontend; \
	npm start;

build-front:
	cd frontend; \
	npm run build

build-image:
	docker build -t codegrove .
