/*
** macro package defines command aliases by text-replacement.
** Macros may be used interactively at the prompt or within batch files,
** they do not define new OSC commands.
**
*/

package macro

import (
	"fmt"
	"strings"
	"strconv"
	"sort"
)

const DEREF = '$'

var macros = make(map[string]*Macro)


/*
** Macro struct defines a command replacement.
** name - The macro's name
** command - The OSC command name
** template - arguments to the expanded command.  
**            Use $n to substitute nth argument into expanded text.
**            
*/
type Macro struct {
	name string
	command string
	template []string

}

func (macro *Macro) String() string {
	acc := fmt.Sprintf("Macro %s --> %s ", macro.name, macro.command)
	for _, t := range macro.template {
		acc += fmt.Sprintf("%s, ", t)
	}
	acc = strings.TrimSpace(acc)
	if len(macro.template) > 0 {
		acc = acc[:len(acc)-1]
	}
	return acc
}

// Define function creates a new Macro definition.
// name - The macro name.
// command - The OSC command the macro is an alias to.
// template - List of arguments to the expanded command.
//    Template contents may either be literal replacement text or have the
//    special syntax '$n' which replaces the nth argument to name into the
//    expanded command.
//
func Define(name string, command string, template []string) {
	m := &Macro{name, command, template}
	macros[name] = m
}


// Delete function deletes the named macro.
//
func Delete(name string) {
	delete(macros, name)
}

// IsMacro returns true iff macro name exists.
//
func IsMacro(name string) bool {
	_, flag := macros[name]
	return flag
}

// Expand function converts macro call to replacement text.
// name - the macro's name.
// args - arguments to name.
//
// Returns expanded macro text.
//
func Expand(name string, args []string) (string, error) {
	var err error
	var acc string
	macro, flag := macros[name]
	if !flag {
		msg := "Macro '%s' is not defined"
		err = fmt.Errorf(msg, name)
		return "", err
	}
	acc = fmt.Sprintf("%s ", macro.command)
	for _, t := range macro.template {
		t = strings.TrimSpace(t)
		if t[0] == DEREF {
			var index int
			index, err = strconv.Atoi(t[1:])
			if err != nil || index < 0 || index >= len(args) {
				msg := "Illegal macro index, '%s', err = %v"
				err = fmt.Errorf(msg, t, err)
				return "", err
			}
			acc += fmt.Sprintf("%s, ", strings.TrimSpace(args[index]))
		} else {
			acc += fmt.Sprintf("%s, ", t)
		}
	}
	acc = strings.TrimSpace(acc)
	if len(macro.template) > 0 {  // trim final comma
		acc = acc[:len(acc)-1]
	}
	return acc, err
}

// ListMacros returns sorted list of defined macros.
//
func ListMacros() []string {
	keys := make([]string, len(macros))
	acc := make([]string, len(macros))
	index := 0
	for key, _ := range macros {
		keys[index] = key
		index++
	}
	sort.Strings(keys)
	for i, key := range keys {
		acc[i] = macros[key].String()
	}
	return acc

}

// Reset deletes all macros
//
func Reset() {
	macros = make(map[string]*Macro)
}
