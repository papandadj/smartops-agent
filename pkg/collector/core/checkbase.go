package core

import (
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/defaults"
	"time"
)

type CheckBase struct {
	checkName      string
	latestWarnings []error
	checkInterval  time.Duration
}

func NewCheckBase(name string) CheckBase {
	return CheckBase{
		checkName:     name,
		checkInterval: defaults.DefaultCheckInterval,
	}
}
func (c *CheckBase) Interval() time.Duration {
	return c.checkInterval
}
func (c *CheckBase) String() string {
	return c.checkName
}
