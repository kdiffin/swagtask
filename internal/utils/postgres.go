package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func Int32ToPgInt4(i int32) pgtype.Int4 {
	var pgi pgtype.Int4
	pgi.Int32 = i
	pgi.Valid = true

	return pgi
}

func StringToPgText(str string) pgtype.Text {
	var pgs pgtype.Text
	if str == "" {
		pgs.Valid = false
	} else {
		pgs.String = str
		pgs.Valid = true
	}
	return pgs
}

func BrowserFormattedtTime(t pgtype.Timestamp) string {

	return t.Time.UTC().Format(time.RFC3339)
}
func PgUUID(str string) pgtype.UUID {
	var pgs pgtype.UUID
	if str == "" {
		pgs.Valid = false
		return pgs
	}

	parsed, err := uuid.Parse(str)
	if err != nil {
		pgs.Valid = false
		return pgs
	}

	copy(pgs.Bytes[:], parsed[:])
	pgs.Valid = true
	return pgs
}
