package main

import (
	"fmt"
	"log"
	"net"
    "encoding/json"
    "strconv"
    "bufio"
    "os"
)

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

func (p Shots_info) String() string {
    s := "Name : " + p.Title + "\nDescription : " + p.Desc + "\nHTML : " + p.Html_url + "\n\n"
    return s
}
func check(e error) {
    if e != nil {
        log.Fatal(e)
        fmt.Println(e)
    }
}

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	    fmt.Print("Enter the text to search: ")
	    reader := bufio.NewScanner(os.Stdin)
	    reader.Scan()
		
	    text := reader.Text()
		json_request := JSON_contain {Action:"query",Status:"0",Value: text}
//		fmt.Println(json_request)
	    encoder := json.NewEncoder(conn)
	    decoder := json.NewDecoder(conn)
	    encoder.Encode(json_request)
	    var newshots []Shots_info
	    decoder.Decode(&newshots)
	    for index := range newshots {
	        fmt.Println("Image #"+strconv.Itoa(index+1)+"\n"+newshots[index].String())
	    }
}