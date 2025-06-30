package load

import (
	"taskgo/pkg/ioc"
	"taskgo/pkg/utils"
	"taskgo/pkg/validate"

	"github.com/go-playground/validator/v10"
)

func InitValidator(c *ioc.Container, registeredRules map[string]validator.Func) {
	err := ioc.Singleton(c, func(c *ioc.Container) (*validate.Validator, error) {
		v := validator.New()

		for tag, rule := range registeredRules {
			if err := v.RegisterValidation(tag, rule); err != nil {
				return nil, err
			}
		}

		return validate.New(v), nil
	})

	if err != nil {
		utils.PrintErr("Failed to load validator module in the ioc container : " + err.Error())
	}
}
