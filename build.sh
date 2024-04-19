go mod tidy
go build -buildmode=plugin -o gostable.so ./staticcheck
go build -o gostable ./golangcilint