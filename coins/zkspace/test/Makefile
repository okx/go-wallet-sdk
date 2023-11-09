all: run

download:
	mkdir -p libs
ifeq ($(shell uname -s),Darwin)
	test -f ./libs/libzks-crypto.a || curl -L https://github.com/zksync-sdk/zksync-crypto-c/releases/download/v0.1.2/zks-crypto-x86_64-apple-darwin.a --output ./libs/libzks-crypto.a
else
	test -f ./libs/libzks-crypto.so || curl -L https://github.com/zksync-sdk/zksync-crypto-c/releases/download/v0.1.2/zks-crypto-linux-x64.so --output ./libs/libzks-crypto.so
endif

build: download
	test -f test || CGO_LDFLAGS="-L./libs" go build

run: build
ifeq ($(shell uname -s),Darwin)
	./test
else
	LD_LIBRARY_PATH=./libs ./test
endif