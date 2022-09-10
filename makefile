.SILENT: clear build start container

clear:
	printf "Cleaning existing builds.\n"
	rm -rf build >> /dev/null

build: clear
	printf "Fetching dependencies...\n"
	go mod vendor >> /dev/null
	printf "Building package...\n"
	CGO_ENABLED=0 go build -o ./build/ocp_quiz . >> /dev/null

container: build
	printf "Building container...\n"
	podman build . -t ocp_quiz:latest >> /dev/null

start: container
	printf "Starting container... \n"
	podman rm -f ocp_quiz || echo "Nothing to clean..."
	podman run -e WEB_PORT="${WEB_PORT}" -e DB_CONN="${DB_CONN}" -p ${WEB_PORT}:${WEB_PORT} -d --name ocp_quiz ocp_quiz:latest >> /dev/null
	printf "Started..."