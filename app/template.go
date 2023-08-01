package app

import (
	"bytes"
	"text/template"
)

func (app *App) initTemplate() {
	tmpl, err := template.ParseFS(FS(), "templates/*.prompt")
	if err != nil {
		panic(err)
	}

	app.promptManager = &PromptManager{tmpl}
}

type PromptManager struct {
	template *template.Template
}

func (pm *PromptManager) Execute(name string, data interface{}) (string, error) {
	buf := &bytes.Buffer{}
	err := pm.template.ExecuteTemplate(buf, name, data)
	return buf.String(), err
}
