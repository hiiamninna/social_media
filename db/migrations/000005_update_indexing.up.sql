CREATE INDEX index_users_id ON users (id);
CREATE INDEX index_users_phone ON users USING HASH (phone);
CREATE INDEX index_users_email ON users USING HASH (email);
CREATE INDEX index_users_total_firend ON users USING BRIN (total_friend);

CREATE INDEX index_friends_user_id ON friends (user_id);
CREATE INDEX index_friends_added_by ON friends (added_by);
CREATE INDEX index_friends_created_at ON friends USING BRIN (created_at);

CREATE INDEX index_posts_id ON posts (id);
CREATE INDEX index_posts_user_id ON posts (user_id);
CREATE INDEX index_posts_post ON posts USING HASH (post);
-- CREATE INDEX index_posts_tags ON posts USING GIST (tags);
CREATE INDEX index_posts_created_at ON posts USING BRIN (created_at);

CREATE INDEX index_comments_id ON comments (id);
CREATE INDEX index_comments_user_id ON comments (user_id);
CREATE INDEX index_comments_post_id ON comments (post_id);
CREATE INDEX index_comments_comment ON comments USING HASH (comment);
CREATE INDEX index_comments_created_at ON comments USING BRIN (created_at);