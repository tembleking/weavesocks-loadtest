# Weavesocks LoadTest

This is a load test for the [Weaveworks' Sock Shop Microservices demo](https://microservices-demo.github.io/).

We were trying to use their [load test](https://microservices-demo.github.io/docs/load-test.html) but since it was
not working, we rewrote the same behaviour in Golang.

# Usage

You can download it with `go get -v github.com/tembleking/weavesocks-loadtest` and execute manually or with
Docker Compose via `docker-compose run --rm weavesocks-loadtest`.

If you download it manually you will require Go 1.11 at least.

Example:

```
$  docker-compose run --rm weavesocks-loadtest -h

Create some fake load in the weavesocks demo application available at https://microservices-demo.github.io

Usage:
weavesocks-loadtest [flags]

Flags:
-c, --clients int       Number of concurrent clients (default 2)
-d, --delay int         Delay before start
-h, --help              help for weavesocks-loadtest
-n, --hostname string   Target host url (eg. http://localhost:8080)
-r, --requests int      Number of requests per client (default 10)
```

# License

GNU GPLv3