
create table first_entity (
  a     varchar(64)     primary key not null,
  b     varchar(64),
  c     integer,
  e     integer
);

create table second_entity (
  x     varchar(64)     primary key not null,
  z     integer
);

create table first_entity_r_second_entity (
  a     varchar(64)     not null references first_entity,
  x     varchar(64)     not null references second_entity,
  primary key(a, x)
);

create table third_entity (
  z     varchar(64)     primary key not null,
  a     varchar(64),
  b     integer,
  c     bigint,
  d     boolean,
  e     double precision,
  f     bytea,
  g     timestamp with time zone,
  h     uuid,
  i     varchar(26)
);

create table fourth_entity (
  x     varchar(64)     primary key not null,
  z     integer
);
