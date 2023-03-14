package altclient

import (
	"net/http"
	"encoding/json"
	"bytes"
)

const ApiURL = "https://rdb.altlinux.org/api/export/branch_binary_packages/"

type Package struct {
	Name      string `json:"name"`
	Epoch     int    `json:"epoch"`
	Version   string `json:"version"`
	Release   string `json:"release"`
	Arch      string `json:"arch"`
	Disttag   string `json:"disttag"`
	Buildtime int64  `json:"buildtime"`
	Source    string `json:"source"`
}

type Branch struct {
	Length  int `json:"length"`
	Arch    map[string]int
	Package []map[string]*Package
}


func NewBranch(br string) (*Branch, error) {
	resp, err := http.Get(ApiURL + br)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	br = Branch{Arch: make(map[string]int)}
	dec := json.NewDecoder(bytes.NewReader(resp.Body))

	for t, err := dec.Token(); t != "packages"; {
		if err != nil {
			return nil, err
		}
		if t == "length" {
			l, err := dec.Token()
			if err != nil {
				return nil, err
			}
			br.Length = int(l)
		}
	}
	dec.Token()
	for dec.More() {
		var pkg Package
		err = dec.Decode(&pkg)
		if err != nil {
			return nil, err
		}
		br.Packages = append(br.Packages, map[string]*Package{pkg.Name: &pkg})
		br.Arch[pkg.Arch] = len(br.Packages) - 1
	}

	return &br, nil
}
