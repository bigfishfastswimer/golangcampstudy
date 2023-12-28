.PHONY: docker
clean:
	@echo "cleaning local build" && rm webook || true
	@go mod tidy
docker: clean
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f  flycash/webook:v0.0.1
	@docker build -t  flycash/webook:v0.0.1 .
docker-k8s: clean
	@GOOS=linux GOARCH=arm go build --tags=k8s -o webook .
	@docker rmi -f  fishereatworld2018/webook:v0.0.1
	@docker build -t  fishereatworld2018/webook:v0.0.1 .