package service

import "errors"

var (
	ErrCannotCreatePVZ = errors.New("cannot create pvz")
	ErrCannotGetPVZ    = errors.New("cannot get pvz")
	ErrUnsupportedCity = errors.New("city is not supported")

	ErrCannotCloseReception  = errors.New("cannot close reception")
	ErrCannotCreateReception = errors.New("cannot create reception")
	ErrNoReceptionInProgress = errors.New("no reception is in progress")
	ErrReceptionInProgress   = errors.New("last reception is in progress")

	ErrCannotAddProduct       = errors.New("cannot add product")
	ErrCannotDeleteProduct    = errors.New("cannot delete product")
	ErrUnsupportedProductType = errors.New("product type is not supported")
	ErrReceptionIsEmpty       = errors.New("reception is empty")
)
