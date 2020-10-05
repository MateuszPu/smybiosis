docker run --name pg-docker -e POSTGRES_PASSWORD=pass -e POSTGRES_USER=user -d -p 5432:5432 postgres
docker run --network="host" --rm -v /home/mateusz/Desktop/Projekty/go/payme/migration:/liquibase/changelog liquibase/liquibase --url="jdbc:postgresql://localhost:5432/postgres?currentSchema=public" --changeLogFile=/changelog/master.xml --username="user" --password="pass" --logLevel=debug update


download sqlboiler
run: sqlboiler psql