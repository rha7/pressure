package printers

import (
	"github.com/rha7/pressure/types"
	"github.com/sirupsen/logrus"
)

// Text //
func Text(results types.Results) {
	logrus.Info("Printing results")
}
