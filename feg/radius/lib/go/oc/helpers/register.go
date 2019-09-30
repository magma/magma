package helpers

import (
	"go.opencensus.io/stats/view"
	"log"
)

// RegisterViews will make sure custom views are exported, promethus exports them diffrently.
// This is currently used only by openflow
func RegisterViews() {
	if err := view.Register(
		LatencyView,
		ErrorCountView,
		SuccessCountView,
		GoroutinesCountView,
		MemFreesCountView,
		MemAllocsCountView,
		MemHeapObjCountView,
		MemHeapReleasedCountView,
		MemPtrLookupsCountView,
	); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}
}
