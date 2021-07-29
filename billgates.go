package main
import "fmt"
import "net/http"
import "net/url"
import "strings"
import "io/ioutil"
import "encoding/json"
import "time"
var validCounties []string =  []string{"ALEXANDRIA", "ARLINGTON", "BAILEYS CROSSROADS", "BURKE", "CHANTILLY", "DALE CITY", "FAIRFAX", "FALLS CHURCH", "HERNDON", "LEESBURG", "MANASSAS", "MANASSAS PARK", "WOODBRIDGE", "ANNANDALE", "VIENNA"}
var m map[string]string
func generateMap(){
  m = make(map[string]string)
  for _, county :=  range validCounties {
    m[county] = "";
  }

}


func updateStates(){
  body := getData()

  var unmarshal interface{}
  if body == nil{
    return
  }

  err := json.Unmarshal(body, &unmarshal)
  if err != nil{
    return
  }

  states := unmarshal.(map[string]interface{})["responsePayloadData"].(map[string]interface{})["data"].(map[string]interface{})["VA"].([]interface{})
  for _, state := range states {

    if  _,ok := m[state.(map[string]interface{})["city"].(string)]; ok{ //check if county being checked is in the map which contains our valid counties
      m[state.(map[string]interface{})["city"].(string)] = state.(map[string]interface{})["status"].(string)

    }

  }

}
func getData() []byte{

  url := "https://www.cvs.com/immunizations/covid-19-vaccine.vaccine-status.VA.json?vaccineinfo"
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return nil
  }
  req.Header.Add("referer", "https://www.cvs.com/immunizations/covid-19-vaccine")
  req.Header.Add("Cookie", "pe=p1")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  defer res.Body.Close()

  if res.StatusCode != 200{
    return nil
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return nil
  }
//  fmt.Println(string(body))
  return body
}
func checkEquality(oldMap map[string]string) bool{

  for _, v:= range validCounties {
    if(oldMap[v] != m[v]){

      return false
    }
  }
  return true
}
func copyMap() map[string]string{
  ret := make(map[string]string)
  for _, v:= range validCounties {
    ret[v] = m[v]
  }
  return ret
}
func postToDisc(content string){

  method := "POST"
  base_url := ""
  body := strings.NewReader("content="+content)
  req,_ := http.NewRequest(method, base_url, body)
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  resp,_ := http.DefaultClient.Do(req)
  bytes,_ := ioutil.ReadAll(resp.Body)
  fmt.Println(string(bytes))
  resp.Body.Close()
}

func main(){
  var oldMap map[string]string
  generateMap()

  updateStates()


  for true {
    time.Sleep(5*time.Minute) //10 minutes

    oldMap = copyMap()
    updateStates()
    if !checkEquality(oldMap){ //we have some changes

      for _, v:= range validCounties {

        if(oldMap[v] != m[v] && m[v] == "Available"){
          content := fmt.Sprintf("<@&829211060980154381> %s is avaliable! \n Schedule an appointment now: \n https://www.cvs.com/vaccine/intake/store/covid-screener/covid-qns", v)
          postToDisc(url.QueryEscape(content))
        }
      }
    }
  }


}
