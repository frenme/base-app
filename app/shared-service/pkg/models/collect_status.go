package models

type CollectStatus int

const (
	Owned CollectStatus = iota
	Wishlist
	Waiting
)
