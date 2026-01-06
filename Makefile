all: build

build: discord-build
	go mod tidy
	go build -o bin/out.exe app/app.go

discord-build: 