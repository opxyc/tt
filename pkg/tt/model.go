package tt

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ActivityStatus string

type TimeChunk struct {
	StartTime *time.Time
	EndTime   *time.Time
}
type ActivityTrack map[int]TimeChunk

type Activity struct {
	gorm.Model
	ID          uint `gorm:"primarykey"`
	Title       string
	Desc        string
	Tags        string
	Status      ActivityStatus // IN_PROGRESS, PAUSED, COMPLETE
	Track       datatypes.JSON
	Duration    float64 // in minutes
	CompletedAt *time.Time
}

/*
Track JSON will store the time track.
Each child object represents chunks of time.
For any task which was completed without taking a break will have only one child (0).
{
	0: {
		startTime: 9:30,
		endTime: 10.00
	},
	1: {
		startTime: 10.10,
		endTime: 10.30
	},
 	...
}
*/
