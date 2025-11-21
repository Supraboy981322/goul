package main
import "net/http"
import "fmt"
import "log"
func main() {
http.HandleFunc("/", listener)
fmt.Println("listening on port 8080")
err := http.ListenAndServe(":8080", nil)
if err != nil {
log.Fatal("http.Lis:  ", err)
}
}
func listener(w http.ResponseWriter, r *http.Request) {
fmt.Println("request recieved")
w.Header().Set("Content-Type", "text/plain")
fmt.Fprintln(w, "goul test")
}
