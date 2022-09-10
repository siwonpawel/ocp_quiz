clear:
	rm -rf build
	podman rmi -f ocp_quiz:latest

build: clear
	go mod vendor
	CGO_ENABLED=0 go build -o ./build/ocp_quiz .

container: build
	podman build . -t ocp_quiz:latest

start: container
	podman run -e WEB_PORT="${WEB_PORT}" -e DB_CONN="${DB_CONN}" -p ${WEB_PORT}:${WEB_PORT} -d ocp_quiz:latest