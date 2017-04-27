package buslogic

import (
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
)

func (w *Workflow) DataTransfer201704() error {
	log.Info("transfer start")
	// User
	log.Info("start to process user")
	legacyUsers, err := w.legacyMongoClient.GetAllLegacyUsers()
	if err != nil {
		log.Errorf("fail to get user: %+v", err)
		return err
	}
	log.Infof("%d users in total", len(legacyUsers))
	for _, u := range legacyUsers {
		switch u.UserType {
		case model.TEACHER:
			teacher := &model.Teacher{
				Username: u.Username,
				Password: u.Password,
				UserType: model.UserTypeTeacher,
				Fullname: u.Fullname,
				Mobile:   u.Mobile,
				Address:  u.Address,
				Gender:   u.Gender,
				Major:    u.Major,
				Academic: u.Academic,
				Aptitude: u.Aptitude,
				Problem:  u.Problem,
			}
			if teacher.Address == "" {
				teacher.Address = "紫荆C楼407室"
			}
			err = w.mongoClient.InsertTeacher(teacher)
			if err != nil {
				log.Errorf("fail to insert teacher %+v, err: %+v", teacher, err)
				return err
			}
		case model.ADMIN:
			admin := &model.Admin{
				Username: u.Username,
				Password: u.Password,
				UserType: model.UserTypeAdmin,
				Fullname: u.Fullname,
				Mobile:   u.Mobile,
			}
			err = w.mongoClient.InsertAdmin(admin)
			if err != nil {
				log.Errorf("fail to insert admin %+v, err: %+v", admin, err)
				return err
			}
		default:
			log.Errorf("fail to process user %+v", u)
		}
	}
	// reservation
	log.Info("start to process reservation")
	legacyReservations, err := w.legacyMongoClient.GetAllLegacyReservations()
	if err != nil {
		log.Errorf("fail to get reservations: %+v", err)
		return err
	}
	log.Infof("%d reservations in total", len(legacyReservations))
	for _, r := range legacyReservations {
		reservation := &model.Reservation{
			StartTime:       r.StartTime,
			EndTime:         r.EndTime,
			TeacherUsername: r.TeacherUsername,
			TeacherFullname: r.TeacherFullname,
			TeacherMobile:   r.TeacherMobile,
			TeacherAddress:  r.TeacherAddress,
			StudentInfo: model.StudentInfo{
				Fullname:        r.StudentInfo.Name,
				Gender:          r.StudentInfo.Gender,
				Username: r.StudentInfo.StudentId,
				School:          r.StudentInfo.School,
				Hometown:        r.StudentInfo.Hometown,
				Mobile:          r.StudentInfo.Mobile,
				Email:           r.StudentInfo.Email,
				Experience:      r.StudentInfo.Experience,
				Problem:         r.StudentInfo.Problem,
			},
			StudentFeedback: model.StudentFeedback{
				Fullname: r.StudentFeedback.Name,
				Problem:  r.StudentFeedback.Problem,
				Choices:  r.StudentFeedback.Choices,
				Score:    r.StudentFeedback.Score,
				Feedback: r.StudentFeedback.Feedback,
			},
			TeacherFeedback: model.TeacherFeedback{
				TeacherFullname: r.TeacherFeedback.TeacherFullname,
				TeacherUsername: r.TeacherFeedback.TeacherUsername,
				StudentFullname: r.TeacherFeedback.StudentFullname,
				Problem:         r.TeacherFeedback.Problem,
				Solution:        r.TeacherFeedback.Solution,
				AdviceToCenter:  r.TeacherFeedback.AdviceToCenter,
			},
		}
		if reservation.TeacherAddress == "" {
			reservation.TeacherAddress = "紫荆C楼407室"
		}
		switch r.Status {
		case model.AVAILABLE:
			reservation.Status = model.ReservationStatusAvailable
		case model.RESERVATED:
			reservation.Status = model.ReservationStatusReservated
		case model.DELETED:
			reservation.Status = model.ReservationStatusDeleted
		default:
			log.Errorf("fail to process reservation %+v", r)
		}
		err = w.mongoClient.InsertReservation(reservation)
		if err != nil {
			log.Errorf("fail to insert reservation %+v, err: %+v", reservation, err)
			return err
		}
	}
	// index
	log.Info("start to process indexes")
	err = w.mongoClient.EnsureAllIndexes()
	if err != nil {
		log.Errorf("fail to ensure all indexes: %+v", err)
		return err
	}

	return nil
}
