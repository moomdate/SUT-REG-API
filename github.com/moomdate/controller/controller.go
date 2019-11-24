package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/moomdate/courseEntity"
)

const (
	headerGroups  = "#F5F5F5"
	acTable       = "body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3)  "
	getCourseName = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > b:nth-child(1) > font:nth-child(1)"
	getCredit     = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(4) > td:nth-child(3) > font:nth-child(1)"
	getDay        = "td:nth-child(4)"
	getTime       = "td:nth-child(5)"
	getRoom       = "td:nth-child(6) "
	getBuilding   = "td:nth-child(7)"
	checkTc       = "td:nth-child(4) > font:nth-child(1)"
	getTc         = "td:nth-child(5) > font:nth-child(1)"
)

//type Courses courseEntity.CourseStruc // use struct

func InitServer() {
	router := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.HandleFunc("/api/{id}/{year}/{semester}", scraping).Methods("GET")

	http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(router))
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
		explodeStr = "00"
	}
	return explodeStr
}
func digCourseCode(ID string, Year string, sem string) string {
	var link string
	var tim string
	if len(ID) < 1 {
		return ""
	}
	mainLink := "http://reg4.sut.ac.th/registrar/class_info_1.asp?coursestatus=O00&facultyid=all&maxrow=1&acadyear=" + Year + "&semester=" + sem + "&coursecode=" + ID
	scrapLink := colly.NewCollector(
		colly.CacheDir("./reg_cacheCourse"),
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
	return out
}
func scraping(w http.ResponseWriter, r *http.Request) {
	var Course courseEntity.CourseStruc
	var courseName, credit string
	var gTemp string
	countInGroup := 0
	getParam := mux.Vars(r)
	pID := getParam["id"]
	pYear := getParam["year"]
	pSemis := getParam["semester"]
	tempCID := pID
	pID = digCourseCode(pID, pYear, pSemis)
	fmt.Println("pID is ", pID)
	baseURL := fmt.Sprintf("http://reg4.sut.ac.th/registrar/class_info_2.asp?backto=home&option=0&courseid=%s&acadyear=%s&semester=%s", pID, pYear, pSemis)
	fmt.Println("base url is ", baseURL)

	bigMC := make(map[string]*courseEntity.GroupBig)
	c := colly.NewCollector(
		colly.DetectCharset(),
		colly.CacheDir("./reg_cache"),
	)
	c.OnHTML(acTable, func(cc *colly.HTMLElement) {
		courseName = cc.ChildText(getCourseName)
		credit = cc.ChildText(getCredit)
		cc.ForEach(" table:nth-child(5) > tbody:nth-child(1) tr", func(_ int, el *colly.HTMLElement) {
			mc2 := make(map[string]courseEntity.Group)
			if strings.ToUpper(el.Attr("bgcolor")) == headerGroups { // checking head of group
				if shouldSumGroups(el.Text) { // sum sec time to group here
					countInGroup++ // section time 2 3 4 ...
					bigMC[gTemp].SecTime[strconv.Itoa(countInGroup)] = courseEntity.Group{
						Day:      el.ChildText(getDay),
						Time:     el.ChildText(getTime),
						Room:     el.ChildText(getRoom),
						Building: el.ChildText(getBuilding),
					}
				} else {
					countInGroup = 0
					gTemp = getGroupNumber(el.Text)
					mc2["0"] = courseEntity.Group{ // main groups
						Day:      el.ChildText(getDay),
						Time:     el.ChildText(getTime),
						Room:     el.ChildText(getRoom),
						Building: el.ChildText(getBuilding),
					}
					countInGroup = 0

					bigMC[gTemp] = &courseEntity.GroupBig{
						SecTime: mc2,
					}
				}
			}
			if el.ChildText(checkTc) == "อาจารย์:" { // อาจารย์
				bigMC[gTemp].Teacher = el.ChildText(getTc)
			} else if el.ChildText(checkTc) == "สอบกลางภาค:" { //mid
				bigMC[gTemp].Mid = el.ChildText(getTc)
			} else if el.ChildText(checkTc) == "สอบประจำภาค:" { //fi
				bigMC[gTemp].Final = el.ChildText(getTc)
			} else if el.ChildText(checkTc) == "หมายเหตุ:" { //fi
				bigMC[gTemp].Note = el.ChildText(getTc)
			}

		}) //end loop
		Course = courseEntity.CourseStruc{
			Name:   courseName,
			ID:     tempCID,
			Credit: credit,
			Groups: bigMC,
		}
	})
	c.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "charset=UTF-8"
	})
	c.Visit(baseURL)
	w.Header().Set("Content-type", "application/json; charset=UTF-8;")
	json.NewEncoder(w).Encode(Course)
}
