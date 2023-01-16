package main

import "errors"

var (
	ErrNotNumber = errors.New("Data is not numeric")
	ErrInvalidCol = errors.New("Invalid column number")
	ErrNoFiles = errors.New("No input file")
	ErrInvalidOp = errors.New("Invalida operation")
)