package main

import (
 "bytes"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "os"
 "sort"
 "strconv"
 "strings"
 "sync"
 "sync/atomic"
 "time"

 "github.com/gorilla/websocket"
 "github.com/openimsdk/protocol/sdkws"
 "google.golang.org/protobuf/proto"
)

var httpClient = &http.Client{Timeout: 15*time.Second, Transport:&http.Transport{MaxIdleConns:128,MaxIdleConnsPerHost:128}}
var reportDir = os.Getenv("REPORTS")
var connected, active, failed, dropped atomic.Int64
var pongCount atomic.Int64
var mu sync.Mutex
var errorsSeen = map[string]int{}
var received = map[string]float64{}
var sentAt = map[string]time.Time{}
var duplicates int
var latencies []float64

func recordError(stage string) {mu.Lock(); errorsSeen[stage]++; mu.Unlock()}
func post(path string, body any, token string) (map[string]any,error) {
 b,_:=json.Marshal(body)
 req,e:=http.NewRequest("POST","http://127.0.0.1:10002"+path,bytes.NewReader(b)); if e!=nil{return nil,e}
 req.Header.Set("Content-Type","application/json"); req.Header.Set("operationID",strconv.FormatInt(time.Now().UnixNano(),10))
 if token!="" {req.Header.Set("token",token)}
 res,e:=httpClient.Do(req); if e!=nil{return nil,e}; defer res.Body.Close()
 var v struct{ErrCode *int; Data json.RawMessage}
 if e=json.NewDecoder(io.LimitReader(res.Body,4<<20)).Decode(&v); e!=nil{return nil,fmt.Errorf("%s response decode failed: %T",path,e)}
 if res.StatusCode!=200{return nil,fmt.Errorf("%s HTTP status=%d",path,res.StatusCode)}
 if v.ErrCode==nil{return nil,fmt.Errorf("%s missing errCode",path)}
 if *v.ErrCode!=0{return nil,fmt.Errorf("%s business code=%d",path,*v.ErrCode)}
 if path!="/auth/get_admin_token" && path!="/auth/get_user_token" {return nil,nil}
 var data map[string]any
 if err:=json.Unmarshal(v.Data,&data);err!=nil{return nil,fmt.Errorf("%s invalid token response data",path)}
 return data,nil
}
func token(v map[string]any) string {s,_:=v["token"].(string);return s}
func snapshot() map[string]any {
 mu.Lock();defer mu.Unlock()
 errs:=map[string]int{};for k,v:=range errorsSeen{errs[k]=v}
 return map[string]any{"connected_total":connected.Load(),"active":active.Load(),"failed":failed.Load(),"dropped":dropped.Load(),"pong_count":pongCount.Load(),"errors":errs}
}
func writeJSON(name string,v any){b,_:=json.MarshalIndent(v,"","  ");_ = os.WriteFile(reportDir+"/"+name,b,0644)}
func worker(id,admin string,done <-chan struct{}, ready *sync.WaitGroup, lifetime *sync.WaitGroup){
 defer lifetime.Done()
 completed:=false
 defer func(){if !completed{ready.Done()}}()
 d,e:=post("/auth/get_user_token",map[string]any{"userID":id,"platformID":5},admin)
 if e!=nil||token(d)==""{failed.Add(1);recordError("user_token");if failed.Load()<=3{fmt.Println("user token failed",e)};return}
 q:=url.Values{"sendID":{id},"token":{token(d)},"platformID":{"5"},"operationID":{id},"isBackground":{"false"},"isMsgResp":{"true"},"sdkType":{"js"}}
 dialer:=websocket.Dialer{HandshakeTimeout:15*time.Second}
 begin:=time.Now()
 c,_,e:=dialer.Dial("ws://127.0.0.1:10001/?"+q.Encode(),nil)
 if e!=nil{failed.Add(1);recordError("websocket_dial");return}
 defer c.Close()
 _=c.SetReadDeadline(time.Now().Add(15*time.Second))
 _,b,e:=c.ReadMessage()
 var ack struct{ErrCode *int}
 if e!=nil||json.Unmarshal(b,&ack)!=nil||ack.ErrCode==nil||*ack.ErrCode!=0{failed.Add(1);recordError("authentication_ack");return}
 mu.Lock();latencies=append(latencies,float64(time.Since(begin).Microseconds())/1000);mu.Unlock()
 connected.Add(1);active.Add(1);completed=true;ready.Done()
 defer active.Add(-1)
 _=c.SetReadDeadline(time.Now().Add(45*time.Second))
 c.SetPongHandler(func(string)error{pongCount.Add(1);return c.SetReadDeadline(time.Now().Add(45*time.Second))})
 stop:=make(chan struct{});defer close(stop)
 go func(){
  t:=time.NewTicker(10*time.Second);defer t.Stop()
  // Close the socket on cancellation to release the read loop.
  for {select{case <-done:c.Close();return;case <-stop:return;case <-t.C:
    if c.WriteControl(websocket.PingMessage,[]byte("ci"),time.Now().Add(5*time.Second))!=nil{c.Close();return}
  }}
 }()
 for {
  kind,b,e:=c.ReadMessage();if e!=nil{select{case <-done:default:dropped.Add(1);recordError("read_or_heartbeat")};return}
  if kind!=websocket.BinaryMessage{continue}
  var env struct{ReqIdentifier int;Data []byte}
  if json.Unmarshal(b,&env)!=nil||env.ReqIdentifier!=2001{continue}
  var push sdkws.PushMessages
  if proto.Unmarshal(env.Data,&push)!=nil{recordError("push_decode");continue}
  for _,batch:=range push.Msgs{for _,msg:=range batch.Msgs{
    var content struct{Content string}
    if json.Unmarshal(msg.Content,&content)!=nil{continue}
    mu.Lock()
    if start,ok:=sentAt[content.Content];ok && msg.RecvID==id{
      if _,seen:=received[content.Content];seen{duplicates++}else{received[content.Content]=float64(time.Since(start).Microseconds())/1000}
    }
    mu.Unlock()
  }}
 }
}
func percentile(a []float64,p float64)float64{if len(a)==0{return 0};sort.Float64s(a);return a[int(float64(len(a)-1)*p)]}
func run() int {
 n,e:=strconv.Atoi(os.Getenv("TARGET"));if e!=nil||n<2{return 2}
 b,e:=os.ReadFile("config/share.yml");if e!=nil{return 2}
 secret:="";for _,l:=range strings.Split(string(b),"\n"){if strings.HasPrefix(l,"secret:"){secret=strings.TrimSpace(strings.TrimPrefix(l,"secret:"))}}
 admin:=""
 for i:=0;i<90;i++{d,e:=post("/auth/get_admin_token",map[string]any{"secret":secret,"userID":"imAdmin"},"");if e==nil{admin=token(d);if admin!=""{break}};time.Sleep(2*time.Second)}
 if admin==""{recordError("server_not_ready");writeJSON("result.json",snapshot());return 1}
 ids:=make([]string,n)
 for i:=range ids{ids[i]=fmt.Sprintf("capacity%06d",i)}
 for i:=0;i<n;i+=200{end:=min(i+200,n);users:=make([]map[string]string,0,end-i);for _,id:=range ids[i:end]{users=append(users,map[string]string{"userID":id,"nickname":id})};if _,e:=post("/user/user_register",map[string]any{"users":users},admin);e!=nil{recordError("register_batch");writeJSON("result.json",snapshot());fmt.Println("registration failed",i,e);return 1}}
 fmt.Println("registered",n)
 done:=make(chan struct{});defer close(done)
 var life sync.WaitGroup
 monitorDone:=make(chan struct{});defer close(monitorDone)
 go func(){t:=time.NewTicker(5*time.Second);defer t.Stop();for{select{case <-monitorDone:return;case <-t.C:
  s:=snapshot();b,_:=json.Marshal(s);fmt.Println(string(b));f,e:=os.OpenFile(reportDir+"/connections.jsonl",os.O_CREATE|os.O_APPEND|os.O_WRONLY,0644);if e==nil{fmt.Fprintln(f,string(b));f.Close()}
 }}}()
 // Start in bounded batches; established sessions remain live while later batches connect.
 for i:=0;i<n;i+=100{end:=min(i+100,n);var ready sync.WaitGroup;ready.Add(end-i);for _,id:=range ids[i:end]{life.Add(1);go worker(id,admin,done,&ready,&life)};ready.Wait()}
 startActive:=active.Load()
 fmt.Println("measurement starts",startActive)
 minActive:=startActive
 sentOK:=0;sentFailed:=0
 // REST publishes, actual authenticated WebSocket recipient push. Not native client-send ACK.
 end:=time.Now().Add(60*time.Second)
 ticker:=time.NewTicker(25*time.Millisecond);defer ticker.Stop()
 for i:=0;time.Now().Before(end);i++{
  <-ticker.C
  if a:=active.Load();a<minActive{minActive=a}
  marker:=fmt.Sprintf("capacity-msg-%d",i)
  mu.Lock();sentAt[marker]=time.Now();mu.Unlock()
  _,e:=post("/msg/send_msg",map[string]any{"sendID":ids[i%n],"recvID":ids[(i+1)%n],"senderNickname":ids[i%n],"sessionType":1,"contentType":101,"content":map[string]string{"content":marker}},admin)
  if e!=nil{sentFailed++;mu.Lock();delete(sentAt,marker);mu.Unlock()}else{sentOK++}
 }
 time.Sleep(10*time.Second)
 if a:=active.Load();a<minActive{minActive=a}
 s:=snapshot()
 mu.Lock()
 delivery:=make([]float64,0,len(received));for _,v:=range received{delivery=append(delivery,v)}
 s["target"]=n;s["min_sampled_active"]=minActive;s["sent_rest_success"]=sentOK;s["sent_rest_failure"]=sentFailed;s["unique_recipient_push"]=len(received);s["duplicate_push"]=duplicates
 s["connection_p95_ms"]=percentile(latencies,.95);s["rest_to_push_p99_ms"]=percentile(delivery,.99)
 pass:=minActive==int64(n)&&failed.Load()==0&&dropped.Load()==0&&sentOK>0&&sentFailed==0&&len(received)==sentOK&&duplicates==0&&pongCount.Load()>=int64(n)
 s["passed"]=pass;mu.Unlock()
 writeJSON("result.json",s);out,_:=json.Marshal(s);fmt.Println(string(out))
 if !pass{return 1};return 0
}
func main(){os.Exit(run())}
