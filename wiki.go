package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

//funcion para verificar si esta enviando y recibiendo informacion
//func handler(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hola me gustan los %s", r.URL.Path[1:])
//}

var validPath = regexp.MustCompile("^/(edit|view|save)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// se elimina esta funcionon porque se hizo una funcion cloures. (makeHandler)
//func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
//	m := validPath.FindStringSubmatch(r.URL.Path)
//	if m == nil {
//		http.NotFound(w, r)
//		return "", errors.New("titulo de pagina no valida")
//	}
//	return m[2], nil
//}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/view/"):]
	//title, erro := getTitle(w, r)
	//if erro != nil {
	//	return
	//}

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplates(w, "view", p)

	//se comenta porque se simplica en la funcion renderTemplates
	//t, _ := template.ParseFiles("view.html")
	//t.Execute(w, p)

	//forma de prueba para mostrar las plantillas o escritura en una pagina.
	//fmt.Fprintf(w, "<h1>%s</h1> <div> %s </div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/edit/"):]
	//title, erro := getTitle(w, r)
	//if erro != nil {
	//	return
	//}

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplates(w, "edit", p)

	// se comenta porque se simplifica
	//t, _ := template.ParseFiles("edit.html")
	//t.Execute(w, p)

}
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/save/"):]
	//title, erro := getTitle(w, r)
	//if erro != nil {
	//	return
	//}

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)

	//renderTemplates(w, "edit", p)

	// se comenta porque se simplifica
	//t, _ := template.ParseFiles("edit.html")
	//t.Execute(w, p)

}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplates(w http.ResponseWriter, tmpl string, p *Page) {
	// se quita para cargar todas las plantilla colocadas en la variable en templates
	//t, err := template.ParseFiles(tmpl + ".html")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	erro := templates.ExecuteTemplate(w, tmpl+".html", p)
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	//p1 := &Page{Title: "TestPage", Body: []byte("Esta es una pagina de muestra.")}
	//p1.save()

	//p2, _ := loadPage("TestPage")
	//fmt.Println(string(p2.Body))

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hola me gustan los %s", "motos")
	//})

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
