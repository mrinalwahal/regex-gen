## Random String Generator Using Regex

Submission for Flanksource job assignment.

### Clone

```
git clone git@github.com:mrinalwahal/regex-gen.git
```

### Build

```
go build -o app
```

### Usage

```
./app "foo(-(boo|bar|woo|quack))"
```

### Characters Not Supported

- Word Boundary
- Word Non-Boundary
- Range Pair List (only partially supported)