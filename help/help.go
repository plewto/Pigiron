package help

// help package provides online help.
// All help documents stored as simple text files in 
//
//     ~/.config/pigiron/resources/help
//

import (
	"fmt"
	"io/ioutil"
	"github.com/plewto/pigiron/pigpath"
)

func Help(topic string) (text string, err error) {
	switch topic {
	case "topics":
		return helpTopics()
	default:
		filename := pigpath.ResourceFilename("help", topic)
		var data []byte
		data, err = ioutil.ReadFile(filename)
		text = string(data)
		return text, err
	}
}
	
func helpTopics() (string, error) {
	dirname := pigpath.ResourceFilename("help")
	topics, err := ioutil.ReadDir(dirname)
	if err != nil {
		errmsg := "Can not accesses help directory\n%s"
		err = fmt.Errorf(errmsg, err)
		return "", err
	}
	acc := "Help topics:\n"
	for _, info := range topics {
		acc += fmt.Sprintf("\t%s\n", info.Name())
	}
	acc += "\n"
	return acc, err
}
