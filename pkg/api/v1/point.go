package v1

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X float64 `json:"lat"`
	Y float64 `json:"lng"`
}

func (p *Point) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "(%f %f)", p.X, p.Y)
	return buf.Bytes(), nil
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v %v)", p.X, p.Y)
}

func (p *Point) Scan(val interface{}) (err error) {
	if bb, ok := val.([]uint8); ok {
		tmp := bb[1 : len(bb)-1]
		coors := strings.Split(string(tmp[:]), ",")
		if p.X, err = strconv.ParseFloat(coors[0], 64); err != nil {
			return err
		}
		if p.Y, err = strconv.ParseFloat(coors[1], 64); err != nil {
			return err
		}
	}
	return nil
}