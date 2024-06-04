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

## Local database setup
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

```bash
make run
```

### Run Tests

```bash
make test

# or 

make test.unit
```

test