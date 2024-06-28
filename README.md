# Event Management System

user management- login / sign up / roles assigned to users  ( organisers , admins and normal users) 

- Event organisers: Should be able to create , edit , update , delete event. 
- Event registration: Normal users can register for an event and since we have limited capacity for events (offline events) this will be checked. 
- Event notification: Once the user registration is complete they should get notified about those events 
- Administration: Admins should be able to deal with users , moderate events etc , see analytics about the events such as popularity/ views etc .
- Event reviews: Users can leave reviews for events they attended etc.

## Setup

```bash
git clone https://github.com/BEOpenSourceCollabs/EventManagementCore
cd EventManagementCore
```

### Local database setup
- You need [Docker](https://www.docker.com/products/docker-desktop/) installed and working.
- Create a `db.env` file in project root with following variables

    ```text
    # Pgadmin credentials
    PGADMIN_DEFAULT_EMAIL=admin@admin.com
    PGADMIN_DEFAULT_PASSWORD=12345

    # Postgresql credentials
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    ```
  you can change these values if you want. If you change postgresql credentials, make sure to change them specially the `OWNER` in `schema.sql`

- Run following command, 
    ```shell
    $ docker compose up
    ```
- Once the containers are running , postgresql will be accessible on `localhost:5432` and Pgadmin will be accessible on `localhost:8888`. 

Next setup the application environment and add the database configuration.

- Create a `.env` file in the project root with the following variables:

  ```text
  DATABASE_HOST=127.0.0.1
  DATABASE_USER=postgres
  DATABASE_PASSWORD=postgres
  DATABASE_NAME=event-mgmt-db
  DATABASE_PORT=5432
  DATABASE_SSL_MODE=disable
  ACCESS_TOKEN_SECRET=test
  REFRESH_TOKEN_SECRET=test
  ```
  Ensure to update these to match your database configuration (these are set in `db.env` for development).

### Run

```shell
make run
```

### Run Tests

```shell
make test

# or 

make test.unit
```

## Development

### Running the application for development

Best practice for running the application is as follows.  
After completing prerequisit steps to setup and configure the application the docker compose can be used to run the development environment.

```shell
docker compose watch
```

> Using `watch` will automatically be syncronized to the running container and restart the application when source is modified.

### Advanced

This section examples the usage of the `docker-dev.sh` script, which builds only the event-mgmt-core image and runs it in a container. You will have to setup and configure the database before running.

  ```shell
  ./docker-dev.sh
  ```

If you want to run the database from `docker-compose.yml` and use that you can do the following:

  ```shell
  # start only db
  docker compose up db

  # optionally you could start both the db and pgadmin
  docker compose up db pgadmin
  ```

  Example: `.env`
  
  ```text
  DATABASE_HOST=event-mgmt-postgres
  DATABASE_USER=postgres
  DATABASE_PASSWORD=postgres
  DATABASE_NAME=event-mgmt-db
  DATABASE_PORT=5432
  DATABASE_SSL_MODE=disable
  ACCESS_TOKEN_SECRET=test
  REFRESH_TOKEN_SECRET=test
  ```
  > Note: the container name 'event-mgmt-postgres' can be used as the host.

Then run the application.

  ```shell
  ./docker-dev.sh
  ```

## Database Migrations

### CLI
To work with database migrations,  [golang-migrate cli](https://pkg.go.dev/github.com/golang-migrate/migrate/v4/cmd/migrate#section-readme) needs to be installed and working. 


### Working with migrations 

The `init.sql` file is responsible for creating an empty database. This file is automatically executed by `docker-compose` during the first run. After the initial setup, database migrations handle subsequent changes to the database schema.

Each migration is versioned and follows a sequential numbering format, starting with `000001`.

Each migration consists of two files - 
* `<SEQUENCE_NUMBER>_<MIGRATION_NAME>.up.sql` : This file contains the SQL statements to apply the changes (migration).
* `<SEQUENCE_NUMBER>_<MIGRATION_NAME>.down.sql` : This file contains the SQL statements to undo the changes (rollback).

To apply all the migrations from first, all the way to the latest

  ```bash
    $ make migrate.all.up
  ```
To view current version of database (current version means current migration # it is at), 
  
  ```bash
  $ make migrate.version
  ```
To create a new migration, 

  ```bash
  $ make migrate.create NAME=name_of_your_migration
  ```
To goto a specific version of database , 
  
  ```bash
  $ make migrate.goto V=version_you_want_to_goto
  ```
   Example: `make migrate.goto V=4`
   - If current version is 2 it will apply #3 and #4 to bring database to version 4.
   - If current version is 6 it will remove #6 and #5 to bring database to version 4.

To undo all the migrations and go to initial state, 

  ```bash
  $ make migrate.all.down
  ```

## Documentation 

Open API documentation with interactive client is available on route [http://localhost:8081/swagger](http://localhost:8081/swagger) and 
the raw json schema can be downloaded on [http://localhost:8081/swagger.json](http://localhost:8081/swagger.json).
