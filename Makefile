protogen-todo:
	protoc -I./lib -I. \
	--go_out=plugins=grpc:. \
	--grpc-gateway_out=logtostderr=true,grpc_api_configuration=./proto/todo/todo.http.yaml:. \
	--swagger_out=logtostderr=true,grpc_api_configuration=./proto/todo/todo.http.yaml:. \
	./proto/todo/todo.proto

	