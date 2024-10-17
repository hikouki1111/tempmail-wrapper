package tempmail_wrapper

type Mail struct {
	ID          string
	From        string
	To          string
	Subject     string
	BodyText    string
	BodyHTML    string
	CreatedAt   string
	Attachments []Attachment
}
