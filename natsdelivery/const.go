package natsdelivery

type ENatsDeliveryStatus = string

const (
	DeliveryStatusPending ENatsDeliveryStatus = "pending"
	DeliveryStatusSending ENatsDeliveryStatus = "sending"
	DeliveryStatusError   ENatsDeliveryStatus = "error"
	DeliveryStatusDone    ENatsDeliveryStatus = "done"
)
