package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 定義一個結構體來表示 MongoDB 中的文件
type User struct {
	Name         string    `bson:"name"`
	Age          int       `bson:"age"`
	Email        string    `bson:"email,omitempty"` // omitempty 表示如果 Email 為空則不包含在 BSON 中
	CreationTime time.Time `bson:"creation_time"`
}

func main() {
	// 設置 MongoDB 連線 URI
	uri := "mongodb://admin:password@192.168.0.236:27017/" // MongoDB 設定

	// 創建一個 MongoDB 客戶端選項
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(time.Hour)

	// 連線到 MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 在程式結束時斷開連線
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Println("成功連線到 MongoDB!")

	// 選擇要操作的資料庫和集合 (Collection)
	// 無指定資料庫時，會自動建立
	database := client.Database("my_database")
	usersCollection := database.Collection("users")

	// 創建 (Create)
	fmt.Println("\n--- 創建使用者 ---")
	newUser := User{Name: "Alice", Age: 30, Email: "alice@example.com", CreationTime: time.Now()}
	insertResult, err := usersCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("插入成功，InsertedID:", insertResult.InsertedID)

	newUsers := []interface{}{
		User{Name: "Bob", Age: 25, CreationTime: time.Now()},
		User{Name: "Charlie", Age: 35, Email: "charlie@example.org", CreationTime: time.Now()},
	}
	insertManyResult, err := usersCollection.InsertMany(context.TODO(), newUsers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("批量插入成功，InsertedIDs:", insertManyResult.InsertedIDs)

	// 讀取 (Read)
	fmt.Println("\n--- 讀取使用者 ---")
	// 讀取單個使用者
	var alice User
	filter := bson.D{{Key: "name", Value: "Alice"}}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&alice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("找不到名為 Alice 的使用者")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("找到使用者: %+v\n", alice)
	}

	// 讀取所有使用者
	fmt.Println("\n--- 所有使用者 ---")
	cursor, err := usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}

	// 更新 (Update)
	fmt.Println("\n--- 更新使用者 ---")
	updateFilter := bson.D{{Key: "name", Value: "Alice"}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 31}}}}
	updateResult, err := usersCollection.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("更新了 %v 個文件\n", updateResult.ModifiedCount)

	// 批量更新
	updateManyFilter := bson.D{{Key: "age", Value: bson.D{{Key: "$lt", Value: 30}}}}
	updateMany := bson.D{{Key: "$inc", Value: bson.D{{Key: "age", Value: 1}}}}
	updateManyResult, err := usersCollection.UpdateMany(context.TODO(), updateManyFilter, updateMany)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("批量更新了 %v 個文件\n", updateManyResult.ModifiedCount)

	// 刪除 (Delete)
	fmt.Println("\n--- 刪除使用者 ---")
	deleteFilter := bson.D{{Key: "name", Value: "Bob"}}
	deleteResult, err := usersCollection.DeleteOne(context.TODO(), deleteFilter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("刪除了 %v 個文件\n", deleteResult.DeletedCount)

	// 批量刪除
	deleteManyFilter := bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 35}}}}
	deleteManyResult, err := usersCollection.DeleteMany(context.TODO(), deleteManyFilter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("批量刪除了 %v 個文件\n", deleteManyResult.DeletedCount)
}
