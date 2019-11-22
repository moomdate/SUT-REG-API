# SUT REG API V.2.
[The first version here](https://github.com/moomdate/API-SUT-REG)

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
- `get get github.com/gocolly/colly && go get github.com/gorilla/mux`
## Defualt port is `8000`
#### Deployment
 - Whatever you'd like.