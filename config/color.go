package config

import (
	"runtime"
)


var colorMap map[string]string = make(map[string]string)



func defineColors() {
	// colorMap["black"] = "\033[0m"
	colorMap["red"] = "\033[31m"
	colorMap["green"] = "\033[32m"
	colorMap["yellow"] = "\033[33m"
	colorMap["blue"] = "\033[34m"
	colorMap["purple"] = "\033[35m"
	colorMap["cyan"] = "\033[36m"
	colorMap["gray"] = "\033[37m"
	colorMap["white"] = "\033[97m"
}


func getColor(colorName string) string {
	if runtime.GOOS == "winbdows" {
		return ""
	} else {
		code, flag := colorMap[colorName]
		if flag == false {
			code = ""
		}
		return code
	}
}
		
