gen: gen_tools gen_clean
	go generate ./...

gen_clean:
	rm -rf client models
	mkdir -p client models

gen_tools:
	go install github.com/go-swagger/go-swagger/cmd/swagger
