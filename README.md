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
  ```
  > Note: the container name 'event-mgmt-postgres' can be used as the host.

Then run the application.

  ```shell
  ./docker-dev.sh
  ```

## Documentation 

Open API documentation with interactive client is available on route [http://localhost:8081/swagger](http://localhost:8081/swagger) and 
the raw json schema can be downloaded on [http://localhost:8081/swagger.json](http://localhost:8081/swagger.json).
