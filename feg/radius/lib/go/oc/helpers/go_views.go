package helpers

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	mGoroutines   = stats.Int64("goroutines", "The number of goroutines", stats.UnitDimensionless)
	mFrees        = stats.Int64("frees", "The number of frees", stats.UnitDimensionless)
	mHeapAllocs   = stats.Int64("heap_allocs", "The number of heap allocations", stats.UnitDimensionless)
	mHeapObjects  = stats.Int64("heap_objects", "The number of objects allocated on the heap", stats.UnitDimensionless)
	mHeapReleased = stats.Int64("heap_released", "The number of objects released from the heap", stats.UnitDimensionless)
	mPtrLookups   = stats.Int64("ptr_lookups", "The number of pointer lookups", stats.UnitDimensionless)
	mStackSys     = stats.Int64("stack_sys", "The memory used by stack spans and OS thread stacks", stats.UnitDimensionless)
)

// Views for the stats quickstart.
var (
	GoroutinesCountView = &view.View{
		Name:        "cpu.goroutines",
		Measure:     mGoroutines,
		Description: "The number of cpu.goroutines",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent}}

	MemFreesCountView = &view.View{
		Name:        "mem_frees",
		Measure:     mFrees,
		Description: "The number of frees",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent}}

	MemAllocsCountView = &view.View{
		Name:        "mem_allocs",
		Measure:     mHeapAllocs,
		Description: "The number of heap allocations",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent}}

	MemHeapObjCountView = &view.View{
		Name:        "mem_heap_objects",
		Measure:     mHeapObjects,
		Description: "The number of objects allocated on the heap",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent}}

	MemHeapReleasedCountView = &view.View{
		Name:        "mem_heap_released",
		Measure:     mHeapReleased,
		Description: "The number of objects released from the heap",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent}}

	MemPtrLookupsCountView = &view.View{
		Name:        "mem_ptr_lookups",
		Measure:     mPtrLookups,
		Description: "The number of pointer lookups",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent}}
)
