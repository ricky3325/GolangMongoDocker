/*package main

import (
    "net/http"
    "encoding/json"
    "log"
)

func main()  {

    http.HandleFunc("/login1", login1)
    http.HandleFunc("/login2", login2)
    http.ListenAndServe("0.0.0.0:8080", nil)
}

type Resp struct {
    Code    string `json:"code"`
    Msg     string `json:"msg"`
}

type  Auth struct {
    Username string `json:"username"`
    Pwd      string   `json:"password"`
}

//post接口接收json數據
func login1(writer http.ResponseWriter,  request *http.Request)  {
    var auth Auth
    if err := json.NewDecoder(request.Body).Decode(&auth); err != nil {
        request.Body.Close()
        log.Fatal(err)
    }
    var result  Resp
    if auth.Username == "admin" && auth.Pwd == "123456" {
        result.Code = "200"
        result.Msg = "登錄成功"
    } else {
        result.Code = "401"
        result.Msg = "賬戶名或密碼錯誤"
    }
    if err := json.NewEncoder(writer).Encode(result); err != nil {
        log.Fatal(err)
    }
}

//接收x-www-form-urlencoded類型的post請求或者普通get請求
func login2(writer http.ResponseWriter,  request *http.Request)  {
    request.ParseForm()
    username, uError :=  request.Form["username"]
    pwd, pError :=  request.Form["password"]

    var result  Resp
    if !uError || !pError {
        result.Code = "401"
        result.Msg = "登錄失敗"
    } else if username[0] == "admin" && pwd[0] == "123456" {
        result.Code = "200"
        result.Msg = "登錄成功"
    } else {
        result.Code = "203"
        result.Msg = "賬戶名或密碼錯誤"
    }
    if err := json.NewEncoder(writer).Encode(result); err != nil {
        log.Fatal(err)
    }
}*/

package main

import (
   "context"
   "fmt"
   "log"
   "net/http"
   "encoding/json"
   "go.mongodb.org/mongo-driver/bson"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   //"time"
)

func main() {
    http.HandleFunc("/login2", login2)
    http.ListenAndServe("0.0.0.0:8080", nil)
	
	
}

type Resp struct {
    Code    string `json:"code"`
    Msg     string `json:"msg"`
}

type  Auth struct {
    Username string `json:"username"`
    Pwd      string   `json:"password"`
}

type Trainer struct {
	Name string
	Age  int
	City string
}

var (
	client     *mongo.Client
	err        error
	db         *mongo.Database
	collection *mongo.Collection
)

var ctx = context.TODO()

//接收x-www-form-urlencoded類型的post請求或者普通get請求
func login2(writer http.ResponseWriter,  request *http.Request)  {
    request.ParseForm()
    username, uError :=  request.Form["username"]
    pwd, pError :=  request.Form["password"]

    var result  Resp
    if !uError || !pError {
        result.Code = "401"
        result.Msg = "登錄失敗"
    } else if username[0] == "admin" && pwd[0] == "0" {
        result.Code = "200"
        result.Msg = "testDB"
		testDB()
	} else if username[0] == "admin" && pwd[0] == "1" {
		result.Code = "200"
		result.Msg = "testInsertOne"
		testInsertOne()
	} else if username[0] == "admin" && pwd[0] == "2" {
		result.Code = "200"
		result.Msg = "testFindOne"
		testFindOne()
    } else if username[0] == "admin" && pwd[0] == "3" {
		result.Code = "200"
		result.Msg = "testDeleteOne"
		testDeleteOne()
    } else if username[0] == "admin" && pwd[0] == "4" {
		result.Code = "200"
		result.Msg = "testUpdateOne"
		testUpdateOne()
    } else {
        result.Code = "203"
        result.Msg = "賬戶名或密碼錯誤"
    }
    if err := json.NewEncoder(writer).Encode(result); err != nil {
        log.Fatal(err)
    }
}

func connectDB(){
	clientOptions := options.Client().ApplyURI("mongodb://rs1:27041/")
    
    // Connect to MongoDB
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    // Check the connection
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")
}//connect to DB

func disconnectDB(){
    err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}//disconnect to DB

func testDB(){
	connectDB()

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println(databases)

	disconnectDB()
}//test DB connect status

func testInsertOne(){
	connectDB()

	collection = client.Database("test").Collection("trainers")

	ash := Trainer{"Ash", 10, "Pallet Town"}
	//misty := Trainer{"Misty", 10, "Cerulean City"}多插入測試資料
	//brock := Trainer{"Brock", 15, "Pewter City"}

	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}

	/*trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)*///多插入方法

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	disconnectDB()
}//testing InsertOne

func testFindOne(){
	connectDB()

	collection = client.Database("test").Collection("trainers")

	// Find a single document
	filter := bson.D{{"name", "Ash"}}

	var result Trainer

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)///// Find a single document

    disconnectDB()
}

func testDeleteOne(){
	connectDB()

	collection = client.Database("test").Collection("trainers")

	deleteResult, err := collection.DeleteOne(context.TODO(), bson.D{{"name", "Ash"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

    disconnectDB()
}

func testUpdateOne(){
	connectDB()

	collection = client.Database("test").Collection("trainers")

	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

    disconnectDB()
}