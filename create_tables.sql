create table ipadic (
  surface     varchar(256),
  original    varchar(256),
  reading     varchar(256)
);

create table neologd (
  surface     varchar(256),
  original    varchar(256),
  reading     varchar(256)
);


create table words (
  surface     varchar(256),
  original    varchar(256) primary key,
  reading     varchar(256)
);
create index reading_index on words(reading);

create table words_queue (
  surface     varchar(256),
  original    varchar(256),
  reading     varchar(256)
);

create table metadata (
  original    varchar(128) primary key,
  minimum_length integer
);






