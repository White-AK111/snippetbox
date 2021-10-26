package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"White-AK111/snippetbox/pkg/models"
)

// Обработчик главной страницы.
// Меняем сигнатуры обработчика home, чтобы он определялся как метод.
// структуры *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/". Если нет, вызывается
	// функция http.NotFound() для возвращения клиенту ошибки 404.
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Получаем последние земетки из БД.
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Заполняем шаблон, отрисовываем страницу.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// Обработчик для отображения содержимого заметки.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение параметра id из URL и попытаемся
	// конвертировать строку в integer используя функцию strconv.Atoi(). Если его нельзя
	// конвертировать в integer, или значение меньше 1, возвращаем ответ
	// 404 - страница не найдена!
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Вызываем метода Get из модели Snipping для извлечения данных для
	// конкретной записи на основе её ID. Если подходящей записи не найдено,
	// то возвращается ответ 404 Not Found (Страница не найдена).
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Заполняем шаблон, отрисовываем страницу.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

// Обработчик для создания новой заметки.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Используем r.Method для проверки, использует ли запрос метод POST или нет.
	// http.MethodPost является строкой и содержит текст "POST".
	if r.Method != http.MethodPost {
		// Если это не так, то
		// Используем метод Header().Set() для добавления заголовка 'Allow: POST' в
		// карту HTTP-заголовков. Первый параметр - название заголовка, а
		// второй параметр - значение заголовка.
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Создаем несколько переменных, содержащих тестовые данные.
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

	// Передаем данные в метод SnippetModel.Insert(), получая обратно
	// ID только что созданной записи в базу данных.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправляем пользователя на соответствующую страницу заметки.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Форма для создания новой заметки..."))
}
