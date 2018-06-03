libnss_etcd.so.2: cmd/libnss_etcd/*.go *.h *.c
	go build -buildmode=c-shared -o $@ ./cmd/libnss_etcd
	chmod +x $@

build: libnss_etcd.so.2

install: libnss_etcd.so.2
	cp $< /usr/lib/$<

clean:
	rm -rf libnss_etcd.so.2 libnss_etcd.so.h

.PHONY: build install clean
