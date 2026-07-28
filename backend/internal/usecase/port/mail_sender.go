package port

import "context"

type IMailSender interface {
	SendPasswordRecoveryCode(ctx context.Context, to, code string) error
}
