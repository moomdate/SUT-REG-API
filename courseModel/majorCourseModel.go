package courseModel

type MajorCourse struct {
	CourseId     int   `json:"courseId"`
	CourseNameEn string `json:"nameEn"`
	CourseNameTh string `json:"nameTh"`
	Credit string `json:"credit"`
}
