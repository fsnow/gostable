go mod tidy
go build -buildmode=plugin -o gostable.so ./plugin
go build -o gostable ./standalone