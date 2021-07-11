package smf

func exError(message string, more ...string) error {
	msg := fmt.Sprintf("ERROR: %s", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nERROR: %s", s)
	}
	return fmt.Errorf(msg)
}

func compoundError(err1 error, message string, more ...string) error {
	msg := fmt.Sprintf("ERROR: %s\n", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nERROR: %s", s)
	}
	msg += fmt.Sprintf("PREVIOUS ERROR:\n%s", err1)
	return fmt.Errorf(msg)
}
