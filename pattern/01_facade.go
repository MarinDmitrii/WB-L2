package pattern

import "fmt"

/*
Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,
а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
Применимость паттерна "фасад":
1. Когда система становится сложной и требует упрощения интерфейса для клиентов.
2. Когда необходимо скрыть детали реализации подсистем от клиентов.

Плюсы паттерна "фасад":
1. Упрощение интерфейса для клиентов.
2. Сокрытие сложности системы.
3. Соблюдение принципа единственной ответственности: Фасад отвечает за координацию работы подсистем, оставляя им отдельные области ответственности

Минусы паттерна "фасад":
1. Может ввести в систему дополнительные слои абстракции, что усложнит обслуживание в будущем.
2. Может привести к созданию god object.
*/

// PaymentService представляет сервис для обработки платежей.
type PaymentService struct {
}

// ProcessPayment моделирует обработку платежа для указанного orderID.
func (ps *PaymentService) ProcessPayment(orderID int) {
	fmt.Printf("PaymentService: Processed payment for order %v\n", orderID)
}

// DeliveryService представляет сервис для обновления статуса доставки.
type DeliveryService struct {
}

// UpdateDeliveryStatus моделирует обновление статуса доставки для указанного orderID.
func (ds *DeliveryService) UpdateDeliveryStatus(orderID int, status string) {
	fmt.Printf("Updated delivery status for order %v - %v\n", orderID, status)
}

// NotificationService представляет сервис для отправки уведомлений.
type NotificationService struct {
}

// SendNotification моделирует отправку подтверждения для указанного orderID.
func (ns *NotificationService) SendNotification(orderID int) {
	fmt.Printf("NotificationService: Sent confirmation for order %v\n", orderID)
}

// OrderFacade представляет фасад, упрощающий взаимодействие с несколькими сервисами.
type OrderFacade struct {
	paymentService      *PaymentService
	deliveryService     *DeliveryService
	notificationService *NotificationService
}

// NewOrderFacade создает новый экземпляр OrderFacade с инициализированными сервисами.
func NewOrderFacade() *OrderFacade {
	return &OrderFacade{
		paymentService:      &PaymentService{},
		deliveryService:     &DeliveryService{},
		notificationService: &NotificationService{},
	}
}

// PlaceOrder упрощает процесс размещения заказа, координируя работу различных сервисов.
func (of *OrderFacade) PlaceOrder(orderID int, status string) {
	of.paymentService.ProcessPayment(orderID)
	of.deliveryService.UpdateDeliveryStatus(orderID, status)
	of.notificationService.SendNotification(orderID)
}

func facadePattern() {
	order := NewOrderFacade()
	order.PlaceOrder(1, "Shipped")
}
