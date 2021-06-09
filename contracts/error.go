package contracts

// Define some error constant
const (
	// ErrFileDimensionQuota when an invalid dimention request comes to the server
	// we check the requested dimension to preven resource waste on the server
	ErrFileDimensionQuota = constError("width or height can't be more than 1000 pixels")
	ErrOTPNotFound        = constError("otp not found")
	ErrUserDeactivated    = constError("user is deactivated")
)

type constError string

func (err constError) Error() string {
	return string(err)
}
