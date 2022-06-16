package registration

import (
	"time"

	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

// GetTime turns Timestamp into Time
func GetTime(timestamp *timestamp.Timestamp) time.Time {
	return time.Unix(timestamp.Seconds, int64(timestamp.Nanos))
}

// GetTimestamp turns Time into Timestamp
func GetTimestamp(time time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.Nanosecond()),
	}
}
