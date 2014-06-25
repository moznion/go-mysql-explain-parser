package mysqlexplainparser

import (
	"testing"

	. "github.com/onsi/gomega"
)

func explain() string {
	return `+----+-------------+-------+-------+---------------+---------+---------+------+------+----------+-------------+
| id | select_type | table | type  | possible_keys | key     | key_len | ref  | rows | filtered | Extra       |
+----+-------------+-------+-------+---------------+---------+---------+------+------+----------+-------------+
|  1 | PRIMARY     | t1    | index | NULL          | PRIMARY | 4       | NULL | 4    | 100.00   |             |
|  2 | SUBQUERY    | t2    | index | a             | a       | 5       | NULL | 3    | 100.00   | Using index |
+----+-------------+-------+-------+---------------+---------+---------+------+------+----------+-------------+`
}

func expected() []map[string]string {
	return []map[string]string{
		{
			"id":            "1",
			"select_type":   "PRIMARY",
			"table":         "t1",
			"type":          "index",
			"possible_keys": "NULL",
			"key":           "PRIMARY",
			"key_len":       "4",
			"ref":           "NULL",
			"rows":          "4",
			"filtered":      "100.00",
			"Extra":         "",
		},
		{
			"id":            "2",
			"select_type":   "SUBQUERY",
			"table":         "t2",
			"type":          "index",
			"possible_keys": "a",
			"key":           "a",
			"key_len":       "5",
			"ref":           "NULL",
			"rows":          "3",
			"filtered":      "100.00",
			"Extra":         "Using index",
		},
	}
}

func TestParseBasicCase(t *testing.T) {
	RegisterTestingT(t)

	Expect(Parse(explain())).To(Equal(expected()))
}

func TestParseNoisyLinesAreOnBothEnds(t *testing.T) {
	RegisterTestingT(t)

	Expect(Parse("\n\n" + explain() + "\n2 rows in set, 1 warning (0.00 sec)\n\n")).To(Equal(expected()))
}

func TestParseMultiByteCharacter(t *testing.T) {
	RegisterTestingT(t)

	explain := `+----+-------------+-----------+-------+---------------+--------+---------+------+------+-------------+
| id | select_type | table     | type  | possible_keys | key    | key_len | ref  | rows | Extra       |
+----+-------------+-----------+-------+---------------+--------+---------+------+------+-------------+
|  1 | SIMPLE      | てすつ    | index | NULL          | ほげ   | 4       | NULL |    5 | Using index |
+----+-------------+-----------+-------+---------------+--------+---------+------+------+-------------+`

	expected := []map[string]string{
		{
			"id":            "1",
			"select_type":   "SIMPLE",
			"table":         "てすつ",
			"type":          "index",
			"possible_keys": "NULL",
			"key":           "ほげ",
			"key_len":       "4",
			"ref":           "NULL",
			"rows":          "5",
			"Extra":         "Using index",
		},
	}

	Expect(Parse(explain)).To(Equal(expected))
}

func explainVertical() string {
	return `*************************** 1. row ***************************
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
        Extra: Using index`
}

func TestParseVerticalBasicCase(t *testing.T) {
	RegisterTestingT(t)

	Expect(ParseVertical(explainVertical())).To(Equal(expected()))
}

func TestParseVerticalNoisyLinesAreOnBothEnds(t *testing.T) {
	RegisterTestingT(t)

	Expect(ParseVertical("\n\n" + explainVertical() + "\n2 rows in set, 1 warning (0.00 sec)\n\n")).To(Equal(expected()))
}

func TestParseVerticalMultiByteCharacter(t *testing.T) {
	RegisterTestingT(t)

	explainVertical := `*************************** 1. row ***************************
           id: 1
  select_type: SIMPLE
        table: てすつ
         type: index
possible_keys: NULL
          key: ほげ
      key_len: 4
          ref: NULL
         rows: 5
        Extra: Using index`

	expected := []map[string]string{
		{
			"id":            "1",
			"select_type":   "SIMPLE",
			"table":         "てすつ",
			"type":          "index",
			"possible_keys": "NULL",
			"key":           "ほげ",
			"key_len":       "4",
			"ref":           "NULL",
			"rows":          "5",
			"Extra":         "Using index",
		},
	}

	Expect(ParseVertical(explainVertical)).To(Equal(expected))
}
