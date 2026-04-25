.PHONY: dev fe-dev be-dev fe-build be-build build lint be-test

be-dev: 
	cd backend && go run .

fe-dev:
	cd frontend && npm run dev

build: be-build fe-build

be-build:
	go build -o bin/chirpy ./backend

fe-build:
	cd frontend && npm run build

be-test:
	cd ./backend && go test ./... 

fe-lint:
	cd frontend && npm run lint





