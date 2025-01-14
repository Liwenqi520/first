package helper

import (
	"log"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)
type ApnsConfig struct {
	PemFile  string `json:"pem_file"`
	Password string `json:"password"`
}
func ApnsPush(Config ApnsConfig, DeviceToken, Topic, PayLoad string) (res interface{}, err error) {
	cert, pemErr := certificate.FromPemFile(Config.PemFile, Config.Password)
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}
	notification := &apns2.Notification{}
	notification.DeviceToken = DeviceToken
	notification.Topic = Topic
	notification.Payload = []byte(PayLoad)
	client := apns2.NewClient(cert).Development()
	res, err = client.Push(notification)
	if err != nil {
		log.Println("Error:", err)
	}
	return
}

