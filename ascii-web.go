package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Content struct {
	Asciiartp string
}

var tpl *template.Template

func main() {
	// 1. we need to create and register a request handler
	// response writer is what the server will respond with from any request
	//http.HandleFunc("/", display)
	fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs)) // handling the CSS

	tpl,_ = template.ParseGlob("static/*.html")
	http.HandleFunc("/", home)

	http.HandleFunc("/ascii-art", display)
	log.Fatal(http.ListenAndServe(":5000", nil))
	// log.fatal: it is going to print to the console and just kill the program inn case of crash
}

func home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.Error(writer, "404 not found.", http.StatusNotFound)
		return
	}

	switch request.Method {
	case "GET":
		err := tpl.ExecuteTemplate(writer, "index.html", nil)
		if err != nil {
			log.Println(err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
		}
		
	default:
		fmt.Fprintf(writer, "Sorry, only GET methods are supported.")
	}
}

func display(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		template, _ := template.ParseFiles("./static/index.html")
		input := r.FormValue("inputtext")
		input = strings.Replace(input, "\r\n", "\\n", -1)
		banner := r.FormValue("bannertype")
		var checkbanner int
		if banner == "standard" || banner == "thinkertoy" || banner == "shadow" {
			checkbanner = 1
		} else{
			checkbanner = 0
		}
		if input == "" || checkbanner == 0 {
			log.Println("input or banner error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}



		page := Content{Asciiartp: asciiart(input, banner)}
		template.Execute(w, page)

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
func asciiart(args string, banner string) string {
	//receiving the argument
	Arg := args

	if Arg[0] == 92 && Arg[1] == 110 {
		return ""
	}
	runes := []rune(Arg)
	lrune := len(runes) - 1
	flag := false

	for _, arune := range Arg {
		if arune < 32 || arune > 126 {
			return ""
		}
	}
	//reading file line by line
	bytes, _ := os.ReadFile(banner + ".txt")
	str := string(bytes)

	line := strings.Split(str, "\n")
	var new []string

	for i := range runes {
		if runes[i] == rune(92) && runes[i+1] == rune(110) {
			flag = true
			if i == lrune-1 {
				new = append(new, string(runes[0:i]))
			} else {
				new = append(new, string(runes[0:i]))
				new = append(new, "\n")
				new = append(new, string(runes[i+2:]))
			}
		}
	}

	newstr := strings.Join(new, "")
	newstr2 := strings.Split(newstr, "\n")
	newstr = ""
	if flag {
		for y := 0; y < len(newstr2); y++ {
			newstr = newstr + Printword([]rune(newstr2[y]), line) + "\n"
		}
	} else {
		return Printword(runes, line)
	}
	return newstr
}

func Printword(runes []rune, line []string) string {
	last := ""
	for k := 0; k <= 7; k++ {
		var print []string
		for i := range runes {
			j := int((runes[i]-32)*9 + 1)
			print = append(print, line[j+k])
		}
		if k == 7 {
			last = last + strings.Join(print, "")
		} else {
			last = last + strings.Join(print, "") + "\n"
		}
	}
	return last
}