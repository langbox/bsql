package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	sql := "SELECT uuid, \"data\" #> '{tags}' AS tags FROM nodes WHERE  \"data\" -> 'tags' ??| array['?'] AND enabled = ?"
	s, _ := replacePlaceholders(sql, mini)
	fmt.Println(s)
}

func mini(buf *bytes.Buffer, i int) error {
	fmt.Fprintf(buf, "$%d", i)
	return nil
}

func replacePlaceholders(sql string, replace func(buf *bytes.Buffer, i int) error) (string, error) {
	buf := &bytes.Buffer{}
	i := 0
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // escape ?? => ?
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]
		} else {
			i++
			buf.WriteString(sql[:p])
			if err := replace(buf, i); err != nil {
				return "", err
			}
			sql = sql[p+1:]
		}
	}

	buf.WriteString(sql)
	return buf.String(), nil
}
