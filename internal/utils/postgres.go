package utils

import "github.com/jackc/pgx/v5/pgtype"

func Int32ToPgInt4(i int32) pgtype.Int4 {
	var pgi pgtype.Int4
	pgi.Int32 = i
	pgi.Valid = true

	return pgi
}

func StringtoPgText(str string) pgtype.Text {
	var pgs pgtype.Text
	if str == "" {
		pgs.Valid = false
	} else {
		pgs.String = str
		pgs.Valid = true
	}
	return pgs
}
