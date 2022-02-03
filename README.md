# time-tracker-rest

Time Tracker Challenge REST implementation with Golang

## Build

`go build .`

## Install

`go install .`

But for you to execute it anywhere you need to add this to your `.bashrc`:
`export PATH=$PATH:$HOME/go/bin` or where your go installation is defined.

## Help

If you run any of the following commands you'll see the help

- `time-tracker-rest`
- `time-tracker-rest --help`

## Comands

- reset: Reset DB Schema
- run: Run REST Server

## Database

In order to run a docker or podman container with MySQL, you
need to run `source ./run.sh`

It will load a couple of shell commands/aliases/functions
to start, stop, check the logs and prompt into the MySQL shell inside the DB.

If you have another MySQL server running, you can provide this server with the following ENV variables:

- DB_TYPE: Database Type (mysql)
- DB_USES: Name of the DB to connecto to (db-test)
- DB_USER: Username (db-user)
- DB_PASS: Password (db-pass)
- DB_HOST: Hostname (localhost)
- DB_PORT: Server Port (3306)
- DB_ARGS: Arguments to pass to the connection (?charset=utf8mb4&parseTime=True&loc=Local)

## Test

There is a Postman Collection inside the project to test the Rest server.
You need to reset the DB before each run.

## Development History

- Created the project
- Copied some of the common package from another project. Just what is needed.
- Copied DB package (dbx) from another project I worked on before and made some modifications
  to exclude some features not needed in this project.
- Created the entities and queries
- Created the Rest Server
- Refactored the project to support Cobra CLI.
- Created the Postman Collection to test the Rest Server.

## Easiest

- Create the project
- Create the REST Server

## Hardest

- Create and Modify the dbx package (1hr)
- Create the Postman Collection (1hr)

## Total Time Spent

Between 3 and 4 hours. Had many distractions (work and family),
But adding all up, less than 4.
