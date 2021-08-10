// Copyright (c) 2020 Alec Randazzo

package writers

type ResultWriter interface {
	Write(data interface{}) error
}
