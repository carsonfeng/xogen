package models

import "time"

// DemoTeamModel1 DemoTeamModel1
type DemoTeamModel1 struct {
	Id              int       `xorm:"not null pk autoincr INT(11)"`
	Uid             int       `xorm:"int(11) unique"`
	Level           string    `xorm:"not null varchar(8) default '' index"`
	Type            string    `xorm:"text"`
	Scenes          string    `xorm:"not null varchar(128) default ''"`
	Team            string    `xorm:"not null varchar(32) default ''"`
	TeamLabel       string    `xorm:"not null varchar(8) default ''"`
	TeamLabelReason string    `xorm:"text"`
	Source          string    `xorm:"not null varchar(16) default '' index"`
	UpdateAt        time.Time `xorm:"updated"`
	CreateAt        time.Time `xorm:"created"`
	ValidAt         time.Time `xorm:"datetime default null index"`
}

// DemoUserModel2 DemoUserModel2
type DemoUserModel2 struct {
	Id    int    `xorm:"not null pk autoincr INT(11)"`
	Uid   int    `xorm:"int(11) unique"`
	Level string `xorm:"not null varchar(8) default '' index"`
}
