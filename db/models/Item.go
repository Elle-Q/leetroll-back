package models

import (
	"time"
)

type Item struct {
	ID         int64
	Cat        *Category
	Name       string
	Desc       string
	Main       string
	Preview    string
	Attachment string
	Type       string
	BLink      string
	Tags       []string
	Price      float64
	Author     string
	DownCnt    int64
	Scores     float64
	Status     string
	CreateTime time.Time
	UpdateTime time.Time
}

func (i *Item) New() *Item {
	i.Cat = &Category{}
	return i
}

func MakeItem() *Item {
	item := Item{}
	item.Cat = &Category{}
	return &item
}
