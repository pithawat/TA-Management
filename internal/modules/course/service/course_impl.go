package service

import (
	"TA-management/internal/modules/course/dto/request"
	courseResponse "TA-management/internal/modules/course/dto/response"
	"TA-management/internal/modules/course/repository"
	"TA-management/internal/modules/shared/dto/response"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type CourseServiceImplementation struct {
	repo        repository.CourseRepository
	redisClient *redis.Client
}

func NewCourseService(repo repository.CourseRepository, redisClient *redis.Client) CourseServiceImplementation {
	return CourseServiceImplementation{repo: repo, redisClient: redisClient}
}

func (s CourseServiceImplementation) GetAllJobPost() (*response.RequestDataResponse, error) {

	courses, err := s.repo.GetAllJobPost()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := response.RequestDataResponse{
		Data:    courses,
		Message: "Success",
	}

	return &response, nil
}

func (s CourseServiceImplementation) GetAllJobPostAllStatus() (*response.RequestDataResponse, error) {

	courses, err := s.repo.GetAllJobPostAllStatus()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := response.RequestDataResponse{
		Data:    courses,
		Message: "Success",
	}

	return &response, nil
}

func (s CourseServiceImplementation) GetAllJobPostByStudentId(studentId int) (*response.RequestDataResponse, error) {
	courses, err := s.repo.GetAllJobPostByStudentId(studentId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := response.RequestDataResponse{
		Data:    courses,
		Message: "Success",
	}

	return &response, nil
}

func (s CourseServiceImplementation) GetAllCourse() (*response.RequestDataResponse, error) {
	// ctx := context.Background()
	// cacheKey := "course:all"

	// Check cache
	// if s.redisClient != nil {
	// 	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	// 	if err == nil {
	// 		fmt.Println("from redis")
	// 		var courses []courseResponse.Course
	// 		if err := json.Unmarshal([]byte(val), &courses); err == nil {
	// 			return &response.RequestDataResponse{
	// 				Data:    courses,
	// 				Message: "Success",
	// 			}, nil
	// 		}
	// 	}
	// }

	courses, err := s.repo.GetAllCourse()
	fmt.Println("from DB")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Set cache
	// if s.redisClient != nil {
	// 	if data, err := json.Marshal(courses); err == nil {
	// 		s.redisClient.Set(ctx, cacheKey, data, 10*time.Minute)
	// 	}
	// }

	response := response.RequestDataResponse{
		Data:    courses,
		Message: "Success",
	}

	return &response, nil
}

func (s CourseServiceImplementation) CreateCourse(body request.CreateCourse) (*response.CreateResponse, error) {
	count, err := s.repo.IsCourseExist(body)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, fmt.Errorf("already have this course")
	}

	id, err := s.repo.CreateCourse(body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Invalidate cache
	if s.redisClient != nil {
		s.redisClient.Del(context.Background(), "course:all")
	}

	return &response.CreateResponse{
		Message: "Created successfully!",
		Id:      id,
	}, nil
}

func (s CourseServiceImplementation) UpdateCourse(body request.UpdateCourse) (response.GeneralResponse, error) {
	err := s.repo.UpdateCourse(body)
	if err != nil {
		fmt.Println(err)
		return response.GeneralResponse{Message: "Update Failed!"}, err
	}
	// Invalidate cache
	if s.redisClient != nil {
		s.redisClient.Del(context.Background(), "course:all")
	}
	return response.GeneralResponse{Message: "Update Successful"}, err
}

func (s CourseServiceImplementation) DeleteCourse(id int) (response.GeneralResponse, error) {
	err := s.repo.DeleteCourse(id)
	if err != nil {
		fmt.Println(err)
		return response.GeneralResponse{Message: "Delete Failed!"}, err
	}
	// Invalidate cache
	if s.redisClient != nil {
		s.redisClient.Del(context.Background(), "course:all")
	}
	return response.GeneralResponse{Message: "Delete Successful"}, err
}

func (s CourseServiceImplementation) CreateJobPost(body request.CreateJobPost) (response.CreateResponse, error) {
	id, err := s.repo.CreateJobPost(body)
	if err != nil {
		fmt.Println(err)
		return response.CreateResponse{
			Message: "Create Job Post Failed!",
		}, err
	}
	return response.CreateResponse{
		Message: "Create Job Post Successfully",
		Id:      id,
	}, nil
}

func (s CourseServiceImplementation) UpdateJobPost(body request.UpdateJobPost) (*response.RequestDataResponse, error) {
	err := s.repo.UpdateJobPost(body)
	if err != nil {
		return nil, err
	}

	jobPost, err := s.repo.GetJobPostByID(body.Id)
	if err != nil {
		return nil, err
	}

	return &response.RequestDataResponse{
		Data:    jobPost,
		Message: "Update Job Post Successful",
	}, nil
}

func (s CourseServiceImplementation) DeleteJobPost(jobPostId int) (response.GeneralResponse, error) {
	err := s.repo.DeleteJobPost(jobPostId)
	if err != nil {
		fmt.Println(err)
		return response.GeneralResponse{Message: "Delete Job Post Failed!"}, err
	}
	return response.GeneralResponse{Message: "Delete Job Post Successful"}, err
}

func (s CourseServiceImplementation) ApplyJobPost(body request.ApplyJobPost) (*response.CreateResponse, error) {
	//check student status on this job
	ok, err := s.repo.CheckStudentJobpostStatus(body)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("already apply to this jobpost")
	}

	taAllocation, err := s.repo.GetTaAllocation(*body.JobPostID)
	if err != nil {
		return nil, err
	}

	allocationCount, err := s.repo.CountTaAllocation(*body.JobPostID)
	if err != nil {
		return nil, err
	}

	if allocationCount >= taAllocation {
		return nil, fmt.Errorf("ta allocation is full")
	}

	tx, err := s.repo.StartDBTx()
	if err != nil {
		return nil, err
	}

	defer s.repo.RollbackTx(tx)
	if body.TranscriptBytes != nil {
		err = s.repo.UpsertTranscript(tx, body)
		if err != nil {
			s.repo.RollbackTx(tx)
			return nil, err
		}
	}

	if body.BankAccountBytes != nil {
		err = s.repo.UpsertBankAccount(tx, body)
		if err != nil {
			s.repo.RollbackTx(tx)
			return nil, err
		}
	}

	if body.StudentCardBytes != nil {
		err = s.repo.UpsertStudentCard(tx, body)
		if err != nil {
			s.repo.RollbackTx(tx)
			return nil, err
		}
	}

	err = s.repo.UpdateStudentData(tx, body)
	if err != nil {
		s.repo.RollbackTx(tx)
		return nil, err
	}

	id, err := s.repo.InsertApplication(tx, body)
	if err != nil {
		s.repo.RollbackTx(tx)
		return nil, err
	}

	err = s.repo.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	return &response.CreateResponse{
		Message: "Apply course successfully",
		Id:      id,
	}, nil
}

func (s CourseServiceImplementation) GetApplicationByStudentId(studentId int) (*response.RequestDataResponse, error) {
	applications, err := s.repo.GetApplicationByStudentId(studentId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &response.RequestDataResponse{
		Data:    applications,
		Message: "GET success",
	}, nil

}

func (s CourseServiceImplementation) GetAllTimeApprovedCoursesByStudentId(studentId int) (*response.RequestDataResponse, error) {
	applications, err := s.repo.GetAllTimeApprovedCoursesByStudentId(studentId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &response.RequestDataResponse{
		Data:    applications,
		Message: "GET success",
	}, nil
}

func (s CourseServiceImplementation) GetApplicationByCourseId(courseId int) (*response.RequestDataResponse, error) {
	applications, err := s.repo.GetApplicationByCourseId(courseId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &response.RequestDataResponse{
		Data:    applications,
		Message: "GET success",
	}, nil
}

func (s CourseServiceImplementation) GetApplicationDetail(applicationId int) (*response.RequestDataResponse, error) {
	application, err := s.repo.GetApplicationDetail(applicationId)
	if err != nil {
		return nil, nil
	}
	return &response.RequestDataResponse{
		Data:    application,
		Message: "GET SUCCESS",
	}, nil

}

func (s CourseServiceImplementation) GetApplicationTranscriptPdf(applicationId int) (*courseResponse.PdfFile, error) {
	applicationPdf, err := s.repo.GetApplicationTranscriptPdf(applicationId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return applicationPdf, nil
}

func (s CourseServiceImplementation) GetApplicationBankAccountPdf(applicationId int) (*courseResponse.PdfFile, error) {
	applicationPdf, err := s.repo.GetApplicationBankAccountPdf(applicationId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return applicationPdf, nil
}

func (s CourseServiceImplementation) GetApplicationStudentCardPdf(applicationId int) (*courseResponse.PdfFile, error) {
	applicationPdf, err := s.repo.GetApplicationStudentCardPdf(applicationId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return applicationPdf, nil
}

func (s CourseServiceImplementation) ApproveApplication(applicationId int) (*response.GeneralResponse, error) {

	courseId, studentId, jobPostId, err := s.repo.GetApproveApplicationData(applicationId)
	if err != nil {
		return nil, fmt.Errorf("fail1 : %v", err)
	}

	//check ta allocation
	taAllocation, err := s.repo.GetTaAllocation(jobPostId)
	if err != nil {
		return nil, fmt.Errorf("fail2 : %v", err)
	}

	allocationCount, err := s.repo.CountTaAllocation(jobPostId)
	if err != nil {
		return nil, fmt.Errorf("fail3 : %v", err)
	}

	if allocationCount >= taAllocation {
		return nil, fmt.Errorf("ta allocation is full")
	}

	tx, err := s.repo.StartDBTx()
	if err != nil {
		return nil, err
	}
	defer s.repo.RollbackTx(tx)

	err = s.repo.UpdateApplicationStatus(tx, applicationId)
	if err != nil {
		s.repo.RollbackTx(tx)
		return nil, err
	}

	err = s.repo.InsertTaCourse(tx, studentId, courseId)
	if err != nil {
		s.repo.RollbackTx(tx)
		return nil, err
	}

	err = s.repo.CommitTx(tx)
	if err != nil {
		return nil, err
	}

	NewCount, err := s.repo.CountTaAllocation(jobPostId)
	if err != nil {
		return nil, err
	}
	if NewCount == taAllocation {
		err = s.repo.UpdateJobPostStatus(jobPostId)
		if err != nil {
			return nil, err
		}
	}

	return &response.GeneralResponse{
		Message: "Approved application Successful",
	}, nil
}

func (s CourseServiceImplementation) GetProfessorCourse(professorId int) (*response.RequestDataResponse, error) {
	courses, err := s.repo.GetProfessorCourse(professorId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := response.RequestDataResponse{
		Data:    courses,
		Message: "Success",
	}

	return &response, nil
}

func (s CourseServiceImplementation) GetApplicationByProfessorId(professorId int) (*response.RequestDataResponse, error) {
	applications, err := s.repo.GetApplicationByProfessorId(professorId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &response.RequestDataResponse{
		Data:    applications,
		Message: "GET success",
	}, nil
}

func (s CourseServiceImplementation) RejectApplication(rq request.RejectApplication) (*response.GeneralResponse, error) {
	err := s.repo.RejectApplication(rq)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &response.GeneralResponse{
		Message: "Rejected application Successful",
	}, nil
}

func (s CourseServiceImplementation) UpdateCourseDiscord(courseId int, roleId string, channelId string, channelName string) (*response.GeneralResponse, error) {
	err := s.repo.UpdateCourseDiscord(courseId, roleId, channelId, channelName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &response.GeneralResponse{
		Message: "Update Discord Successful",
	}, nil
}

func (s CourseServiceImplementation) SoftDeleteExpiredData() error {
	return s.repo.SoftDeleteExpiredData()
}

func (s CourseServiceImplementation) GetTermHistory() (*response.RequestDataResponse, error) {
	result, err := s.repo.GetExpiredSemesters()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s CourseServiceImplementation) GetHistoryCourses(semesterID int) (*response.RequestDataResponse, error) {
	result, err := s.repo.GetCoursesBySemesterID(semesterID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
