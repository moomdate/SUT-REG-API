package services

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"reg-api/courseModel"
	"strconv"
	"strings"
	"time"
)

const ( //child access
	headerGroups     = "#F5F5F5"
	acTable          = "body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3)  "
	getStatus        = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(5) > td:nth-child(3) > font:nth-child(1)"
	getCourseNameEn  = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > b:nth-child(1) > font:nth-child(1)"
	getCourseNameTh  = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(2) > font:nth-child(1)"
	getBelongTo      = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(3) > td:nth-child(3) > font:nth-child(1)	"
	getCredit        = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(4) > td:nth-child(3) > font:nth-child(1)"
	getDay           = "td:nth-child(4)"                     //get Date
	getTime          = "td:nth-child(5)"                     //get Time
	getRoom          = "td:nth-child(6) "                    //get room
	getBuilding      = "td:nth-child(7)"                     //get building
	checkTc          = "td:nth-child(4) > font:nth-child(1)" //check teacher
	getTc            = "td:nth-child(5) > font:nth-child(1)" //get teacher
	getYearSemester  = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(8) > td:nth-child(2) > font:nth-child(3)"
	getCourseId      = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > b:nth-child(1) > font:nth-child(1)"
	getDescription   = "table:nth-child(6) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > font:nth-child(3)"
	getCourseList    = "body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > font:nth-child(5) > font:nth-child(2) > font:nth-child(2) > font:nth-child(3) > div:nth-child(3) > table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > font:nth-child(1) > table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(1) > tbody:nth-child(1) > tr"
	getOpenAmount    = "td:nth-child(9)"
	getReserveAmount = "td:nth-child(10)"
	getRemainAmount  = "td:nth-child(11)"
)
const baseURLCourseInMajor = "http://reg5.sut.ac.th/registrar/program_info_1.asp?programid="

func ClearCache(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	cacheType := getParam["type"]
	cacheFolder := ""
	switch cacheType {
	case "dig":
		cacheFolder = "digCode"
		break
	case "major":
		cacheFolder = "majorCourse"
	default:
		cacheFolder = "course"
		break
	}
	log.Print("clear reg cache")
	err := os.RemoveAll("./reg_cache/" + cacheFolder)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

// check difference day of group
func shouldSumGroups(strIn string) (getnumber bool) {
	//data := strings.Split(strIn, " ") not work
	status := false
	splitStrToArr := strings.Fields(strIn)
	if len(splitStrToArr[0]) >= 3 { // number of group 0-99 [type string]
		status = true
	}
	return status
}
func FindCourseAll(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	cid := getParam["cid"]
	acadyears := getParam["year"]
	semesters := getParam["semester"]
	fmt.Print(" =====", cid+" "+acadyears+" - "+semesters)
	scrapLink := colly.NewCollector(
		//colly.CacheDir("./reg_cache/allCourseVersion"),
	)

	var courseVersionList = []courseModel.CourseVersion{}

	scrapLink.OnHTML("body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > font:nth-child(6) > font:nth-child(1) > font:nth-child(1) > font:nth-child(1) > table:nth-child(2) > tbody:nth-child(1)", func(el *colly.HTMLElement) {
		el.ForEach("tr", func(_ int, elTr *colly.HTMLElement) {
			if elTr.Attr("valign") == "TOP" {
				elTr.ForEach(" td:nth-child(2) > font:nth-child(1) > a:nth-child(1)", func(_ int, elTb *colly.HTMLElement) {
					split := strings.Split(elTb.Text, " - ")
					courseD := courseModel.CourseVersion{
						CourseID: split[0],
						Version:  split[1],
						AliasID:  subCourse(elTb.Attr("href")),
					}
					courseVersionList = append(courseVersionList, courseD)

				})
			}
		})
	})
	scrapLink.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"
	})
	scrapLink.Request("POST",
		"http://reg5.sut.ac.th/registrar/class_info_1.asp?avs782309057=22&backto=home",
		strings.NewReader(fmt.Sprintf("coursestatus=O00&facultyid=all&maxrow=50&acadyear=%s&semester=%s&CAMPUSID=&LEVELID=&coursecode=%s&coursename=&cmd=2", acadyears, semesters, cid)),
		nil,
		http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}})

	uniqueVersion := filterUniqVersion(courseVersionList)

	w.Header().Set("Content-type", "application/json; charset=UTF-8;")

	json.NewEncoder(w).Encode(uniqueVersion)
}

func filterUniqVersion(courseVersionList []courseModel.CourseVersion) []courseModel.CourseVersion {
	uniqueVersion := []courseModel.CourseVersion{}
	for i := range courseVersionList {
		if !Contains(uniqueVersion, courseVersionList[i].Version) {
			uniqueVersion = append(uniqueVersion, courseVersionList[i])
		}
	}
	return uniqueVersion
}

func Contains(a []courseModel.CourseVersion, version string) bool {
	for _, n := range a {
		if version == n.Version {
			return true
		}
	}
	return false
}

func getGroupNumber(str string) string {
	explodeAr := strings.Fields(str)
	explodeStr := explodeAr[0]
	if len(explodeStr) > 3 {
		explodeStr = "0"
	}
	realn, err := strconv.Atoi(explodeStr)
	if err != nil { // need to remove zero subtract the number
		fmt.Println("it error ", err)
	}
	arlOut := strconv.Itoa(realn)
	return arlOut
}
func digCourseCode(ID string, Year string, sem string) string {
	var link string
	var tim string
	if len(ID) < 1 {
		return ""
	}
	mainLink := "http://reg5.sut.ac.th/registrar/class_info_1.asp?coursestatus=O00&facultyid=all&maxrow=1&acadyear=" + Year + "&semester=" + sem + "&coursecode=" + ID
	scrapLink := colly.NewCollector(
		colly.CacheDir("./reg_cache/digCode"),
	)
	scrapLink.SetRequestTimeout(5 * time.Second)
	scrapLink.OnHTML("a[href]", func(el *colly.HTMLElement) {

		link = el.Attr("href")
		if strings.Contains(link, "courseid") {
			tim = subCourse(link)
		}
	})
	scrapLink.Visit(mainLink)
	return tim
}

//===========================
//filter coursecode from link
// return "Course ID"
//===========================
func subCourse(inputStr string) string {
	var out string
	tim := strings.Split(inputStr, "&")
	for _, mam := range tim {
		if strings.Contains(mam, "courseid") { // see
			out = strings.Split(mam, "=")[1] // get number
			break
		}
	}
	return out
}

type Amount struct {
	open      int
	reserve   int
	remaining int
}

func splitCourseName(text string) (string, string) {
	thaiIndex := 0
	for i, r := range text {
		if r > 3500 {
			thaiIndex = i
			break
		}
	}
	firstIndex := 0
	if text[0] == 194 { // remove double space
		firstIndex = 2
	}
	nameEn := text[firstIndex : thaiIndex-2] // remove double space
	nameTh := text[thaiIndex:len(text)]
	return nameEn, nameTh
}
func getNumber(strNumber string) int {

	tempStr := strNumber[2:len(strNumber)]
	result, err := strconv.Atoi(tempStr)
	if err != nil {
		log.Fatal("program id is err")
	}
	return result
}

func ImportFormReg(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	stdid := getParam["stdid"]
	acadyears := getParam["acadyear"]
	semesters := getParam["semester"]
	if strings.ToUpper(stdid[0:1]) == "B" {
		stdid = "1" + stdid[1:]
	} else if strings.ToUpper(stdid[0:1]) == "M" {
		stdid = "2" + stdid[1:]
	} else {
		stdid = "3" + stdid[1:]
	}
	URL := "http://reg5.sut.ac.th/registrar/learn_time.asp?f_cmd=2&studentid=" + stdid + "&acadyear=" + acadyears + "&maxsemester=3&rnd=43673.6771527778&firstday=22/7/2562&semester=" + semesters
	log.Print("# url:", URL)
	c := colly.NewCollector()
	var cid, gid, cname, ver string

	var Data []courseModel.CourseDetail
	c.OnHTML(getCourseList, func(e *colly.HTMLElement) {
		print("===", e.Text)

		if e.Index != 0 {
			e.ForEach("td", func(_ int, el *colly.HTMLElement) {
				if el.Index == 0 {
					cid = strings.TrimSpace(el.Text)
					datacid := strings.Split(cid, "-")
					cid = datacid[0]
					ver = datacid[1]
				}
				if el.Index == 1 {
					for i := 0; i < len(el.Text); i++ {
						if int(el.Text[i]) > 160 {
							s := strings.Split(el.Text, string([]rune(el.Text)[i]))
							cname = s[0]
							break
						}
					}
				}
				if el.Index == 2 {
					gid = el.Text
				}
			})
			courseD := courseModel.CourseDetail{
				Name:     cname,
				Group:    gid,
				CourseID: cid,
				Version:  ver,
			}
			Data = append(Data, courseD)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=utf-8"
	})

	c.Visit(URL)
	courses := &courseModel.Course{
		Acadyear: acadyears,
		Semester: semesters,
		Data:     Data,
	}
	w.Header().Set("Content-type", "application/json; charset=UTF-8;")
	json.NewEncoder(w).Encode(courses)
}
func GetMajor(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	id := getParam["id"]
	//mainLink := "test"

	var majorList = []courseModel.MajorModel{}
	var major = courseModel.MajorModel{}
	requestParams := "facultyid=" + id + "&f_cmd="
	if id != "10000" {
		requestParams = requestParams + "&Levelid=1&f_cmd=&Acadyear=-1"
	}
	getCredit := false
	scrapLink := colly.NewCollector(
		//colly.CacheDir("./reg_cache/majorList"),
	)
	scrapLink.SetRequestTimeout(5 * time.Second)
	scrapLink.OnHTML("body > table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > table:nth-child(6) > tbody:nth-child(1)", func(el *colly.HTMLElement) {
		el.ForEach("tr", func(_ int, elTr *colly.HTMLElement) {
			elTr.ForEach("td", func(_ int, elTr2 *colly.HTMLElement) {
				if elTr2.Attr("width") == "20" { // program id
					tempStr := elTr2.Text
					major.ProgramId = getNumber(tempStr)
				} else if elTr2.Attr("valign") == "TOP" && len(elTr2.Text) > 5 { // course
					getCredit = true
					major.Course = elTr2.Text[0 : len(elTr2.Text)-1]
				} else if elTr2.Attr("valign") == "TOP" && elTr2.Attr("bgcolor") == "#FFFFDE" && elTr2.Attr("align") == "CENTER" && getCredit == true {
					getCredit = false
					credit, err := strconv.Atoi(elTr2.Text)
					if err != nil {
						log.Fatal("program id is err")
					}
					major.Credit = credit
					majorList = append(majorList, major)
					fmt.Println(elTr2.Text)
				}
			})
		})
	})
	scrapLink.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"
	})

	header := http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}, "User-Agent": []string{"application/reg1"}}
	scrapLink.Request("POST",
		"http://reg5.sut.ac.th/registrar/program_info.asp",
		strings.NewReader(requestParams),
		nil,
		header)

	w.Header().Set("Content-type", "application/json; charset=UTF-8;")
	json.NewEncoder(w).Encode(majorList)
}

func CourseMajor(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	id := getParam["id"]
	var majorCourseList []courseModel.MajorCourse
	var courseTemp = courseModel.MajorCourse{}
	state := "courseId" // ready

	mainLink := baseURLCourseInMajor + id
	scrapLink := colly.NewCollector(
		colly.CacheDir("./reg_cache/majorCourse"),
	)
	scrapLink.SetRequestTimeout(5 * time.Second)
	scrapLink.OnHTML("body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > table:nth-child(4) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table", func(el *colly.HTMLElement) {
		el.ForEach("tbody", func(_ int, eltr *colly.HTMLElement) {
			eltr.ForEach("tr", func(_ int, elFont *colly.HTMLElement) {

				if elFont.Attr("bgcolor") == "#FFFFDE" {
					elFont.ForEach("td", func(_ int, elFontEl *colly.HTMLElement) {

						if elFontEl.Attr("valign") == "TOP" && state == "courseId" {
							tempText := strings.Split(elFontEl.ChildText("font"), " ")
							tempCourseIn, err := strconv.Atoi(tempText[0])
							if err != nil {
								log.Fatal("get course id err")
							}
							courseTemp.CourseId = tempCourseIn
						} else if elFontEl.Attr("valign") == "TOP" && state == "credit" {
							courseTemp.Credit = elFontEl.ChildText("font")
							majorCourseList = append(majorCourseList, courseTemp)
							courseTemp = courseModel.MajorCourse{}
							state = "courseId"
						} else if elFontEl.Attr("valign") == "CENTER" {
							nameEn, nameTh := splitCourseName(elFontEl.ChildText("font"))
							courseTemp.CourseNameEn = nameEn
							courseTemp.CourseNameTh = nameTh
							state = "credit"
						}
					})
				}
			})
		})
	})
	scrapLink.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"
	})
	scrapLink.Visit(mainLink)
	w.Header().Set("Content-type", "application/json; charset=UTF-8;")
	json.NewEncoder(w).Encode(majorCourseList)
}
func getDateFormCourse(text string) {
	seeDash := 0
	if len(text) > 0 {
		for i, r := range text {
			if string(r) == "-" {
				seeDash = i
				break
			}
		}
		dateTimeStr := text[0 : seeDash+7] // 6 ก.ย. 2562 เวลา 12:00 - 14:00
		seeDash = 0
		textTemp := text
		patternLen := len(dateTimeStr)
		for i := 0; i >= 0; i++ {
			if len(textTemp) <= patternLen {
				break
			}
			textOriLen := len(textTemp)
			textTemp = textTemp[patternLen+1 : textOriLen]
		}
	}
	//fmt.Println(text)
}
func subDescription(text string) string {
	split := strings.Split(text, "หมายเหตุเรียน")
	return split[0]
}
func ScrapingCourseDetail(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	pID := getParam["id"]
	pYear := getParam["year"]
	pSemis := getParam["semester"]

	var Course courseModel.CourseStruc
	var Group []courseModel.Group

	var courseNameEn, courseNameTh, belongTo, status, credit, description string
	var sectionTimeTemp []courseModel.SectionTime
	var gTemp int
	var tempDetail courseModel.Group
	var amount Amount
	year, _ := strconv.Atoi(pYear)
	semester, _ := strconv.Atoi(pSemis)

	var tempCID string
	newPId := digCourseCode(pID, pYear, pSemis)

	if newPId == "" {
		newPId = pID
	}

	baseURL := fmt.Sprintf("http://reg5.sut.ac.th/registrar/class_info_2.asp?backto=home&option=0&courseid=%s&acadyear=%s&semester=%s", newPId, pYear, pSemis)
	log.Println("request URL: ", baseURL)

	c := colly.NewCollector(
		colly.DetectCharset(),
		colly.CacheDir("./reg_cache/course"),
	)

	c.OnHTML(acTable, func(cc *colly.HTMLElement) {
		courseNameEn = cc.ChildText(getCourseNameEn)
		courseNameTh = cc.ChildText(getCourseNameTh)
		belongTo = cc.ChildText(getBelongTo)
		status = cc.ChildText(getStatus)
		tempCID = cc.ChildText(getCourseId)
		description = subDescription(cc.ChildText(getDescription))

		fmt.Println("----------- course name en", courseNameEn)
		fmt.Println("----------- course name th", courseNameTh)
		fmt.Println("----------- course name th", description)

		credit = cc.ChildText(getCredit)

		cc.ForEach(" table:nth-child(5) > tbody:nth-child(1) tr", func(_ int, el *colly.HTMLElement) {
			// mc2 := make(map[string]courseModel.Group)
			if strings.ToUpper(el.Attr("bgcolor")) == headerGroups { // checking head of group
				if shouldSumGroups(el.Text) { // sum sec time to group here
					sectionTimeTemp = append(sectionTimeTemp, courseModel.SectionTime{
						Day:      el.ChildText(getDay),
						Time:     el.ChildText(getTime),
						Room:     el.ChildText(getRoom),
						Building: el.ChildText(getBuilding),
					})
					fmt.Printf("groups t:%d \r\n", gTemp)

				} else {
					gTemp, _ = strconv.Atoi(getGroupNumber(el.Text))
					sectionTimeTemp = append(sectionTimeTemp, courseModel.SectionTime{
						Day:      el.ChildText(getDay),
						Time:     el.ChildText(getTime),
						Room:     el.ChildText(getRoom),
						Building: el.ChildText(getBuilding),
					})
					if len(el.ChildText(getOpenAmount)) > 0 {
						amount.open, _ = strconv.Atoi(el.ChildText(getOpenAmount))
					}
					if len(el.ChildText(getReserveAmount)) > 0 {
						amount.reserve, _ = strconv.Atoi(el.ChildText(getReserveAmount))
					}
					if len(el.ChildText(getRemainAmount)) > 0 {
						amount.remaining, _ = strconv.Atoi(el.ChildText(getRemainAmount))
					}

					fmt.Printf("groups b:%d \r\n", gTemp)
				}
			}
			if el.Attr("align") == "left" { // line hr in html tag
				// push group hear
				Group = append(Group, courseModel.Group{
					SecTime:   sectionTimeTemp,
					Group:     gTemp,
					Teacher:   tempDetail.Teacher,
					Final:     tempDetail.Final,
					Mid:       tempDetail.Mid,
					Note:      tempDetail.Note,
					Open:      amount.open,
					Reserved:  amount.reserve,
					Remaining: amount.remaining,
				})
				sectionTimeTemp = []courseModel.SectionTime{}
				tempDetail = courseModel.Group{}
				amount = Amount{}
			}
			if el.ChildText(checkTc) == "อาจารย์:" { // อาจารย์
				tempDetail.Teacher = el.ChildText(getTc)
			} else if el.ChildText(checkTc) == "สอบกลางภาค:" { //mid
				tempDetail.Mid = el.ChildText(getTc)
			} else if el.ChildText(checkTc) == "สอบประจำภาค:" { //fi
				tempDetail.Final = el.ChildText(getTc)
			} else if el.ChildText(checkTc) == "หมายเหตุ:" { //fi
				tempDetail.Note = el.ChildText(getTc)
			}

		}) //end loop

		Course = courseModel.CourseStruc{
			NameEn:      courseNameEn,
			NameTh:      courseNameTh,
			BelongTo:    belongTo,
			Status:      status,
			ID:          tempCID,
			Credit:      credit,
			Description: description,
			Groups:      Group,
			Year:        year,
			Semester:    semester,
		}
	})
	c.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"
	})
	c.Visit(baseURL)
	w.Header().Set("Content-type", "application/json; charset=UTF-8;")
	json.NewEncoder(w).Encode(Course)
}
