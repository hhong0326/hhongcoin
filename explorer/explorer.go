package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/hhong0326/hhongcoin/blockchain"
)

var templates *template.Template

const (
	templateDir string = "explorer/templates/"
)

type homeData struct {
	PageTitle string // template public/private에 영향
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// // data := homeData{"Home", blockchain.GetBlockChain().AllBlocks()}

	// templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		blockchain.BlockChain().AddBlock()

		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}

}
func Start(port int) {
	handler := http.NewServeMux()
	//load all that html templates
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	//update
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
