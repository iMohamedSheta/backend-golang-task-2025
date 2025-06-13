package bootstrap

import (
	"sync"
	"taskgo/internal/rules"
	"taskgo/internal/tasks"

	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
)

/*
This file is used to register all the custom stuff in the application
like the validations rules, ...etc
*/

// Global registered tasks
var registeredTasks map[string]asynq.Handler
var once sync.Once

// Register new validations rules
var registeredRules = map[string]validator.Func{
	// Add your custom validation rules here
	"unique_db":      rules.UniqueDB,
	"exists_db":      rules.ExistsDB,
	"egyptian_phone": rules.EgyptianPhone,
}

// registerTaskHandlers defines all individual task handlers
func registerTaskHandlers() map[string]asynq.Handler {
	return map[string]asynq.Handler{
		tasks.TypeProcessPayment: tasks.NewProcessPaymentHandler(),
		tasks.TypeInventoryCheck: tasks.NewInventoryCheckHandler(),
		//...
	}
}
