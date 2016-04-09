use admin;
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
