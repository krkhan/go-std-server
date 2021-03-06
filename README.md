# go-std-server

Example of a golang HTTP server built entirely using the standard library

![no-deps](https://user-images.githubusercontent.com/697967/131789704-a2d6eb65-9d44-4d33-86d2-74c4e91e45f1.gif)

## Building

```
$ make help
```
```
Usage:
  make <target>

Targets:
  Build:
    build               Build your project and put the output binary in out/bin/
    run                 Execute server binary with default arguments
    clean               Remove build related file
    vendor              Copy of all packages needed to support builds and tests in the vendor directory
    watch               Run the code with cosmtrek/air to have automatic reload on changes
  Test:
    test                Run the tests of the project
    coverage            Run the tests of the project and export the coverage
  Lint:
    lint                Run all available linters
    lint-dockerfile     Lint your Dockerfile
    lint-go             Use golintci-lint on your project
    lint-yaml           Use yamllint on the yaml file of your projects
  Docker:
    docker-build        Use the dockerfile to build the container
    docker-release      Release the container with tag latest and version
  Help:
    help                Show this help.
```

## Testing
```
$ make test
```
```
go clean -testcache
go test -v -race ./...
=== RUN   TestHTTPServer
=== RUN   TestHTTPServer/stats_request_without_any_traffic
    main_test.go:40: Making GET call to http://127.0.0.1:40162/stats
2021/09/07 14:58:15 Handled GET request from 127.0.0.1:41512 for URL '/stats', response: OK
=== RUN   TestHTTPServer/valid_hash_request
    main_test.go:46: Making POST call to http://127.0.0.1:40162/hash with 20 bytes in body
2021/09/07 14:58:15 Handled POST request from 127.0.0.1:41512 for URL '/hash', response: OK
=== RUN   TestHTTPServer/multiple_passwords
    main_test.go:46: Making POST call to http://127.0.0.1:40162/hash with 44 bytes in body
2021/09/07 14:58:15 Multiple (2) passwords provided
2021/09/07 14:58:15 Handled POST request from 127.0.0.1:41512 for URL '/hash', response: Bad Request
--- PASS: TestHTTPServer (3.02s)
    --- PASS: TestHTTPServer/stats_request_without_any_traffic (0.01s)
    --- PASS: TestHTTPServer/valid_hash_request (0.00s)
    --- PASS: TestHTTPServer/multiple_passwords (0.00s)
=== RUN   TestDelay
2021/09/07 14:58:15 Handled POST request from 127.0.0.1:54428 for URL '/hash', response: OK
2021/09/07 14:58:15 Key not found: 1
2021/09/07 14:58:15 Handled GET request from 127.0.0.1:54428 for URL '/hash/1', response: Bad Request
2021/09/07 14:58:20 Handled GET request from 127.0.0.1:54428 for URL '/hash/1', response: OK
--- PASS: TestDelay (5.02s)
PASS
ok      github.com/krkhan/go-std-server 8.080s
=== RUN   TestRouter
=== RUN   TestRouter/get_request_without_any_parameters
2021/09/07 14:58:12 Handled GET request from  for URL '/', response: OK
=== RUN   TestRouter/post_request_without_any_parameters
2021/09/07 14:58:12 Handled POST request from  for URL '/profilePicture', response: OK
=== RUN   TestRouter/get_request_with_a_numeric_parameter
2021/09/07 14:58:12 Handled GET request from  for URL '/profile/12345', response: OK
--- PASS: TestRouter (0.00s)
    --- PASS: TestRouter/get_request_without_any_parameters (0.00s)
    --- PASS: TestRouter/post_request_without_any_parameters (0.00s)
    --- PASS: TestRouter/get_request_with_a_numeric_parameter (0.00s)
PASS
ok      github.com/krkhan/go-std-server/router  0.029s
=== RUN   TestSha512DigestStore
--- PASS: TestSha512DigestStore (10.03s)
PASS
ok      github.com/krkhan/go-std-server/store   10.056s
```

You can get a coverage report via:

```
$ make coverage
```

```
go test -cover -covermode=count -coverprofile=profile.cov ./...
ok      github.com/krkhan/go-std-server 8.539s  coverage: 57.3% of statements
ok      github.com/krkhan/go-std-server/router  0.013s  coverage: 71.0% of statements
ok      github.com/krkhan/go-std-server/store   10.009s coverage: 100.0% of statements
go tool cover -func profile.cov
github.com/krkhan/go-std-server/main.go:26:             handleError                     100.0%
github.com/krkhan/go-std-server/main.go:32:             postHash                        66.7%
github.com/krkhan/go-std-server/main.go:60:             getHash                         78.6%
github.com/krkhan/go-std-server/main.go:81:             getStats                        83.3%
github.com/krkhan/go-std-server/main.go:106:            ServeHTTP                       100.0%
github.com/krkhan/go-std-server/main.go:112:            startHttpServer                 72.7%
github.com/krkhan/go-std-server/main.go:148:            main                            0.0%
github.com/krkhan/go-std-server/router/router.go:16:    NewRoute                        100.0%
github.com/krkhan/go-std-server/router/router.go:49:    NewLoggingResponseWriter        100.0%
github.com/krkhan/go-std-server/router/router.go:53:    WriteHeader                     0.0%
github.com/krkhan/go-std-server/router/router.go:58:    Serve                           72.0%
github.com/krkhan/go-std-server/router/router.go:93:    GetParam                        100.0%
github.com/krkhan/go-std-server/store/store.go:19:      AddDigest                       100.0%
github.com/krkhan/go-std-server/store/store.go:31:      GetDigest                       100.0%
total:                                                  (statements)                    65.1%
```

## Launching the server

```
$ make run
```

Or, if you want to launch the executable manually:

```
$ go-std-server [listen address (default ":8080")]
```
```
2021/09/07 15:01:39 Launching HTTP server on :8080
2021/09/07 15:01:46 Handled POST request from 127.0.0.1:58096 for URL '/hash', response: OK
2021/09/07 15:01:47 Handled POST request from 127.0.0.1:58098 for URL '/hash', response: OK
2021/09/07 15:01:50 Key not found: 1
2021/09/07 15:01:50 Handled GET request from 127.0.0.1:58100 for URL '/hash/1', response: Bad Request
2021/09/07 15:01:51 Key not found: 2
2021/09/07 15:01:51 Handled GET request from 127.0.0.1:58102 for URL '/hash/2', response: Bad Request
2021/09/07 15:01:57 Handled GET request from 127.0.0.1:58104 for URL '/hash/1', response: OK
2021/09/07 15:01:58 Handled GET request from 127.0.0.1:58106 for URL '/hash/2', response: OK
^C2021/09/07 15:02:00 Received signal 'interrupt', shutting down HTTP server
2021/09/07 15:02:00 HTTP server terminated successfully
```

## Client Examples

### Request/queue a new password digest
```
$ curl --data password=angryMonkey http://127.0.0.1:8080/hash
```
```
1
```

### Get digest value
```
$ curl http://127.0.0.1:8080/hash/1
```

You may get:
```
Key not found: 1
```

Or, if the digest has been processed (default time: 5 seconds) you'll get its value:
```
ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==
```

Another example:
```
$ curl --data password=angryMonkey2 http://127.0.0.1:8080/hash
```
```
2
```

Sleep with one eye open, gripping your pillow tight:
```
$ sleep 5
```

Exit light:
```
$ curl http://127.0.0.1:8080/hash/2
```
```
k9w4qLUoXzJEUp6TXTL59Vhjq1g600F0Va9v/VLkNeegC7Oro7kh/AIMU20+RlnG4fBDdfmv9qY4NHc5rF7YTw==
```

### Get statistics for the `POST /hash` endpoint
```
curl http://127.0.0.1:8080/stats
```
```
{"total": 4, "average": 113}
```
Please note that since the hashes are calculated asynchronously, the average is just for the time servicing the HTTP request that triggered the SHA512 calculation. In other words, the hash performance itself is not part of the stats. (Not sure about requirements here, it can easily be tweaked to measure hash performance instead.)

### Shutdown via GET request
```
curl --verbose http://127.0.0.1:8080/shutdown
```
```
*   Trying 127.0.0.1:8080...
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET /shutdown HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.74.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 02 Sep 2021 06:46:50 GMT
< Content-Length: 0
<
* Connection #0 to host 127.0.0.1 left intact
```

While on the server side:
```
2021/09/01 23:46:48 Launching HTTP server on :8080
2021/09/01 23:46:50 Received shutdown request, terminating self
2021/09/01 23:46:50 HTTP server terminated successfully
```

## Todo

* Add security around digests (return digests to only those who queued them)
* Do not return a key unless the digest has been processed and stored; while making sure the relevant digest isn't returned unless required delay has elapsed
* Better error reporting
  * Return something along the lines of "digest being calculated" when it's queued, instead of just saying "key not found" for all cases
  * Use appropriate HTTP status codes (instead of blanket-returning 400 in case of trouble)

