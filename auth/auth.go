package auth

import (
	"encoding/hex"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mostafasolati/leviathan/contracts"
	"github.com/mostafasolati/leviathan/models"
	"github.com/mostafasolati/leviathan/utils"
)

// NewAuthService creates a new IAccountingService.
func NewAuthService(
	config contracts.IConfigService,
	logger contracts.ILogger,
	userService contracts.IUserService,
	notification contracts.INotificationService,
) contracts.IAuth {
	return &auth{
		config:       config,
		logger:       logger,
		userService:  userService,
		notification: notification,
	}
}

// FindOTPOfAUser implements IAccounting.FindOTPOfAUser
func (s *auth) FindOTPOfAUser(phone string) (string, error) {
	otp, ok := s.otpMap.Load(phone)
	if !ok || otp.(otpType).expireAt.Before(time.Now()) {
		return "", contracts.ErrOTPNotFound
	}
	return otp.(otpType).code, nil
}

// SendOTP sends otp to user phone number
func (s *auth) SendOTP(phone string, app string) error {
	phone = utils.NormalizePhoneNumber(phone)

	// Todo: validate phone number
	user, _ := s.userService.FindByPhone(phone)
	if user != nil && user.DeletedAt() != nil {
		return contracts.ErrUserDeactivated
	}

	otp := s.generateOTP(phone)
	// todo: send it via an event
	s.notification.SendSMS(user.Phone, otp)
	return nil
}

// NoSendOTP generates OTP for a phone number but doesn't send it.
func (s *auth) NoSendOTP(phone string) string {
	return s.generateOTP(phone)
}

// LoginByOTP either register or login user by authenticating the otp sent in previous step
func (s *auth) LoginByOTP(phone, code string, guestID int) (token *models.Token, err error) {
	phone = utils.NormalizePhoneNumber(phone)
	validator := NewValidation()

	validator.PhoneNumber(phone)
	if err := validator.Validate(); err != nil {
		return nil, err
	}

	otp, ok := s.otpMap.Load(phone)
	if !ok || otp.(otpType).code != code || otp.(otpType).expireAt.Before(time.Now()) {
		return nil, contracts.ErrOTPIsIncorrect
	}

	user, err := s.storage.FindByPhone(phone)
	if err != nil {
		switch err {
		case contracts.ErrUserNotFound:
			// register user if not found in database
			user = &models.User{
				Phone: phone,
			}
			err = s.storage.Create(user)
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

	// create s new refresh token for user if not exists or expired
	if user.RefreshTokenExpiry.Before(time.Now()) {
		bs := make([]byte, 32)
		rand.Read(bs)
		user.RefreshToken = hex.EncodeToString(bs)
		user.RefreshTokenExpiry = time.Now().Add(6 * 30 * 24 * time.Hour)

		if err = s.storage.Update(user); err != nil {
			return nil, err
		}
	}

	accessToken, err := s.createJwtToken(user)
	if err != nil {
		return nil, err
	}

	if guestID != 0 {
		// Remove the guest, as it is no longer needed.
		if err = s.guestStorage.Remove(guestID); err != nil {
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
func (s *auth) RefreshToken(refreshToken string) (accessToken string, err error) {
	user, err := s.storage.FindByRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return s.createJwtToken(user)
}

// ParseToken validates the access token and extracts its claims.
func (s *auth) ParseToken(accessToken string) (*models.UserClaims, error) {
	claims, err := s.parseJwtToken(accessToken)
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

type auth struct {
	config       contracts.IConfigService
	otpMap       sync.Map
	logger       contracts.ILogger
	userService  contracts.IUserService
	notification contracts.INotificationService
}

/********** Helper functions ***********/

// A helper function to generate a random number
// and store in a map of generated otps
func (s *auth) generateOTP(phone string) string {

	// search for existing otp in map
	if otp, ok := s.otpMap.Load(phone); ok {
		// if the existing code expired delete it otherwise return it
		if otp.(otpType).expireAt.Before(time.Now()) {
			s.otpMap.Delete(phone)
		} else {
			return otp.(otpType).code
		}
	}

	// in case of not found otp in map
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(99999-10000) + 10000
	otp := strconv.Itoa(n)

	if !s.config.IsProduction() && phone == "09120000000" {
		otp = "00000"
	}

	s.otpMap.Store(phone, otpType{code: otp, expireAt: time.Now().Add(10 * time.Minute)})

	return otp
}

func (s *auth) createJwtToken(user *models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(s.config.Int(CPJWTExpiration)) * time.Minute).Unix(),
		},

		ID:    user.ID,
		Name:  user.FullName(),
		Roles: []string{"user"},
		Phone: user.Phone,
	})

	// Sign and get the complete encoded token as s string using the secret
	tokenString, err := token.SignedString([]byte(s.config.String(CPAccountingSecret)))

	return tokenString, err
}

func (s *auth) parseJwtToken(accessToken string) (*models.UserClaims, error) {
	var claims models.UserClaims
	token, err := jwt.ParseWithClaims(accessToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.config.String(CPAccountingSecret)), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(*models.UserClaims), nil
}
