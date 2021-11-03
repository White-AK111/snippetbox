-- Сайт SnippetBox. Хранение заметок пользователя с возможностью добавления заметок и напоминаний об истечении их сроков через различные каналы связи.
-- use-case
-- регистрация\идентификация\аутентификация\авторизация пользователя на форме входа;
-- отображение всех заметок пользователя, с не истёкшим сроком, на главной странице;
-- редактирование заметки, на форме работы с заметкой;
-- удаление заметки, на форме работы с заметкой;
-- создание новой заметки, на форме работы с заметкой;
-- настройка уведомления по заметке, на форме работы с заметкой;
-- получение уведомления по заметке по выбранному каналу связи (почта, TG, СМС, etc..).

-- таблица пользователей
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
    name VARCHAR(100) NOT NULL, -- отображаемое имя
	login VARCHAR(100) UNIQUE NOT NULL CONSTRAINT login_length CHECK (char_length(login) >= 3), -- логин, длинна больше или равна 3
	email VARCHAR(100) UNIQUE NOT NULL CONSTRAINT email_length CHECK (char_length(email) >= 4), -- почта, длинна больше или равна 4
	hashed_password BYTEA, -- хэш пароля
	created TIMESTAMPTZ NOT NULL, -- дата\время создания
	confirmed BOOLEAN DEFAULT FALSE -- флаг указывающий что УЗ пользователя подтверждена с помощью email
);

-- таблица заметок
DROP TABLE IF EXISTS snippets CASCADE;
CREATE TABLE snippets (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
	user_id INT NOT NULL REFERENCES users (id), -- внешний ключ на пользователя
    title VARCHAR(100) NOT NULL, -- заголовок
	content TEXT NOT NULL, -- содержимое
	created TIMESTAMPTZ NOT NULL, -- дата\время создания
	expires TIMESTAMPTZ NOT NULL, -- дата\время истечения
	changed TIMESTAMPTZ NOT NULL, -- дата\время изменения
	deleted BOOLEAN DEFAULT FALSE -- флаг удаления
);

-- справочник типов каналов связи
DROP TABLE IF EXISTS channels_types CASCADE;
CREATE TABLE channels_types (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
	name VARCHAR(100) UNIQUE NOT NULL, -- наименование канала (ex: Telegram)
	deleted BOOLEAN DEFAULT FALSE -- флаг удаления
);

-- таблица каналов связи пользователей
DROP TABLE IF EXISTS users_channels CASCADE;
CREATE TABLE users_channels (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
	user_id INT NOT NULL REFERENCES users (id), -- внешний ключ на пользователя
	channel_type_id INT NOT NULL REFERENCES channels_types (id), -- внешний ключ на тип канала
	address VARCHAR(100) NOT NULL CONSTRAINT address_length CHECK (char_length(address) >= 3), -- адрес доставки (ex: @white_ak111), длинна больше или равна 
	deleted BOOLEAN DEFAULT FALSE -- флаг удаления
);

-- таблица уведомлений\напоминаний
DROP TABLE IF EXISTS notifications CASCADE;
CREATE TABLE notifications (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
    snippet_id INT NOT NULL REFERENCES snippets (id), -- внешний ключ на заметку
	user_channel_id INT NOT NULL REFERENCES users_channels (id), -- внешний ключ на канал связи пользователя
	created TIMESTAMPTZ NOT NULL, -- дата\время создания
	time_to_send TIMESTAMPTZ NOT NULL, -- дата\время назначенной отправки уведомления
	sended  BOOLEAN DEFAULT FALSE, -- флаг отправки уведомления
	deleted BOOLEAN DEFAULT FALSE -- флаг удаления
);

-- таблица логов
DROP TABLE IF EXISTS logs CASCADE;
CREATE TABLE logs (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
    log_body JSON NOT NULL, -- сообщение лога в формате json
	created TIMESTAMPTZ NOT NULL -- -- дата\время создания
);

--делаем составной индекс на login + hashed_password
CREATE INDEX CONCURRENTLY users_login_password_idx on users (login, hashed_password);
--DROP INDEX CONCURRENTLY users_login_password_idx;

--делаем покрывающий составной индекс
CREATE INDEX CONCURRENTLY snippets_id_user_id_expires_deleted_changed_idx on snippets (id, user_id, expires, deleted, changed)
INCLUDE (title, content, created, expires, changed);
--DROP INDEX CONCURRENTLY snippets_id_user_id_expires_deleted_changed_idx;

--делаем покрывающий составной индекс
CREATE INDEX CONCURRENTLY notifications_sended_deleted_idx on notifications (sended, deleted)
INCLUDE (id, snippet_id, user_channel_id, created, time_to_send);
--DROP INDEX CONCURRENTLY notifications_sended_deleted_idx;