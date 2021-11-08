-- create tables
CREATE TABLE users (
                       id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
                       name VARCHAR(100) NOT NULL, -- отображаемое имя
                       login VARCHAR(100) UNIQUE NOT NULL CONSTRAINT login_length CHECK (char_length(login) >= 3), -- логин, длинна больше или равна 3
                       email VARCHAR(100) UNIQUE NOT NULL CONSTRAINT email_length CHECK (char_length(email) >= 4), -- почта, длинна больше или равна 4
                       hashed_password BYTEA, -- хэш пароля
                       created TIMESTAMPTZ NOT NULL, -- дата\время создания
                       confirmed BOOLEAN DEFAULT FALSE -- флаг указывающий что УЗ пользователя подтверждена с помощью email
);

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

CREATE TABLE channels_types (
                                id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
                                name VARCHAR(100) UNIQUE NOT NULL, -- наименование канала (ex: Telegram)
                                deleted BOOLEAN DEFAULT FALSE -- флаг удаления
);

CREATE TABLE users_channels (
                                id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
                                user_id INT NOT NULL REFERENCES users (id), -- внешний ключ на пользователя
                                channel_type_id INT NOT NULL REFERENCES channels_types (id), -- внешний ключ на тип канала
                                address VARCHAR(100) NOT NULL CONSTRAINT address_length CHECK (char_length(address) >= 3), -- адрес доставки (ex: @white_ak111), длинна больше или равна
                                deleted BOOLEAN DEFAULT FALSE -- флаг удаления
);

CREATE TABLE notifications (
                               id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
                               snippet_id INT NOT NULL REFERENCES snippets (id), -- внешний ключ на заметку
                               user_channel_id INT NOT NULL REFERENCES users_channels (id), -- внешний ключ на канал связи пользователя
                               created TIMESTAMPTZ NOT NULL, -- дата\время создания
                               time_to_send TIMESTAMPTZ NOT NULL, -- дата\время назначенной отправки уведомления
                               sended  BOOLEAN DEFAULT FALSE, -- флаг отправки уведомления
                               deleted BOOLEAN DEFAULT FALSE -- флаг удаления
);

CREATE TABLE logs (
                      id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- id, первичный ключ
                      log_body JSON NOT NULL, -- сообщение лога в формате json
                      created TIMESTAMPTZ NOT NULL -- -- дата\время создания
);