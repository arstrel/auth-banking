# Banking project

This project consist of the following services:

- [REST api](https://github.com/arstrel/rest-banking)
- [Auth](https://github.com/arstrel/auth-banking)

# auth-banking

Auth service for banking app

## Role-Based Access Control (RBAC)

We have 2 roles in this app. This information is stored in users database table

- Admin role can use all endpoints
- User role can "Get customer by ID" and "Make a transaction" ony

---

# auth-banking

Auth service for banking app.
Auth process has 6 steps as shown below

![Auth flow](https://filedn.com/lTTdn1W2IjNme17D5yWuF74/Resources/go-banking-auth.png)

Learning objectives:

- Implement Authentication & Authorization in Golang
- Work with JWT Tokens and Role-Based Access Control (RBAC)
