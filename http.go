package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"html/template"
	"log"
	"net"
	"net/http"
)

var decoder = schema.NewDecoder()

type homePageData struct {
	Title string
	// Neighbors Neighbors
	Data map[string]string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.New("index").Parse(index_tpl)
	if err != nil {
		log.Println("Could not parse template", err)
		fmt.Fprintln(w, "Problem parsing template", err)
		return
	}

	data := new(homePageData)
	data.Title = "NTP Pool Stats"

	// fmt.Fprintf(w, "%s\t%s\t%v\n", neighbor, data.State, data.Updates)
	// fmt.Printf("TMPL %s %#v\n", tmpl, tmpl)

	tmpl.Execute(w, data)

}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)

	_ = r.ParseForm()
	ip := net.ParseIP(vars["ip"])

	if ip == nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Bad IP address")
		return
	}

	type Options struct {
		Points int    `schema:"points"`
		Type   string `schema:"type"`
	}

	log.Println("Looking for IP", ip)

	srv := getServer(&ip)
	if srv == nil {
		w.WriteHeader(404)
		return
	}

	log.Println("looking for data for server", srv.Id)

	history, err := srv.GetData()
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error fetching recent server data", err)
		fmt.Fprintln(w, "Could not fetching recent server data", err)
		return
	}

	scores := history.History

	options := new(Options)
	decoder.Decode(options, r.Form)

	if options.Points > 0 {
		log.Println("Sampling points/method", options.Points, options.Type, len(scores))
		switch options.Type {
		case "", "sample":
			scores = scores.Sample(options.Points)
		case "worst":
			scores = scores.WorstOffset(options.Points)
		default:
			w.WriteHeader(500)
			fmt.Fprintln(w, "invalid 'sample' parameter")
			return
		}

		log.Println("Now has", len(scores))
		history.History = scores
	}

	history.Server = map[string]string{"ip": ip.String()}

	js, err := json.Marshal(history)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error generating json", err)
		fmt.Fprintln(w, "Could not generate JSON")
		return
	}

	w.Header().Add("Cache-Control", "max-age=10400")

	fmt.Fprint(w, string(js), "\n")
}

func httpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/data/{ip:[0-9a-f.:]+}", ApiHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8085", nil)
}

const index_tpl = `<!DOCTYPE html>
<html>
	<head><title>NTP Pool Stats</title>
		<link href="http://st.pimg.net/cdn/libs/bootstrap/2/css/bootstrap.min.css" rel="stylesheet">
		<style>
			html,
			body {
			  margin: 10px;
			  margin-top: 20px;
			}
		</style>
	</head>
	<body>

	<h1>{{.Title}}</h1>

	<p>This is an internal service used by the <a href="http://www.ntppool.org/">NTP Pool</a>.</p>

	<script src="http://st.pimg.net/cdn/libs/jquery/1.8/jquery.min.js"></script>
	<script src="http://st.pimg.net/cdn/libs/underscore/1/underscore-min.js"></script>
	<script>
		"use strict";
		(function ($) {
	})(jQuery);
	</script>

	</body>
</html>
`
