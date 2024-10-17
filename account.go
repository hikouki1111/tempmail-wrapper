package tempmail_wrapper

import (
	"bytes"
	"encoding/json"
	"log"

	requests "github.com/RabiesDev/request-helper"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
)

type Account struct {
	Email string
	Token string
}

func NewAccount() (*Account, error) {
	payload := map[string]interface{}{
		"max_name_length": 10,
		"min_name_length": 10,
	}
	req := requests.Post(newAccountURL, bytes.NewReader(lo.Must(json.Marshal(payload))))
	requests.SetHeaders(req, jsonHeader)
	body, _, err := requests.DoAndReadString(client, req)

	return &Account{
		Email: gjson.Get(body, "email").String(),
		Token: gjson.Get(body, "token").String(),
	}, err
}

func (a *Account) Delete() error {
	payload := map[string]interface{}{
		"token": a.Token,
	}
	req := requests.Delete(controlURL(a.Email), bytes.NewReader(lo.Must(json.Marshal(payload))))
	requests.SetHeaders(req, jsonHeader)
	_, err := requests.Do(client, req)

	return err
}

type attachmentJson struct {
	ID         string `json:"id"`
	HasPreview bool   `json:"has_preview"`
	Name       string `json:"name"`
	Size       int    `json:"size"`
}

type emailJson struct {
	ID          string           `json:"id"`
	From        string           `json:"from"`
	To          string           `json:"to"`
	Subject     string           `json:"subject"`
	BodyText    string           `json:"body_text"`
	BodyHTML    string           `json:"body_html"`
	CreatedAt   string           `json:"created_at"`
	Attachments []attachmentJson `json:"attachments"`
}

func (a *Account) GetMailbox() ([]Mail, error) {
	req := requests.Get(mailboxURL(a.Email))
	body, _, err := requests.DoAndReadString(client, req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var emailJsons []emailJson
	err = json.Unmarshal([]byte(body), &emailJsons)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var mailSlice []Mail
	for _, eJ := range emailJsons {
		mail := Mail{
			ID:          eJ.ID,
			From:        eJ.From,
			To:          eJ.To,
			Subject:     eJ.Subject,
			BodyText:    eJ.BodyText,
			BodyHTML:    eJ.BodyHTML,
			CreatedAt:   eJ.CreatedAt,
			Attachments: make([]Attachment, 0),
		}

		if len(eJ.Attachments) > 0 {
			var attachments []Attachment
			for _, aJ := range eJ.Attachments {
				attachments = append(attachments, Attachment{
					Name: aJ.Name,
					URL:  attachmentURL(aJ.ID),
					Size: aJ.Size,
				})
			}
			mail.Attachments = attachments
		}

		mailSlice = append(mailSlice, mail)
	}

	return mailSlice, nil
}
