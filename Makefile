# Go parameters
GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOTEST=${GOCMD} test
GOGET=${GOCMD} get

BINARY_NAME=base


all: clean build 

SRC=$(ls "./src/*.go")

build: src/*.go
	# building binary ...
	$(GOBUILD)  -o $(BINARY_NAME) -v src/*.go 
	# compressing ...
	# upx -9 ${BINARY_NAME}
	# done

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

mods:
	# MOD: test
	go build --buildmode=plugin -o plugin/test.so src/modules/test/test.go
