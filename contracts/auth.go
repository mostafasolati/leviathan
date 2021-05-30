package contracts

import "github.com/mostafasolati/leviathan/models"

type IAuth interface {
	// SendOTP : User send his/her phone number and asks for an otp code
	SendOTP(phone string, app string) error

	// NoSendOTP generates OTP for a phone number but doesn't send it.
	NoSendOTP(phone string) string

	// LoginByOTP : User send his/her phone and received otp from previous step
	// and asks for login. We will register the user if he/she is not registered
	// before. Then we will send a pair of access token and refresh token to the
	// user
	//
	// When a guest attempts to login, guestID will be non-zero and represents
	// the ID of the guest.
	LoginByOTP(phone, code string, guestID int) (token *models.Token, err error)

	// RefreshToken : Jwt lifetime is short. most of the time 30 minutes. in order to
	// prevent hackers to steal other people tokens and send request.
	// so every time jwt expires the client send it's refresh token to us
	// and we give him another access token
	RefreshToken(refreshToken string) (accessToken string, err error)

	// ParseToken validates the access token and extracts its claims.
	//
	// It returns ErrUnauthorized if the access token is invalid or expired.
	ParseToken(accessToken string) (*models.UserClaims, error)
}
