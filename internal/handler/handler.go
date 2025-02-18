package handler

import (
	"html/template"
	"net/http"

	"github.com/is-web-y26/m3302-milovatskiy/internal/data"
)

func HandleCommon(
	responseWriter http.ResponseWriter,
	request *http.Request,
	handlerTemplate string,
) {
	template, err := template.ParseFiles(
		"templates/layout.html",
		"templates/header.html",
		"templates/footer.html",
		handlerTemplate,
	)
	if err != nil {
		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	if err := template.ExecuteTemplate(
		responseWriter,
		"layout",
		data.Data,
	); err != nil {
		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}
