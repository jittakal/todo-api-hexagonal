## Steps Followed After Cloning Github Repository

```
$ # clone githib.com/jittakal/todo-api-hexagonal
# cd todo-api-hexagonal
```

```
$ go install github.com/swaggo/swag/cmd/swag@latest
```

```
$ export GOPATH=/Users/jittakal/go
$ export PATH=$PATH:$GOPATH/bin
$ swag init -g ./internal/adapters/handlers/http.go
```

```
$ go run cmd/main.go
$ # OR
$ go run cmd/main.go --handler=http --repository=in-memory
```

```
protoc -I. --go_out=plugins=grpc:. task.proto
or
protoc -I. --go_out=plugins=grpc:. task.proto

```


```
$ cd to_project_root_folder
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_
out=. --go-grpc_opt=paths=source_relative api/proto/task.proto
```
