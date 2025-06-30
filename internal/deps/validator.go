package deps

import (
	"taskgo/pkg/ioc"
	"taskgo/pkg/validate"
)

func Validator() *validate.Validator {
	v, err := ioc.AppMake[*validate.Validator]()
	if err != nil {
		panic(err)
	}
	return v
}
