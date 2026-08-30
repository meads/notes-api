## 🗺️ Roadmap / TODO

### Backend
- [x] Register route to create new users, create tokens and session.
- [x] Users routes to retrieve and modify database users.
- [x] Login route to authenticate existing users, create tokens and session.
- [x] Sessions route to retrieve and modify sessions. 
- [x] Notes route to retrieve and modify notes db records.
- [x] Notes table for persisting notes records.
- [x] Users table for persisting uesrs records.
- [x] Sessions table for persisting sessions records.
- [x] Reliable means of hashing and comparing passwords during Registration and Login.
- [x] JWT middleware to protect certain routes by token validation. 
- [x] Database first design to generate data layer calls from sql statements.
- [x] Generate mocks for interfaces allowing mocks for unit testing code that modifies database state.
- [x] Generation of access & refresh tokens for custom claims using JWT library v5.
- [ ] Config struct to store environment variables in main package, validate and provide ease of use.
- [x] Near 100% code coverage for unit testing of all backend code.
- [x] Convert sqlc.yaml to version 2 and convert all TIMESTAMP columns to TIMESTAMPZ to handle NULL types
- [x] Create an internal api layer that becomes the source to interact with backing data stores.
- [x] Refactor handlers data calls to api calls.
- [x] Refactor handler tests to use mocks of the api.
- [x] Unit test api layer.
- [x] Create error types that stop wrapping at the handler boundary. So end users don't see system errors.
- [x] Access and Refresh token duration from environment variables

### Frontend
- [x] UI to register users.
- [x] UI to login existing users.
- [ ] UI to change user password.
- [ ] UI to view sessions and tokens.
- [x] UI to create, edit, list and delete notes.
- [x] UI to intercept 401 token expired calls
 

