package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/White-AK111/snippetbox/pkg/models"

	"github.com/go-redis/cache/v8"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SnippetModel - struct for pool connections sql.DB
type SnippetModel struct {
	DB    *pgxpool.Pool
	CTX   context.Context
	Cache *cache.Cache
}

// InsertSnippet insert snippet into DB
func (a *SnippetModel) InsertSnippet(snippet *models.Snippet) (models.SnippetID, error) {
	const sql = `
insert into snippets (user_id, title, content, created, expires, changed, deleted) values
	($1, $2, $3, $4, $5, $6, $7)
returning id;
`
	var id models.SnippetID
	err := a.DB.QueryRow(a.CTX, sql,
		snippet.UserId,
		snippet.Title,
		snippet.Content,
		snippet.Created,
		snippet.Expires,
		snippet.Changed,
		snippet.Deleted,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert snippet: %w", err)
	}

	key := string(id) + "#" + string(snippet.UserId)

	err = a.Cache.Set(&cache.Item{
		Ctx:   a.CTX,
		Key:   key,
		Value: snippet,
		TTL:   time.Hour,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to insert snippet: %w", err)
	}
	fmt.Printf("INFO\tStore in cache\n")

	return id, nil
}

// InsertUser insert user into DB
func (a *SnippetModel) InsertUser(user *models.User) (models.UserID, error) {
	const sql = `
INSERT into users (name, login, email, hashed_password, created, confirmed) VALUES 
($1, $2, $3, $4, $5, $6)
returning id;
`
	var id models.UserID
	err := a.DB.QueryRow(a.CTX, sql,
		user.Name,
		user.Login,
		user.Email,
		user.HashedPassword,
		user.Created,
		user.Confirmed,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

// GetSnippet get snippet by id
func (a *SnippetModel) GetSnippet(id uint, userId uint) (*models.Snippet, error) {

	snippet := &models.Snippet{}

	key := string(id) + "#" + string(userId)
	err := a.Cache.Get(a.CTX, key, snippet)

	switch err {
	case nil:
		fmt.Printf("INFO\tGET from cache\n")
		return snippet, nil
	case cache.ErrCacheMiss:
		const sql = `
SELECT id, user_id, title, content, created, expires, changed, deleted  
FROM snippets 
WHERE id = $1 
AND user_id = $2 
AND expires > current_timestamp  
AND deleted = FALSE;
`
		row := a.DB.QueryRow(a.CTX, sql, id, userId)
		err = row.Scan(&snippet.Id, &snippet.UserId, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires, &snippet.Changed, &snippet.Deleted)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}

		err = a.Cache.Set(&cache.Item{
			Ctx:   a.CTX,
			Key:   key,
			Value: snippet,
			TTL:   time.Hour,
		})

		if err != nil {
			return nil, err
		}

		fmt.Printf("INFO\tGET from DB\n")
		return snippet, nil
	default:
		return nil, err
	}
}

// LatestSnippets get last user snippets
func (a *SnippetModel) LatestSnippets(userId uint, limit uint) ([]*models.Snippet, error) {
	const sql = `
SELECT id, user_id, title, content, created, expires, changed, deleted  
FROM snippets 
WHERE user_id = $1 
AND expires > current_timestamp 
AND deleted = FALSE
ORDER BY changed DESC LIMIT $2;
`
	var snippets []*models.Snippet

	rows, err := a.DB.Query(a.CTX, sql, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		snippet := &models.Snippet{}

		err = rows.Scan(&snippet.Id, &snippet.UserId, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires, &snippet.Changed, &snippet.Deleted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		snippets = append(snippets, snippet)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}

	return snippets, nil
}

// GetUserByLogin get snippet by id
func (a *SnippetModel) GetUserByLogin(login string) (*models.User, error) {
	const sql = `
SELECT id, name, login, email, hashed_password, created, confirmed
FROM users
WHERE login = $1;
`
	row := a.DB.QueryRow(a.CTX, sql, login)
	user := &models.User{}

	err := row.Scan(&user.Id, &user.Name, &user.Login, &user.Email, &user.HashedPassword, &user.Created, &user.Confirmed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

// GetUserByLoginAndPassHash get snippet by id
func (a *SnippetModel) GetUserByLoginAndPassHash(login string, passHash string) (*models.User, error) {
	const sql = `
SELECT id, name, login, email, hashed_password, created, confirmed
FROM users
WHERE login = $1
AND hashed_password = decode($2, 'hex');
`
	row := a.DB.QueryRow(a.CTX, sql, login, passHash)
	user := &models.User{}

	err := row.Scan(&user.Id, &user.Name, &user.Login, &user.Email, &user.HashedPassword, &user.Created, &user.Confirmed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

// GetNotSendedNotifications get not send notifications
func (a *SnippetModel) GetNotSendedNotifications() ([]*models.Notification, error) {
	const sql = `
SELECT id, snippet_id, user_channel_id, created, time_to_send, sended, deleted
FROM notifications
WHERE sended = FALSE
AND deleted = FALSE;
`
	var notifications []*models.Notification

	rows, err := a.DB.Query(a.CTX, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		notification := &models.Notification{}

		err = rows.Scan(&notification.Id, &notification.SnippetId, &notification.UserChannelId, &notification.Created, &notification.TimeToSend, &notification.Sended, &notification.Deleted)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		notifications = append(notifications, notification)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}

	return notifications, nil
}
