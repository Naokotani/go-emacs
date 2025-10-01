SRC_PATH=cmd
BIN_PATH=bin
PORT=8080
BIN=go-emacs
export CONFIG_PATH=/home/naokotani/code/go/go-emacs/config.toml

build:
	@mkdir -p ${SRC_PATH} ${BIN_PATH}
	@go build -C ./${SRC_PATH} -o ../${BIN_PATH}/${BIN} 2>&1

run: build
	@./${BIN_PATH}/${BIN}

clean:
	@rm -rf ./${BIN_PATH}/${BIN}
	@rm ./${OUTPUT_PATH}/*

serve: run
	@echo "Starting test server on port ${PORT}"
	@python3 -m http.server ${PORT} --directory www
