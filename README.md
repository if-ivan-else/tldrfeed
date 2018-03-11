# TLDRFEED

TL;DR;

`tldrfeed` is a simple JSON feed reader service.

tldrfeed was developed as part of a programming challenge / exercise provided by the Confluent Cloud team.

## Problem

We want to make a feed reader system. We will have 3 entities in the system: Users, Feeds, Articles. It should
support the following operations:

1.  Subscribe/Unsubscribe a User to a Feed
2.  Add Articles to a Feed
3.  Get all Feeds a Subscriber is following
4.  Get Articles from the set of Feeds a Subscriber is following

## Requirements

1.  Write a service with HTTP endpoints that allow items 1-4 from above
2.  It should handle multiple concurrent clients
3.  It should persist data across restarts
4.  Supply a README explaining your choices and how to run/test your service

### Clarifications

a) Service implementation assumes we are building a backend service and authN is provided by an upstream service.

b) Routes to create Users and Feeds should be provided for convenience of exercising the service functionality.

c) Feeds and Articles have the following required fields:

* `Feed`: `id`, `name`, list of `Articles`
* `Article`: `id`, `title`, `body`

d) Format of the API is JSON

e) There is currently no need for pagination

f) Sorting should be implemented in Feeds

* Assuming by most recent published time

g) There is no need for tracking User's read progress

h) There is no need for expiration of old Articles

i) HTTP is sufficient for exposed endpoints (no TLS)

k) No need to implement server-side push of updates (no websockets, etc.)

## Implementation

`tldrfeed` is built using Go, with the standard net/http library fullfilling most of the needs for API implementation.

Data persistence is built using MongoDB as the backend.

### Dependencies

Dependencies are managed by the dep vendoring manager.

Key dependencies in the project are:

* [gorilla/mux](https://github.com/gorilla/mux) - HTTP Routing
* [codegangsta/negroni](https://github.com/codegangsta/negroni) - Flexible HTTP middlewre (logging, etc.)
* [unrolled/render](https://github.com/unrolled/render) - HTTP/JSON rendering
* [globalsign/mgo](https://github/globalsign/mgo) - MongoDB for Go
* [spf13/cobra](https://github.com/spf13/cobra) - CLI parsing
* [spf13/viper](https://github.com/spf13/viper) - env var binding to flags
* [dghubble/sling](https://github.com/dghubble/sling) - simplified JSON REST client implementation
* [asaskevich/govalidator](https://github.com/asaskevich/govalidator) - annotation based JSON struct validation

### Package Layout

There are three top-level packages at the top of the project:

* `cmd` - CLI parsing and command implementation (can be used outside of this project)
* `api` - externally consumable APIs (REST API Client and related JSON api), ideally this would be consumed programmatically
* `internal` - internal package not meant to be used outside of the project, majority of the service implementation resides here

The internal package is broken up like so:

* `internal/buildinfo` - build version
* `internal/db` - DB/persistence interface and its implementations
* `internal/db/mock` - mock implementation of the db.Repository interface
* `internal/db/mongo` - MongoDB implementation of the db.Repository interface
* `internal/service` - implementation of the REST HTTP service, complete with routing and request validation

## Building and Testing

`tldrfeed` can built using the standard Go tools however there is a Makefile at the top of the project for convenience:

* `make vendor` - ensure vendored dependencies
* `make build` - build the project
* `make install` - build and install the service
* `make test` - run tests
* `make cloc` - run code stats
* `make todo` - list remaining TODOs
* `make lint` - lint for subtle and not so subtle coding errors

To get started:

```bash
make vendor install
```

The Makefile is a flavor of one found in [apex/up](https://github.com/apex/up)

## Running

`tldrfeed` needs a MongoDB database instance.

To run mongo using docker:

```bash
docker run -d --name mongo -p 27017:27017 mvertes/alpine-mongo
```

By default the service binds to port 8080 but that can be overriden by specifying the --port parameter.

To run the `tldrfeed` service (see build and install steps above):

```bash
tldrfeed server -d 0.0.0.0:27017
```

### Running in Docker

To build and run the service in docker:

```bash
make docker
docker run -p 8080:8080 --net host -e DB_URL="0.0.0.0:27017/db" tldrfeed
```

Note that in the above docker commane we are taking advantage of the host network in order to connect to MongoDB.

You can also run both the service and its DB via `docker-compose`:

```bash
docker-compose up -d
```

## Testing

To run the unit tests:

```bash
make test
```

Note that a test MongoDB instance is created as part of the test run.

Ultimately we should be testing at every level, however due to time limitations the primary areas of focus for test so far:

* Testing HTTP handlers, JSON validation and routing
* Testing persistence layer (db.Repository implementation)

The above areas apear to be the more complicated so they were tested first.

### API Testing using HTTPie

A couple of examples exercising the `tldrfeed` API using [httpie](https://httpie.org):

1.  Creating Users

```bash
http --json POST localhost:8080/api/v1/users name=boris
HTTP/1.1 201 Created
Content-Length: 60
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 21:52:07 GMT

{
    "id": "66a7854c-6657-4b85-9b0e-9b065a1b79d1",
    "name": "boris"
}
```

2.  Listing Users

```bash
http --json GET localhost:8080/api/v1/users
HTTP/1.1 200 OK
Content-Length: 62
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 21:56:54 GMT

[
    {
        "id": "66a7854c-6657-4b85-9b0e-9b065a1b79d1",
        "name": "boris"
    }
]
```

3.  Creating Feeds

```bash
http --json POST localhost:8080/api/v1/feeds name="boris' blog"
HTTP/1.1 201 Created
Content-Length: 66
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 21:55:30 GMT

{
    "id": "50b217e2-c5a2-44df-b6f2-c3e624557566",
    "name": "boris' blog"
}
```

4.  Listing Feeds

```bash
http --json GET localhost:8080/api/v1/feeds
HTTP/1.1 200 OK
Content-Length: 68
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 21:57:32 GMT

[
    {
        "id": "50b217e2-c5a2-44df-b6f2-c3e624557566",
        "name": "boris' blog"
    }
]
```

5.  Subscribing a User to a Feed

```bash
http --json POST localhost:8080/api/v1/users/66a7854c-6657-4b85-9b0e-9b065a1b79d1/feeds feed_id="50b217e2-c5a2-44df-b6f2-c3e624557566"
HTTP/1.1 202 Accepted
Content-Length: 114
Content-Type: text/plain; charset=UTF-8
Date: Sat, 03 Feb 2018 22:04:34 GMT

Successfully subscribed User '66a7854c-6657-4b85-9b0e-9b065a1b79d1' to Feed '50b217e2-c5a2-44df-b6f2-c3e624557566'
```

6.  Listing User's Feeds

```bash
http --json GET localhost:8080/api/v1/users/66a7854c-6657-4b85-9b0e-9b065a1b79d1/feeds
HTTP/1.1 200 OK
Content-Length: 68
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 22:19:43 GMT

[
    {
        "id": "50b217e2-c5a2-44df-b6f2-c3e624557566",
        "name": "boris' blog"
    }
]
```

7.  Publishing Articles to a Feed

```bash
http --json POST localhost:8080/api/v1/feeds/50b217e2-c5a2-44df-b6f2-c3e624557566/articles title="Morning News" body="Eating borsch with sauerkraut"
HTTP/1.1 201 Created
Content-Length: 45
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 22:22:27 GMT

{
    "id": "0500dd05-ab8f-4a48-ae97-e7dccd95cc3c"
}
http --json POST localhost:8080/api/v1/feeds/50b217e2-c5a2-44df-b6f2-c3e624557566/articles title="Evening News" body="Drinking kvas"
HTTP/1.1 201 Created
Content-Length: 45
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 22:23:15 GMT

{
    "id": "225b8bd3-1a81-45ad-ac1b-1535212a75db"
}
```

8.  Viewing Articles in a Feed

```bash
http --json GET localhost:8080/api/v1/feeds/50b217e2-c5a2-44df-b6f2-c3e624557566/articles
HTTP/1.1 200 OK
Content-Length: 285
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 22:24:08 GMT

[
    {
        "body": "Drinking kvas",
        "id": "225b8bd3-1a81-45ad-ac1b-1535212a75db",
        "published_at": "2018-02-03T22:23:15.715Z",
        "title": "Evening News"
    },
    {
        "body": "Eating borsch with sauerkraut",
        "id": "0500dd05-ab8f-4a48-ae97-e7dccd95cc3c",
        "published_at": "2018-02-03T22:22:27.175Z",
        "title": "Morning News"
    }
]
```

9.  Viewing Articles for a User in a subscribed Feed

```bash
http --json GET localhost:8080/api/v1/users/66a7854c-6657-4b85-9b0e-9b065a1b79d1/feeds/50b217e2-c5a2-44df-b6f2-c3e624557566/articles
HTTP/1.1 200 OK
Content-Length: 285
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 22:25:00 GMT

[
    {
        "body": "Drinking kvas",
        "id": "225b8bd3-1a81-45ad-ac1b-1535212a75db",
        "published_at": "2018-02-03T22:23:15.715Z",
        "title": "Evening News"
    },
    {
        "body": "Eating borsch with sauerkraut",
        "id": "0500dd05-ab8f-4a48-ae97-e7dccd95cc3c",
        "published_at": "2018-02-03T22:22:27.175Z",
        "title": "Morning News"
    }
]
```

9.  Viewing Articles for a User in all subscribed Feeds

```bash
http --json GET localhost:8080/api/v1/users/66a7854c-6657-4b85-9b0e-9b065a1b79d1/feeds
HTTP/1.1 200 OK
Content-Length: 143
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 22:29:46 GMT

[
    {
        "id": "50b217e2-c5a2-44df-b6f2-c3e624557566",
        "name": "boris' blog"
    },
    {
        "id": "a491159d-3607-42ee-b0ea-37ed3282aac9",
        "name": "Olga's High Fashion"
    }
]
http --json GET localhost:8080/api/v1/users/66a7854c-6657-4b85-9b0e-9b065a1b79d1/articles
HTTP/1.1 200 OK
Content-Length: 530
Content-Type: application/json; charset=UTF-8
Date: Sat, 03 Feb 2018 22:32:00 GMT

[
    {
        "body": "Christian Louboutin is a French fashion designer whose high-end stiletto footwear incorporates shiny, red-lacquered soles",
        "id": "6889a4f7-94bb-4fd9-8e71-759bf2f3a50a",
        "published_at": "2018-02-03T22:31:44.251Z",
        "title": "Louboutin Shoes"
    },
    {
        "body": "Drinking kvas",
        "id": "225b8bd3-1a81-45ad-ac1b-1535212a75db",
        "published_at": "2018-02-03T22:23:15.715Z",
        "title": "Evening News"
    },
    {
        "body": "Eating borsch with sauerkraut",
        "id": "0500dd05-ab8f-4a48-ae97-e7dccd95cc3c",
        "published_at": "2018-02-03T22:22:27.175Z",
        "title": "Morning News"
    }
]
```

### Testing using CLI

`tldrfeed` implements the following CLI commands in addition to the `server` command used to run the service:

```bash
tldrfeed help create
Create users, feeds or articles in tldrfeed

Usage:
  tldrfeed create [flags]
  tldrfeed create [command]

Available Commands:
  article     Create new article
  feed        Create new feed
  user        Create new user

Flags:
  -h, --help         help for create
  -u, --url string   tldrfeed service URL (default "http://localhost:8080")

Use "tldrfeed create [command] --help" for more information about a command.
```

Examples:

* Listing Articles for a User

```bash
tldrfeed list articles --user 66a7854c-6657-4b85-9b0e-9b065a1b79d1
2018/02/03 21:31:21 Articles for user 66a7854c-6657-4b85-9b0e-9b065a1b79d1 in all Feeds
{ID:6889a4f7-94bb-4fd9-8e71-759bf2f3a50a Title:Louboutin Shoes Body:Christian Louboutin is a French fashion designer whose high-end stiletto footwear incorporates shiny, red-lacquered soles PublishedTime:2018-02-03 22:31:44.251 +0000 UTC}
{ID:225b8bd3-1a81-45ad-ac1b-1535212a75db Title:Evening News Body:Drinking kvas PublishedTime:2018-02-03 22:23:15.715 +0000 UTC}
{ID:0500dd05-ab8f-4a48-ae97-e7dccd95cc3c Title:Morning News Body:Eating borsch with sauerkraut PublishedTime:2018-02-03 22:22:27.175 +0000 UTC}
```

## TODO

There are still quite a few unfinished things in this project:

* Ensuring uniqueness of user and feed names upon create
* More and better (configurable?) logging (sirupsen/logrus)
* More unit testing (internal/db/mongo and internal/service are the only packages that have coverage)
* E2E testing
* Profiling and perf testing
* More extensible data model
* ..
