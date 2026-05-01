package service

import "errors"

var ErrBookingOverlap = errors.New("these dates are already booked")
