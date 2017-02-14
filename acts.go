// package name
package acts

// import statement
import (
  // from Go standard libraries
  "net/http"
  "html/template"
  "time"

  // App Engine specific
  "appengine"
  "appengine/datastore"
)

// Act struct that we will fill with form data and put in the datastore.
// Date is used only for sorting query results
type Act struct {
  Title string
  Description string
  FocusArea string
  Date time.Time
}

// init is called before the application starts
func init() {
  // register handler functions for the routes we use
  http.HandleFunc("/", root)
  http.HandleFunc("/submit", submit)
}

// Helper function to create new datastore keys
func actsOfPeaceKey(c appengine.Context) *datastore.Key {
  return datastore.NewKey(c, "ActsOfPeace", "default_acts_of_peace", 0, nil)
}

func root(w http.ResponseWriter, r *http.Request) {
  // Context holds information about the request
  // 'Package context defines the Context type, which carries deadlines, cancelation signals, and other request-scoped values across API boundaries and between processes'
  c := appengine.NewContext(r)
  // Query for Act entities.
  // Using 'Ancestor' speeds up the query. It limits the results to the
  // specified key and its descendants. In this case, we are only looking
  // at Acts that have the actsOfPeaceKey as their ancestor
  // Order by Date, descending
  q := datastore.NewQuery("Act").Ancestor(actsOfPeaceKey(c)).Order("-Date")

  // declare an empty slice called acts
  var acts []Act
  // GetAll runs the query and puts the results in the acts slice
  if _, err := q.GetAll(c, &acts); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  // Execute the template with the acts slice as the data context
  if err := actsListTemplate.Execute(w, acts); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// Parse the acts list template
var actsListTemplate, _ = template.ParseFiles("templates/acts.html")

func submit(w http.ResponseWriter, r *http.Request) {
  // Only insert in the datastore if the request method was POST
  if r.Method == "POST" {
    c := appengine.NewContext(r)

    // Put form data into an Act struct. &Act is a pointer to this new struct
    act := &Act{
      Title: r.FormValue("title"),
      Description: r.FormValue("description"),
      FocusArea: r.FormValue("focusArea"),
      Date: time.Now(),
    }
    // Create a new datastore Incomplete Key for the Act entity
    actKey := datastore.NewIncompleteKey(c, "Act", actsOfPeaceKey(c))
    // Put the struct into the datastore with the key we just created
    _, err := datastore.Put(c, actKey, act)
    // Report any errors with the db transaction
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
  // Redirect to the root path
  http.Redirect(w, r, "/", http.StatusFound)
}
