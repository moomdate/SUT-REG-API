package controller

import (
	"../courseModel"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const port  =  8080
const ( //child access
	headerGroups  = "#F5F5F5"
	acTable       = "body > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(3)  "
	getCourseName = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > b:nth-child(1) > font:nth-child(1)"
	getCredit     = "table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(2) > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(4) > td:nth-child(3) > font:nth-child(1)"
	getDay        = "td:nth-child(4)" //get Date
	getTime       = "td:nth-child(5)" //get Time
	getRoom       = "td:nth-child(6) " //get room
	getBuilding   = "td:nth-child(7)" //get building
	checkTc       = "td:nth-child(4) > font:nth-child(1)" //check teacher
	getTc         = "td:nth-child(5) > font:nth-child(1)" //get teacher
)

//type Courses courseEntity.CourseStruc // use struct

func InitServer() {
	router := mux.NewRouter()

	router.HandleFunc("/api/{id}/{year}/{semester}", scraping).Methods("GET")
	mcors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	handler := mcors.Handler(router)
	fmt.Print("server port:",port)
	http.ListenAndServe(":8080"  , (handler))
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
	//fmt.Println("course id", out)
	return out
}
func scraping(w http.ResponseWriter, r *http.Request) {
	var Course courseModel.CourseStruc
	var courseName, credit string
	var gTemp string
	countInGroup := 0
	getParam := mux.Vars(r)
	pID := getParam["id"]
	pYear := getParam["year"]
	pSemis := getParam["semester"]
	tempCID := pID
	pID = digCourseCode(pID, pYear, pSemis)
	baseURL := fmt.Sprintf("http://reg3.sut.ac.th/registrar/class_info_2.asp?backto=home&option=0&courseid=%s&acadyear=%s&semester=%s", pID, pYear, pSemis)

	bigMC := make(map[string]*courseModel.GroupBig)
	c := colly.NewCollector(
		colly.DetectCharset(),
		colly.CacheDir("./reg_cache"),
	)
	c.OnHTML(acTable, func(cc *colly.HTMLElement) {
		courseName = cc.ChildText(getCourseName)
		credit = cc.ChildText(getCredit)
		cc.ForEach(" table:nth-child(5) > tbody:nth-child(1) tr", func(_ int, el *colly.HTMLElement) {
			mc2 := make(map[string]courseModel.Group)
			if strings.ToUpper(el.Attr("bgcolor")) == headerGroups { // checking head of group
				if shouldSumGroups(el.Text) { // sum sec time to group here
					countInGroup++ // section time 2 3 4 ...
					bigMC[gTemp].SecTime[strconv.Itoa(countInGroup)] = courseModel.Group{
						Day:      el.ChildText(getDay),
						Time:     el.ChildText(getTime),
						Room:     el.ChildText(getRoom),
						Building: el.ChildText(getBuilding),
					}
				} else {
					countInGroup = 0
					gTemp = getGroupNumber(el.Text)
					mc2["0"] = courseModel.Group{ // main groups
						Day:      el.ChildText(getDay),
						Time:     el.ChildText(getTime),
						Room:     el.ChildText(getRoom),
						Building: el.ChildText(getBuilding),
					}
					countInGroup = 0

					bigMC[gTemp] = &courseModel.GroupBig{
						SecTime: mc2,
					}
					bigMC[gTemp].Group = gTemp
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
		Course = courseModel.CourseStruc{
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
