PBFLAGS=--go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
PBC?=protoc
PB=$(patsubst %.proto, %_grpc.pb.go, $(wildcard controller/pb/*.proto))
PB_GRPC=$(patsubst %.proto, %.pb.go, $(wildcard controller/pb/*.proto))


all: ${PB} ${PB_GRPC} .dep

%_grpc.pb.go %.pb.go: %.proto
	${PBC} $< ${PBFLAGS}

.dep:
	go get
	touch .dep
clean:
	rm ${PB} ${PB_GRPC}