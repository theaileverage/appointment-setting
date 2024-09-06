APIs
Encore lets you define APIs as regular Go functions with an annotation. API calls are made as regular function calls. Read the docs to learn more.
Defining APIs
import "context"

// PingParams is the request data for the Ping endpoint.
type PingParams struct {
    Name string
}

// PingResponse is the response data for the Ping endpoint.
type PingResponse struct {
    Message string
}

//encore:api public method=POST path=/ping
func Ping(ctx context.Context, params *PingParams) (*PingResponse, error) {
    msg := "Hello, " + params.Name
    return &PingResponse{Message: msg}, nil
}

Calling APIs
import "encore.app/hello" // import service

//encore:api public
func MyOtherAPI(ctx context.Context) error {
    resp, err := hello.Ping(ctx, &hello.PingParams{Name: "World"})
    if err == nil {
        log.Println(resp.Message) // "Hello, World!"
    }
    return err
}

Hint: Import the service package and call the API endpoint using a regular function call.

Raw endpoints
import "net/http"

// A raw endpoint operates on standard HTTP requests.
// It's great for things like Webhooks, WebSockets, and GraphQL.
//encore:api public raw
func Webhook(w http.ResponseWriter, req *http.Request) {
    // ... operate on the raw HTTP request ...
}

GraphQL
Encore supports GraphQL servers through raw endpoints. We recommend using [gqlgen](https://gqlgen.com/).

An example of using GraphQL with Encore can be found [here](https://github.com/encoredev/examples/tree/main/graphql).