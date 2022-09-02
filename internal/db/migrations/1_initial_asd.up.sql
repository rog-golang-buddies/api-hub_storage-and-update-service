create table api_spec_docs (
    id serial primary key ,
    title text,
    description text,
    type text,
    md5sum text,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    fetched_at timestamp
);

create table groups (
    id serial primary key,
    name text,
    description text,
    api_spec_docs_id int references api_spec_docs(id) on delete cascade on update cascade
);

create table api_methods (
  id serial primary key,
  path text,
  name text,
  description text,
  type text,
  parameters text,
  request_body text,
  api_spec_docs_id int references api_spec_docs(id) on delete cascade on update cascade,
  groups_id int references groups(id) on delete cascade on update cascade
);

create table external_docs (
    id serial primary key,
    description text,
    url text,
    api_methods_id int unique references api_methods(id) on delete cascade on update cascade
);

create table servers (
    id serial primary key,
    url text,
    description text,
    api_methods_id int references api_methods(id) on delete cascade on update cascade
);
