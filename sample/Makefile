.PHONY: db-start
db-start:
	docker run --rm -d \
	    --name sqltest-db \
	    --network host \
	    -e POSTGRES_DB=test \
	    postgres:9.5

.PHONY: db-stop
db-stop:
	docker stop \
	    sqltest-db

.PHONY: sql
sql:
	docker exec -ti \
	    sqltest-db \
	    psql -h localhost -p 5432 test postgres

