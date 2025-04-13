package repository

import "errors"

var (
	ErrUnsupportedCity        = errors.New("city is not supported")
	ErrUnsupportedProductType = errors.New("product type is not supported")
	ErrReceptionInProgress    = errors.New("last reception is in progress")
	ErrNoReceptionInProgress  = errors.New("no reception is in progress")
	ErrReceptionIsEmpty       = errors.New("reception is empty")
)
