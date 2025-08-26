package models

type CollectStatus string

const (
	Owned    CollectStatus = "owned"
	Wishlist CollectStatus = "wishlist"
	Waiting  CollectStatus = "waiting"
)
