DROP INDEX index_users_id;
DROP INDEX index_users_phone;
DROP INDEX index_users_email;
DROP INDEX index_users_total_firend;

DROP INDEX index_friends_user_id;
DROP INDEX index_friends_added_by;
DROP INDEX index_friends_created_at;

DROP INDEX index_posts_id;
DROP INDEX index_posts_user_id;
DROP INDEX index_posts_post;
-- DROP INDEX index_posts_tags;
DROP INDEX index_posts_created_at;

DROP INDEX index_comments_id;
DROP INDEX index_comments_user_id;
DROP INDEX index_comments_post_id;
DROP INDEX index_comments_comment;
DROP INDEX index_comments_created_at;