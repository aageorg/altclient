package altclient

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-version"
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

// Map where key is an architecture name (e.g. "aarch64"),
// value is an index of the corresponding map in packages slice

	arch     map[string]int        `json:"-"`

// Slice of maps with packages combined by name of architecture.
// The key is a package name, the value - pointer to a Package object

	packages []map[string]*Package `json:"-"`
}

// Returns lists of supported architectures

func (br *Branch) GetArchs() (archs []string) {
	for key, _ := range br.arch {
		archs = append(archs, key)
	}
	return
}

// Returns list of packages with a specified archtecture

func (br *Branch) getPackages(a string) map[string]*Package {
	if _, ok := br.arch[a]; !ok {
		return nil
	}
	return br.packages[br.arch[a]]
}

// Returns list of out of missing packages in comparing branch

func (br *Branch) GetMissing(br_to_compare *Branch, a string) []Package {
	from_comparing := br_to_compare.getPackages(a)
	var pkgs []Package
	if from_comparing == nil {
		for name, _ := range br.packages[br.arch[a]] {
			pkgs = append(pkgs, *br.packages[br.arch[a]][name])
		}
	} else {
		for name, _ := range br.packages[br.arch[a]] {
			if _, ok := from_comparing[name]; !ok {
				pkgs = append(pkgs, *br.packages[br.arch[a]][name])
			}
		}
	}
	return pkgs
}

// Returns list of out of date packages in comparing branch

func (br *Branch) GetOutOfDate(br_to_compare *Branch, a string) []Package {
	from_comparing := br_to_compare.getPackages(a)
	var pkgs []Package
	if from_comparing == nil {
		return pkgs
	} else {
		for name, _ := range br.packages[br.arch[a]] {
			v1, err := version.NewVersion(br.packages[br.arch[a]][name].Version)
			if err != nil {
				continue
			}
			if _, ok := from_comparing[name]; ok {
				v2, err := version.NewVersion(from_comparing[name].Version)
				if err != nil {
					continue
				}
				if v2.LessThan(v1) {
					pkgs = append(pkgs, *from_comparing[name])
				}
			}
		}
	}
	return pkgs
}

// returns a pointer to a new Branch object

func NewBranch(br string) (*Branch, error) {
	resp, err := http.Get(ApiURL + br)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	branch := Branch{Name: br, arch: make(map[string]int)}
	dec := json.NewDecoder(resp.Body)

	for t, err := dec.Token(); t != "packages"; {
		if err != nil {
			return nil, err
		}

		if t == "validation_message" {
			dec.Token()
			t, err = dec.Token()
			if err != nil {
				return nil, err
			}
			var message string

			message = t.(string)
			for dec.More() {
				t, err = dec.Token()
				if err != nil {
					return nil, err
				}
				message += ", " + t.(string)
			}
			return nil, errors.New(message)
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
		if _, ok := branch.arch[pkg.Arch]; !ok {
			branch.packages = append(branch.packages, map[string]*Package{})
			branch.arch[pkg.Arch] = len(branch.packages) - 1
		}
		branch.packages[branch.arch[pkg.Arch]][pkg.Name] = &pkg
	}
	return &branch, nil
}
