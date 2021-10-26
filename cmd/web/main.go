package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"White-AK111/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// Создаем структуру `application` для хранения зависимостей всего веб-приложения.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Создаем новый флаг командной строки, значение по умолчанию: ":4000".
	// Добавляем небольшую справку, объясняющая, что содержит данный флаг.
	// Значение флага будет сохранено в переменной addr.
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	//флаг из командной строки для настройки MySQL подключения.
	dsn := flag.String("dsn", "web:pA55w0rD@/snippetbox?parseTime=true", "Название MySQL источника данных")

	// Вызываем функцию flag.Parse() для извлечения флага из командной строки.
	// Она считывает значение флага из командной строки и присваивает его содержимое
	// переменной. Вам нужно вызвать ее *до* использования переменной addr
	// иначе она всегда будет содержать значение по умолчанию ":4000".
	// Если есть ошибки во время извлечения данных - приложение будет остановлено.
	flag.Parse()

	// Используем log.New() для создания логгера для записи информационных сообщений. Для этого нужно
	// три параметра: место назначения для записи логов (os.Stdout), строка
	// с префиксом сообщения (INFO или ERROR) и флаги, указывающие, какая
	// дополнительная информация будет добавлена. Флаги
	// соединяются с помощью оператора OR |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Создаем логгер для записи сообщений об ошибках таким же образом, но используем stderr как
	// место для записи и используем флаг log.Lshortfile для включения в лог
	// названия файла и номера строки где обнаружилась ошибка.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Создаём пул соединений с помощью функции openDB(). Передаем в нее полученный
	// источник данных (DSN) из флага командной строки.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Откладываем вызов db.Close(), чтобы пул соединений был закрыт
	// до выхода из функции main().
	defer db.Close()

	// Инициализируем кэш шаблонов
	templateCache, err := newTemplateCache("../../ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	//Инициализируем новую структуру с зависимостями приложения.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Значение, возвращаемое функцией flag.String(), является указателем на значение
	// из флага, а не самим значением. Убраем ссылку на указатель.

	// Инициализируем новую структуру http.Server. Устанавливаем поля Addr и Handler, так
	// что сервер использует тот же сетевой адрес и маршруты, что и раньше, и назначаем
	// поле ErrorLog, чтобы сервер использовал наш логгер
	// при возникновении проблем.
	srv := &http.Server{
		Addr:           *addr,
		ErrorLog:       errorLog,
		Handler:        app.routes(),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	// Вызываем метод ListenAndServe() от структуры http.Server
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// Функция openDB() обертывает sql.Open() и возвращает пул соединений sql.DB
// для заданной строки подключения (DSN).
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
