# usql

A universal command-line interface for PostgreSQL, MySQL, Oracle Database,
SQLite3, Microsoft SQL Server, [and other SQL databases][Database Support].

[Installing][] | [Using][] | [Commands][] | [Building][] | [Database Support][] | [Releases][]

[Installing]: #installing (Installing)
[Using]: #using (Using)
[Commands]: #commands (Commands)
[Building]: #building (Building)
[Database Support]: #database-support (Database Support)
[Releases]: https://github.com/xo/usql/releases (Releases)

## Overview

`usql` provides a simple way of working with SQL databases via a command-line
inspired by PostgreSQL's `psql` tool and has a few additional features that
`psql` does not, such as syntax highlighting and context-based completion.

Database administrators and developers that would prefer to work with
non-PostgreSQL databases with a tool like `psql`, will find `usql` intuitive,
easy-to-use, and a great replacement for the command-line clients/tools
available for other databases.

## Installing

`usql` can be installed by [via Release][], [via Go][], or [via Homebrew][]:

[via Release]: #installing-via-release
[via Go]: #installing-via-go
[via Homebrew]: #installing-via-homebrew-osx

### Installing via Release

1. [Download a release for your platform][Releases]
2. Extract the `usql` or `usql.exe` file from the `.tar.bz2` or `.zip` file
3. Move the executable to somewhere on your `$PATH` (Linux/OSX) or `%PATH%` (Windows)

### Installing via Go

`usql` can be installed in the usual Go fashion:

```sh
# install usql with basic database support (includes PosgreSQL, MySQL, SQLite3, and MS SQL drivers)
$ go get -u github.com/xo/usql
```

Support for additional databases can be specified with [build tags][Database Support]:

```sh
# install usql with most drivers (excludes drivers requiring CGO)
$ go get -u -tags most github.com/xo/usql

# install usql with all drivers (includes drivers requiring CGO, namely Oracle and ODBC drivers)
$ go get -u -tags all github.com/xo/usql
```

### Installing via Homebrew (OSX)

`usql` is available in the [`xo/xo`][xo-tap] tap, and can be installed in the
usual way with the [`brew` command][homebrew]:

```sh
# add tap
$ brew tap xo/xo

# install usql with "most" drivers
$ brew install usql
```

Additional support for [Oracle and ODBC databases][Database Support] can be
installed by passing `--with-*` parameters during install:

```sh
# install usql with oracle and odbc support
$ brew install --with-oracle --with-odbc usql
```

Please note that Oracle support requires using the [`xo/xo` tap's][xo-tap]
`instantclient-sdk` formula. Any other `instantclient-sdk` formulae or older
versions of the Oracle Instant Client SDK [should be uninstalled][xo-tap-notes]
prior to attempting the above:

```sh
# uninstall the instantclient-sdk formula
$ brew uninstall InstantClientTap/instantclient/instantclient-sdk

# remove conflicting tap
$ brew untap InstantClientTap/instantclient
```

## Using

After [installing][Installing], `usql` can be used similarly to the following:

```sh
# connect to a postgres database
$ usql postgres://booktest@localhost/booktest

# connect to an oracle database
$ usql oracle://user:pass@host/oracle.sid

# connect to a postgres database and run script.sql
$ usql pg://localhost/ -f script.sql
```

### Connecting to a Database

`usql` opens a database connection by [parsing a URL][dburl] and passing the
resulting connection string to the [database driver][Database Support].
Database connection strings (aka "data source name" or DSNs) have the same
parsing rules as URLs, and can be passed to `usql` via command-line, or to the
`\connect` or `\c` [commands][Commands].

`usql` connection strings look like the following:

```txt
   driver+transport://user:pass@host/dbname?opt1=a&opt2=b
   driver:/path/to/file
   /path/to/file
```

Where the above are:

| Component          | Description                                                               |
|--------------------|---------------------------------------------------------------------------|
| driver             | driver name or alias                                                      |
| transport          | `tcp`, `udp`, `unix` or driver name <i>(for ODBC and ADODB)</i>           |
| user               | username                                                                  |
| pass               | password                                                                  |
| host               | hostname                                                                  |
| dbname<sup>*</sup> | database name, instance, or service name/ID                               |
| ?opt1=a&...        | database driver options (see respective SQL driver for available options) |
| /path/to/file      | a path on disk                                                            |

<i><sup><b>*</b></sup> for Microsoft SQL Server, the syntax to supply an
instance and database name is `/instance/dbname`, where `/instance` is
optional. For Oracle databases, `/dbname` is the unique database ID (SID).</i>

#### Driver Aliases

`usql` supports the same driver names and aliases from the [`dburl`][dburl]
package. Most databases have at least one or more alias - please refer to the
[`dburl` documentation][dburl-schemes] for all supported aliases.

##### Short Aliases

All database drivers have a two character short form that is usually the first
two letters of the database driver. For example, `pg` for `postgres`, `my` for
`mysql`, `ms` for `mssql`, `or` for `oracle`, or `sq` for `sqlite3`.

#### Passing Driver Options

Driver options are specified as standard URL query options in the form of
`?opt1=a&obt2=b`. Please refer to the [relevant database driver's
documentation][Database Support] for available options.

#### Paths on Disk

If a URL does not have a `driver:` scheme, `usql` will check if it is a path on
disk. If the path exists, `usql` will attempt to use an appropriate database
driver to open the path.

If the specified path is a Unix Domain Socket, `usql` will attempt to open it
using the MySQL driver. If the path is a directory, `usql` will attempt to open
it using the PostgreSQL driver. If the path is a regular file, `usql` will
attempt to open the file using the SQLite3 driver.

### Connection String Examples

The following are example connection strings and additional ways to connect to
databases with `usql`:

```sh
# connect to a postgres database
$ usql pg://user:pass@localhost/dbname
$ usql pgsql://user:pass@localhost/dbname
$ usql postgres://user:pass@localhost:port/dbname
$ psql /var/run/postgresql

# connect to a mysql database
$ usql my://user:pass@localhost/dbname
$ usql mysql://user:pass@localhost:port/dbname
$ usql /var/run/mysqld/mysqld.sock

# connect to a mssql (Microsoft SQL) database
$ usql ms://user:pass@localhost/dbname
$ usql mssql://user:pass@localhost:port/dbname

# connect to a mssql (Microsoft SQL) database using Windows domain authentication
$ runas /user:ACME\wiley /netonly "usql mssql://host/dbname/"

# connect to a oracle database
$ usql or://user:pass@localhost/dbname
$ usql oracle://user:pass@localhost:port/dbname

# connect to a sqlite database that exists on disk
$ usql dbname.sqlite3

# NOTE: when the "<driver>://" or "<driver>:" scheme is omitted, the file must
# already exist on disk.
#
# if the file does not yet exist, the URL must incorporate file:, sq:, sqlite3:,
# or any other recognized sqlite3 driver alias to force usql to create a new,
# empty database at the specified path:
$ usql sq://path/to/dbname.sqlite3
$ usql sqlite3://path/to/dbname.sqlite3
$ usql file:/path/to/dbname.sqlite3

# connect to a adodb ole resource (windows only)
$ usql adodb://Microsoft.Jet.OLEDB.4.0/myfile.mdb
$ usql "adodb://Microsoft.ACE.OLEDB.12.0/?Extended+Properties=\"Text;HDR=NO;FMT=Delimited\""
```

### Executing Queries and Commands

`usql` provides a command intrepreter that intreprets `\ ` [commands][Commands]
and sends SQL queries to a database similar to `psql`:

```sh
$ usql sqlite://example.sqlite3
Connected with driver sqlite3 (SQLite3 3.17.0)
Type "help" for help.

sq:example.sqlite3=> create table test (test_id int, name string);
CREATE TABLE
sq:example.sqlite3=> insert into test (test_id, name) values (1, 'hello');
INSERT 1
sq:example.sqlite3=> select * from test;
  test_id | name
+---------+-------+
        1 | hello
(1 rows)

sq:example.sqlite3=> select * from test
sq:example.sqlite3-> \p
select * from test
sq:example.sqlite3-> \g
  test_id | name
+---------+-------+
        1 | hello
(1 rows)

sq:example.sqlite3=> \c postgres://booktest@localhost
error: pq: 28P01: password authentication failed for user "booktest"
Enter password:
Connected with driver postgres (PostgreSQL 9.6.6)
pg:booktest@localhost=> select * from authors;
  author_id |      name
+-----------+----------------+
          1 | Unknown Master
          2 | blah
          3 | aoeu
(3 rows)

pg:booktest@localhost=>
```

## Commands

`usql` recognizes backslash (`\ `) commands similar to `psql`. Currently
available commands:

```sh
$ usql
Type "help" for help.

(not connected)=> \?
General
  \q                    quit usql
  \copyright            show usql usage and distribution terms
  \drivers              display information about available database drivers
  \g [FILE] or ;        execute query (and send results to file or |pipe)
  \gexec                execute query and execute each value of the result
  \gset [PREFIX]        execute query and store results in usql variables

Help
  \? [commands]         show help on backslash commands
  \? options            show help on usql command-line options
  \? variables          show help on special variables

Query Buffer
  \e [FILE] [LINE]      edit the query buffer (or file) with external editor
  \p                    show the contents of the query buffer
  \r                    reset (clear) the query buffer
  \w FILE               write query buffer to file

Input/Output
  \echo [STRING]        write string to standard output
  \i FILE               execute commands from file
  \ir FILE              as \i, but relative to location of current script

Transaction
  \begin                begin a transaction
  \commit               commit current transaction
  \rollback             rollback (abort) current transaction

Connection
  \c URL                connect to database with url
  \c DRIVER PARAMS...   connect to database with SQL driver and parameters
  \Z                    close database connection
  \password [USERNAME]  change the password for a user
  \conninfo             display information about the current database connection

Operating System
  \cd [DIR]             change the current working directory
  \setenv NAME [VALUE]  set or unset environment variable
  \! [COMMAND]          execute command in shell or start interactive shell

Variables
  \prompt [TEXT] NAME   prompt user to set internal variable
  \set [NAME [VALUE]]   set internal variable, or list all if no parameters
  \unset NAME           unset (delete) internal variable
```

The `usql` project's goal is to support all standard `psql` commands. Pull
Requests would be greatly appreciated!

## Building

When building `usql` with [Go][go-project], only drivers for PostgreSQL, MySQL,
SQLite3 and Microsoft SQL Server will be enabled by default. Other databases
can be enabled by specifying the build tag for their [database driver][Database Support].
Additionally, the `most` and `all` build tags include most, and all SQL
drivers, respectively:

```sh
# install all drivers
$ go get -u -tags all github.com/xo/usql

# install with most drivers (same as all but excludes Oracle/ODBC)
$ go get -u -tags most github.com/xo/usql

# install with base drivers and Oracle/ODBC support
$ go get -u -tags 'oracle odbc' github.com/xo/usql
```

For every build tag `<driver>`, there is also the `no_<driver>` build tag
disabling the driver:

```sh
# install all drivers excluding avatica and couchbase
$ go get -u -tags 'all no_avatica no_couchbase' github.com/xo/usql
```

### Release Builds

[Release builds][Releases] are built with the `most` build tag. [Additional
SQLite3 build tags](contrib/build-release.sh) are also specified for releases.

### Using as a Package

An effort has been made to keep `usql`'s packages modular, and reusable by
other developers wishing to leverage the `usql` code base. As such, it is
possible to build a SQL command-line interface (e.g, for use by some other
project as an "official" client) using the core `usql` source tree.

Please refer to [main.go](main.go) to see how `usql` puts together its
packages. `usql`'s code is also well-documented -- please refer to the [GoDoc
listing][godoc] to see the various APIs available.

## Database Support

`usql` works with all Go standard library compatible SQL drivers supported by
[`github.com/xo/dburl`][dburl].

The databases supported, the respective build tag, and the driver used by `usql` are:

| Driver               | Build Tag  | Driver Used                                                                      |
|----------------------|------------|----------------------------------------------------------------------------------|
| Microsoft SQL Server | mssql      | [github.com/denisenkom/go-mssqldb][d-mssql]                                      |
| MySQL                | mysql      | [github.com/go-sql-driver/mysql][d-mysql]                                        |
| PostgreSQL           | postgres   | [github.com/lib/pq][d-postgres]                                                  |
| SQLite3              | sqlite3    | [github.com/mattn/go-sqlite3][d-sqlite3]                                         |
| Oracle               | oracle     | [gopkg.in/rana/ora.v4][d-oracle]                                                 |
|                      |            |                                                                                  |
| MySQL                | mymysql    | [github.com/ziutek/mymysql/godrv][d-mymysql]                                     |
| PostgreSQL           | pgx        | [github.com/jackc/pgx/stdlib][d-pgx]                                             |
|                      |            |                                                                                  |
| Apache Avatica       | avatica    | [github.com/Boostport/avatica][d-avatica]                                        |
| Cassandra            | cassandra  | [github.com/MichaelS11/go-cql-driver][d-cassandra]                               |
| ClickHouse           | clickhouse | [github.com/kshvakov/clickhouse][d-clickhouse]                                   |
| Couchbase            | couchbase  | [github.com/couchbase/go_n1ql][d-couchbase]                                      |
| Cznic QL             | ql         | [github.com/cznic/ql][d-ql]                                                      |
| Firebird SQL         | firebird   | [github.com/nakagami/firebirdsql][d-firebird]                                    |
| Microsoft ADODB      | adodb      | [github.com/mattn/go-adodb][d-adodb]                                             |
| ODBC                 | odbc       | [github.com/alexbrainman/odbc][d-odbc]                                           |
| Presto               | presto     | [github.com/prestodb/presto-go-client/presto][d-presto]                          |
| SAP HANA             | hdb        | [github.com/SAP/go-hdb/driver][d-hdb]                                            |
| VoltDB               | voltdb     | [github.com/VoltDB/voltdb-client-go/voltdbclient][d-voltdb]                      |
|                      |            |                                                                                  |
| Google Spanner       | spanner    | github.com/xo/spanner (not yet public)                                           |
|                      |            |                                                                                  |
| **MOST DRIVERS**     | most       | all drivers excluding Oracle and ODBC (requires CGO and additional dependencies) |
| **ALL DRIVERS**      | all        | all drivers                                                                      |

## Related Projects

* [dburl][dburl] - a Go package providing a standard, URL style mechanism for parsing and opening database connection URLs
* [xo][xo] - a command-line tool to generate Go code from a database schema

## TODO

`usql` aims to eventually provide a drop-in replacement for PostgreSQL's `psql`
command. This is on-going -- an attempt has been made in good-faith to provide
support for the most frequently used aspects/features of `psql`. Compatability
(where possible) with `psql`, takes general development priority.

##### General
1. Fix multiline behavior to mimic psql properly (on arrow up/down through history)
2. PAGER
3. \qecho + \o support
4. fix table output / formatting
5. context-based completion (WIP)
6. \encoding and environment/command line options to set encoding of input (to
    convert to utf-8 before feeding to SQL driver) (how important is this ... ?)
7. better --help support/output cli, man pages

##### Command Processing + `psql` compatibility

1. the \j* commands (WIP)
2. \watch
3. \errverbose
4. formatting settings (\pset, \a, etc)
5. all \d* commands from `psql` (WIP, need to finish work extracting introspection code from `xo`)
6. remaining `psql` cli parameters

##### Not important / "Nice to haves" / extremely low priority:

1. correct operation of interweaved -f/-c commands, ie: -f 1 -c 1 -c 2 -f 2 -f 3 -c 3 runs in the specified order

##### Testing

1. test suite for databases, doing a minimal set of SELECT, INSERT, UPDATE, DELETE

##### Future Database Support

1. InfluxDB
2. CSV via SQLite3 vtable
3. Google Sheets via SQLite3 vtable
4. Atlassian JIRA JQL (why not? lol)
5. Google Spanner

[dburl]: https://github.com/xo/dburl
[dburl-schemes]: https://github.com/xo/dburl#protocol-schemes-and-aliases
[godoc]: https://godoc.org/github.com/xo/usql
[go-project]: https://golang.org/project
[homebrew]: https://brew.sh/
[xo]: https://github.com/xo/xo
[xo-tap]: https://github.com/xo/homebrew-xo
[xo-tap-notes]: https://github.com/xo/homebrew-xo#oracle-notes

[d-mssql]: https://github.com/denisenkom/go-mssqldb
[d-mysql]: https://github.com/go-sql-driver/mysql
[d-postgres]: https://github.com/lib/pq
[d-sqlite3]: https://github.com/mattn/go-sqlite3
[d-oracle]: https://gopkg.in/rana/ora.v4
[d-mymysql]: https://github.com/ziutek/mymysql
[d-pgx]: https://github.com/jackc/pgx
[d-avatica]: https://github.com/Boostport/avatica
[d-cassandra]: https://github.com/MichaelS11/go-cql-driver
[d-clickhouse]: https://github.com/kshvakov/clickhouse
[d-couchbase]: https://github.com/couchbase/go_n1ql
[d-ql]: https://github.com/cznic/ql
[d-firebird]: https://github.com/nakagami/firebirdsql
[d-adodb]: https://github.com/mattn/go-adodb
[d-odbc]: https://github.com/alexbrainman/odbc
[d-presto]: https://github.com/prestodb/presto-go-client
[d-hdb]: https://github.com/SAP/go-hdb
[d-sqlago]: https://github.com/a-palchikov/sqlago
[d-voltdb]: https://github.com/VoltDB/voltdb-client-go
[d-spanner]: https://github.com/xo/spanner
