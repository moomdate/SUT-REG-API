# SUT REG API V.2.
[The first version here](https://github.com/moomdate/API-SUT-REG)
```json

{

    "Name": "CHEMISTRY LABORATORY II",
    "courseID": "102104",
    "Credit": "1 (0-3-0)",
    "Groups": {
        "01": {
            "SecTime": {
                "0": {
                    "Day": "พุธ",
                    "Time": "13:00-16:00",
                    "Room": "F2-111",
                    "Building": "F2"
                }
            },
            "Teacher": "ผู้ช่วยศาสตราจารย์ ดร.ธนพร แม่นยำ",
            "Mid": "27 ธ.ค. 2562 เวลา 15:00 - 17:00 อาคาร B ห้อง N",
            "Final": "18 ก.พ. 2563 เวลา 09:00 - 12:00 อาคาร B ห้อง N",
            "Note": "เริ่มเรียนสัปดาห์ที่ 2"
        },
        "02": {
            "SecTime": {
                "0": {
                    "Day": "พุธ",
                    "Time": "13:00-16:00",
                    "Room": "F2-112",
                    "Building": "F2"
                }
            },
            "Teacher": "ผู้ช่วยศาสตราจารย์ ดร.ธนพร แม่นยำ อาจารย์ ดร.กมลวัช งามเชื้อ",
            "Mid": "27 ธ.ค. 2562 เวลา 15:00 - 17:00 อาคาร B ห้อง N",
            "Final": "18 ก.พ. 2563 เวลา 09:00 - 12:00 อาคาร B ห้อง N",
            "Note": "เริ่มเรียนสัปดาห์ที่ 2"
        }
    }

}

```
### What's new.
- Move from php to Golang.
- Remove json file. use url filter instead of json file.
- None update Json file.
- Caching reg.sut.ac.th. able to use when can't access Reg Page.
- Faster than before.
- Easily updated.

### How to use.
 
#### `GET` http://moomdate.tech/api/:courseid/:acadyear/:semester
`http://moomdate.tech/api/:courseid/:acadyear/:semester`
#### PATH VARIABLES
##### courseid = Course code such as `52310`
##### acadyear = Acadmic Year in Buddhist year such as `2556` `2557`
##### semester = Semester of course such as `1` `2` `3`


#### Prepare to use
- `$ go get github.com/gocolly/colly && go get github.com/gorilla/mux`
## Defualt port is `8000`
#### Deployment
 - Whatever you'd like.