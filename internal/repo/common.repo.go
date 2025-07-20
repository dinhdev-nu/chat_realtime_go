package repo

import (
	"database/sql"
	"strings"
	"time"
	"unicode"
)

// Const
const LIMIT_USER = 6

// HandKeySearch
func RemoveVietnameseAccent(keySearch string) string {
	accent := map[rune]string{
		'á': "a", 'à': "a", 'ả': "a", 'ã': "a", 'ạ': "a", 'â': "a", 'ấ': "a", 'ầ': "a", 'ậ': "a", 'ă': "a",
		'é': "e", 'è': "e", 'ẻ': "e", 'ẽ': "e", 'ẹ': "e", 'ê': "e", 'ế': "e", 'ề': "e", 'ệ': "e",
		'í': "i", 'ì': "i", 'ỉ': "i", 'ĩ': "i", 'ị': "i",
		'ó': "o", 'ò': "o", 'ỏ': "o", 'õ': "o", 'ọ': "o", 'ô': "o", 'ố': "o", 'ồ': "o", 'ộ': "o", 'ơ': "o",
		'ú': "u", 'ù': "u", 'ủ': "u", 'ũ': "u", 'ụ': "u", 'ư': "u", 'ứ': "u", 'ừ': "u", 'ự': "u",
		'ý': "y", 'ỳ': "y", 'ỷ': "y", 'ỹ': "y", 'ỵ': "y",
		'đ': "d",
	}
	var result strings.Builder
	for _, char := range strings.ToLower(keySearch) {
		if r, ok := accent[char]; ok {
			result.WriteString(r)
		} else if unicode.IsLetter(char) || unicode.IsSpace(char) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func NullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func NullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: i != 0,
	}
}

func NullUint64(i uint64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(i),
		Valid: i != 0,
	}
}

func NullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func NullBool(b bool) bool {
	return b
}

func GetNameRoom(roomName string) [2]string {
	parts := strings.Split(roomName, " ")
	name2 := parts[2] + " " + parts[1] + " " + parts[0]

	return [2]string{roomName, name2}
}

func ConvertInt64ToUint64(i int64) uint64 {
	if i < 0 {
		return 0
	}
	return uint64(i)
}
