package game

import "bytes"

type Ans struct{ value int } // 0, 1, 2, encoded backwards
func (ans Ans) String() string {
	var b bytes.Buffer
	v := ans.value
	for i := 0; i < 5; i++ {
		b.WriteByte(byte(v%3) + '0')
		v /= 3
	}
	return b.String()
}

func (ans Ans) Equals(other Ans) bool {
	return ans.value == other.value
}

func FromBytes(ans [5]byte) Ans {
	val := 0
	l := len(ans) - 1
	for i := range ans {
		val = val*3 + int(ans[l-i])
	}
	return Ans{val}
}
