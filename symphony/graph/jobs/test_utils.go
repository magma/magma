package jobs

import (
	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/symphony/graph/ent"
)

//TestJobsResolver contains data for jobs resolver
type TestJobsResolver struct {
	drv        dialect.Driver
	client     *ent.Client
	jobsRunner jobs
}
