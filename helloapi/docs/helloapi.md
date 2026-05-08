### 1. N/A

1. route definition

- Url: /api/dept/count
- Method: GET
- Request: `DeptCountRequest`
- Response: `DeptCountResponse`

2. request definition



```golang
type DeptCountRequest struct {
	Dept string `json:"dept"`
}
```


3. response definition



```golang
type DeptCountResponse struct {
	Dept string `json:"dept"`
	Count int `json:"count"`
}
```

### 2. N/A

1. route definition

- Url: /hello/:name
- Method: GET
- Request: `Request`
- Response: `Response`

2. request definition



```golang
type Request struct {
	Name string `path:"name"`
}
```


3. response definition



```golang
type Response struct {
	Message string `json:"message"`
}
```

