package dispacher

import "log"

type Notifier struct {
	NotificationChan chan []byte
}

func NewNotifier(notificationChan chan []byte) *Notifier {
	return &Notifier{
		NotificationChan: notificationChan,
	}
}

func (d *Notifier) Send(msg []byte) {
	if d.NotificationChan == nil {
		log.Fatal("Error: NotificationChan es nil")
	}
	d.NotificationChan <- msg
}

func (d *Notifier) Consumer() <-chan []byte {
	if d.NotificationChan == nil {
		log.Fatal("Error: NotificationChan es nil")
	}
	return d.NotificationChan
}

func (d *Notifier) Close() {
	close(d.NotificationChan)
}
