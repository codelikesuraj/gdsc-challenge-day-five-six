# Days 5-6: Here comes the CRUD operations

Well done. At this stage of the challenge, you should have a local server connected to a database with a book model. Now we need to perform Create, Read, Update, and Delete Operations on these books.

Create API endpoints that will implement this behavior. These API endpoints are expected to follow proper naming conventions.

## Setup
- Navigate to the root of this repo.
- Run the command ```go run ./main.go``` to start the server.
- Visit the following url endpoints:
    |METHOD|DESCRIPTION|ENDPOINT|SAMPLE BODY|
    |------|-----------|--------|----|
    |GET   |Get all books|http://127.0.0.1:3000/books|-|
    |GET   |Get a book   |http://127.0.0.1:3000/books/{id}|-|
    |POST  |Create a book|http://127.0.0.1:3000/books|{"author":"book_author","title":"book_title","published_at":"YYYY-MM-DD"}|
    |PUT   |Update a book|http://127.0.0.1:3000/books/{id}|{"author":"book_author","title":"book_title","published_at":"YYYY-MM-DD"}|
    |DELETE|Delete a book|http://127.0.0.1:3000/books/{id}|-|
