ifneq (,$(wildcard .env))
    include .env
    export
endif

DB_CONTAINER = mariadb
SQL_FILE = dev/maria.sql


db-setup:
	docker-compose exec -T $(DB_CONTAINER) \
	mariadb -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) $(MYSQL_DATABASE) < $(SQL_FILE)

aws:
	dev/aws-setup.sh