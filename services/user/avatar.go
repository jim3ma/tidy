package user

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	mod "github.com/jim3mar/tidy/models/user"
	"github.com/jim3mar/tidy/utilities/cache"
	"gopkg.in/mgo.v2/bson"
)

func (ur *UserResource) QueryAvatar(uid bson.ObjectId) (string, error) {
	/// query cache first
	/// TBD
	///
	conn := cache.Pool.Get()
	defer conn.Close()

	id := fmt.Sprintf("user_avatar:%s", uid.Hex())
	ava, err := redis.String(conn.Do("GET", id))
	if len(ava) > 8 {
		return ava, nil
	}

	query := ur.CollUser.Find(
		bson.M{
			"_id": uid,
		},
	)
	if c, err := query.Count(); c == 0 {
		return "", err
	}

	var user mod.User
	err = query.One(&user)
	if err != nil {
		return "", err
	}
	/// update cache
	/// TBD
	///
	conn.Do("SET", id, user.Portrait)
	return user.Portrait, nil
}
