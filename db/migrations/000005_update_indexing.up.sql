CREATE INDEX index_hash_user_id ON users USING HASH (id);
CREATE INDEX index_hash_user_phone ON users USING HASH (phone);
CREATE INDEX index_hash_user_email ON users USING HASH (email);

CREATE INDEX index_friend_user_id ON friends (user_id);
CREATE INDEX index_friend_added_by ON friends (added_by);

CREATE INDEX index_hash_post_id ON posts USING HASH (id);

CREATE INDEX index_comment_id ON comments (id);
CREATE INDEX index_hash_comment_post_id ON comments USING HASH (post_id);
