PORT ?= 7676
NODE_IMAGE = node:24-alpine
NODE = docker run --rm -v $(CURDIR)/web:/app -w /app $(NODE_IMAGE)

.PHONY: build ui run dev test docker up

build: ui
	go build -o bin/caern ./cmd/caern

ui: web/node_modules
	$(NODE) npm run build

web/node_modules: web/package.json web/package-lock.json
	$(NODE) npm ci
	@touch $@

run: build
	CAERN_ADDR=:$(PORT) bin/caern

dev: web/node_modules
	go build -o bin/caern ./cmd/caern
	@CAERN_ADDR=:$(PORT) bin/caern & pid=$$!; trap "kill $$pid" EXIT INT TERM; \
	docker run --rm -it -v $(CURDIR)/web:/app -w /app -p 5173:5173 \
		-e CAERN_API=http://host.docker.internal:$(PORT) $(NODE_IMAGE) npm run dev

test: web/node_modules
	go test ./...
	$(NODE) npm run check

docker:
	docker build -t caern:latest .

up:
	docker compose up --build -d
