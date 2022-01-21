package game

import "bytes"

type Ans struct{ bytes [5]byte } // 0, 1, 2
func (ans Ans) String() string {
	var b bytes.Buffer
	for i := range ans.bytes {
		b.WriteByte(ans.bytes[i] + '0')
	}
	return b.String()
}
