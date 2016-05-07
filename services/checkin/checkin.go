package checkin

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/checkin"
	//muser "github.com/jim3mar/tidy/models/user"
	svcuser "github.com/jim3mar/tidy/services/user"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CheckInResource struct {
	Mongo        *mgo.Session
	CollCI       *mgo.Collection
	CollUser     *mgo.Collection
	UserResource *svcuser.UserResource
}

func (cr *CheckInResource) Init(session *mgo.Session) {
	db := viper.GetString("mongo.db")
	cr.Mongo = session
	cr.CollCI = cr.Mongo.DB(db).C("checkin")
	cr.CollUser = cr.Mongo.DB(db).C("user")
}

// CheckIn do checkin task for special user id
// Method: POST
func (cr *CheckInResource) CheckIn(c *gin.Context) {
	now := time.Now()
	content := c.PostForm("content")
	uidString := c.PostForm("uid")
	username := c.PostForm("user_name")
	//log.Printf("Username: %s", username)
	img := c.PostForm("img")
	log.Print("Checkin user_id: " + uidString)
	uid := bson.ObjectIdHex(uidString)
	userinfo, err := cr.UserResource.QueryUserInfoByID(uidString)
	log.Printf("User info: %s", userin)
	if err != nil {
		panic(err)
	}
	if !userinfo.CanCheckIn() {
		log.Printf("user_id: %s, already checkin", uidString)
		c.JSON(http.StatusForbidden, "Already checkin today")
		return
	}
	ciData := &mod.CheckIn{
		Id_:         bson.NewObjectId(),
		UserId:      uid,
		UserName:    username,
		Content:     content,
		CreateAt:    now,
		CreateDay:   now.Day(),
		CreateMonth: int(now.Month()),
		CreateYear:  now.Year(),
		CreateHour:  now.Hour(),
		CreateMin:   now.Minute(),
		CreateSec:   now.Second(),
		Timestamp:   now.Unix(),
		Images:      strings.Split(img, "|"),
	}
	//log.Printf("Checkin content: %s", *ciData)
	err = cr.CollCI.Insert(ciData)
	if err != nil {
		panic(err)
	}
	err = cr.CollUser.Update(
		bson.M{
			"_id": uid,
		},
		bson.M{
			"$inc": bson.M{
				"continuous": 1,
			},
			"$set": bson.M{
				"last_checkin": *ciData,
			},
		})

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, now.Unix())
}

// ListCheckIn type
const (
	ListPersonal = iota
	ListAll
	ListSpecial
)

// ListCheckIn return all checkin records
// Method: GET
func (cr *CheckInResource) ListCheckIn(c *gin.Context) {
	//col := cr.Mongo.DB("tidy").C("checkin")
	//uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
	//objectId := bson.ObjectIdHex(id)
	//log.Print("user_id: " + uid)
	var ci []mod.CheckIn
	//col.Find(bson.M{"user_id": uid}).All(&ci)
	tp, err := strconv.Atoi(c.DefaultQuery("type", strconv.Itoa(ListPersonal)))
	if err != nil {
		tp = ListPersonal
	}
	switch tp {
	case ListPersonal:
		uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
		cr.CollCI.Find(bson.M{"user_id": uid}).All(&ci)
	case ListAll:
		cr.CollCI.Find(bson.M{}).All(&ci)
	case ListSpecial:
		spUID := bson.ObjectIdHex(c.DefaultQuery("special_uid", ""))
		cr.CollCI.Find(bson.M{"user_id": spUID}).All(&ci)
	default:
		uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
		cr.CollCI.Find(bson.M{"user_id": uid}).All(&ci)
	}
	//col.Find(nil).All(&ci)
	//log.Printf("%s", ci)
	c.JSON(http.StatusOK, ci)
}

func (cr *CheckInResource) UploadImg(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	//file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	//log.Println(header)
	//filename split
	fns := strings.Split(header.Filename, ".")
	log.Println(fns)
	fileext := "png"
	if l := len(fns); l >= 2 {
		fileext = fns[l-1]
	}
	guid := bson.NewObjectId().Hex()
	filename := guid + "." + fileext
	log.Println(filename)
	//fmt.Println(header.Filename)
	out, err := os.Create("./tmp/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, struct {
		GUID string `json:"guid"`
		Ext  string `json:"ext"`
	}{
		GUID: guid,
		Ext:  fileext,
	})
}
