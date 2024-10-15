package dispacher

type Notifier struct {
	NotificationChan chan []byte
}

func NewNotifier(notificationChan chan []byte) *Notifier {
	return &Notifier{
		NotificationChan: notificationChan,
	}
}

func (d *Notifier) Send(msg []byte) {
	d.NotificationChan <- msg
}

func (d *Notifier) Consumer() <-chan []byte {
	return d.NotificationChan
}

func (d *Notifier) Close() {
	close(d.NotificationChan)
}
