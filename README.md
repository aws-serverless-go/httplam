# httplam
httplam is a library that bind [AWS Lambda](https://aws.amazon.com/lambda/), [AWS API Gateway](https://aws.amazon.com/api-gateway/) to `net/http`

# Install

```shell
go get -u github.com/aws-serverless-go/httplam
```

# Examples
## `net/http`

```go
package main

import (
	"github.com/aws-serverless-go/httplam"
	"net/http"
)

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("welcome"))
	})
	if httplam.IsLambdaRuntime() {
		httplam.StartLambdaWithAPIGateway(m)
	} else {
		http.ListenAndServe(":3000", m)
	}
}
```

## `github.com/julienschmidt/httprouter`

```go
package main

import (
	"fmt"
	"github.com/aws-serverless-go/httplam"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	if httplam.IsLambdaRuntime() {
		httplam.StartLambdaWithAPIGateway(router)
	} else {
		log.Fatal(http.ListenAndServe(":8080", router))
	}
}
```

## `github.com/go-chi/chi`

```go
package main

import (
	"github.com/aws-serverless-go/httplam"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	if httplam.IsLambdaRuntime() {
		httplam.StartLambdaWithAPIGateway(r)
	} else {
		http.ListenAndServe(":3000", r)
	}
}
```

# License
[MIT License](LICENSE)