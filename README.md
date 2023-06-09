# Client for ALT Linux Repository API

## Dependencies
altclient uses [go-version](http://github.com/hashicorp/go-version).
library for parsing and comparing package versions.

## Installation and Usage
Installation can be done with a normal `go get`:
```
$ go get github.com/aageorg/altclient
```

Type Branch represents a structure of the specified repository and contains unexported fields

Creation of a new branch object from the ALTRepo API. 
```go
// Returns *altclient.Branch. If cannot retrieve a list of packages, returns error
br1, err := altclient.NewBranch("p10")
br2, err := altclient.NewBranch("p9")
```
How to get a list of supporting architectures:
```go
br1, err := altclient.NewBranch("p10")
...
archs := br1.GetArchs() // Returns a []string slice

for _, a := range archs {
    fmt.Println(a)
}
```
Get a list of the packages, which are presenting in the first branch but missing from the second one
```go
br1, err := altclient.NewBranch("p10")
...
pkgs := br1.GetMissing(br2, "aarch64") // returns a slice []altclient.Package
```

Get a list of the packages from the second branch with the older versions than in the first one
```go
br1, err := altclient.NewBranch("p10")
...
pkgs := br1.GetOutOfDate(br2, "aarch64") // returns a slice []altclient.Package
```

## Issues and Contributing
If you find an issue with this library, please report an issue. If you'd like, we welcome any contributions. Fork this library and submit a pull request.