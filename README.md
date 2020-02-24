# DBX – DataBase eXtensions
DBX is a convenience layer built on [`database/sql`](https://golang.org/pkg/database/sql/) and [`sqlx`](https://github.com/jmoiron/sqlx) which aims to strike a very particular balance between automating SQL drudgery and telling you how to model or query your data.

## DBX is not an ORM (but it does some ORM stuff)
Among other things, DBX:

* Marshals fields from your structs into and out of a database,
* Provides mechanisms for storing and fetching graphs of entities,
* Makes very specific aspects of writing SQL easier.

Some of the things DBX **does not do** include:

* Requiring you to interact with a `Functional().Meta().Query().Language()`,
* Dictating any particular manner of modeling your data,
* Generating source code or SQL (mostly).

Oh yeah and it's only really tested with Postgres, cause that's the only relational database I care about. Sorry.

## It's Exampletown
Ok, so you've decided to completely rewrite your application to use DBX. Excellent decision. But where do you start? Let's talk about it.

### Act One, in which we meet the entity

Let's say we have a simple entity in our Go program that we want to persist as a row in Postgres. It's modeled as the following struct.
```go
type User struct {
  Id        string    `db:"id,pk"
  Username  string    `db:"username"
  Password  string    `db:"password"
  Notes     string    `db:"notes,omitempty"
  Created   time.Time `db:"created_at,omitempty"
}
```

We need to get this struct into this database table, which we have to write ourself, on purpose:
```sql
CREATE TABLE users (
  id          varchar(32)   primary key,
  username    varchar(32)   not null,
  password    varchar(64)   not null, -- bcrypt, because we're not dumb
  notes       text,
  created_at  timestamp with time zone not null default now()
);
```

Well, good news for us. Because _this is exactly the sort of thing DBX was created for!_ That was extremely lucky.

### Act Two, in which we persist the entity
What we need for this job is a _persister_. This is the high-level concept that deals with converting a struct to and from its database representation. Let's make one.

```go
pst := persist.New(
  db, // make a *sqlx.DB somehow, I'm not your mom
  entity.DefaultFieldMapper(),
  registry.DefaultRegistry(),
  ident.AlphaNumeric(16),
)
```
We'll discuss some of those parameters later, but for the moment, note the last one. A persister sometimes needs to generate primary keys in order to insert new entities. `ident.AlphaNumeric(16)` returns a function that generates random alpha-numeric strings 16 characters long.

There are a few common generators in the `ident` package that will create UUIDs, ULIDs, and random strings. If those don't meet your needs you can easily write your own.

Ok, we have our persister now. Let's store an instance of our `User` type.

```go
user := &User{
  Username: "cooldude",
  Password: "some long hash",
}

err := pst.Store("users", user, nil)
if err != nil {
  panic(err)
}
```

Now we have a database row like this:

```
| id | username | password | notes | created_at |
+----+----------+----------+-------+------------+
|jnjIYRgmCIC0oCUE | cooldude | some long hash | NULL | 2020-02-20 15:24:38.743665+00 |
```

As a special treat for us, before it persisted our entity, DBX used the identifier generator function we passed into our persister to create a new primary key for this record because it didn't have one. (If it did already have one, DBX would have performed an `UPDATE` using that key instead of an `INSERT`.)

So once the `Store` call succeeds we can reference `user.Id`, which will be populated with the persisted record's primary key.

### Act Three, in which we restore the entity
Alright, let's fetch it back now.

```go
dup := &User{}
err = pst.Fetch("users, &dup, user.Id)
if err != nil {
  panic(err)
}

assertEqual(user, dup) // "Cha, brah"
```
