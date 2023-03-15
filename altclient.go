// package altclient
package main

import (
	"encoding/json"
	"net/http"
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
	Name     string                `json:"-"`
	Length   int                   `json:"length"`
	Arch     map[string]int        `json:"-"`
	Packages []map[string]*Package `json:"-"`
}

type Diff struct {
	Branch            string
	Missing          []Package
	Redundant        []Package
	OutOfDate        []Package
}

// Returns lists of supported architectures

func (br *Branch) GetArchs() (archs []string) {
	for key, _ := range br.Arch {
		archs = append(archs, key)
	}
	return
}

// Returns list of packages with a specified archtecture

func(br *Branch) getPackages(arch string) map[string]*Package {
if _,ok:=br.Arch[arch]; !ok {
return nil
}
return br.Packages[br.Arch[arch]]
}

// Returns list of the missing packages in comparing branch

func (br *Branch) GetMissing(br_to_compare *Branch, arch string) ([]Package, error) {
	from_comparing:= br_to_compare.getPackages(arch)
	var pkgs []Package
	if from_comparing == nil {
		for name, _ := range br.Packages[br.Arch[arch]] {
			pkgs = append(pkgs, *br.Packages[br.Arch[arch]][name])
		}
	} else {
		for name, _ := range br.Packages[br.Arch[arch]] {
			if _, ok := from_comparing[name]; !ok {
				pkgs = append(pkgs, *br.Packages[br.Arch[arch]][name])
			}
		}
	}
	return pkgs, nil
}



func NewBranch(br string) (*Branch, error) {
	resp, err := http.Get(ApiURL + br)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	branch := Branch{Name: br, Arch: make(map[string]int)}
	dec := json.NewDecoder(resp.Body)

	for t, err := dec.Token(); t != "packages"; {
		if err != nil {
			return nil, err
		}
		if t == "length" {
			l, err := dec.Token()
			if err != nil {
				return nil, err
			}
			branch.Length = int(l.(float64))
		}
		t, err = dec.Token()
	}
	dec.Token()
	for dec.More() {
		var pkg Package
		err = dec.Decode(&pkg)
		if err != nil {
			return nil, err
		}
		if _, ok := branch.Arch[pkg.Arch]; !ok {
			branch.Packages = append(branch.Packages, map[string]*Package{})
			branch.Arch[pkg.Arch] = len(branch.Packages) - 1
		}
		branch.Packages[branch.Arch[pkg.Arch]][pkg.Name] = &pkg
	}

	return &branch, nil
}

func main() {

}
