[![Build Status](https://travis-ci.org/moznion/go-mysql-explain-parser.svg?branch=master)](https://travis-ci.org/moznion/go-mysql-explain-parser)

[https://godoc.org/github.com/moznion/go-mysql-explain-parser](https://godoc.org/github.com/moznion/go-mysql-explain-parser)

go-mysql-explain-parser
=======================

go-mysql-explain-parser is the parser for result of EXPLAIN of MySQL.

This package is port of [MySQL::Explain::Parser](http://search.cpan.org/~moznion/MySQL-Explain-Parser/lib/MySQL/Explain/Parser.pm) from Perl to Go.

## Getting Started

```
go get github.com/moznion/go-mysql-explain-parser
```

## Synopsis

```go
import (
	"github.com/moznion/go-mysql-explain-parser"
)

func main() {

	explain := `
+----+-------------+-------+-------+---------------+---------+---------+------+------+----------+-------------+
| id | select_type | table | type  | possible_keys | key     | key_len | ref  | rows | filtered | Extra       |
+----+-------------+-------+-------+---------------+---------+---------+------+------+----------+-------------+
|  1 | PRIMARY     | t1    | index | NULL          | PRIMARY | 4       | NULL | 4    | 100.00   |             |
|  2 | SUBQUERY    | t2    | index | a             | a       | 5       | NULL | 3    | 100.00   | Using index |
+----+-------------+-------+-------+---------------+---------+---------+------+------+----------+-------------+
`
	mysqlexplainparser.Parse(explain)

	explainVertical := `
*************************** 1. row ***************************
           id: 1
  select_type: PRIMARY
        table: t1
         type: index
possible_keys: NULL
          key: PRIMARY
      key_len: 4
          ref: NULL
         rows: 4
     filtered: 100.00
        Extra:
*************************** 2. row ***************************
           id: 2
  select_type: SUBQUERY
        table: t2
         type: index
possible_keys: a
          key: a
      key_len: 5
          ref: NULL
         rows: 3
     filtered: 100.00
        Extra: Using index
`
	mysqlexplainparser.ParseVertical(explainVertical)
}
```

## Functions

- `func Parse(explain string) []map[string]string`

Returns the result of parsed EXPLAIN

- `func ParseVertical(explain string) []map[string]string`

Returns the result of parsed EXPLAIN as vertical

## See Also

- [MySQL::Explain::Parser](http://search.cpan.org/~moznion/MySQL-Explain-Parser/lib/MySQL/Explain/Parser.pm)
- [http://dev.mysql.com/doc/en/explain-output.html](http://dev.mysql.com/doc/en/explain-output.html)
- [http://dev.mysql.com/doc/en/explain-extended.html](http://dev.mysql.com/doc/en/explain-extended.html)

## LICENSE

MIT
