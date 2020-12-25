# SizeOf
**SizeOf** helps to know size of any data structure in the memory. 

## API 
* SizeOf returns total size in bytes that the object allocate in memory.
`
    func SizeOf(v interface{}) int {
`
* SizeOfVerbose returns total size in bytes that the object allocate in memory. The second returned value is a detailed report of space usage.
`
    func SizeOfVerbose(v interface{}) (int, SpaceUsageReport) {
`
* PrintReport prints SpaceUsageReport to the given io.Writer.
`
    func PrintReport(r *SpaceUsageReport, w io.Writer) 
`
* MemHumanReadableValue converts bytes to human readable string. Is behaves like the -h option in 'du' command.
`
    func MemHumanReadableValue(bytes int) string {
`

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
| ------------ |:---------------------:| ----------------------------------------------------------------------------:|
| *            | \___type              | Datatype of the object                                                       |
| Everyove has | \__object-kind        | Kind of the object array/map/slice/struct and etc...                         |
|              | \__size               | String with human a readable size value                                      |
| Array        | \_length              | Lenght of the array                                                          |
|              | \_count-each-key      | Shows that size calculation individualy provided for each underlining object |
| Map          | \_length              | Lenght of the array                                                          |
|              | \_size-structure      | Memory for holding size structure of the slice                               |
|              | \_count-each-key      | Shows that size calculation individualy provided for each underlining object |
| Pointer      | point-to              | Type of the object to whitch pointer points to                               |
|              | already-taken         | Size already taken to account somewhere else                                 |
| Slice        | \_length              | Lenght of the slice                                                          |
|              | \_capacity            | Capacity of the slice                                                        |
|              | \_cap-len:length      | Difference between capacity and length                                       |
|              | \_cap-len:size        | Allocation for that difference                                               |
|              | \_size-slice-overhead | Memory for holding size structure of the slice                               |
| String       | length                | Length of the string                                                         |
| Struct       | \____field            | Field name of the structure                                                  |


## Donations
 If you want to support this project, please consider donating:
 * PayPal: https://paypal.me/MaxFe
