CREATE TABLE app.users
(
    id         UUID PRIMARY KEY,
    email      TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE app.chats
(
    id              UUID PRIMARY KEY,
    title           TEXT        NOT NULL,
    user_id         UUID        NOT NULL REFERENCES app.users (id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_message_at TIMESTAMPTZ
);

CREATE TABLE app.messages
(
    id         UUID PRIMARY KEY,
    chat_id    UUID        NOT NULL REFERENCES app.chats (id),
    role       TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)