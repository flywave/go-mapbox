package sprite

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidSubFile = errors.New("invalid sub file")
)

type StyleFile struct {
	Type       string
	Ext        string
	PixelRatio int
}

func (f *StyleFile) UnmarshalText(text []byte) error {
	sf, err := ParseStyleFile(string(text))
	if err != nil {
		return err
	}
	*f = *sf
	return nil
}

func (f StyleFile) MarshalText() (text []byte, err error) {
	return []byte(f.String()), nil
}

func (f StyleFile) String() string {
	if f.PixelRatio <= 1 {
		return f.Type + f.Ext
	}
	return f.Type + "@" + strconv.FormatInt(int64(f.PixelRatio), 10) + "x" + f.Ext
}

var reSubFile = regexp.MustCompile(`([^@]+)(@([0-9]+)x)?(\.[A-Za-z0-9]+)$`)

func ParseStyleFile(s string) (*StyleFile, error) {
	if !reSubFile.MatchString(s) {
		return nil, ErrInvalidSubFile
	}

	matched := reSubFile.FindAllStringSubmatch(s, -1)

	m := &StyleFile{
		Type:       strings.ToLower(matched[0][1]),
		Ext:        strings.ToLower(matched[0][4]),
		PixelRatio: 1,
	}

	if matched[0][3] != "" {
		if i, err := strconv.ParseInt(matched[0][3], 10, 64); err == nil {
			m.PixelRatio = int(i)
		}
	}

	return m, nil
}
