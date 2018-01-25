package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "io"
    "os"
    "encoding/json"
    "strconv"
	"database/sql"
     _ "github.com/mattn/go-sqlite3"
     "net"
)

var acc_token = "access_token=0fe8ce5727b53609c8a90de22acdf00e3b7dd1dc10c4b1044723389fba527c09" //jimmyler API token

type JSON_contain struct {
	Action   string `json:"action"`
	Status   string `json:"status"`
	Value    string `json:"value"`
}

type Shots_info struct {
	Id          int	   `json:"id"`
	Title       string `json:"title"`
	Desc	    string `json:"description"`
	Html_url    string `json:"html_url"`
	Img_url     Links  `json:"images"`
}
type Links struct {
	Hidpi  string
	Normal string
	Teaser string
}
type json_info struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
}
func check(e error) {
    if e != nil {
        log.Fatal(e)
        fmt.Println(e)
    }
}

func listShots () {
	url := "https://api.dribbble.com/v2/"
    client := &http.Client{}
    r, _ := http.NewRequest("GET", url+"user/shots?"+acc_token, nil)
    resp, _ := client.Do(r)
    f, err := ioutil.ReadAll(resp.Body)    
    check(err)
//    fmt.Println(string(f))
    var result []Shots_info
	err = json.Unmarshal(f, &result)
    for i := range result {
//    	fmt.Println(string(result[i].Title)+" , id = "+strconv.Itoa(result[i].Id)+" , desc = "+string(result[i].Desc))
    	getShot(strconv.Itoa(result[i].Id))
    }
}

func getShot (img_id string) {
	url := "https://api.dribbble.com/v2/"
    client := &http.Client{}
    r, _ := http.NewRequest("GET", url+"/shots/"+img_id+"?"+acc_token, nil)
    resp, _ := client.Do(r)
    f, err := ioutil.ReadAll(resp.Body)
    check(err)
//    fmt.Println(string(f))
    var result Shots_info
	err = json.Unmarshal(f, &result)
	check(err)
//	fmt.Println(result.Img_url.Normal)
	db, _ := sql.Open("sqlite3", "sqlite/image.db")
	tx, _ := db.Begin()
	insert_sql := "Insert Into shots_db (id,title,desc,img_url,html_url,filename) values("+strconv.Itoa(result.Id)+",'"+result.Title+"','"+result.Desc+"','"+result.Img_url.Normal+"','"+result.Html_url+"','"+result.Title+".jpg')"
	tx.Exec(insert_sql)
	tx.Commit()
	defer db.Close()
	response, e := http.Get(result.Img_url.Normal)
	check(e)
	defer response.Body.Close()
	file, err := os.Create( "downloaded_shots/"+result.Title + ".jpg")
	check(err)
	_, err = io.Copy(file, response.Body)
	check(err)
	file.Close()
//	fmt.Println(".")
}

func main() {
	listShots()
	fmt.Println("Launching server...")
	listener, _ := net.Listen("tcp", "127.0.0.1:8081")
	for {
        conn, err := listener.Accept()
        check(err)
        encoder := json.NewEncoder(conn)
        decoder := json.NewDecoder(conn)
		
		var json_req JSON_contain
		var dt []Shots_info
		decoder.Decode(&json_req)
		if json_req.Action == "query" {
			db, _ := sql.Open("sqlite3", "sqlite/image.db")
			q, err := db.Query("SELECT title,desc,html_url FROM shots_db WHERE title LIKE '%" + json_req.Value + "%' OR desc Like '%" + json_req.Value + "%'")
//			fmt.Println("SELECT title,desc,html_url FROM shots_db WHERE title LIKE '%" + json_req.Value + "%' OR desc Like '%" + json_req.Value + "%'")
			check(err)
			for q.Next() {
				var title string
				var desc string
				var html_url string
				q.Scan(&title, &desc, &html_url)
				tmp_Shots_info := Shots_info{
				    Title:	title,
				    Desc: 	desc,
				    Html_url: html_url}
//			    fmt.Println(tmp_Shots_info)
			    dt = append(dt, tmp_Shots_info)
			}
		}
		encoder.Encode(dt)
		conn.Close() // we're finished
	}
}