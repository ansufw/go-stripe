package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates
var emailTemplateFS embed.FS

func (app *application) SendMail(from, to, subject, tmpl string, data interface{}) error {

	templateToRender := fmt.Sprintf("templates/%s.html.gohtml")

	t, err := template.New("email-html").ParseFS(emailTemplateFS)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	formattedMessage := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.gohtml")

	t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	plainMessaes := tpl.String()

	return nil
}
