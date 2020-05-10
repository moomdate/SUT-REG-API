package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../courseModel"
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
	getOpenAmount    = "td:nth-child(9)"
	getReserveAmount = "td:nth-child(10)"
	getRemainAmount  = "td:nth-child(11)"
)

//type Courses courseEntity.CourseStruc // use struct

func InitServer() {
	router := mux.NewRouter()

	router.HandleFunc("/major/{id}", courseMajor).Methods("GET")
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
	firstIndex :=0
	if text[0]==194{ // remove double space
		firstIndex = 2
	}

	nameEn := text[firstIndex : thaiIndex-2] // remove double space
	nameTh := text[thaiIndex:len(text)]
	return nameEn, nameTh
}
func courseMajor(w http.ResponseWriter, r *http.Request) {
	getParam := mux.Vars(r)
	id := getParam["id"]
	var majorCourseList = []courseModel.MajorCourse{}
	var courseTemp = courseModel.MajorCourse{}
	mainLink := "http://reg4.sut.ac.th/registrar/program_info_1.asp?programid=" + id
	scrapLink := colly.NewCollector(
		colly.CacheDir("./reg_cache/majorCourse"),
	)
	scrapLink.SetRequestTimeout(5 * time.Second)
	scrapLink.OnHTML("body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3) > table:nth-child(4) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table", func(el *colly.HTMLElement) {
		el.ForEach("tbody", func(_ int, eltr *colly.HTMLElement) {
			eltr.ForEach("tr", func(_ int, elFont *colly.HTMLElement) {

				if elFont.Attr("bgcolor") == "#FFFFDE" {
					elFont.ForEach("td > font", func(_ int, elFontEl *colly.HTMLElement) {
						if len(elFontEl.Text) == 10 {

							tempText := strings.Split(elFontEl.Text, " ")
							tempCourseIn, err := strconv.Atoi(tempText[0])
							if err != nil {
								log.Fatal("get course id err")
							}
							courseTemp.CourseId = tempCourseIn
						} else if len(elFontEl.Text) == 9 {
							courseTemp.Credit = elFontEl.Text
							majorCourseList = append(majorCourseList, courseTemp)
							courseTemp = courseModel.MajorCourse{}

						} else {
							nameEn, nameTh := splitCourseName(elFontEl.Text)
							courseTemp.CourseNameEn = nameEn
							courseTemp.CourseNameTh = nameTh
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

func scraping(w http.ResponseWriter, r *http.Request) {
	var Course courseModel.CourseStruc
	var Group []courseModel.Group

	var courseNameEn, courseNameTh, belongTo, status, credit string
	var sectionTimeTemp []courseModel.SectionTime
	var gTemp int
	var tempDetail courseModel.Group
	//countInGroup := 0
	var amount Amount
	getParam := mux.Vars(r)
	pID := getParam["id"]
	pYear := getParam["year"]
	pSemis := getParam["semester"]
	tempCID := pID
	pID = digCourseCode(pID, pYear, pSemis)
	baseURL := fmt.Sprintf("http://reg3.sut.ac.th/registrar/class_info_2.asp?backto=home&option=0&courseid=%s&acadyear=%s&semester=%s", pID, pYear, pSemis)
	fmt.Println(baseURL)
	//var secTime []*courseModel.SectionTime;

	//bigMC := make([]*courseModel.Group, 100)

	c := colly.NewCollector(
		colly.DetectCharset(),
		colly.CacheDir("./reg_cache/course"),
	)

	c.OnHTML(acTable, func(cc *colly.HTMLElement) {
		courseNameEn = cc.ChildText(getCourseNameEn)
		courseNameTh = cc.ChildText(getCourseNameTh)
		belongTo = cc.ChildText(getBelongTo)
		status = cc.ChildText(getStatus)
		fmt.Println("----------- course name en", courseNameEn)
		fmt.Println("----------- course name th", courseNameTh)
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
				tempDetail.Mid = el.ChildText(getTc)
			} else if el.ChildText(checkTc) == "สอบประจำภาค:" { //fi
				tempDetail.Final = el.ChildText(getTc)
			} else if el.ChildText(checkTc) == "หมายเหตุ:" { //fi
				tempDetail.Note = el.ChildText(getTc)
			}

		}) //end loop

		Course = courseModel.CourseStruc{
			NameEn:   courseNameEn,
			NameTh:   courseNameTh,
			BelongTo: belongTo,
			Status:   status,
			ID:       tempCID,
			Credit:   credit,
			Groups:   Group,
		}
	})
	c.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"
	})
	c.Visit(baseURL)
	w.Header().Set("Content-type", "application/json; charset=UTF-8;")
	json.NewEncoder(w).Encode(Course)
}
