package notification

import (
	"fmt"

	kn "github.com/kavenegar/kavenegar-go"
	"github.com/mostafasolati/leviathan/contracts"
)

type kavenegar struct {
	api    *kn.Kavenegar
	config contracts.IConfigService
}

type ApikeyType string

func NewKavenegar(config contracts.IConfigService) contracts.INotificationService {

	return &kavenegar{
		api:    kn.New(config.String("notification.sms.api_key")),
		config: config,
	}
}

func (k *kavenegar) Receive() (map[string]string, error) {
	ret := make(map[string]string, 0)
	messages, err := k.api.Message.Receive(k.config.String("notification.sms.phone"), true)
	if err != nil {
		return nil, err
	}
	for _, message := range messages {
		ret[message.Sender] = message.Message
	}

	return ret, nil
}

func (k *kavenegar) SendSMS(phone string, message string) error {
	_, err := k.api.Message.Send("1000596446", []string{phone}, message, nil)
	if err != nil {
		fmt.Println("KAVENEGAR ERR:", phone, err)
	}
	return err
}
