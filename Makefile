SRC_PATH=cmd
BIN_PATH=bin
BIN=go-emacs
export OUTPUT_PATH=./www

build:
	@mkdir -p ${SRC_PATH} ${BIN_PATH}
	@mkdir -p ${OUTPUT_PATH}
	@go build -C ./${SRC_PATH} -o ../${BIN_PATH}/${BIN}

run: build
	@./${BIN_PATH}/${BIN}

clean:
	@rm -rf ./${BIN_PATH}/${BIN}
	@rm ./${OUTPUT_PATH}/*
