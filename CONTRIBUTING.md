# Contributing Guidelines


## install for contributors

make sure you have the following in your .bashrc or .bash_profile 

```bash
# Adds the main 'go' compiler command to your path (probably already set)
export PATH="$PATH:/usr/local/go/bin"

# export the GOPATH variable explicitly
export GOPATH="$HOME/go"

# Adds user installed binaries 'go install' tools are in your path
export PATH="$PATH:$GOPATH/bin"
```
docker
  - why: To run the api and database in containers. Install docker desktop.
  - how: [typical docker install steps](https://www.docker.com/get-started/)

sqlc
  - why: to generate data access code for the database
  - how: 

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

mockgen
  - why: to generate mocks of the Querier interface for mocking data calls in http handler unit tests
  - how: 

```bash  
go install go.uber.org/mock/mockgen@latest
```

migrate
  - why: to control the version and state of the database in use
      Dual Files: Generates distinct .up.sql (to apply changes) and .down.sql (to roll back changes) schema files.
      Safety Tracking: Creates a tracking table (schema_migrations) in your database to log the exact current version and flag failed ("dirty") states.
  - how: 

```bash  
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## environment variables

Create a local .env file in the root directory and specify the values below. 
The .env file is excluded from the project via .gitignore file.
    
```env
# api database connection string
DATABASE_URL=postgresql://firstly:firstly@db:5432/firstly?sslmode=disable

# used for jwt signing. generate using command -> openssl rand -hex 32
SECRET=
# set the domain field of access/refresh token cookies
COOKIE_DOMAIN=localhost

# used for configuring gin router cors middleware
ALLOW_ORIGINS="http://localhost:3000"

# variables for db service
POSTGRES_USER=firstly
POSTGRES_PASSWORD=firstly
POSTGRES_DB=firstly

ACCESS_TOKEN_DURATION="15s"
REFRESH_TOKEN_DURATION="2m"

```

## generate

Go module housekeeping is performed. Generates the data access code from the 
db project using sqlc and generates the mocks used in tests for the data access code using 
mockgen.

```bash
$ make generate
```

To generate new dml sql scripts use golang-migrate/migrate tool. This will add files to the
db/migration directory following the migrate tools naming convention. 

```bash
$ migrate create -ext sql -dir db/migration create_example_table
```

## test

Run the unit tests for the entire application.

```bash
$ make test
```

## test coverage

Run the unit tests for the entire application and generate a coverage report.
```bash
$ make test-cover
```

## verify

Quickly run both generating code and test recipes
```bash
$ make verify
```

## build

Build the docker containers with the appropriate files and environment 
configurations. 

```bash
$ docker compose build
```

## run local

Run local instances of the docker compose containers

```bash
$ docker compose up
```

## cleanup local

Remove all networks and containers and start fresh.
```bash
$ docker compose down
```

## deploy

```bash
# TODO: choose another container orchestration platform and define deploy step
$ make deploy
```

