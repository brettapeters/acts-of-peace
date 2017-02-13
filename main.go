package acts

import (
  "net/http"
  "html/template"

  "appengine"
  "appengine/datastore"
)

type Act struct {
  Title string
  Description string
  FocusArea string
}

func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/submit", submit)
}

func actsOfPeaceKey(c appengine.Context) *datastore.Key {
  return datastore.NewKey(c, "ActsOfPeace", "default_acts_of_peace", 0, nil)
}

func root(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  q := datastore.NewQuery("Act").Ancestor(actsOfPeaceKey(c))

  var acts []Act
  if _, err := q.GetAll(c, &acts); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if err := actsListTemplate.Execute(w, acts); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

var actsListTemplate, _ = template.ParseFiles("templates/acts.html")

func submit(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    c := appengine.NewContext(r)
    g := Act{
      Title: r.FormValue("title"),
      Description: r.FormValue("description"),
      FocusArea: r.FormValue("focusArea"),
    }

    key := datastore.NewIncompleteKey(c, "Act", actsOfPeaceKey(c))
    _, err := datastore.Put(c, key, &g)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }

  http.Redirect(w, r, "/", http.StatusFound)
}
