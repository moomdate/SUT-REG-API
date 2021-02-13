package courseModel

type CourseStruc struct {
	NameEn      string  `json:"nameEn"`
	NameTh      string  `json:"nameTh"`
	BelongTo    string  `json:"belongTo"`
	Status      string  `json:"status"`
	ID          string  `json:"courseID"`
	Credit      string  `json:"credit"`
	Description string  `json:"description"`
	Year        int     `json:"year""`
	Semester    int     `json:"semester"`
	Groups      []Group `json:"groups"`
}

// default is Groups []GroupDetail
type Group struct {
	SecTime   []SectionTime `json:"sectionTime"`
	Group     int           `json:"group"`
	Open      int           `json:"openAmount"`
	Reserved  int           `json:"reservedAmount"`
	Remaining int           `json:"remainingAmount"`
	Teacher   string        `json:"teacher"`
	Mid       string        `json:"mid"`
	Final     string        `json:"final"`
	Note      string        `json:"note"`
}

// should change to day eiei
type SectionTime struct {
	Day      string `json:"day"`
	Time     string `json:"time"`
	Room     string `json:"room"`
	Building string `json:"building"`
}
type Course struct {
	Acadyear string         `json:"acadyear"`
	Semester string         `json:"semester"`
	Data     []CourseDetail `json:"courses"`
}
type CourseDetail struct {
	Name     string `json:"name"`
	Group    string `json:"group"`
	CourseID string `json:"courseID"`
	Version  string `json:"version"`
}
type CourseVersion struct {
	CourseID string `json:"courseID"`
	Version  string `json:"version"`
	AliasID  string `json:"aliasID"`
}
