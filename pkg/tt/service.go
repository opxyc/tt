package tt

type ListFilters struct {
	Status []ActivityStatus
	Date   string // in the format YYYY-mm-dd
}

type TTService interface {
	Start(title, desc, tags string) (*Activity, error)
	Stop(activityID string) (*Activity, error)
	List(filters *ListFilters) ([]Activity, error)
	Pause(activityID string) (*Activity, error)
	Resume(activityID string) (*Activity, error)
	Delete(activiyID string) (*Activity, error)
}
