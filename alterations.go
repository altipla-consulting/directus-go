package directus

type Alterations[T any, PK string | int64] struct {
	Create []T  `json:"create,omitempty"`
	Update []T  `json:"update,omitempty"`
	Delete []PK `json:"delete,omitempty"`
}
