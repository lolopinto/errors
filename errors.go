package errorsutil

import (
	"strings"
	"sync"
)

func CoalesceErr(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var errors []error
	for _, err := range errs {
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return &ErrorList{errs: errors}
}

// ErrorList encompasses a list of errors. It's also an error
type ErrorList struct {
	errs []error
}

func (e *ErrorList) Error() string {
	var sb strings.Builder
	for idx, e2 := range e.errs {
		sb.WriteString(e2.Error())
		if idx != 0 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// SyncError abstracts out collecting errors in goroutines
// Makes it safe to add each new error in Append
// This doesn't implement the error interface so that it doesn't break the if err != nil
// paradigm in golang. To check if there's an error, call the Err() method and check if that's nil
type SyncError struct {
	m    sync.Mutex
	errs []error
}

// Append takes an error and adds a new instance to its collection of errors
func (e *SyncError) Append(err error) {
	if err == nil {
		return
	}
	e.m.Lock()
	defer e.m.Unlock()
	e.errs = append(e.errs, err)
}

// Err returns something which satisfies the error interface
func (e *SyncError) Err() error {
	return CoalesceErr(e.errs...)
}
