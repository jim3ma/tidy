package checkin

import (
	"io"
	//"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
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
	CollThumb    *mgo.Collection
	CollFavor    *mgo.Collection
	CollComments *mgo.Collection
	CollSComment *mgo.Collection
	UserResource *svcuser.UserResource
}

func (cr *CheckInResource) Init(session *mgo.Session) {
	db := viper.GetString("mongo.db")
	cr.Mongo = session
	cr.CollCI = cr.Mongo.DB(db).C("checkin")
	cr.CollUser = cr.Mongo.DB(db).C("user")
	cr.CollThumb = cr.Mongo.DB(db).C("ci_thumb")
	cr.CollFavor = cr.Mongo.DB(db).C("ci_favor")
	cr.CollComments = cr.Mongo.DB(db).C("ci_comments")
	cr.CollSComment = cr.Mongo.DB(db).C("ci_comment")
}

func (cr *CheckInResource) canEditCheckIn(uid bson.ObjectId, cid bson.ObjectId) bool {
	var ci mod.CheckIn
	cr.CollCI.Find(
		bson.M{
			"_id": cid,
		},
	).One(&ci)
	if ci.UserID == uid {
		return true
	}
	return false
}

// EditCheckIn will add a new checkin to replace old checkin
// Method: PUT
func (cr *CheckInResource) EditCheckIn(c *gin.Context) {
	// old checkin id
	cid := bson.ObjectIdHex(c.PostForm("cid"))

	content := c.PostForm("content")
	username := c.PostForm("user_name")
	img := c.PostForm("img")
	uid := bson.ObjectIdHex(c.PostForm("uid"))
	var pub bool
	if c.PostForm("pub") == "false" {
		pub = false
	} else {
		pub = true
	}

	log.Info("Checkin user_id: " + uid.Hex())

	userinfo, err := cr.UserResource.QueryUserInfoByID(uid.Hex())
	log.Infof("User info: %+v", userinfo)
	if err != nil {
		panic(err)
	}
	var ci mod.CheckIn
	cr.CollCI.Find(
		bson.M{
			"_id":     cid,
			"deleted": false,
		},
	).One(&ci)
	if ci.UserID != uid {
		log.Infof("user_id: %s, the owner of checkin isn't this user", uid)
		c.JSON(http.StatusForbidden, "Error user for this checkin")
		return
	}

	now := time.Now()
	ciData := &mod.CheckIn{
		ID:          bson.NewObjectId(),
		UserID:      uid,
		UserName:    username,
		Content:     content,
		CreateAt:    ci.CreateAt,
		CreateDay:   ci.CreateDay,
		CreateMonth: ci.CreateMonth,
		CreateYear:  ci.CreateYear,
		CreateHour:  ci.CreateHour,
		CreateMin:   ci.CreateMin,
		CreateSec:   ci.CreateSec,
		Timestamp:   ci.Timestamp,
		Images:      strings.Split(img, "|"),
		Deleted:     false,
		Public:      pub,
	}
	//log.Infof("Checkin content: %s", *ciData)
	// insert updated checkin
	err = cr.CollCI.Insert(ciData)
	if err != nil {
		panic(err)
	}
	log.Debug("edit checkin record: %+v", *ciData)

	// tag old checkin deleted
	err = cr.CollCI.Update(
		bson.M{
			"_id": cid,
		},
		bson.M{
			"$set": bson.M{
				"deleted": true,
			},
		})

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, now.Unix())
}

// CommentCheckIn add a new comment for special checkin
// Method: POST
func (cr *CheckInResource) CommentCheckIn(c *gin.Context) {
	uid := bson.ObjectIdHex(c.PostForm("uid"))
	cid := bson.ObjectIdHex(c.PostForm("cid"))
	content := c.PostForm("content")
	username := c.PostForm("user_name")

	now := time.Now()
	commendID := bson.NewObjectId()
	comment := &mod.SingleComment{
		ID:          commendID,
		UserID:      uid,
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
	}

	err := cr.CollSComment.Insert(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("upsert checkin comment error: %s", err)
		return
	}

	change, err := cr.updateComments(cid, uid, MgoCollAddSet, comment.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update checkin comments error: %s", err)
		return
	}
	log.Infof("Update comments, cid: %s, uid: %s, change info: %s", cid, uid, change)
	c.JSON(http.StatusOK, gin.H{"error": "none"})
}

// UnCommentCheckIn mark a new comment for special checkin
// Method: DELETE
func (cr *CheckInResource) UnCommentCheckIn(c *gin.Context) {
	uid := bson.ObjectIdHex(c.Query("uid"))
	cid := bson.ObjectIdHex(c.Query("cid"))
	commentID := bson.ObjectIdHex(c.Query("comment_id"))

	// check whether can delete the comment
	count, err := cr.CollSComment.Find(
		bson.M{
			"_id": cid,
			"uid": uid,
		}).Count()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	change, err := cr.updateComments(cid, uid, MgoCollPull, commentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update checkin comments error: %s", err)
		return
	}
	log.Infof("Update comments, cid: %s, uid: %s, comment id: %s, change info: %s", cid, uid, commentID, change)
	c.JSON(http.StatusOK, gin.H{"error": "none"})
}

func (cr *CheckInResource) updateComments(cid bson.ObjectId, uid bson.ObjectId, action string, commentID bson.ObjectId) (*mgo.ChangeInfo, error) {
	change, err := cr.CollComments.Upsert(
		bson.M{
			"_id": cid,
		},
		bson.M{
			action: bson.M{
				"comment_ids": commentID,
			},
		})
	return change, err
}

// ThumbCheckIn update checkin and thumb collection
// Method: POST
func (cr *CheckInResource) ThumbCheckIn(c *gin.Context) {
	uid := bson.ObjectIdHex(c.PostForm("uid"))
	cid := bson.ObjectIdHex(c.PostForm("cid"))

	///
	/// TBD check if user had already thumb the ci
	///
	err := cr.updateCIThumb(cid, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update checkin error: %s", err)
		return
	}

	change, err := cr.updateThumb(cid, uid, MgoCollAddSet)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update ci_thumb error: %s", err)
		return
	}

	log.Infof("Update thumb, cid: %s, uid: %s, change info: %s", cid, uid, change)
	c.JSON(http.StatusOK, gin.H{"error": "none"})
}

// UnThumbCheckIn update checkin and thumb collection
// Method: DELETE
func (cr *CheckInResource) UnThumbCheckIn(c *gin.Context) {
	uid := bson.ObjectIdHex(c.Query("uid"))
	cid := bson.ObjectIdHex(c.Query("cid"))

	///
	/// TBD check if user had already thumb the ci
	///
	err := cr.updateCIThumb(cid, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update checkin error: %s", err)
		return
	}

	change, err := cr.updateThumb(cid, uid, MgoCollPull)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update ci_thumb error: %s", err)
		return
	}

	log.Infof("Update thumb, cid: %s, uid: %s, change info: %s", cid, uid, change)
	c.JSON(http.StatusOK, gin.H{"error": "none"})
}

func (cr *CheckInResource) updateCIFavor(cid bson.ObjectId, val int) error {
	err := cr.CollCI.Update(
		bson.M{
			"_id": cid,
		},
		bson.M{
			"$inc": bson.M{
				"favor_count": val,
			},
		})
	return err
}

func (cr *CheckInResource) updateCIThumb(cid bson.ObjectId, val int) error {
	err := cr.CollCI.Update(
		bson.M{
			"_id": cid,
		},
		bson.M{
			"$inc": bson.M{
				"thumb_count": val,
			},
		})
	return err
}

// Mongo collection update type
const (
	MgoCollAddSet = "$addToSet"
	MgoCollInc    = "$inc"
	MgoCollPop    = "$pop"
	MgoCollPush   = "$push"
	MgoCollPull   = "$pull"
	MgoCollSet    = "$set"
)

func (cr *CheckInResource) updateFavor(cid bson.ObjectId, uid bson.ObjectId, action string) (*mgo.ChangeInfo, error) {
	change, err := cr.CollFavor.Upsert(
		bson.M{
			"_id": uid,
		},
		bson.M{
			action: bson.M{
				"cids": cid,
			},
		})
	return change, err
}

func (cr *CheckInResource) updateThumb(cid bson.ObjectId, uid bson.ObjectId, action string) (*mgo.ChangeInfo, error) {
	change, err := cr.CollThumb.Upsert(
		bson.M{
			"_id": uid,
		},
		bson.M{
			action: bson.M{
				"cids": cid,
			},
		})
	return change, err
}

// FavorCheckIn update checkin and ci_favor collection
// Method: POST
func (cr *CheckInResource) FavorCheckIn(c *gin.Context) {
	uid := bson.ObjectIdHex(c.PostForm("uid"))
	cid := bson.ObjectIdHex(c.PostForm("cid"))

	///
	/// TBD check if user had already favor the ci
	///
	err := cr.updateCIFavor(cid, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update checkin error: %s", err)
		return
	}

	change, err := cr.updateFavor(cid, uid, MgoCollAddSet)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update ci_favor error: %s", err)
		return
	}

	log.Infof("Update favor, cid: %s, uid: %s, change info: %s", cid, uid, change)
	c.JSON(http.StatusOK, gin.H{"error": "none"})
}

// UnFavorCheckIn update checkin and ci_favor collection
// Method: DELETE
func (cr *CheckInResource) UnFavorCheckIn(c *gin.Context) {
	uid := bson.ObjectIdHex(c.Query("uid"))
	cid := bson.ObjectIdHex(c.Query("cid"))

	///
	/// TBD check if user had already favor the ci
	///
	err := cr.updateCIFavor(cid, -1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update checkin error: %s", err)
		return
	}

	change, err := cr.updateFavor(cid, uid, MgoCollPull)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Debugf("update ci_favor error: %s", err)
		return
	}

	log.Infof("Update favor, cid: %s, uid: %s, change info: %s", cid, uid, change)
	c.JSON(http.StatusOK, gin.H{"error": "none"})
}

// CheckIn do checkin task for special user id
// Method: POST
func (cr *CheckInResource) CheckIn(c *gin.Context) {
	now := time.Now()
	content := c.PostForm("content")
	uidString := c.PostForm("uid")
	username := c.PostForm("user_name")
	//log.Infof("Username: %s", username)
	img := c.PostForm("img")
	var pub bool
	if c.PostForm("pub") == "false" {
		pub = false
	} else {
		pub = true
	}
	log.Info("Checkin user_id: " + uidString)
	uid := bson.ObjectIdHex(uidString)
	userinfo, err := cr.UserResource.QueryUserInfoByID(uidString)
	log.Infof("User info: %+v", userinfo)
	if err != nil {
		panic(err)
	}
	if !userinfo.CanCheckIn() {
		log.Infof("user_id: %s, already checkin", uidString)
		c.JSON(http.StatusForbidden, "Already checkin today")
		return
	}
	ciData := &mod.CheckIn{
		ID:          bson.NewObjectId(),
		UserID:      uid,
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
		Public:      pub,
		//FavorCount:  0,
		//ThumbCount:  0,
	}
	//log.Infof("Checkin content: %s", *ciData)
	err = cr.CollCI.Insert(ciData)
	if err != nil {
		panic(err)
	}
	err = cr.CollUser.Update(
		bson.M{
			"_id": uid,
		},
		bson.M{
			//"$inc": bson.M{
			//	"continuous": 1,
			//},
			"$set": bson.M{
				"last_checkin": *ciData,
				"continuous":   userinfo.CalcContinuous(),
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
	ListAllPublic
	ListSpecialPersonal
	ListPersonalWithCID
)

func (cr *CheckInResource) changePublic(c *gin.Context, pub bool) {
	uid := bson.ObjectIdHex(c.PostForm("uid"))
	cid := bson.ObjectIdHex(c.PostForm("cid"))
	err := cr.CollCI.Update(
		bson.M{
			"_id":     cid,
			"user_id": uid,
		},
		bson.M{
			"$set": bson.M{
				"public": pub,
			},
		})

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, "success!")
}

// MarkCIPublic change the checkin to public
// Method: PUT
func (cr *CheckInResource) MarkCIPublic(c *gin.Context) {
	cr.changePublic(c, true)
}

// MarkCIPrivate change the checkin to private
// Method: PUT
func (cr *CheckInResource) MarkCIPrivate(c *gin.Context) {
	cr.changePublic(c, false)
}

// DeleteCheckIn tag the checkin deleted
// Method: DELETE
func (cr *CheckInResource) DeleteCheckIn(c *gin.Context) {
	uid := bson.ObjectIdHex(c.Query("uid"))
	cid := bson.ObjectIdHex(c.Query("cid"))
	err := cr.CollCI.Update(
		bson.M{
			"_id":     cid,
			"user_id": uid,
		},
		bson.M{
			"$set": bson.M{
				"deleted": true,
			},
		})

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, "success")
}

// ListCheckIn return all checkin records
// Method: GET
func (cr *CheckInResource) ListCheckIn(c *gin.Context) {
	//col := cr.Mongo.DB("tidy").C("checkin")
	//uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
	//objectId := bson.ObjectIdHex(id)
	//log.Info("user_id: " + uid)
	var ci []mod.CheckIn
	//col.Find(bson.M{"user_id": uid}).All(&ci)
	tp, err := strconv.Atoi(c.DefaultQuery("type", strconv.Itoa(ListPersonal)))
	if err != nil {
		tp = ListPersonal
	}
	timestamp, tserr := strconv.ParseInt(
		c.DefaultQuery(
			"timestamp",
			strconv.FormatInt(time.Now().Unix(), 10),
		),
		10, 64)
	if tserr != nil {
		timestamp = time.Now().Unix()
	}
	count, cnterr := strconv.Atoi(c.DefaultQuery("count", "32"))
	if cnterr != nil {
		count = 32
	}
	var queryM bson.M
	switch tp {
	case ListPersonal:
		uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
		//log.Info(timestamp)
		//log.Info(count)
		queryM = bson.M{
			"user_id": uid,
			"timestamp": bson.M{
				"$lt": timestamp,
			},
			"deleted": false,
		}
	case ListPersonalWithCID:
		uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
		cid := bson.ObjectIdHex(c.DefaultQuery("cid", ""))
		//log.Info(timestamp)
		//log.Info(count)
		queryM = bson.M{
			"_id":     cid,
			"user_id": uid,
			"timestamp": bson.M{
				"$lt": timestamp,
			},
			"deleted": false,
		}
	case ListAllPublic:
		queryM = bson.M{
			"timestamp": bson.M{
				"$lt": timestamp,
			},
			"deleted": false,
			"public":  true,
		}
	case ListSpecialPersonal:
		spUID := bson.ObjectIdHex(c.DefaultQuery("special_uid", ""))
		queryM = bson.M{
			"user_id": spUID,
			"timestamp": bson.M{
				"$lt": timestamp,
			},
			"deleted": false,
			"public":  true,
		}
	default:
		uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
		queryM = bson.M{
			"user_id": uid,
			"timestamp": bson.M{
				"$lt": timestamp,
			},
			"deleted": false,
			"public":  true,
		}
	}
	cr.CollCI.Find(queryM).Limit(count).All(&ci)
	//col.Find(nil).All(&ci)
	//log.Infof("%s", ci)

	// update personal thumbed and favored infomation
	if tp != ListPersonal {
		uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
		cr.updateQueriedCIs(uid, ci)
	}
	c.JSON(http.StatusOK, ci)
}

func (cr *CheckInResource) updateQueriedCIs(uid bson.ObjectId, ci []mod.CheckIn) {
	///////////////////////////////////////
	/// TBD
	///////////////////////////////////////
	var thumb mod.Thumb
	var favor mod.Favorite
	cr.CollThumb.Find(
		bson.M{
			"_id": uid,
		},
	).One(&thumb)
	cr.CollFavor.Find(
		bson.M{
			"_id": uid,
		},
	).One(&favor)

	//log.Debugf("before updating: %+v", ci)
	for i := range ci {
		///////////////////////////////////////
		/// TBD
		///////////////////////////////////////
		for _, tID := range thumb.CheckinIDs {
			if ci[i].ID == tID {
				ci[i].Thumbed = true
			}
			//ci[i].Thumbed = true
		}
		for _, fID := range favor.CheckinIDs {
			if ci[i].ID == fID {
				ci[i].Favored = true
			}
			//ci[i].Favored = true
		}
	}
	//log.Debugf("after updating: %+v", ci)

}

func (cr *CheckInResource) UploadImg(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	//file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}
	//log.Debug(header)
	//filename split
	fns := strings.Split(header.Filename, ".")
	log.Debugf("Header.Filename: %s", fns)
	fileext := "png"
	if l := len(fns); l >= 2 {
		fileext = fns[l-1]
	}
	guid := bson.NewObjectId().Hex()
	filename := guid + "." + fileext
	log.Debug(filename)
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
