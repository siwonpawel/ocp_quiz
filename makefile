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
	docker build . -t ocp_quiz:latest >> /dev/null

start: container
	printf "Starting container... \n"
	docker rm -f ocp_quiz || echo "Nothing to clean..."
	docker run -e WEB_PORT="${WEB_PORT}" -e DB_CONN="${DB_CONN}" -p ${WEB_PORT}:${WEB_PORT} -d --name ocp_quiz ocp_quiz:latest >> /dev/null
	printf "Started...\n"