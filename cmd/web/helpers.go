package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера".
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник clientError отправляет определенный код состояния и соответствующее описание
// пользователю. Мы будем использовать это в следующий уроках, чтобы отправлять ответы вроде 400 "Bad
// Request", когда есть проблема с пользовательским запросом.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Обработка (рендер) шаблона, возврат обработанного шаблона в ответ на запрос.
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]

	if !ok {
		app.serverError(w, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

// Функция добавления данных по умолчнию в шаблон.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {

	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()
	//td.Flash = app.session.PopString(r, "flash")
	return td
}
