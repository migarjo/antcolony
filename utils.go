package main

//ApplicationError Errors that should be caught and handled gracefully
type ApplicationError struct {
	message string
}

func (e ApplicationError) Error() string {
	return e.message
}
