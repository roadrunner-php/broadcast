clean:
	rm -rf rr-jobbroadcast
install: all
	cp rr-broadcast /usr/local/bin/rr-broadcastt
uninstall: 
	rm -f /usr/local/bin/rr-broadcasts
test:
	composer update
	go test -v -race -cover