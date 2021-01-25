# Tracker

---

## About
_Coming soon_

---

## How to run

### Running the server
You can launch the server by running the following command from the root of the project.

```shell
go run cmd/server/server.go
```

A default host.conf file will be created after the first launch. This contains details about the host of the server which can be modified as required.

### Running the scraper
Before running the scraper you should setup your databases using the schema located at `schema/schema.sql`
You can launch the scraper by running the following command from the root of the project.

```shell
go run cmd/scraper/scraper.go
```

Note: without adding entries to the `tracker/shows` table, the crawler will have nothing to do.

---
