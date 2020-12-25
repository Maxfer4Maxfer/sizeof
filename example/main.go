package main

import (
	"fmt"
	"os"

	"github.com/Maxfer4Maxfer/sizeof"
)

type a struct {
	a int
	b string
	c []*a
}

func main() {
	a1 := a{a: 1, b: "one"}
	a2 := &a{a: 2, b: "two"}
	a3 := &a{a: 3, b: "three"}
	a4 := a{a: 4, b: "four"}

	all1 := a{1000, "all", []*a{&a1}}
	all2 := a{1000, "all", []*a{&a1, a2}}
	all3 := a{1000, "all", []*a{&a1, a2, a3}}
	all4 := a{1000, "all", []*a{&a1, a2, a3, &a4}}

	allS := []a{all1, all2, all3, all4}

	all := make(map[string]interface{}, 100)

	all["a-1"] = a1
	all["a-2"] = a2
	all["all-one"] = all1
	all["all-two"] = all2
	all["all-tree"] = all3
	all["all-four"] = all4
	all["all-slice"] = allS

	size, report := sizeof.SizeOfVerbose(all, sizeof.ExtendedReport())

	fmt.Printf("total %d bytes or %s\n", size, sizeof.MemHumanReadableValue(size))

	sizeof.PrintReport(&report, os.Stdout)
}
