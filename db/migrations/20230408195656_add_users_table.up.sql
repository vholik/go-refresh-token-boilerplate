CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "email" varchar NOT NULL UNIQUE,
    "password" varchar NOT NULL
)

