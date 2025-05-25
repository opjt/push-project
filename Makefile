ifneq (,$(wildcard .env))
    include .env
    export
endif

DB_CONTAINER = mariadb
SQL_FILE = dev/maria.sql


db-setup:
	docker-compose exec -T $(DB_CONTAINER) \
	mariadb -u$(MARIA_USER) -p$(MARIA_PASSWORD) $(MARIA_DATABASE) < $(SQL_FILE)

aws:
	dev/aws-setup.sh