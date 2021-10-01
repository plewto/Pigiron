package piglog

import (
	"fmt"
	"os"
	"log"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/pigpath"
)

var (
	logfile string
	file *os.File
)


func init() {
	if config.GlobalParameters.EnableLogging {
		logfile = pigpath.SubSpecialDirectories(config.GlobalParameters.Logfile)
		var err error
		_ = os.Remove(logfile)
		file, err = os.OpenFile(logfile, os.O_WRONLY | os.O_CREATE, 0666)
	
		if err != nil {
			errmsg := "ERROR: Can not open log file: %s, logging disabled\n\n"
			fmt.Printf(errmsg, logfile)
			config.GlobalParameters.EnableLogging = false
		} else {
			log.SetOutput(file)
		}
	}
}


func Logfile() string {
	return logfile
}

func Close() {
	if file != nil {
		fmt.Printf("Closing logfile: %s\n", logfile)
		file.Close()
	}
}


func Log(text ...string) {
	if config.GlobalParameters.EnableLogging {
		for _, s := range text {
			log.Print(s)
		}
	}
}

func Print(s string) {
	if config.GlobalParameters.EnableLogging {
		log.Print(s)
	}
}
