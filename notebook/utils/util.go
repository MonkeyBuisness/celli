package utils

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Close closes the io.Closer instance and logs an error if the close did not complete as expected.
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		logrus.Warnf("could not close stream: %v", err)
	}
}
