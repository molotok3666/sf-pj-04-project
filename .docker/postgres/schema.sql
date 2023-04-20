DROP TABLE IF EXISTS news;

CREATE TABLE news (
    id   bigserial
        constraint news_pk
            primary key,
    guid TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time BIGINT NOT NULL,
    link TEXT NOT NULL
);

CREATE UNIQUE INDEX news_guid_uindex
    on news (guid);

CREATE TABLE comments
(
    id         bigserial
        constraint comments_pk
            primary key,
    news_id  BIGINT NOT NULL
        constraint comments_news_id_fk
            references news (id)
            on update cascade on delete cascade,
    content    TEXT NOT NULL,
    parent_id  BIGINT default 0
);

CREATE UNIQUE INDEX comments_id_uindex
    on comments (id);

CREATE TABLE request_logs
(
    id bigserial
        constraint request_logs_pk
            primary key,
    ip          VARCHAR(50) not null,
    timestamp   INTEGER     not null,
    status_code VARCHAR(3)  not null,
    request_id  TEXT        not null
);

create unique index request_logs_id_uindex
    on request_logs (id);
