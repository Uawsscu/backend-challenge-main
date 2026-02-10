<style>
  @page {
    size: A4;
    margin: 15mm;
  }
  body {
    font-family: 'Inter', 'Helvetica', 'Arial', sans-serif;
    line-height: 1.5;
    color: #1a1a1a;
  }
  h1, h2, h3 {
    color: #111;
    margin-top: 1em;
  }
  table { 
    width: 100% !important; 
    border-collapse: collapse; 
    margin-bottom: 24px;
    table-layout: fixed;
  }
  th, td { 
    border: 1px solid #e0e0e0; 
    padding: 10px 12px; 
    text-align: left; 
    font-size: 8.5pt; 
    word-wrap: break-word;
    vertical-align: top;
  }
  th { 
    background-color: #f7f7f7; 
    font-weight: 600;
    text-transform: uppercase;
    font-size: 8pt;
    letter-spacing: 0.02em;
  }
  code { 
    font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
    font-size: 8pt; 
    background-color: #f5f5f5;
    padding: 2px 4px;
    border-radius: 3px;
  }
  pre code {
    white-space: pre-wrap !important; 
    word-break: break-all !important;
    display: block;
    padding: 16px;
    border: 1px solid #e0e0e0;
    background-color: #fafafa;
    border-radius: 4px;
    line-height: 1.4;
  }
  table tr th:nth-child(1):nth-last-child(2), table tr td:nth-child(1):nth-last-child(2) { width: 25%; }
  table tr th:nth-child(2):nth-last-child(1), table tr td:nth-child(2):nth-last-child(1) { width: 75%; }
  table tr th:nth-child(1):nth-last-child(3), table tr td:nth-child(1):nth-last-child(3) { width: 25%; }
  table tr th:nth-child(2):nth-last-child(2), table tr td:nth-child(2):nth-last-child(2) { width: 35%; }
  table tr th:nth-child(3):nth-last-child(1), table tr td:nth-child(3):nth-last-child(1) { width: 40%; }
</style>

# User Management API Specification

## Authentication
Every request in this module requires an **Authorization** header:
`Authorization: Bearer <accessToken>`

---

## 1. Create User
| Field | Value |
| :--- | :--- |
| **Method** | `POST` |
| **URL** | `{{host}}/api/v1/users` |
| **Description** | Admin or system tool creates a new user |

### Header Attributes
| Header | Require | Type | Description | Example Value |
| :--- | :--- | :--- | :--- | :--- |
| `Authorization` | true | String | Bearer <accessToken> | `Bearer eyJhbGci...` |

### Request Body
| Attribute | Require | Type | Description |
| :--- | :--- | :--- | :--- |
| `name` | true | String | User's full name |
| `email` | true | String | Unique email address |
| `password` | true | String | Login password |

### Example Response (201 Created)
```json
{
  "id": "uuid-string",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "createdAt": "2024-02-11T00:00:00Z"
}
```

---

## 2. List Users
| Field | Value |
| :--- | :--- |
| **Method** | `GET` |
| **URL** | `{{host}}/api/v1/users` |
| **Description** | Retrieve all users in the system |

### Header Attributes
| Header | Require | Type | Description | Example Value |
| :--- | :--- | :--- | :--- | :--- |
| `Authorization` | true | String | Bearer <accessToken> | `Bearer eyJhbGci...` |

### Example Response (200 OK)
```json
[
  {
    "id": "uuid-user-1",
    "name": "User 1",
    "email": "user1@example.com",
    "createdAt": "2024-02-11T00:00:00Z"
  },
  {
    "id": "uuid-user-2",
    "name": "User 2",
    "email": "user2@example.com",
    "createdAt": "2024-02-11T00:00:00Z"
  }
]
```

---

## 3. Get User by ID
| Field | Value |
| :--- | :--- |
| **Method** | `GET` |
| **URL** | `{{host}}/api/v1/users/{id}` |
| **Description** | Retrieve specific user details |

### Header Attributes
| Header | Require | Type | Description | Example Value |
| :--- | :--- | :--- | :--- | :--- |
| `Authorization` | true | String | Bearer <accessToken> | `Bearer eyJhbGci...` |

### Example Response (200 OK)
```json
{
  "id": "uuid-string",
  "name": "User Name",
  "email": "user@example.com",
  "createdAt": "2024-02-11T00:00:00Z"
}
```

---

## 4. Update User
| Field | Value |
| :--- | :--- |
| **Method** | `PUT` |
| **URL** | `{{host}}/api/v1/users/{id}` |
| **Description** | Update user name or email |

### Header Attributes
| Header | Require | Type | Description | Example Value |
| :--- | :--- | :--- | :--- | :--- |
| `Authorization` | true | String | Bearer <accessToken> | `Bearer eyJhbGci...` |

### Request Body
| Attribute | Require | Type | Description |
| :--- | :--- | :--- | :--- |
| `name` | true | String | Updated name |
| `email` | true | String | Updated email |

### Example Response (200 OK)
```json
{
  "id": "uuid-string",
  "name": "Updated Name",
  "email": "updated@example.com",
  "createdAt": "2024-02-11T00:00:00Z"
}
```

---

## 5. Delete User
| Field | Value |
| :--- | :--- |
| **Method** | `DELETE` |
| **URL** | `{{host}}/api/v1/users/{id}` |
| **Description** | Remove a user from the system |

### Header Attributes
| Header | Require | Type | Description | Example Value |
| :--- | :--- | :--- | :--- | :--- |
| `Authorization` | true | String | Bearer <accessToken> | `Bearer eyJhbGci...` |

### Example Response (200 OK)
```json
{
  "message": "User deleted successfully"
}
```
