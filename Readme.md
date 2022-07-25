# Database Schema
![db schema](assets/db%20schema.png)

# Go Image
[Postgres Official Docker Image](https://hub.docker.com/_/postgres)

`docker pull postgres:12-alpine`

`docker run --name some-postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:12-alpine`

to see your running containers: 
`docker ps`


# Go-Lang Migrate
[Migrate repo](https://github.com/golang-migrate/migrate)

To install the CLI tool:

`brew install golang-migrate`

`mkdir -p db/migration`
`migrate create -ext sql -dir db/migration -seq init_schema`

2 files are generated as:

1. Upstream file is used to do changes to the schema. (1->2->3)
2. Downstream is used to revert the changes done to the schema by Upstream.(1<-2<-3)


