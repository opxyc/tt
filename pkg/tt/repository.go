package tt

type Repository interface {
	Find(activityID uint) (*Activity, error)
	Store(activity *Activity) error
	List(filters *ListFilters) ([]Activity, error)
	Update(activityID uint, activity *Activity) error
	Delete(activityID uint) error
}
