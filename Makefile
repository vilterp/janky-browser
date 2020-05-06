run:
	go run -v main.go

serve-samples:
	# npm install -g http-server
	cd testdata && http-server
