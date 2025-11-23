package app_err

import (
	er "errors"
)

var (
	ServerErr           = er.New("something went wrong, try later")
	ParseErr            = er.New("cannot parse data, please check provided data")
	DocumentNotFoundErr = er.New("document not found")
	UpdateDocumentErr   = er.New("document update error")
)
