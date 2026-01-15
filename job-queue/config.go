package main

import "time"

const (
	MaxAttempts	=	3;
	BaseBackoff	=	500*time.Millisecond;
	MaxBackoff	=	5*time.Second;
)