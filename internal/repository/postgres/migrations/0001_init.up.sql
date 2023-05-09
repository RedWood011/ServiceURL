
create table if not exists schema_migrations
(
    version bigint  not null,
    dirty   boolean not null,
    primary key (version)
    );

CREATE TABLE IF NOT EXISTS public.urls
(
    id            serial primary key,
    user_id       varchar,
    original_url  varchar,
    short_url     varchar,
    is_deleted    boolean  NOT NULL DEFAULT FALSE,
    CONSTRAINT original_unique UNIQUE (original_url)
);
