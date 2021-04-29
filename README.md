Sample Go REST API
==================
Facilitating `alert` and `tag` creation/queries.

Build status (master): [![Build Status](https://travis-ci.org/konrads/go-tagged-articles.svg?branch=master)](https://travis-ci.org/konrads/go-tagged-articles)

Specification
-------------
Implement GET/POST endpoints for persistence/retrieval of tagged articles such as:
```
{
  "id": "1",
  "title": "latest science shows that potato chips are better for you than sugar",
  "date" : "2016-09-22",
  "body" : "some text, potentially containing simple markup about how potato chips are great",
  "tags" : ["health", "fitness", "science"]
}
```
Implement GET end point for tag specific information:
```
{
  "tag" : "health",
  "count" : 17,
    "articles" :
      [
        "1",
        "7"
      ],
    "related_tags" :
      [
        "science",
        "fitness"
      ]
}
```

Assumptions
-----------
From the tag specific information, the following has been assumed:
* count = count of all tags of associated articles, for given selection criteria (tag and date)
* articles = list of article_ids for given selection criteria
* related_tags = unique set of tags for the given selection criteria, excluding the current tag

Design
------
The API contains a single API endpoint, [restapi.go](cmd/restapi/restapi.go), which serves REST API via [handler.go](pkg/handler/handler.go). This connects up to [postgres.go](pkg/db/postges.go) db. [postgres.go](pkg/db/postges.go) implements an interface, for the purpose of decoupling and ease of testing.

Database mode consists of a single table, as Postgres supports queries within embedded arrays.

Rest API sits on top of the [gin](https://github.com/gin-gonic/gin) framework, and persistence was done with [pq](https://github.com/lib/pq) postgres library.

Restapi endpoint results in typical HTTP errors, such as 404/400/500.

Docker (and docker-compose) were utilized for end-to-end testing.

Unit test
---------
Tests validating Go codebase, with the use of mocks. NOTE: needs more error testing.
```
go test ./... -v
```

End-to-end test
---------------
Start up postgres and restapi via `docker-compose`, run tests via `curl`, validate with `jq`. NOTE: needs more error testing.
This test aids in overall interactions, from REST calls up to db persistence.
```
# ensure you have installed: docker-compose, jq, curl
end2end_test.sh
```

TODOs
-----
* add failure scenarios, in unit and end-to-end testing
* add end-to-end tests to CICD
* should the landscape get more complex, consider splitting up the db into its own gRPC fronted microservice
* for production readiness consider security, rate limiting, telemetry, timeouts
