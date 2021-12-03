package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Error represents generic model for error.
type Error struct {
	base    error
	details []string
}

// Notebook Error.
var (
	ErrReadMarkdownSource = Error{
		base: errors.New("could not read markdown source"),
	}
	ErrMarshalCommentPayload = Error{
		base: errors.New("could not marshal comment payload"),
	}
	ErrRenderNotebook = Error{
		base: errors.New("could not render notebook data"),
	}
	ErrReadNotebookSource = Error{
		base: errors.New("could not read notebook source"),
	}
	ErrParseNotebookContent = Error{
		base: errors.New("could not parse notebook content"),
	}
	ErrCreateTemplateContent = Error{
		base: errors.New("could not create template content"),
	}
)

// New creates a new copy of Error.
func (e Error) New(details ...string) Error {
	return Error{
		base:    e.base,
		details: details,
	}
}

// Error returns human-readable error message.
func (e Error) Error() string {
	err := e.base.Error()
	if len(e.details) != 0 {
		err = fmt.Sprintf("%s: %s", err, strings.Join(e.details, "; "))
	}

	return err
}

// Is compares target error with the current one.
func (e Error) Is(target error) bool {
	err, ok := target.(Error)
	if !ok {
		return false
	}

	return err.base == e.base
}
