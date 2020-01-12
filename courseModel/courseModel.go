package courseModel

type CourseStruc struct {
	Name   string `json:"Name"`
	ID     string `json:"courseID"`
	Credit string `json:"Credit"`
	Groups map[string]*GroupBig
}

//defualt is Groups []GroupDetail
type GroupBig struct {
	SecTime map[string]Group
	Group   string `json:group`
	Teacher string `json:"Teacher"`
	Mid     string `json:Mid`
	Final   string `json:Final`
	Note    string `json:Note`
}

// should change to day eiei
type Group struct {
	Day      string `json:Day`
	Time     string `json:Time`
	Room     string `json:Room`
	Building string `json:Building`
}
