# Secure Request

Secure request is a request signature method based on (but not fully compliant with) OAuth 1.0 (http://oauth.net/core/1.0a/#signing_process).

# Configuration

Both sides, client and server must know the applications which will communicate using this signature. Both main classes constructors receives a `app_secrets` parameter which contains a mapping of application key (or name) and its secret (salt for the hashing algorithm).

# Usage

## Client-side

You can Sign requests yourself.

```go
err := securerequest.Sign(request)
``` 

or you can you use SigningRoundTripper with your http.Client

```go
// Take your existing http.Client and change it's Transport
client.Transport = securerequest.NewSigningRoundTripper(nil)
// it now automatically signs requests
client.Do(req)
``` 

# Server-side

Validate requests manually

```go
valid, err := securerequest.Validate(request)
```

or use ValidatingHandler as a middleware for your existing handler

```go
http.Handle("/foo", fooHandler)
// becomes
http.Handle("/foo", securerequest.NewValidatingHandler(fooHandler))
``` 

If the authorization fails, it will respond with 403 status and a plain text body: "Unauthorized access.".
