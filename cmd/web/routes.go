package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	// Используется функция http.NewServeMux() для инициализации нового рутера, затем
	mux := http.NewServeMux()
	// Используем методы из структуры в качестве обработчиков маршрутов.
	// функцию "home" регистрируется как обработчик для URL-шаблона "/".
	mux.HandleFunc("/", app.home)
	// Регистрируем два новых обработчика и соответствующие URL-шаблоны в маршрутизаторе servemux
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Инициализируем FileServer, он будет обрабатывать
	// HTTP-запросы к статическим файлам из папки "./ui/static".
	// Используем настраиваемую файловую систему.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("../../ui/static/")})
	// Используем функцию mux.Handle() для регистрации обработчика для
	// всех запросов, которые начинаются с "/static/". Убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

// Структура настраиваемой файловой системуы
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Если файл index.html не существует, то метод вернет ошибку os.ErrNotExist
// (которая, в свою очередь, будет преобразована через http.FileServer в ответ 404 страница не найдена).
// Также вызываем метод Close() для закрытия только, что открытого index.html файла, чтобы избежать утечки файлового дескриптора.
// Во всех остальных случаях мы просто возвращаем файл и даем http.FileServer сделать то, что он должен.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() && err != nil {
		index := "index.html"
		if _, err := nfs.fs.Open(index); err == nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
		} else {
			return nil, err
		}
	}

	return f, nil
}
