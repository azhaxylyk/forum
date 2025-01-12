CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,                 -- Уникальный идентификатор уведомления
    user_id TEXT,                        -- ID пользователя, для которого создается уведомление
    action_by TEXT,                      -- ID пользователя, совершившего действие (например, поставившего лайк)
    action_type TEXT,                    -- Тип действия (например, "like", "dislike", "comment")
    target_id TEXT,                      -- ID объекта действия (например, поста или комментария)
    target_type TEXT,                    -- Тип объекта действия (например, "post" или "comment")
    is_read INTEGER DEFAULT 0,           -- Флаг, прочитано ли уведомление (0 — не прочитано, 1 — прочитано)
    created_at DATETIME,                 -- Время создания уведомления
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (action_by) REFERENCES users(id)
);
