package user

import "html/template"

type Service struct {
	indexTpl *template.Template
	showTpl  *template.Template
}

func New(tplDir string) *Service {
	return &Service{
		indexTpl: template.Must(template.ParseFiles()),
	}
}
