package dispacher

type Dispacher struct {
	NotificationChan chan []byte
}

func NewDispacher(notificationChan chan []byte) *Dispacher {
	return &Dispacher{
		NotificationChan: notificationChan,
	}
}

func (d *Dispacher) Notifique(msg []byte) {
	d.NotificationChan <- msg
}

func (d *Dispacher) Consumer() <-chan []byte {
	return d.NotificationChan
}
