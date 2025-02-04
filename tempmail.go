package tempmail_wrapper

import (
    "fmt"
    "net/http"
)

var (
    client = &http.Client{}

    baseURL       = "https://api.internal.temp-mail.io/api/v3"
    newAccountURL = fmt.Sprintf("%s/email/new", baseURL)
    controlURL    = func(email string) string { return fmt.Sprintf("%s/email/%s", baseURL, email) }
    mailboxURL    = func(email string) string { return fmt.Sprintf("%s/messages", controlURL(email)) }
    attachmentURL = func(id string) string { return fmt.Sprintf("%s/attachment/%s", baseURL, id) }
)
