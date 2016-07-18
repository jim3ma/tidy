package user

//import (
//"time"

//ci "github.com/jim3mar/tidy/models/checkin"
//"gopkg.in/mgo.v2/bson"
//)

type Setting struct {
	IMGUploadJS string `bson:"img_uploadjs" json:"img_uploadjs"`
	//OldPassword    string        `bson:",omitempty" json:"old_password"`
	//NewPassword    string        `bson:",omitempty" json:"new_password"`
	Gender int `bson:"gender" json:"gender"`
}
