CREATE TABLE "refresh_tokens" (
    id serial NOT NULL UNIQUE,
    user_id int REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    token varchar(255) NOT NULL UNIQUE,
    expires_at timestamp NOT NULL
)