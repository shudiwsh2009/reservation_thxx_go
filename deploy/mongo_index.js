use appointment
show collections
db.user.createIndex({"username": 1})
db.user.createIndex({"fullname": 1})
db.user.createIndex({"mobile": 1})
db.appointment.createIndex({"studentInfo.studentId": 1})
db.appointment.createIndex({"startTime": 1, "status_GO": 1})
