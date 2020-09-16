# poker-planning

Simple poker planning web app with Go templating front-end

## Development

You need to generate a TLS certificate to use HTTPS. `generate_cert.go` can be used for development purposes.  
In the root of the project repository :     
```
$ mkdir tls
$ go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```
