package utils

import "github.com/sirupsen/logrus"

// Error handling function
func check(e error, m string) {
	if e != nil {
		logrus.WithError(e).Fatal(m)
	}
}

// Error handling function
func Check(e error, m string) {
	check(e, m)
}
