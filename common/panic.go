package common

import "github.com/op/go-logging"

func HandleDecodePanic(log *logging.Logger) {
	if r := recover(); r != nil {
		log.Warningf("recovered from decode panic")
	}
}
