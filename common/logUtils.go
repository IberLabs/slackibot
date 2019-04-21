package common

import (
	"github.com/google/logger"
	"fmt"
)

/**
	Main log/print method all in one
 */
func Display(txt string, error bool, logIt bool) {
	if logIt {
		if error {
			txt = "ERROR " + txt
			logger.Errorln(txt)
		} else {
			logger.Infoln(txt)
		}
	}
	fmt.Println(txt)
}
