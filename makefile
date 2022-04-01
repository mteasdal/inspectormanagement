compile: ./cmd/InspectorManager/main.go
	 go build -o ./cmd/InspectorManager/inspector ./cmd/InspectorManager/main.go

build: compile
	 mv ./cmd/InspectorManager/inspector ./bin/inspector


testinspector: ./inspector/inspector2_test.go
	go test ./inspector