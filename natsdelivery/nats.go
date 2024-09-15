package natsdelivery

import (
	"encoding/json"

	"github.com/kowiste/boilerplate/src/messaging"
	"gorm.io/gorm"
)

type NatsDelivery struct {
	db     *gorm.DB
	broker messaging.IMessaging
	need   chan bool
}

// Product represents a product model/
type message struct {
	gorm.Model
	Topic  string              `json:"topic"`
	Event  string              `json:"event"`
	Data   []byte              `json:"data"`
	Status ENatsDeliveryStatus `json:"status"`
}

func New(db *gorm.DB, broker messaging.IMessaging) *NatsDelivery {
	return &NatsDelivery{
		db:     db,
		broker: broker,
		need:   make(chan bool),
	}
}

func (nd *NatsDelivery) Init() (err error) {
	err = nd.db.AutoMigrate(&message{})
	if err != nil {
		return
	}

	go nd.deliver()
	nd.need <- true
	return
}

func (nd *NatsDelivery) Send(topic, event string, data any) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = nd.db.Save(&message{
		Topic:  topic,
		Event:  event,
		Data:   b,
		Status: DeliveryStatusPending,
	}).Error
	nd.need <- true
	return
}

func (nd *NatsDelivery) deliver() {
	var messages []message
	var err error
	for range nd.need {
		err = nd.db.Find(&messages).
			Where("status != ? OR status != ?", DeliveryStatusDone, DeliveryStatusSending).Error
		if err != nil {
			nd.need <- true
		}
		for i := range messages {
			go nd.sendMessage(messages[i])
		}

	}
}
func (nd NatsDelivery) sendMessage(msg message) {
	var err error
	defer func() {
		if err != nil {
			nd.need <- true
		}
	}()
	msg.Status = DeliveryStatusSending
	err = nd.db.Updates(msg).Error
	if err != nil {
		return
	}
	err = nd.broker.Send(msg.Topic, msg.Event, msg.Data)
	msg.Status = DeliveryStatusDone
	if err != nil {
		msg.Status = DeliveryStatusError
	}
	nd.db.Updates(msg)

}
