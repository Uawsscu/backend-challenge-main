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

# Auth API Specification

## 1. Register User
| Field | Value |
| :--- | :--- |
| **Method** | `POST` |
| **URL** | `{{host}}/api/v1/auth/register` |
| **Description** | Register a new user to the system |

### Request Body
| Attribute | Require | Type | Description |
| :--- | :--- | :--- | :--- |
| `name` | true | String | User's full name |
| `email` | true | String | Unique email address |
| `password` | true | String | Login password |

### Example Request
```json
{
  "name": "Lottery User",
  "email": "user@example.com",
  "password": "password123"
}
```

### Example Response (201 Created)
```json
{
  "id": "uuid-string",
  "name": "Lottery User",
  "email": "user@example.com",
  "createdAt": "2024-02-11T00:00:00Z"
}
```

---

## 2. Login
| Field | Value |
| :--- | :--- |
| **Method** | `POST` |
| **URL** | `{{host}}/api/v1/auth/login` |
| **Description** | Authenticate user and receive tokens |

### Request Body
| Attribute | Require | Type | Description |
| :--- | :--- | :--- | :--- |
| `email` | true | String | Registered email |
| `password` | true | String | Login password |

### Example Response (200 OK)
```json
{
  "accessToken": "eyJhbG...",
  "refreshToken": "eyJhbG..."
}
```

---

## 3. Logout
| Field | Value |
| :--- | :--- |
| **Method** | `POST` |
| **URL** | `{{host}}/api/v1/auth/logout` |
| **Description** | Invalidate the current access token |

### Header Attributes
| Header | Require | Type | Description |
| :--- | :--- | :--- | :--- |
| `Authorization` | true | String | Bearer <accessToken> |

### Example Response (200 OK)
```json
{
  "message": "Logged out successfully"
}
```

---

## 4. Refresh Token
| Field | Value |
| :--- | :--- |
| **Method** | `POST` |
| **URL** | `{{host}}/api/v1/auth/refresh` |
| **Description** | Get a new access token using a refresh token |

### Request Body
| Attribute | Require | Type | Description |
| :--- | :--- | :--- | :--- |
| `refreshToken` | true | String | Valid refresh token |

### Example Response (200 OK)
```json
{
  "accessToken": "eyJhbG...",
  "refreshToken": "eyJhbG..."
}
```
