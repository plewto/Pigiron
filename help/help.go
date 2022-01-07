package help

// help package provides online help.
// All help documents stored as simple text files in 
//
//     ~/.config/pigiron/resources/help
//

import (
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/plewto/pigiron/pigpath"
)


// Help function returns documentation for topic.
//
func Help(topic string) (text string, err error) {
	tokens := strings.Split(topic, " ")
	if len(tokens) < 1 {
		return
	}
	head := strings.TrimSpace(tokens[0])
	switch head {
	case "topics":
		filter := ""
		if len(tokens) > 1 {
			filter = strings.TrimSpace(tokens[1])
		}
		return helpTopics(filter)
	default:
		filename := pigpath.ResourceFilename("help", head)
		var data []byte
		data, err = ioutil.ReadFile(filename)
		text = string(data)
		return text, err
	}
}
	
func helpTopics(filter string) (string, error) {
	dirname := pigpath.ResourceFilename("help")
	topics, err := ioutil.ReadDir(dirname)
	if err != nil {
		errmsg := "Can not accesses help directory\n%s"
		err = fmt.Errorf(errmsg, err)
		return "", err
	}
	acc := "Help topics:\n"
	for _, info := range topics {
		name := fmt.Sprintf("%s", info.Name())
		if strings.Contains(name, filter) {
			acc += fmt.Sprintf("\t%s\n", name)
		}
	}
	acc += "\n"
	return acc, err
}
