package prefix

type Prefix string

const (
	SegmentPrefix           Prefix = "s:"
	TreePrefix              Prefix = "t:"
	DictionaryPrefix        Prefix = "d:"
	DimensionPrefix         Prefix = "i:"
	ExemplarDataPrefix      Prefix = "v:"
	ExemplarTimestampPrefix Prefix = "t:"
)

func (p Prefix) String() string { return string(p) }

func (p Prefix) Bytes() []byte { return []byte(p) }

func (p Prefix) Key(k string) []byte { return []byte(string(p) + k) }

func (p Prefix) Trim(k []byte) ([]byte, bool) {
	if len(k) > len(p) {
		return k[len(p):], true
	}
	return nil, false
}
