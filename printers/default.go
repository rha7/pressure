package printers

import (
	"github.com/rha7/pressure/apptypes"
	"github.com/sirupsen/logrus"
)

// Default //
func Default(logger *logrus.Logger, summary apptypes.Summary) (string, error) {
	return Text(logger, summary)
}
