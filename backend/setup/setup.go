package setup

import (
	"encoding/base64"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/pquerna/otp/totp"
)

func GenerateAdminSetupConstants(
	password string,
	passwordSettings *common.PasswordHashSettings,
	isDev bool,
	core common.CoreService,
) (
	envVars *common.AdminAuthEnvVars,
	totpURL string,
	wrappedErr common.WrappedError,
) {
	issuer := "Cryptic Stash"
	if isDev {
		issuer = "Cryptic Stash (Development)"
	}
	key, stdErr := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: common.AdminUsername,
	})
	if stdErr != nil {
		return nil, "", ErrWrapperGenerateAdminSetupConstants.Wrap(
			ErrWrapperTotp.Wrap(stdErr),
		)
	}
	totpSecret := key.Secret()

	salt := core.GenerateSalt()
	encryptionKey := core.HashPassword(password, salt, passwordSettings)

	return &common.AdminAuthEnvVars{
			ADMIN_PASSWORD_HASH: base64.StdEncoding.EncodeToString(encryptionKey),
			ADMIN_PASSWORD_SALT: base64.StdEncoding.EncodeToString(salt),
			ADMIN_TOTP_SECRET:   totpSecret,
		},
		key.URL(),
		nil
}

func CheckTotpCode(totpCode string, totpSecret string) bool {
	return totp.Validate(totpCode, totpSecret)
}
