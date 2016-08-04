package system

import (
	"bytes"
	"strconv"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// MailTemplate is used for sending email to users
type MailTemplate struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Type    int           `bson:"type" json:"type"`
	Subject string        `bson:"subject" json:"subject"`
	Content string        `bson:"content" json:"content"`
	Version int           `bson:"version" json:"version"`
}

// GenerateContent execute template and return content and error
func (m *MailTemplate) GenerateContent(data interface{}) (string, error) {
	tmpl := template.Must(template.New(strconv.Itoa(m.Type)).Parse(m.Content))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		log.Errorf("Execute template error: %s", err)
		log.Errorf("template type: %d, content: %s, version: %d", m.Type, m.Content, m.Version)
		return "", err
	}
	return buf.String(), nil
}
