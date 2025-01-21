all: ccallbacks

ccallbacks: libclibrary.a goprog.go cfunc.go
	go build -o $@ *.go

libclibrary.a: libclibrary.o
	ar cru libclibrary.a clibrary.o

libclibrary.o: clibrary.h clibrary.c
	gcc -c clibrary.c

clean:
	rm -rf *.a
	rm -rf *.o
	rm -rf ccallbacks