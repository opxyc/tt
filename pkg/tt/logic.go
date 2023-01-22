package tt

import (
	"encoding/json"
	"errors"
	"math"
	"time"
)

var ErrOtherActivityAlreadyInProgress = errors.New("other activity in progress")

type ttService struct {
	repo Repository
}

func (tts *ttService) Start(title, desc, tags string) (*Activity, error) {
	activitiesInProgress, err := tts.repo.List(&ListFilters{Status: []ActivityStatus{StatusInProgress}})
	if err != nil {
		return nil, err
	}

	if len(activitiesInProgress) > 0 {
		return nil, ErrOtherActivityAlreadyInProgress
	}

	now := time.Now()
	tc := TimeChunk{StartTime: &now}
	track := ActivityTrack{0: tc}
	trackJSON, err := json.Marshal(track)
	if err != nil {
		return nil, err
	}

	newActivity := &Activity{
		Title:    title,
		Desc:     desc,
		Tags:     tags,
		Status:   StatusInProgress,
		Track:    trackJSON,
		Duration: 0,
	}

	err = tts.repo.Store(newActivity)
	if err != nil {
		return nil, err
	}

	return newActivity, nil
}

func (tts *ttService) Stop(activityID uint) (*Activity, error) {
	activity, err := tts.repo.Find(activityID)
	if err != nil {
		return nil, err
	}

	track := ActivityTrack{}
	err = json.Unmarshal(activity.Track, &track)
	if err != nil {
		return nil, err
	}

	lastTimeChunk := track[len(track)-1]
	now := time.Now()
	track[len(track)-1] = TimeChunk{lastTimeChunk.StartTime, &now}
	trackJSON, err := json.Marshal(track)
	if err != nil {
		return nil, err
	}

	duration := 0.0
	for _, chunk := range track {
		duration += math.Round(chunk.EndTime.Sub(*chunk.StartTime).Minutes())
	}

	activity.Track = trackJSON
	activity.CompletedAt = &now
	activity.Duration = duration
	activity.Status = StatusComplete
	err = tts.repo.Update(activityID, activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (tts *ttService) List(filters *ListFilters) ([]Activity, error) {
	// TODO modify List method to accept `OR` conditions like status = COMPLETED OR status = PAUSED
	return tts.repo.List(filters)
}

func (tts *ttService) Pause(activityID uint) (*Activity, error) {
	activity, err := tts.repo.Find(activityID)
	if err != nil {
		return nil, err
	}

	track := ActivityTrack{}
	err = json.Unmarshal(activity.Track, &track)
	if err != nil {
		return nil, err
	}

	lastTimeChunk := track[len(track)-1]
	now := time.Now()
	track[len(track)-1] = TimeChunk{lastTimeChunk.StartTime, &now}
	trackJSON, err := json.Marshal(track)
	if err != nil {
		return nil, err
	}

	activity.Track = trackJSON
	activity.Status = StatusPaused
	err = tts.repo.Update(activityID, activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (tts *ttService) Resume(activityID uint) (*Activity, error) {
	activity, err := tts.repo.Find(activityID)
	if err != nil {
		return nil, err
	}

	track := ActivityTrack{}
	err = json.Unmarshal(activity.Track, &track)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	// add a new time chunk
	track[len(track)] = TimeChunk{StartTime: &now}
	trackJSON, err := json.Marshal(track)
	if err != nil {
		return nil, err
	}

	activity.Track = trackJSON
	activity.Status = StatusInProgress
	err = tts.repo.Update(activityID, activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (tts *ttService) Delete(activityID uint) (*Activity, error) {
	activity, err := tts.repo.Find(activityID)
	if err != nil {
		return nil, err
	}

	err = tts.repo.Delete(activityID)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func NewTTService(repo Repository) TTService {
	return &ttService{repo: repo}
}
