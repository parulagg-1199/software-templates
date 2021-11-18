## List of API Endpoints ##

## 1. To check health of service ##

```
GET /healthz
```

## 2. To create an account in postgres ##

```
POST /v1/signup
Header:-
Content-Type: application/json

Body:-
email     string    required
password  string    required
```

## 3. Account Login##

```
POST /v1/login
Header:-
Content-Type: application/json

Body:-
email     string    required
password  string    required
```

## 4. Save user details in mysql##

```
POST /v1/userDetail
Header:-
Content-Type: application/json
Authorization: Bearer <token>


Body:-
name     string    required
contact  string    required
```

## 5. Get user details from mysql##

```
GET /v1/userDetail
Header:-
Authorization: Bearer <atoken>
```

## 6. Save user information in mongodb##

```
POST /v1/userInfo
Header:-
Content-Type: application/json
Authorization: Bearer <token>


Body:-
{
  "address":{
    "city":"string required",
    "state":"string required",
    "country":"string required",
    "zipCode":"string",
  },
  "interests":[{
    "type":"string required",
    "name":"string required"
    }]
}
```

## 7. Get user information from mongodb##

```
GET /v1/userInfo
Header:-
Authorization: Bearer <token>
```

## 8. Delete account all details and information by admin only from all databases##

```
DELETE /v1/deleteAccount
Header:-
Authorization: Bearer <admin_token>

Query Data:-
email string required
```

## 9. Save Key value in redis database##

```
POST /v1/userInfo
Header:-
Content-Type: application/json
Authorization: Bearer <token>


Body:-
key   string  required
value string  required
```

## 10. Check key in redis##

```
GET /v1/userInfo
Header:-
Authorization: Bearer <token>

Query Data:-
key string required
```

## 11. Delete key from redis##

```
DELETE /v1/deleteAccount
Header:-
Authorization: Bearer <token>

Query Data:-
key string required
```
