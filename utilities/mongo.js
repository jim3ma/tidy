use tidy;
db.createRole(
   {
     role: "tidyOpRole",
     privileges: [
       {
           resource: { db: "tidy", collection: "" },
           actions: [
               "collStats",
               "convertToCapped",
               "createCollection",
               "dbHash",
               "dbStats",
               "dropCollection",
               "createIndex",
               "dropIndex",
               "emptycapped",
               "find",
               "insert",
               "killCursors",
               "listIndexes",
               "listCollections",
               "remove",
               "renameCollectionSameDB",
               "update" ]
       }
     ],
     roles: []
   }
);

db.createUser(
    {
        user:"tidy",
        pwd:"111111",
        roles:[
            "tidyOpRole"
        ]
    }
);

// create Index

db.user.ensureIndex({"email": 1});

db.checkin.ensureIndex({"timestamp": -1});
