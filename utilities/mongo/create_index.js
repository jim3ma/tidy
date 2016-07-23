use tidy;

// create Index
db.user.ensureIndex({"email": 1});

db.checkin.ensureIndex({"timestamp": -1});
// db.ci_comment.ensureIndex({"checkin_id": 1});
// db.ci_like.ensureIndex({"checkin_id": 1});

db.wechat.ensureIndex({"openid": 1});
