create extension if not exists "uuid-ossp";

create table if not exists users (
  id uuid primary key default uuid_generate_v4(),
  name text not null,
  username varchar(20) unique not null,
  bio text,
  password text not null,
  profile_picture text,
  created_at date default now()
);

create table if not exists attachments (
  id uuid primary key default uuid_generate_v4(),
  document_type text not null,
  url text not null
);

create table if not exists follows (
  follower_id uuid not null references users(id),
  following_id uuid not null references users(id),
  primary key(follower_id, following_id)
);

create table if not exists tweets (
  id uuid primary key default uuid_generate_v4(),
  content text not null,
  author_id uuid not null references users(id),
  attachment uuid references attachments(id),
  created_at date default now(),
  updated_at date
);