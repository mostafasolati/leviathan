package notification

import (
	"fmt"

	kn "github.com/kavenegar/kavenegar-go"
	"github.com/mostafasolati/leviathan/contracts"
)

type kavenegar struct {
	api *kn.Kavenegar
}

type ApikeyType string

func NewKavenegar(apikey ApikeyType) contracts.INotificationService {
	return &kavenegar{
		api: kn.New(apikey),
	}
}

func (k *kavenegar) Receive() (map[string]string, error) {
	ret := make(map[string]string, 0)
	messages, err := k.api.Message.Receive("100077070", true)
	if err != nil {
		return nil, err
	}
	for _, message := range messages {
		ret[message.Sender] = message.Message
	}

	return ret, nil
}

func (k *kavenegar) SendSMS(reciever []string, message string) error {
	_, err := k.api.Message.Send("1000596446", reciever, message, nil)
	if err != nil {
		fmt.Println("KAVENEGAR ERR:", reciever, err)
	}
	return err
}
