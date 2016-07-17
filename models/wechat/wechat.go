package wechat

import (
	//"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

const (
	LanguageZhCN = "zh_CN" // 简体中文
	LanguageZhTW = "zh_TW" // 繁体中文
	LanguageEN   = "en"    // 英文
)

const (
	SexUnknown = iota // 未知
	SexMale           // 男性
	SexFemale         // 女性
)

type WeChatUserInfo struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	UserId   bson.ObjectId `bson:"user_id" json:"user_id"`
	OpenId   string        `bson:"openid" json:"openid"`     // 用户的唯一标识
	Nickname string        `bson:"nickname" json:"nickname"` // 用户昵称
	Sex      int           `bson:"sex" json:"sex"`           // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	City     string        `bson:"city" json:"city"`         // 普通用户个人资料填写的城市
	Province string        `bson:"province" json:"province"` // 用户个人资料填写的省份
	Country  string        `bson:"country" json:"country"`   // 国家, 如中国为CN

	// 用户头像，最后一个数值代表正方形头像大小有0、46、64、96、132数值可选，0代表640*640正方形头像，
	// 用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	HeadImageURL string `bson:"headimgurl,omitempty" json:"headimgurl,omitempty"`

	// 用户特权信息，json 数组，如微信沃卡用户为chinaunicom
	Privilege []string `bson:"privilege,omitempty" json:"privilege,omitempty"`
	// 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	UnionId string `bson:"unionid,omitempty" json:"unionid,omitempty"`
}
