package mail

import (
	"encoding/json"
	"github.com/lowitea/jeevez/internal/config"
	mRand "math/rand"
	"net/http"
	"net/url"
)

const (
	createMailboxEndpoint         = "/api/mail/createMailbox"
	changeMailboxSettingsEndpoint = "/api/mail/changeMailboxSettings"
	forwardListAddMailboxEndpoint = "/api/mail/forwardListAddMailbox"
	//MailboxForwardOption          = "forward_and_delete"
	MailboxForwardOption = "forward"
	tempMailboxPwdLen    = 32
)

var tempMailboxNameRunes = []rune("abcdefghijklmnopqrstuvwxyz")

type createMailboxData struct {
	Domain          string `json:"domain"`
	Mailbox         string `json:"mailbox"`
	MailboxPassword string `json:"mailbox_password"`
}

type changeMailboxSettingsData struct {
	Domain            string `json:"domain"`
	Mailbox           string `json:"mailbox"`
	SpamFilterStatus  int    `json:"spam_filter_status"`
	SpamFilter        int    `json:"spam_filter"`
	ForwardMailStatus string `json:"forward_mail_status"`
}

type forwardListAddMailboxData struct {
	Domain         string `json:"domain"`
	Mailbox        string `json:"mailbox"`
	ForwardMailbox string `json:"forward_mailbox"`
}

func genPwd(length int) string {
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + digits + specials
	buf := make([]byte, length)
	buf[0] = digits[mRand.Intn(len(digits))]
	buf[1] = specials[mRand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[mRand.Intn(len(all))]
	}
	mRand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

func CreateMail(domain, mName string) (mBox string) {
	send := func(endpoint string, data interface{}) (*http.Response, error) {
		iData, _ := json.Marshal(data)
		rURL, _ := url.Parse(config.Cfg.Mail.Host)
		rURL.Path = endpoint
		params := url.Values{
			"login":         {config.Cfg.Mail.Login},
			"passwd":        {config.Cfg.Mail.Password},
			"input_format":  {"json"},
			"output_format": {"json"},
			"input_data":    {string(iData)},
		}
		rURL.RawQuery = params.Encode()
		return http.Get(rURL.String())
	}

	// создаём новый ящик
	_, _ = send(createMailboxEndpoint, createMailboxData{
		Domain:          domain,
		Mailbox:         mName,
		MailboxPassword: genPwd(tempMailboxPwdLen),
	})

	// настраиваем созданный ящик
	_, _ = send(changeMailboxSettingsEndpoint, changeMailboxSettingsData{
		Domain:            domain,
		Mailbox:           mName,
		SpamFilterStatus:  1,
		SpamFilter:        50,
		ForwardMailStatus: MailboxForwardOption,
	})

	// добавляем адрес переадресации
	_, _ = send(forwardListAddMailboxEndpoint, forwardListAddMailboxData{
		Domain:         domain,
		Mailbox:        mName,
		ForwardMailbox: config.Cfg.Admin.Email,
	})

	return mName + "@" + domain
}
