# Can(Can)

This package is a simple implementation of a can/cannot method for authorization.

```go
if can.Can(someAccount, "update", someResource) {
  // do something...
} else {
  // show an error maybe?
}

```

Or, if you only want to see if something is not allowed:

```go
if can.Not(someAccount, "update", someResource) {
  // do something...
}

```
