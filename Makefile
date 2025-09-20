SRC_PATH=src
BIN_PATH=bin
BIN=go-emacs

build:
	@mkdir -p ${SRC_PATH} ${BIN_PATH}
	@go build -C ./${SRC_PATH} -o ../${BIN_PATH}/${BIN}

run: build
	@./${BIN_PATH}/${BIN}

clean:
	@rm -rf ./${BIN_PATH}/${BIN}
