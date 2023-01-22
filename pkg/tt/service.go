package tt

type ListFilters struct {
	Status []ActivityStatus
	Date   string // in the format YYYY-mm-dd
}

type TTService interface {
	Start(title, desc, tags string) (*Activity, error)
	Stop(activityID uint) (*Activity, error)
	List(filters *ListFilters) ([]Activity, error)
	Pause(activityID uint) (*Activity, error)
	Resume(activityID uint) (*Activity, error)
	Delete(activiyID uint) (*Activity, error)
}
