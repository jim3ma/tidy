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
	return m.generateHelp(data, m.Content)
}

// GenerateSubject execute template and return subject and error
func (m *MailTemplate) GenerateSubject(data interface{}) (string, error) {
	return m.generateHelp(data, m.Subject)
}

func (m *MailTemplate) generateHelp(data interface{}, field string) (string, error) {
	tmpl := template.Must(template.New("subject_" + strconv.Itoa(m.Type)).Parse(field))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		log.Errorf("Execute template error: %s", err)
		log.Errorf("template type: %d, field: %s, version: %d", m.Type, field, m.Version)
		return "", err
	}
	return buf.String(), nil
}
