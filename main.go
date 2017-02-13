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

var actsListTemplate = template.Must(template.New("list").Parse(`
<!DOCTYPE html>
<html>
  <head>
    <link rel="stylesheet" href="/css/styles.css" />
    <title>Acts of Peace</title>
  </head>
  <body>
    <table>
      <thead>
        <tr>
          <td>Title</td>
          <td>Description</td>
          <td>Focus Area</td>
        </tr>
      </thead>
      <tbody>
      {{range .}}
        <tr>
          <td>{{.Title}}</td>
          <td>{{.Description}}</td>
          <td>{{.FocusArea}}
        </tr>
      {{end}}
      </tbody>
    </table>
    <form action="/submit" method="post">
      <div>
        <label for="title">Title</label>
        <input type="text" name="title" />
      </div>
      <div>
        <label for="description">Description</label>
        <textarea name="description" rows="6" cols="60"></textarea>
      </div>
      <div>
        <label for="focusArea">Focus Area</label>
        <select name="focusArea">
          <option disabled selected value> -- select a focus area -- </option>
          <option value="Education and Community Development">Education and Community Development</option>
          <option value="Protecting the Environment">Protecting the Environment</option>
          <option value="Alleviating Extreme Poverty">Alleviating Extreme Poverty</option>
          <option value="Global Health and Wellness">Global Health and Wellness</option>
          <option value="Non-proliferation & Disarmament">Non-proliferation & Disarmament</option>
          <option value="Human Rights For All">Human Rights For All</option>
          <option value="Ending Racism & Hate">Ending Racism & Hate</option>
          <option value="Advancing Women and Children">Advancing Women and Children</option>
          <option value="Clean Water For Everyone">Clean Water For Everyone</option>
          <option value="Conflict Resolution">Conflict Resolution</option>
        </select>
      </div>
      <div><input type="submit" value="Submit"></div>
    </form>
  </body>
</html>
`))

func submit(w http.ResponseWriter, r *http.Request) {
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
  http.Redirect(w, r, "/", http.StatusFound)
}
