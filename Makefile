GC=go

all:
	$(GC) build -o bin/fzip fzip.go

install:
	cp bin/* /usr/local/bin/
