package contracts

type INotificationService interface {
	SendSMS(phone, message string) error
}
