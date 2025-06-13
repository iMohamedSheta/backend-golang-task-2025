package bootstrap

import (
	"taskgo/internal/tasks"

	"github.com/hibiken/asynq"
)

// GetRegisteredTaskHandlers returns all task handlers including the configured orchestrator
func GetRegisteredTaskHandlers() map[string]asynq.Handler {
	once.Do(func() {
		// Create the orchestrator first
		orchestrator := tasks.NewChainOrchestrator()

		// Get all individual task handlers
		individualHandlers := registerTaskHandlers()

		// Auto-register all individual handlers with the orchestrator
		for taskType, handler := range individualHandlers {
			orchestrator.RegisterHandler(taskType, handler)
		}

		// Initialize the global registeredTasks with individual handlers
		registeredTasks = make(map[string]asynq.Handler)
		for taskType, handler := range individualHandlers {
			registeredTasks[taskType] = handler
		}

		// Add the configured orchestrator to the map
		registeredTasks[tasks.TypeChainOrchestrator] = orchestrator
	})

	return registeredTasks
}
