## Currency Converter

REST API for converting currency using the [Currency Layer](https://currencylayer.com) service.

### Running
There are three environment variables used by this application

* `PORT` - The port the application will expose for requests
* `CURRENCY_HOST` - The service host used by the application (i.e https://apilayer.net)
* `CURRENCY_ACCESSKEY` - The access key used in the service requests

To run the application:

```
$ go run main.go
```

### Testing

To test the application and check the coverage per package (aka directory) using the `go test`:

```
$ go test -cover ./...
```

To generate a profile from test, so you can check coverage per function or even create an `.html` file to see which lines were covered by the tests:

``` 
$ go test -cver ./... -coverprofile=cover.out
```

Checking coverage by function

```
$ go tool cover -func=cover.out
```

Generating `.html` of the coverage
```
$ go tool cover -html=cover.out -o coverage.html
```

Now the `coverage.html` file can be open in the browser.
