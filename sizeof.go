package sizeof

import (
	"fmt"
	"reflect"
	"runtime"
)

const (
	lengthThresholdForGC   = 4000000
	chidlrenThresholdForGC = 100
)

// MemHumanReadableValue converts bytes to human readable string
// like the -h option in 'du' command.
func MemHumanReadableValue(bytes int) string {
	var helper func(v float64, units []string) string

	helper = func(v float64, units []string) string {
		if len(units) == 1 {
			if v-float64(int64(v)) > 0 {
				return fmt.Sprintf("%.2f%s", v, units[0])
			}

			return fmt.Sprintf("%.0f%s", v, units[0])
		}

		if int64(v/1024) != 0 {
			return helper(v/1024, units[1:])
		}

		return helper(v, units[:1])
	}

	return helper(float64(bytes), []string{"b", "K", "M", "G"})
}

// SpaceUsageReport holds detail information of space usage.
type SpaceUsageReport struct {
	Values   map[string]string
	Children []*SpaceUsageReport
}

func newSpaceUsageReport() *SpaceUsageReport {
	return &SpaceUsageReport{
		Values:   make(map[string]string),
		Children: make([]*SpaceUsageReport, 0),
	}
}

func (sur *SpaceUsageReport) addValue(key string, value string) {
	if sur == nil {
		return
	}

	sur.Values[key] = value
}

func (sur *SpaceUsageReport) addChild() *SpaceUsageReport {
	if sur == nil {
		return nil
	}

	child := newSpaceUsageReport()

	sur.Children = append(sur.Children, child)

	return child
}

func (sur *SpaceUsageReport) childrenLength() int {
	if sur == nil {
		return 0
	}

	return len(sur.Children)
}

type sizeOfCalculator struct {
	ptrTraversal   map[uintptr]struct{}
	extendedReport bool
}

func (soc *sizeOfCalculator) sizeOf(
	rv reflect.Value, report *SpaceUsageReport,
) (int, *SpaceUsageReport) {
	size := 0

	if rv.IsZero() {
		report.addValue("__object-kind", "zero-value")

		return size, report
	}

	report.addValue("___type", rv.Type().String())
	report.addValue("__object-kind", rv.Type().Kind().String())

	switch rv.Type().Kind() {
	case reflect.Ptr:
		size = soc.sizePointer(rv, report)
	case reflect.Interface:
		size = soc.sizeInterface(rv, report)
	case reflect.String:
		size = soc.sizeString(rv, report)
	case reflect.Map:
		size = soc.sizeMap(rv, report)
	case reflect.Slice:
		size = soc.sizeSlice(rv, report)
	case reflect.Array:
		size = soc.sizeArray(rv, report)
	case reflect.Struct:
		size = soc.sizeStruct(rv, report)
	case reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64,
		reflect.Float32, reflect.Float64, reflect.Func,
		reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Invalid, reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		size = int(rv.Type().Size())
	default:
		size = int(rv.Type().Size())
	}

	report.addValue("__size", MemHumanReadableValue(size))

	if report.childrenLength() > chidlrenThresholdForGC {
		runtime.GC()
	}

	return size, report
}

func (soc *sizeOfCalculator) sizeSlice(
	rv reflect.Value, report *SpaceUsageReport,
) int {
	soc.ptrTraversal[rv.Pointer()] = struct{}{}

	report.addValue("_length", fmt.Sprintf("%d", rv.Len()))
	report.addValue("_capacity", fmt.Sprintf("%d", rv.Cap()))

	size := int(rv.Type().Size())

	if rv.Cap() > rv.Len() {
		tl := rv.Cap() - rv.Len()
		ds := defaultSize(rv.Type().Elem())
		s := ds * tl

		size += s

		report.addValue("cap-len:length", fmt.Sprintf("%d", tl))
		report.addValue("cap-len:size", MemHumanReadableValue(s))
	}

	report.addValue("size-slice-overhead", MemHumanReadableValue(size))

	size += soc.sizeArray(rv, report)

	return size
}

func (soc *sizeOfCalculator) sizeArray(
	rv reflect.Value, report *SpaceUsageReport,
) int {
	report.addValue("_length", fmt.Sprintf("%d", rv.Len()))

	size := 0

	if isTypeWithFloatingValueSize(rv.Type().Elem()) {
		report.addValue("count-each-key", "yes")

		for i := 0; i < rv.Len(); i++ {
			var child *SpaceUsageReport

			if soc.extendedReport {
				child = report.addChild()
			}

			s, _ := soc.sizeOf(rv.Index(i), child)
			size += s

			if (i+1)%lengthThresholdForGC == 0 {
				runtime.GC()
			}
		}
	} else if rv.Len() != 0 {
		s, _ := soc.sizeOf(rv.Index(0), report.addChild())
		size += s * rv.Len()
	}

	return size
}

func (soc *sizeOfCalculator) sizeMap(
	rv reflect.Value, report *SpaceUsageReport,
) int {
	soc.ptrTraversal[rv.Pointer()] = struct{}{}

	report.addValue("_length", fmt.Sprintf("%d", rv.Len()))

	size := 0
	// + 8 bytes (:pointer)
	size += int(rv.Type().Size())
	// + 48 bytes (:hmap)
	size += 48 + 24
	// + len(map) / 8 (:bucket's size) * 8 bytes (:bmap)
	size += rv.Len() / 8 * 8

	report.addValue("size-structure", MemHumanReadableValue(size))

	keys := rv.MapKeys()

	if isTypeWithFloatingValueSize(rv.Type().Key()) {
		report.addValue("count-each-key", "yes")

		for i := 0; i < len(keys); i++ {
			var child *SpaceUsageReport

			if soc.extendedReport {
				child = report.addChild()
			}

			s, _ := soc.sizeOf(keys[i], child)
			size += s

			if (i+1)%lengthThresholdForGC == 0 {
				runtime.GC()
			}
		}
	} else if rv.Len() != 0 {
		s, _ := soc.sizeOf(keys[0], report.addChild())
		size += s * len(keys)
	}

	iter := rv.MapRange()

	if isTypeWithFloatingValueSize(rv.Type().Elem()) {
		report.addValue("count-each-value", "yes")

		i := 0

		for iter.Next() {
			var child *SpaceUsageReport

			if soc.extendedReport {
				child = report.addChild()
			}

			s, _ := soc.sizeOf(iter.Value(), child)
			size += s

			i++

			if i%lengthThresholdForGC == 0 {
				runtime.GC()
			}
		}
	} else if rv.Len() != 0 {
		iter.Next()
		s, _ := soc.sizeOf(iter.Value(), report.addChild())
		size += s * len(keys)
	}

	return size
}

func (soc *sizeOfCalculator) sizeStruct(
	rv reflect.Value, report *SpaceUsageReport,
) int {
	size := int(rv.Type().Size())

	for i := 0; i < rv.NumField(); i++ {
		rChild := report.addChild()

		rChild.addValue("____field", rv.Type().Field(i).Name)

		s, _ := soc.sizeOf(rv.Field(i), rChild)

		size += s
	}

	return size
}

func (soc *sizeOfCalculator) sizeString(
	rv reflect.Value, report *SpaceUsageReport,
) int {
	report.addValue("length", fmt.Sprintf("%d", len(rv.String())))

	return len(rv.String()) + int(rv.Type().Size())
}

func (soc *sizeOfCalculator) sizePointer(
	rv reflect.Value, report *SpaceUsageReport,
) int {
	report.addValue("point-to", reflect.Indirect(rv).Type().String())

	size := int(rv.Type().Size())

	if _, ok := soc.ptrTraversal[rv.Pointer()]; ok {
		report.addValue("already-taken", "yes")

		return 0
	}

	soc.ptrTraversal[rv.Pointer()] = struct{}{}

	if !rv.IsNil() {
		s, _ := soc.sizeOf(reflect.Indirect(rv), report.addChild())
		size += s
	}

	return size
}

func (soc *sizeOfCalculator) sizeInterface(
	rv reflect.Value, report *SpaceUsageReport,
) int {
	size := int(rv.Type().Size())

	s, _ := soc.sizeOf(rv.Elem(), report.addChild())
	size += s

	return size
}

func defaultSize(rt reflect.Type) int {
	switch rt.Kind() {
	case reflect.Struct:
		size := 0

		for i := 0; i < rt.NumField(); i++ {
			size += defaultSize(rt.Field(i).Type)
		}

		return size
	case reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64,
		reflect.Float32, reflect.Float64, reflect.Func,
		reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Invalid, reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer,
		reflect.Array, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.String:
		return int(rt.Size())
	default:
		return int(rt.Size())
	}
}

func isTypeWithFloatingValueSize(rt reflect.Type) bool {
	switch rt.Kind() {
	case reflect.Slice:
		return true
	case reflect.Array, reflect.Ptr:
		return isTypeWithFloatingValueSize(rt.Elem())
	case reflect.Interface:
		return true
	case reflect.String:
		return true
	case reflect.Map:
		return isTypeWithFloatingValueSize(rt.Key()) ||
			isTypeWithFloatingValueSize(rt.Elem())
	case reflect.Struct:
		for i := 0; i < rt.NumField(); i++ {
			if isTypeWithFloatingValueSize(rt.Field(i).Type) {
				return true
			}
		}

		return false
	case reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64,
		reflect.Float32, reflect.Float64, reflect.Func,
		reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Invalid, reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		return false
	default:
		return false
	}
}

// SizeOf returns total size in bytes that the object allocate in memory.
func SizeOf(v interface{}) int {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	soc := sizeOfCalculator{
		ptrTraversal: make(map[uintptr]struct{}),
	}

	size, _ := soc.sizeOf(rv, nil)

	return size
}

// Options is a struct for options for this component.
type Options struct {
	extendedReport bool
}

func defaultOptions() Options {
	return Options{
		extendedReport: false,
	}
}

// Option is the func interface to assign options.
type Option func(*Options)

// ExtendedReport includes to the report every object in each slice and map.
func ExtendedReport() Option {
	return func(o *Options) {
		o.extendedReport = true
	}
}

// SizeOfVerbose returns total size in bytes that the object allocate in memory.
// The second returned value is a detailed report of space usage.
func SizeOfVerbose(v interface{}, opts ...Option) (int, SpaceUsageReport) {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	options := defaultOptions()

	for _, o := range opts {
		o(&options)
	}

	soc := sizeOfCalculator{
		ptrTraversal:   make(map[uintptr]struct{}),
		extendedReport: options.extendedReport,
	}

	size, report := soc.sizeOf(rv, newSpaceUsageReport())

	return size, *report
}
