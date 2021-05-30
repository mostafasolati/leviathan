package auth

import (
	"bastek7/contracts"
	"bastek7/lib/event"
	"bastek7/lib/utils"
	"bastek7/models"
	"encoding/hex"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// CPAccountingSecret is the secret key to encode and decode JWT tokens
	CPAccountingSecret = "accounting.secret"

	// CPGuestExpireTimeout is the timeout (seconds) at-which guests expire.
	CPGuestExpireTimeout = "accounting.guest.expire-timeout"

	// CPJWTExpiration is the jwt expiration time in minute
	CPJWTExpiration = "accounting.jwt-expiration"
)

// NewAuthService creates a new IAccountingService.
func NewAuthService(
	config contracts.IConfigService,
	logger contracts.ILogger,
) contracts.IAccountingService {
	service := &accounting{
		config: config,
		logger: logger,
	}

	return service
}

// FindOTPOfAUser implements IAccounting.FindOTPOfAUser
func (a *accounting) FindOTPOfAUser(phone string) (string, error) {
	otp, ok := a.otpMap.Load(phone)
	if !ok || otp.(otpType).expireAt.Before(time.Now()) {
		return "", contracts.ErrOTPNotFound
	}
	return otp.(otpType).code, nil
}

// SendOTP sends otp to user phone number
func (a *accounting) SendOTP(phone string, app string) error {
	phone = utils.NormalizePhoneNumber(phone)

	validator := NewValidation()
	validator.PhoneNumber(phone)
	if err := validator.Validate(); err != nil {
		return err
	}

	user, _ := a.storage.FindByPhone(phone)
	if user != nil && user.DeletedAt != nil {
		return contracts.ErrUserDeactivated
	}

	otp := a.generateOTP(phone)
	event.Fire(&models.OTPEvent{Phone: phone, OTP: otp, App: app})
	return nil
}

// NoSendOTP generates OTP for a phone number but doesn't send it.
func (a *accounting) NoSendOTP(phone string) string {
	return a.generateOTP(phone)
}

// LoginByOTP either register or login user by authenticating the otp sent in previous step
func (a *accounting) LoginByOTP(phone, code string, guestID int) (token *models.Token, err error) {
	phone = utils.NormalizePhoneNumber(phone)
	validator := NewValidation()

	validator.PhoneNumber(phone)
	if err := validator.Validate(); err != nil {
		return nil, err
	}

	otp, ok := a.otpMap.Load(phone)
	if !ok || otp.(otpType).code != code || otp.(otpType).expireAt.Before(time.Now()) {
		return nil, contracts.ErrOTPIsIncorrect
	}

	user, err := a.storage.FindByPhone(phone)
	if err != nil {
		switch err {
		case contracts.ErrUserNotFound:
			// register user if not found in database
			user = &models.User{
				Phone: phone,
			}
			err = a.storage.Create(user)
			if err != nil {
				return nil, err
			}

			event.Fire(&models.UserCreatedEvent{
				UserID:  user.ID,
				GuestID: guestID,
			})
		default:
			return nil, err
		}
	}

	// create a new refresh token for user if not exists or expired
	if user.RefreshTokenExpiry.Before(time.Now()) {
		bs := make([]byte, 32)
		rand.Read(bs)
		user.RefreshToken = hex.EncodeToString(bs)
		user.RefreshTokenExpiry = time.Now().Add(6 * 30 * 24 * time.Hour)

		if err = a.storage.Update(user); err != nil {
			return nil, err
		}
	}

	accessToken, err := a.createJwtToken(user)
	if err != nil {
		return nil, err
	}

	if guestID != 0 {
		// Remove the guest, as it is no longer needed.
		if err = a.guestStorage.Remove(guestID); err != nil {
			return nil, err
		}
	}

	event.Fire(&models.UserLoggedInEvent{
		UserID:  user.ID,
		GuestID: guestID,
	})

	return &models.Token{
		AccessToken:  accessToken,
		RefreshToken: user.RefreshToken,
	}, nil
}

// RefreshToken return a new access token based on user refresh token
func (a *accounting) RefreshToken(refreshToken string) (accessToken string, err error) {
	user, err := a.storage.FindByRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return a.createJwtToken(user)
}

// ParseToken validates the access token and extracts its claims.
func (a *accounting) ParseToken(accessToken string) (*models.UserClaims, error) {
	claims, err := a.parseJwtToken(accessToken)
	if err != nil {
		return nil, contracts.ErrUnauthorized
	}

	return claims, nil
}

/********** types **********/
type otpType struct {
	code     string
	expireAt time.Time
}

type accounting struct {
	storage       contracts.IAccountingStorage
	guestStorage  contracts.IGuestStorage
	config        contracts.IConfigService
	messageBroker contracts.IMessageBroker
	otpMap        sync.Map
	acl           contracts.IACL
	logger        contracts.ILogger
}

/********** Helper functions ***********/

// A helper function to generate a random number
// and store in a map of generated otps
func (a *accounting) generateOTP(phone string) string {

	// search for existing otp in map
	if otp, ok := a.otpMap.Load(phone); ok {
		// if the existing code expired delete it otherwise return it
		if otp.(otpType).expireAt.Before(time.Now()) {
			a.otpMap.Delete(phone)
		} else {
			return otp.(otpType).code
		}
	}

	// in case of not found otp in map
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(99999-10000) + 10000
	otp := strconv.Itoa(n)

	if !a.config.IsProduction() && phone == "09120000000" {
		otp = "00000"
	}

	a.otpMap.Store(phone, otpType{code: otp, expireAt: time.Now().Add(10 * time.Minute)})

	return otp
}

func (a *accounting) createJwtToken(user *models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(a.config.Int(CPJWTExpiration)) * time.Minute).Unix(),
		},

		ID:    user.ID,
		Name:  user.FullName(),
		Roles: []string{"user"},
		Phone: user.Phone,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(a.config.String(CPAccountingSecret)))

	return tokenString, err
}

func (a *accounting) parseJwtToken(accessToken string) (*models.UserClaims, error) {
	var claims models.UserClaims
	token, err := jwt.ParseWithClaims(accessToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(a.config.String(CPAccountingSecret)), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(*models.UserClaims), nil
}
