package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func CreateUUID() string {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
}

// GenerateHTML fucntion
func GenerateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

// parse HTML templates
// pass in a list of file names, and get a template
func ParseTemplateFiles(filenames ...string) *template.Template {
	var files []string
	var t *template.Template
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("./templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return t
}
