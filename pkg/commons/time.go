package commons

import (
	"time"
)

func MalaysiaTimeNow() time.Time {
	//init the loc
	loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	//set timezone,
	now := time.Now().In(loc)
	return now
}

// MalaysiaTime gets Malaysia time
func MalaysiaTime(t time.Time) time.Time {
	// Load required location
	location, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	return t.In(location)
}

// DaysElapsed find number of days elapsed given two days
func DaysElapsed(from time.Time, to time.Time) int64 {
	duration := (MalaysiaTime(to).Sub(MalaysiaTime(from))).Hours() / 24
	return int64(duration)
}

// // TimeToMilli converts time to millisecond
// func TimeToMilli(time time.Time) int64 {
// 	return MalaysiaTime(time).UnixNano() / 1000000
// }

// // MilliToTime converts millisecond to time
// func MilliToTime(milli int64) time.Time {
// 	return MalaysiaTime(time.Unix(0, milli*int64(time.Millisecond)))
// }

// TimeToMilli converts time to millisecond
func TimeToMilli(time time.Time) int64 {
	mt := MalaysiaTime(time)
	return mt.Unix()*1000 + int64(mt.Nanosecond()/1000000)
}

// MilliToTime converts millisecond to time
func MilliToTime(milli int64) time.Time {
	// fmt.Printf("\t%v\n", time.Unix(milli/1000, 0))
	return MalaysiaTime(time.Unix(milli/1000, (milli%1000)*int64(time.Millisecond)))
}

// DateStringToTime converts date string to time
func DateStringToTime(date string) (time.Time, error) {
	t, err := time.Parse("02012006", date)
	if err != nil {
		return time.Now(), err
	}
	t = t.Add(-8 * time.Hour)

	return MalaysiaTime(t), nil
}

// TimeToDateString timestamp to date string (ddMMyyyy)
func TimeToDateString(t time.Time) string {
	return MalaysiaTime(t).Format("02012006")
}
