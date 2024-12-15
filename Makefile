.PHONY: refresh-db-test

refresh-db:
	docker rm -f mysql-employees-system || true
	docker run --name mysql-employees-system \
		-e MYSQL_ROOT_PASSWORD=root \
		-e MYSQL_DATABASE=employees-system \
		-d \
		-v ./db/init-scripts:/docker-entrypoint-initdb.d \
		-p 3398:3306 \
		mysql:8 \
		--default-time-zone='-03:00'
		@echo "Database refreshed"

refresh-db-test:
	docker rm -f mysql-employees-system-test || true
	docker run --name mysql-employees-system-test \
		-e MYSQL_ROOT_PASSWORD=test \
		-e MYSQL_DATABASE=employees-system \
		-d \
		-v ./db/init-scripts:/docker-entrypoint-initdb.d \
		-p 3377:3306 \
		mysql:8 \
		--default-time-zone='-03:00'
		@echo "Database test refreshed"
