package attacker

import (
	"github.com/rha7/pressure/types"
	"github.com/sirupsen/logrus"
)

// Perform //
func Perform(logger *logrus.Logger, spec *types.TestSpec) (types.Results, error) {
	var results types.Results
	for i := int64(0); i < spec.ConcurrentThreads; i++ {
		go processor(logger)
	}
	return results, nil
}
