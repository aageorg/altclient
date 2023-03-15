# Client for ALT Linux Repository API

## Dependencies
altclient uses [go-version](http://github.com/hashicorp/go-version).
library for parsing and comparing package versions.

## Installation and Usage
Installation can be done with a normal `go get`:
```
$ go get github.com/aageorg/altclient
```

The only useful type is Branch. It represents structure of specified repository and contains unexported fields

####func NewBranch(branch_name string) *Branch, error
Returns a new branch object from the ALTRepo API. If cannot retrieve a list of packages, returns error

```go
br1, err := altclient.NewBranch("p10")
br2, err := altclient.NewBranch("p9")
```

####func (*Branch) GetArchs() []string
Returns a list of supported architectures
```go
br1, err := altclient.NewBranch("p10")
...
archs := br1.GetArchs()
for _, a := range archs {
    fmt.Println(a)
}
```

####func (*Branch) GetMissing(br *Branch, arch string) []Package
Returns a list of the packages, which are presenting in the first branch but missing from the second one
```go
br1, err := altclient.NewBranch("p10")
...
pkgs := br1.GetMissing(br2, "aarch64")
```

####func (*Branch) GetOutOfDate(br *Branch, arch string) []Package
Returns a list of packages from the second branch with the older versions than in the first one
```go
br1, err := altclient.NewBranch("p10")
...
pkgs := br1.GetOutOfDate(br2, "aarch64") 
```

## Issues and Contributing
If you find an issue with this library, please report an issue. If you'd like, we welcome any contributions. Fork this library and submit a pull request.