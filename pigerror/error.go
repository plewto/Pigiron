// pigerror package implemented a composable error structure.
//

package pigerror

import "fmt"


// pigError implements the Error interface with composable messages.
//
type pigError struct {
	list []string
}

func (err *pigError) Error() string {
	if len(err.list) == 0 {
		return "Error <no-message>"
	}
	acc := ""
	for i, s := range err.list {
		acc += fmt.Sprintf("[%2d] ERROR: %s\n", i, s)
	}
	return acc
}

// Merge appends the messages from another pigError into this instance.
//
func (err *pigError) Merge(other *pigError) {
	err.list = append(err.list, other.list...)
}


// Add adds a new message to the pigError message list.
//
func (err *pigError) Add(message string) {
	err.list = append(err.list, message)
}

// AddError adds the message of another error into the pigError message list.
//
func (err *pigError) AddError(other error) {
	err.Add(fmt.Sprintf("%s", other))
}

// Print displays the pigError messages to the terminal.
//
func (err *pigError) Print() {
	fmt.Println(err.Error())
}

// PigError returns a pointer to a new instance of pigError.
// Initial messages may be empty.
//
func PigError(messages ...string) *pigError {
	list := make([]string, 0, len(messages))
	for _, s := range messages {
		list = append(list, s)
	}
	return &pigError{list}
}


