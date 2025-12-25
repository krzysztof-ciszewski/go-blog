ALTER TABLE posts ADD COLUMN author_id UUID;

CREATE INDEX idx_posts_author_id ON posts(author_id);

ALTER TABLE posts DROP COLUMN author;
ALTER TABLE posts ADD CONSTRAINT fk_posts_author_id FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL;