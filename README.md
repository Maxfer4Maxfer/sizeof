# SizeOf
**SizeOf** helps to know size of any data structure in the memory. 

## API 
* SizeOf returns total size in bytes that the object allocate in memory.  
`
    func SizeOf(v interface{}) int {
`
* SizeOfVerbose returns total size in bytes that the object allocate in memory. The second returned value is a detailed report of space usage.  
`
    func SizeOfVerbose(SizeOfVerbose(v interface{}, opts ...Option) (int, SpaceUsageReport) {
`
* PrintReport prints SpaceUsageReport to the given io.Writer.  
`
    func PrintReport(r *SpaceUsageReport, w io.Writer) 
`
* MemHumanReadableValue converts bytes to human readable string. Is behaves like the -h option in 'du' command.  
`
    func MemHumanReadableValue(bytes int) string {
`

Options SizeOfVerbose:
 * ExtendedReport includes to the report every object in each slice and map.

## SpaceUsageReport 
SpaceUsageReport holds detail information of space usage.

`
    type SpaceUsageReport struct {
        Values   map[string]string
        Children []*SpaceUsageReport
    }
`

Possible valuues of the SpaceUsageReport:
| Go Type      | Report's field        | Description  |
| ------------ | --------------------- | ---------------------------------------------------------------------------- |
| *            | \___type              | Datatype of the object                                                       |
| Everyove has | \__object-kind        | Kind of the object array/map/slice/struct and etc...                         |
|              | \__size               | String with human a readable size value                                      |
| Array        | \_length              | Lenght of the array                                                          |
|              | \_count-each-key      | Shows that size calculation individualy provided for each underlining object |
| Map          | \_length              | Lenght of the array                                                          |
|              | \_size-structure      | Memory for holding size structure of the slice                               |
|              | \_count-each-key      | Shows that size calculation individualy provided for each stored key         |
|              | \_count-each-value    | Shows that size calculation individualy provided for each stored value       |
| Pointer      | point-to              | Type of the object to whitch pointer points to                               |
|              | already-taken         | Size already taken to account somewhere else                                 |
| Slice        | \_length              | Lenght of the slice                                                          |
|              | \_capacity            | Capacity of the slice                                                        |
|              | \_cap-len:length      | Difference between capacity and length                                       |
|              | \_cap-len:size        | Allocation for that difference                                               |
|              | \_size-slice-overhead | Memory for holding size structure of the slice                               |
| String       | length                | Length of the string                                                         |
| Struct       | \____field            | Field name of the structure                                                  |

## Example 
`
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
`

## Donations
 If you want to support this project, please consider donating:
 * PayPal: https://paypal.me/MaxFe
