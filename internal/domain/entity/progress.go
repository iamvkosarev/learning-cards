package entity

type Mark int32

const (
	MARK_NULL = Mark(iota)
	MARK_A
	MARK_B
	MARK_C
	MARK_D
	MARK_E
)

type GroupProgress struct {
	Id   GroupId
	Mark Mark
}

type CardProgress struct {
	Id   CardId
	Mark Mark
}
