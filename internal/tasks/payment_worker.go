package tasks

import (
	"log"
	"taskgo/internal/database/models"
	"taskgo/internal/services"
)

type OrderTask struct {
	Order *models.Order
}

var paymentQueue = make(chan OrderTask, 100)

func StartPaymentWorkerPool(service *services.PaymentService, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go paymentWorker(i, service)
	}
}

func paymentWorker(id int, service *services.PaymentService) {
	for task := range paymentQueue {
		log.Printf("[Worker %d] 🏁 Starting payment task\n", id)
		err := service.ProcessPayment(task.Order)
		if err != nil {
			log.Printf("[Worker %d] ❌ Payment failed: %v\n", id, err)
		} else {
			log.Printf("[Worker %d] ✅ Payment task completed\n", id)
		}
	}
}

func EnqueueOrder(task OrderTask) {
	paymentQueue <- task
}
