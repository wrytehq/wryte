package validator

type SetupForm struct {
	Name            string `form:"name" validate:"required,min=2,max=100"`
	Email           string `form:"email" validate:"required,email"`
	Password        string `form:"password" validate:"required,min=6,max=72"`
	ConfirmPassword string `form:"confirmPassword" validate:"omitempty,eqfield=Password"`
}

type LoginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

type ForgotPasswordForm struct {
	Email string `form:"email" validate:"required,email"`
}

type ResetPasswordForm struct {
	Password        string `form:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `form:"confirmPassword" validate:"required,eqfield=Password"`
	Token           string `form:"token" validate:"required"`
}

type ChangePasswordForm struct {
	CurrentPassword string `form:"currentPassword" validate:"required"`
	NewPassword     string `form:"newPassword" validate:"required,min=8,max=72"`
	ConfirmPassword string `form:"confirmPassword" validate:"required,eqfield=NewPassword"`
}
