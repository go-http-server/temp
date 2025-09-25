package mailer

import (
	"testing"

	"github.com/go-http-server/core/utils"
	"github.com/stretchr/testify/require"
)

func TestSendMailWithTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	env, err := utils.LoadEnviromentVariables("./../../../")
	require.NoError(t, err)
	require.NotEmpty(t, env)

	sender := NewGmailSender(env.EmailUsernameSender, env.EmailAddressSender, env.EmailPasswordSender)
	require.NotNil(t, sender)
	receiver := UserReceive{
		Username:     utils.RandomString(6),
		EmailAddress: "21A100100257@students.hou.edu.vn",
		Code:         utils.RandomCode(),
		Fullname:     "Pham Hai Nam",
		AttachFiles:  []string{"../../../docker-compose.yml"},
	}
	err = sender.SendWithTemplate("[Go core] Kích hoạt tài khoản", "../../../templates/verify_account.html", receiver)
	require.NoError(t, err)
}
