# Dynamis

Dynamis is a lightweight tool for working with
[DynamoDB](https://aws.amazon.com/dynamodb/) via the
[aws-sdk](https://github.com/aws/aws-sdk-go). Its goal is to reduce some of the
repetitive and cumbersome aspects of the API, rather than to abstract it.

Status: Dynamis is well tested, and seems to have a stable API, but is minimal
in its scope. Please open an issue or pull request if you need funtionality
that seems to be missing.

# Example

Following is a simple example of writing and reading fields to DynamoDB, and
how dynamis saves you some headaches.

```golang
// Initialize the item to write to DynamoDB
item := map[string]*dynamodb.AttributeValue{}

// Initialize a dyanmis.ValueWriter for convenience.
w := dynamo.NewValueWriter(item)
w.Str("user_id", u.UserID.String())

// Name could be empty, which DynamoDB doesn't allow. Dynamis
// doesn't add it to the map if the value is empty.
w.Str("name", u.Name)

// See below to define a custom reader to turn this back into a date.
w.Str("start_date", u.StartDate.Format(time.RFC822)

// Use the item like normal.
resp, err := db.PutItem(&dynamodb.PutItemInput{
  TableName:           aws.String(tableName),
  Item:                item,
})
```

Now read the record out. Dynamis lets us access fields without an excess of
error handling, preferring to return a type's zero value over errors.

```golang
resp, err := db.GetItem(&dynamodb.GetItemInput{
  TableName: aws.String(tableName),
  Key: map[string]*dynamodb.AttributeValue{
    "user_id": {
      S: aws.String(id.String()),
    },
  },
})

// Initialize a ValueReader of the response item.
r := dynamo.NewValueReader(resp.Item)

// Define a custom reader for start_date to turn it back into a time.
r.Def("start_date", func(vr dynamis.ValueReader) interface{} {
  startDate, err := time.Parse(time.RFC822, vr.Str("start_date"))
  if err != nil {
    return Date{}
  }
  return startDate
})

// Now read from the item. ValueReader handles missing keys, missing values, 
// etc, and returns a zero value instead.
user := &User{
  UserID: UserID(r.Str("user_id")),

  // Name could be missing but that's ok.
  Name: r.Str("name"),

  // We know that the Def() of start_date returns a Time, so typecasting is safe.
  StartDate: r.Get("start_date").(time.Time),
}
```


## Author

Ryan Carver (ryan@ryancarver.com / @rcarver)

## License

MIT
