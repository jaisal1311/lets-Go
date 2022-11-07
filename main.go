package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// models
type Course struct {
	CourseId string  `json:"courseId"`
	Name     string  `json:"name"`
	Author   *Author `json:"author"`
	Price    int     `json:"price"`
}

type Author struct {
	Name    string `json:"name"`
	Website string `json:"website"`
}

// fake DB

var courses []Course

// helper func

func (c *Course) IsEmpty() bool {
	return c.Name == ""
}

// controller func

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getCourseById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, val := range courses {
		if val.CourseId == params["id"] {
			json.NewEncoder(w).Encode(val)
			return
		}
	}
	json.NewEncoder(w).Encode("404 not found")
}

func addCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please pass content of body")
		return
	}

	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("Please pass valid body")
		return
	}

	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)

	json.NewEncoder(w).Encode(course)

}

func deleteCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println(params)

	for index, val := range courses {
		if val.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode(true)

			return
		}
	}
	json.NewEncoder(w).Encode(false)

}

// main router func

func main() {
	r := mux.NewRouter()

	//seeding
	courses = append(courses, Course{CourseId: "2", Name: "ReactJS", Price: 2099, Author: &Author{Name: "John Smilga", Website: "web.dev"}})
	courses = append(courses, Course{CourseId: "4", Name: "MERN Stack", Price: 10099, Author: &Author{Name: "FCC", Website: "go.dev"}})

	r.HandleFunc("/", serveRoot).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getCourseById).Methods("GET")
	r.HandleFunc("/course", addCourse).Methods("POST")
	r.HandleFunc("/course/{id}", deleteCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4000", r))
}
