package tt

type Repository interface {
	Find(activityID string) (*Activity, error)
	Store(activity *Activity) error
	List(filters *ListFilters) ([]Activity, error)
	Update(activityID string, activity *Activity) error
	Delete(activityID string) error
}
