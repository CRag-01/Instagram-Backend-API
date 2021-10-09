# Instagram Backend HTTP REST API using GO Lang and Mongo DB

Project for Appointy Summer Internship . Project built within 25 hrs, with no prior knowledge of GO.

### Objectives - To be completed
- Create User Endpoint
- Fetch User Endpoint
- Create Post Endpoint
- Fetch Post Endpoint
- Fetch Post of User Endpoint
#### Additional Tasks 
- Password Encryption
- Server Thread Safety
- Pagination
- Unit Testing

### Progress
- [X] Create User Endpoint

Routed using - http:localhost:5000/users 

(Found in /EndPoints/user.go )

- [X] Fetch User Endpoint

Routed using - http:localhost:5000/user/{id} 

(Found in /EndPoints/user.go )

- [X] Create Post Endpoint

Routed using - http:localhost:5000/posts 

(Found in /EndPoints/post.go )

- [X] Fetch Post Endpoint

Routed using - http:localhost:5000/post/{id} 

(Found in /EndPoints/post.go )

- [X] Fetch Post of User Endpoint

Routed using - http:localhost:5000/posts/users/{id} 

(Found in /EndPoints/post.go )

#### Additional Tasks
- [x] Password Encryption
- [X] Server Thread Safety

Implemented using syncs

- [X] Pagination

Implemented using limit in URL

(Found in /Pagination/pagination.go)

- []Unit Testing

Could've implemented if not for the time constraint. :(
