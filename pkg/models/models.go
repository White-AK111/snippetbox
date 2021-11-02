package models

import (
	"errors"
	"github.com/jackc/pgtype"
)

var ErrNoRecord = errors.New("models: no records found")

type (
	SnippetID      int
	UserID         int
	ChannelTypeID  int
	UserChannelID  int
	NotificationID int
	LogID          int
)

type (
	Address string
	Email   string
)

// Snippet data model
type Snippet struct {
	Id      SnippetID          `pgx:"id"`      // id, первичный ключ
	UserId  UserID             `pgx:"user_id"` // внешний ключ на пользователя
	Title   string             `pgx:"title"`   // заголовок
	Content string             `pgx:"content"` // содержимое
	Created pgtype.Timestamptz `pgx:"created"` // дата\время создания
	Expires pgtype.Timestamptz `pgx:"expires"` // дата\время истечения
	Changed pgtype.Timestamptz `pgx:"changed"` // дата\время изменения
	Deleted bool               `pgx:"deleted"` // флаг удаления
}

// User data model
type User struct {
	Id             UserID             `pgx:"id"`              // id, первичный ключ
	Name           string             `pgx:"name"`            // отображаемое имя
	Login          string             `pgx:"login"`           // логин, длинна больше или равна 3
	Email          Email              `pgx:"email"`           // почта, длинна больше или равна 4
	HashedPassword pgtype.Bytea       `pgx:"hashed_password"` // хэш пароля
	Created        pgtype.Timestamptz `pgx:"created"`         // дата\время создания
	Confirmed      bool               `pgx:"confirmed"`       // флаг указывающий что УЗ пользователя подтверждена с помощью email
}

// ChannelType data model
type ChannelType struct {
	Id      ChannelTypeID `pgx:"id"`      // id, первичный ключ
	Name    string        `pgx:"name"`    // наименование канала (ex: Telegram)
	Deleted bool          `pgx:"deleted"` // флаг удаления
}

// UserChannel data model
type UserChannel struct {
	Id            UserChannelID `pgx:"id"`              // id, первичный ключ
	UserId        UserID        `pgx:"user_id"`         // внешний ключ на пользователя
	ChannelTypeId ChannelTypeID `pgx:"channel_type_id"` // внешний ключ на тип канала
	Address       Address       `pgx:"address"`         // адрес доставки (ex: @white_ak111)
	Deleted       bool          `pgx:"deleted"`         // флаг удаления
}

// Notification data model
type Notification struct {
	Id            NotificationID     `pgx:"id"`              // id, первичный ключ
	SnippetId     SnippetID          `pgx:"snippet_id"`      // внешний ключ на заметку
	UserChannelId UserChannelID      `pgx:"user_channel_id"` // внешний ключ на канал связи пользователя
	Created       pgtype.Timestamptz `pgx:"created"`         // дата\время создания
	TimeToSend    pgtype.Timestamptz `pgx:"time_to_send"`    // дата\время назначенной отправки уведомления
	Sended        bool               `pgx:"sended"`          // флаг отправки уведомления
	Deleted       bool               `pgx:"deleted"`         // флаг удаления
}

// Log data model
type Log struct {
	Id      LogID              `pgx:"id"`       // id, первичный ключ
	LogBody pgtype.JSON        `pgx:"log_body"` // сообщение лога в формате json
	Created pgtype.Timestamptz `pgx:"created"`  // дата\время создания
}
