use tidy;

// create Index
db.user.ensureIndex({"email": 1});
db.checkin.ensureIndex({"timestamp": -1});
db.wechat.ensureIndex({"openid": 1});
