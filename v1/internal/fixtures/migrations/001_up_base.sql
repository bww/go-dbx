
create table test_entity (
  a     varchar(64)     primary key not null,
  b     varchar(64),
  c     integer
);

create table another_entity (
  x     varchar(64)     primary key not null,
  z     integer
);

create table test_entity_r_another_entity (
  a     varchar(64)     not null references test_entity,
  x     varchar(64)     not null references another_entity,
  primary key(a, x)
);
