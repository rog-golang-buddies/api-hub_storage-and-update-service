create table api_spec_docs
(
    id          bigserial primary key,
    title       text,
    description text,
    type        text,
    md5sum      text,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    fetched_at  timestamp with time zone
);

create index idx_api_spec_docs_deleted_at
    on api_spec_docs (deleted_at);

create table groups
(
    id              bigserial primary key,
    name            text,
    description     text,
    api_spec_doc_id int references api_spec_docs (id) on delete cascade on update cascade
);

create table api_methods
(
    id              bigserial primary key,
    path            text,
    name            text,
    description     text,
    type            text,
    parameters      text,
    request_body    text,
    api_spec_doc_id int references api_spec_docs (id) on delete cascade on update cascade,
    group_id        int references groups (id) on delete cascade on update cascade
        check (num_nonnulls(api_spec_doc_id, group_id) = 1)
);

create table external_docs
(
    id            bigserial primary key,
    description   text,
    url           text,
    api_method_id int unique references api_methods (id) on delete cascade on update cascade
);

create table servers
(
    id            bigserial primary key,
    url           text,
    description   text,
    api_method_id int references api_methods (id) on delete cascade on update cascade
);
