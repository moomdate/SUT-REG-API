package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reg-api/courseModel"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const port = 8081
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
	getCourseId      = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > b:nth-child(1) > font:nth-child(1)"
	getDescription   = "table:nth-child(6) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > font:nth-child(3)"
	getOpenAmount    = "td:nth-child(9)"
	getReserveAmount = "td:nth-child(10)"
	getRemainAmount  = "td:nth-child(11)"
	Baseroot         = "body > table > tbody > tr:nth-child(1) > td:nth-child(3) > font > b > div > table > tbody > tr > td > table > tbody > tr > td > font > table > tbody > tr > td > table > tbody"
	TermDetail       = "body > table > tbody > tr:nth-child(1) > td:nth-child(3) > table:nth-child(3) > tbody > tr:nth-child(7) > td:nth-child(2) > font > font"
)
const baseURLCourseInMajor = "http://reg4.sut.ac.th/registrar/program_info_1.asp?programid="

//type Courses courseEntity.CourseStruc // use struct

func InitServer() {
	router := mux.NewRouter()
	router.HandleFunc("/api/find/{cid}/{year}/{semester}", findCourseAll).Methods("GET")
	router.HandleFunc("/api/course/major/{id}", courseMajor).Methods("GET")
	router.HandleFunc("/api/major/list/{id}", getMajor).Methods("GET")
	router.HandleFunc("/api/import/{stdid}/{acadyear}/{semester}", importFormReg).Methods("GET")
	router.HandleFunc("/api/{id}/{year}/{semester}", scraping).Methods("GET")

	mcors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in producti
		Debug: true,
	})
	handler := mcors.Handler(router)
	fmt.Print("server port:", port)
	http.ListenAndServe(":8081", (handler))
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
func findCourseAll(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	cid := getParam["cid"]
	acadyears := getParam["year"]
	semesters := getParam["semester"]
	fmt.Print(" =====", cid+" "+acadyears+" - "+semesters)
	scrapLink := colly.NewCollector(
		colly.CacheDir("./reg_cache/majorList"),
	)
	scrapLink.OnResponse(func(r *colly.Response) {

	})
	scrapLink.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"

	})
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
	mainLink := "http://reg3.sut.ac.th/registrar/class_info_1.asp?coursestatus=O00&facultyid=all&maxrow=1&acadyear=" + Year + "&semester=" + sem + "&coursecode=" + ID
	scrapLink := colly.NewCollector(
		colly.CacheDir("./reg_cache/digCode"),
	)
	scrapLink.SetRequestTimeout(5 * time.Second)
	scrapLink.OnHTML("a[href]", func(el *colly.HTMLElement) {

		link = el.Attr("href")
		fmt.Println(link)
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
		//fmt.Println(mam)
		if strings.Contains(mam, "courseid") { // see
			out = strings.Split(mam, "=")[1] // get number
			//fmt.Println("see out is :", out)
			break
		}
	}
	//fmt.Println("course id", out)
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
			//fmt.Println(i, r, string(r))
			thaiIndex = i
			break
		}
	}
	firstIndex := 0
	if text[0] == 194 { // remove double space
		firstIndex = 2
	}
	//fmt.Println("errr------------>",text,firstIndex,thaiIndex-2)
	nameEn := text[firstIndex : thaiIndex-2] // remove double space
	nameTh := text[thaiIndex:len(text)]
	return nameEn, nameTh
}
func getNumber(strNumber string) int {
	//for i,r := range strNumber {
	//	fmt.Print(i,r,string(i))
	//}
	tempStr := strNumber[2:len(strNumber)]
	result, err := strconv.Atoi(tempStr)
	if err != nil {
		log.Fatal("program id is err")
	}
	return result
}

func importFormReg(w http.ResponseWriter, r *http.Request) {
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
	c := colly.NewCollector()
	var cid, gid, cname, ver string

	var Data []courseModel.CourseDetail
	c.OnHTML("body > table > tbody > tr:nth-child(1) > td:nth-child(3) > font > b > div > table > tbody > tr > td > table > tbody > tr > td > font > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr", func(e *colly.HTMLElement) {
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
func getMajor(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	id := getParam["id"]
	//mainLink := "test"
	scrapLink := colly.NewCollector(
		colly.CacheDir("./reg_cache/majorList"),
	)
	scrapLink.OnResponse(func(r *colly.Response) {

	})
	scrapLink.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"

	})
	var majorList = []courseModel.MajorModel{}
	var major = courseModel.MajorModel{}
	getCredit := false
	scrapLink.OnHTML("body > table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > table:nth-child(6) > tbody:nth-child(1)", func(el *colly.HTMLElement) {
		el.ForEach("tr", func(_ int, elTr *colly.HTMLElement) {
			elTr.ForEach("td", func(_ int, elTr2 *colly.HTMLElement) {
				if elTr2.Attr("width") == "20" { // program id
					tempStr := elTr2.Text
					//fmt.Println("string:",tempStr)
					major.ProgramId = getNumber(tempStr)
					//fmt.Println(elTr2.Text)
				} else if elTr2.Attr("valign") == "TOP" && len(elTr2.Text) > 5 { // course
					getCredit = true
					major.Course = elTr2.Text[0 : len(elTr2.Text)-1]
					//fmt.Println(elTr2.Text)
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
	scrapLink.Request("POST",
		"http://reg4.sut.ac.th/registrar/program_info.asp",
		strings.NewReader("facultyid="+id+"&f_cmd=&Levelid=1&f_cmd=&Acadyear=-1&f_cmd="),
		nil,
		http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}})

	w.Header().Set("Content-type", "application/json; charset=UTF-8;")
	json.NewEncoder(w).Encode(majorList)
}

func courseMajor(w http.ResponseWriter, r *http.Request) {
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
							log.Println("[courseOfMajor] CourseId :", elFontEl.ChildText("font"))
							tempText := strings.Split(elFontEl.ChildText("font"), " ")
							tempCourseIn, err := strconv.Atoi(tempText[0])
							if err != nil {
								log.Fatal("get course id err")
							}
							courseTemp.CourseId = tempCourseIn
						} else if elFontEl.Attr("valign") == "TOP" && state == "credit" {
							log.Println("[courseOfMajor] Credit :", elFontEl.ChildText("font"))
							courseTemp.Credit = elFontEl.ChildText("font")
							majorCourseList = append(majorCourseList, courseTemp)
							courseTemp = courseModel.MajorCourse{}
							state = "courseId"
						} else if elFontEl.Attr("valign") == "CENTER" {
							nameEn, nameTh := splitCourseName(elFontEl.ChildText("font"))
							log.Println("[courseOfMajor] Course name:", nameEn, nameTh)
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
			//fmt.Println("--->",i, r, string(r))
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
			position := strings.Index(textTemp, dateTimeStr) // หา index ของ word match
			if len(textTemp) <= patternLen {
				break
			}
			textScope := text[0:patternLen] // 0, length

			fmt.Println("text scope->", textScope)

			textOriLen := len(textTemp)
			//cutAt := len(textScope)
			fmt.Println("text org->", textTemp)
			textTemp = textTemp[patternLen+1 : textOriLen]
			position = strings.Index(textTemp, dateTimeStr)
			Building := textTemp[0:position]
			fmt.Println("new text->", Building)

			fmt.Println("new position:", position)
			fmt.Println("xxxxxx--", textTemp[0:position])
		}
		fmt.Printf("=============================================")
	}
	//fmt.Println(text)
}
func subDescription(text string) string {
	split := strings.Split(text, "หมายเหตุเรียน")
	return split[0]
}
func scraping(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	log.Println("host Requeter:", host)

	var Course courseModel.CourseStruc
	var Group []courseModel.Group

	var courseNameEn, courseNameTh, belongTo, status, credit, description string
	var sectionTimeTemp []courseModel.SectionTime
	var gTemp int
	var tempDetail courseModel.Group
	//countInGroup := 0
	var amount Amount
	getParam := mux.Vars(r)
	pID := getParam["id"]
	pYear := getParam["year"]
	pSemis := getParam["semester"]
	var tempCID string
	newPId := digCourseCode(pID, pYear, pSemis)

	if newPId == "" {
		newPId = pID
	}

	baseURL := fmt.Sprintf("http://reg3.sut.ac.th/registrar/class_info_2.asp?backto=home&option=0&courseid=%s&acadyear=%s&semester=%s", newPId, pYear, pSemis)
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
				fmt.Println("====================== ,", gTemp)
				fmt.Printf("%+v\n", sectionTimeTemp)
				//fmt.Println(tempDetail)
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
				//getDateFormCoursegetDateFormCourse(el.ChildText(getTc))
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
		}
	})
	c.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"
	})
	c.Visit(baseURL)
	w.Header().Set("Content-type", "application/json; charset=UTF-8;")
	json.NewEncoder(w).Encode(Course)
}
