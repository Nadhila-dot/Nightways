# Backend / Server for This project


Use go for the backend and use docker to deploy.
The docker images and etc will be in github. 

Examples to define a route.

```go
func json(){
    Route.Get("/json", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "This is a JSON response",
            "status":  "success",
        })
    })
}
``` 

Now this route is wrapped around a function thus you then register it in the main.go
```go
func main() {
     init() // <-- Your go function for route /json
}
```
A dynamic solution to Route loading. 


