SHELL=bash
PBFLAGS=--go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
PBC?=protoc
PB=$(patsubst %.proto, %_grpc.pb.go, $(wildcard controller/pb/*.proto))
PB_GRPC=$(patsubst %.proto, %.pb.go, $(wildcard controller/pb/*.proto))


all: ${PB} ${PB_GRPC} .dep		## install dependencies and compile protobuf
pb: ${PB_GRPC} ${PB}			## compile protobuf and gRPC stuff

.pbc: 
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	touch .pbc

%_grpc.pb.go %.pb.go: %.proto .pbc
	${PBC} $< ${PBFLAGS}
dep: .dep 						## get dependencies. make pb first if errors occurred
.dep: .pbc
	go get
	touch .dep
clean: 							## remove protobuf artifacts and .dep
	rm -f ${PB} ${PB_GRPC}
	rm -f .dep .pbc


# Self Documenting Makefile from https://www.client9.com/self-documenting-makefiles/
help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
	printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)

.DEFAULT_GOAL=help
.PHONY: all help clean
