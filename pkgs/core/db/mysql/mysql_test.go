package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)

var dbCli *gorm.DB

type user struct {
	ID     int    `json:"id" gorm:"column:id"`
	Name   string `json:"name" gorm:"column:name"`
	UserID int    `json:"user_id" gorm:"column:user_id"`
	Card   card
}

type card struct {
	ID     int    `json:"id" gorm:"column:id"`
	Email  string `json:"email" gorm:"column:email"`
	UserID int    `json:"user_id" gorm:"column:user_id"`
}

type User struct {
	ID        int `json:"id" gorm:"column:id"`
	Profile   Profile
	ProfileID int `json:"profile_id" gorm:"column:profile_id"`
}

type User1 struct {
	ID        int       `json:"id" gorm:"column:id"`
	Profile   []Profile `gorm:"FOREIGNKEY:TID;ASSOCIATION_FOREIGNKEY:ID"`
	ProfileID int       `json:"profile_id" gorm:"column:profile_id"`
}

type Profile struct {
	ID   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
	TID  int    `json:"tid" gorm:"column:tid"`
}

func init() {
	nDB := MysqlFactory{}
	dbCli = nDB.Create("192.168.38.140:33061").NewDBClient().Debug().LogMode(true)
}

func TestMysqlDB_NewDBClient(t *testing.T) {
	nDB := MysqlFactory{}
	nDB.Create("192.168.38.140:33061").NewDBClient().Debug().LogMode(true)
}

func TestFineOne(t *testing.T) {
	u := &User{
		ProfileID: 2,
	}
	//p := &Profile{}
	//db := dbCli.Model(&u).Related(&u.Profile).Find(&u)
	// db.Model(&user).Related(&profile)
	dbCli.Find(&u)
	db := dbCli.Model(&u).Association("Profile").Find(&u.Profile)
	if db.Error != nil {
		t.Fatal(db.Error)
	}
	t.Logf("%+v", u)

	// 关联查询所有记录
	var list []User1
	err := dbCli.Debug().Table("users").Preload("Profile").Find(&list).Error
	if err != nil {
		fmt.Println("[Preload] Query all failed:", err)
		return
	}
	fmt.Printf("[Preload] List one:%+v\n", list)

}
