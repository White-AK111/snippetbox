-- получить данные пользователя при его аутентификации
SELECT id, login, name, email, confirmed
FROM users
WHERE login = 'mike111'
AND hashed_password = decode('013d7d16d7ad4fefb61bd95b765c8ceb', 'hex');

-- получить TOP10 заметок пользователя user_id = 1, срок заметки не истёк, заметка не удалена
SELECT id, title, content, created, expires, changed 
FROM snippets 
WHERE user_id = 1 
AND expires > current_timestamp 
AND deleted = FALSE
ORDER BY changed DESC LIMIT 10;

-- получить данные по заметке с id = 1 пользователя user_id = 1, срок заметки не истёк, заметка не удалена
SELECT id, title, content, created, expires, changed 
FROM snippets 
WHERE id = 1 
AND user_id = 1 
AND expires > current_timestamp  
AND deleted = FALSE;

-- получить все не отправленные и не удалённые уведомления\напоминания 
SELECT id, snippet_id, user_channel_id, created, time_to_send
FROM notifications
WHERE sended = FALSE
AND deleted = FALSE;