# golang-simple-serverless

## overview ##
In this repo, I used [go](https://go.dev/) and [Serverless Framework](https://www.serverless.com/) to implement a simple serverless service.
* note: Learn more about serverless, you can check my medium: [認識無服務器-Serverless](https://medium.com/@barney30818/%E8%AA%8D%E8%AD%98%E7%84%A1%E6%9C%8D%E5%8B%99%E5%99%A8%E6%9E%B6-serverless-a63216517c29)

This project simulates like a registration system,customer can call API endpoint and type member info in payload,<br>
then it will trigger the lambda function and store the info in DynamoDB.<br>
At meantime,this lambda function calls SQS service and SQS will trigger another lambda function to send gmail.
## architecture ## 

![go-simple-serverless-arch](https://user-images.githubusercontent.com/69409373/179388416-c538c064-9a46-429c-a775-48a4bbbf41d7.jpg)

 ## API ## 
  
 **createMember**

```
POST /createMember
```

 **example request payload**
 
 ```
 {
   "Id":"barney30818",
   "password":"aaaaaaaa",
   "name":"Barney",
   "email":"barney308188@gmail.com"
}
```

