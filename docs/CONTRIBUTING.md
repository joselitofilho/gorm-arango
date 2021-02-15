# Contributing

By participating to this project, you agree to abide our [code of
conduct](CODE_OF_CONDUCT.md).

## Setup your machine

`gorm-arango` is written in [Go](https://golang.org/).

Prerequisites:

- [Go 1.15+](https://golang.org/doc/install)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

Clone `gorm-arango` anywhere:

```sh
git clone git@github.com:joselitofilho/gorm-arango.git
```

Create a docker container:

```sh
docker-compose build  # only first time
docker-compose up
```

Install the build and lint dependencies:
```sh
go get -v ./...
```

A good way of making sure everything is all right is running the test suite:

```sh
go test -v -count=1 ./...
```

Check the coverage of code:
```sh
go test -coverpkg=./... -coverprofile="coverage.out" -v -count=1 ./...
go tool cover -html="coverage.out" -o="coverage.html"
```

## Runtime configuration

| ENV key                        | Default value            | Description                                                         |
| ------------------------------ | ------------------------ | ------------------------------------------------------------------- |
| LOGRUS_LEVEL                   | "info"                   | Logging level (trace, debug, info, etc).                            |
| ARANGODB_URI                   | "http://localhost:8529"  | URI where to find the ArangoDB server (including protocol and port) |
| ARANGODB_DATABASE              | "test"                   | ArangoDB database                                                   |
| ARANGODB_USER                  | "user"                   | ArangoDB user                                                       |
| ARANGODB_PASSWORD              | "password"               | ArangoDB user password                                              |

## Create a commit

Commit messages should be well formatted, and to make that "standardized", we
are using Conventional Commits.

You can follow the documentation on
[their website](https://www.conventionalcommits.org).

## Submit a pull request

Push your branch to `gorm-arango` repository and open a pull request against the
main branch.

## Credits

### Contributors

Thank you to all the people who have already contributed to gorm-arango!

[<img src="https://avatars.githubusercontent.com/u/1815812?s=64&v=4" width="32" height="32"/>](https://github.com/ricardogpsf)
[<img src="https://avatars.githubusercontent.com/u/64505737?s=64&v=4" width="32" height="32"/>](https://github.com/LucasSaraiva019)