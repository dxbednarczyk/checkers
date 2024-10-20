run filename:
    go run cmd/checkers/checkers.go {{filename}}

gen:
    go run cmd/drivegen/drivegen.go

build:
    go build -o checkers cmd/checkers/checkers.go

logs:
    go run cmd/checkers/checkers.go logs/xld.log
    go run cmd/checkers/checkers.go logs/xld-error.log
    go run cmd/checkers/checkers.go logs/xld-noacc.log

lint:
    gofumpt -l -w .

    go vet ./...
    go mod tidy
    go clean

    golangci-lint -c .golangci-lint.yml run
