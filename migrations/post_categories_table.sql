CREATE TABLE IF NOT EXISTS post_categories (
	post_id TEXT,
	category_id TEXT,
	PRIMARY KEY (post_id, category_id),
	FOREIGN KEY (post_id) REFERENCES posts(id),
	FOREIGN KEY (category_id) REFERENCES categories(id)
);