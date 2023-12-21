# snippetbox

A web application which lets people paste and share snippets of text. A bit like Pastebin or GitHubs Gists.
This application was developed by following the [Let's Go book by Alex Edwards](https://lets-go.alexedwards.net/).

## What I learned

This was my very first web application and I learned a ton about backend-development by making this project.
The most important takeaways for me was the following:

- REST-principles.
- SQL database management.
- User authentication.
- Protection against common vulnerabilites.
- Password hashing.
- Session management.
- How to setup a HTTPS server.

## Commands

### To download dependencies:

```console
foo@bar:~$ go mod download
```

### To verify dependencies:

```console
foo@bar:~$ go mod verify
```

### To generate TLS key and cert:

```console
foo@bar:~$ go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

```

### To run the server:

```console
foo@bar:~$ go run ./cmd/web
```

### To list available command-line flags:

```console
foo@bar:~$ go run ./cmd/web -help
```
