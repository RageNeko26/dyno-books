package helper

import (
	"fmt"
	"log"

	"dyno-books/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Declare table schema
type Data struct {
  UserEmail string `json:"user_email"`
  // Relation one-to-many. Because single user can have lot of books.
  Book []Book
}


type Book struct {
  Title string `json:"title"`
  Author string `json:"author"`
}

func AddNote(data *Data) {
  av, err := dynamodbattribute.MarshalMap(data)

  if err != nil {
    log.Fatal(err)
  }

  input := &dynamodb.PutItemInput{
    Item: av,
    TableName: aws.String("books"),
  }

  _, err = config.Dyno.PutItem(input)

  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Successfully add data")
  

}

func FindBook(email string) {
  result, err := config.Dyno.GetItem(&dynamodb.GetItemInput{
    // Define table name
    TableName: aws.String("books"),
    Key: map[string]*dynamodb.AttributeValue{
      "user_email": {
        S: aws.String(email),
      },
    },
  })

  if err != nil {
    log.Fatal(err)
  }

  if result.Item == nil {
    log.Fatal("Data is Not found")
  }

  var data Data
  // Decode raw data into struct object
  err = dynamodbattribute.UnmarshalMap(result.Item, &data) 

  if err != nil {
    log.Fatal(err)
  } 

  // Because property of "Book" data is in array 
  // We need to loop over it to get the data

  for index, hasil := range data.Book {
    fmt.Printf("Data #%v \n", index)
    fmt.Println("Title:", hasil.Title)
    fmt.Println("Author:", hasil.Author)
  }
}
