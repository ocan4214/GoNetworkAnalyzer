package RabbitMQInterfaces

type ISubscriber interface {
	HandleMessage()
}
