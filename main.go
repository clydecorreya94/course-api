package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

type Author struct {
	AuthorName string `json:"authorname"`
	Website    string `json:"website"`
}

var courses []Course

func (c *Course) IsEmpty() bool {
	return c.CourseName == "" && c.CourseId == ""
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello... Welcome to Courses APIs</h2>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all courses!")
	w.Header().Set("Content-Type", "content-application")

	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one course!")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("The searched element is not available!!")
	return
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating new course!")
	w.Header().Set("Content-Type", "applicatiob/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("The body is nil")
	}

	var course Course

	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("The json body is empty")
		return
	}

	rand.Seed(time.Now().Unix())
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)

}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating course!")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	if params["id"] != "" {
		for index, course := range courses {
			if course.CourseId == params["id"] {
				courses = append(courses[:index], courses[index+1:]...)

				var course Course
				_ = json.NewDecoder(r.Body).Decode(&course)
				course.CourseId = params["id"]
				courses = append(courses, course)

				json.NewEncoder(w).Encode(courses)
				return
			}
		}
	} else {
		json.NewEncoder(w).Encode("The Course id is not found")
		return
	}
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete course!")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode(courses)

			return
		}

	}
}

func main() {
	fmt.Println("Welcome to the Courses API")
	courses = append(courses,
		Course{
			CourseId:    "1",
			CourseName:  "Java",
			CoursePrice: 100,
			Author: &Author{
				AuthorName: "James",
				Website:    "java.com",
			},
		},
		Course{
			CourseId:    "2",
			CourseName:  "Go",
			CoursePrice: 100,
			Author: &Author{
				AuthorName: "Govind",
				Website:    "golang.com",
			},
		},
		Course{
			CourseId:    "3",
			CourseName:  "C++",
			CoursePrice: 100,
			Author: &Author{
				AuthorName: "Steven",
				Website:    "cpp.com",
			},
		},
		Course{
			CourseId:    "4",
			CourseName:  "C",
			CoursePrice: 100,
			Author: &Author{
				AuthorName: "Spilberg",
				Website:    "c.com",
			},
		},
	)

	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/courses/{id}", deleteOneCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4000", r))
}
