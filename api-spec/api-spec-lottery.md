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

# Lottery Search API Specification

## Authentication
Every request in this module requires an **Authorization** header:
`Authorization: Bearer <accessToken>`

---

## 1. Search Lottery
| Field | Value |
| :--- | :--- |
| **Method** | `GET` |
| **URL** | `{{host}}/api/v1/lotteries/search` |
| **Description** | Search for lottery tickets by number pattern |

### Header Attributes
| Header | Require | Type | Description | Example Value |
| :--- | :--- | :--- | :--- | :--- |
| `Authorization` | true | String | Bearer <accessToken> | `Bearer eyJhbGci...` |

### Query Parameters
| Parameter | Require | Type | Description | Example Value |
| :--- | :--- | :--- | :--- | :--- |
| `pattern` | true | String | 6-character pattern (e.g. 123***, ****56) | `****23` |

### Example Request
`GET {{host}}/api/v1/lotteries/search?pattern=****23`

### Example Response (200 OK)
```json
{
    "count": 10,
    "pattern": "****23",
    "results": [
        {
            "ID": "698b6e4cd9666be7d11fffc2",
            "Number": "004223",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.062Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d120266f",
            "Number": "014123",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.339Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d1201f67",
            "Number": "012323",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.339Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d12024df",
            "Number": "013723",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.339Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d11ff8ba",
            "Number": "002423",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.062Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d1201ab7",
            "Number": "011123",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.339Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d120102a",
            "Number": "008423",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.063Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d11ff91e",
            "Number": "002523",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.062Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d11fff5e",
            "Number": "004123",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.062Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        },
        {
            "ID": "698b6e4cd9666be7d1200602",
            "Number": "005823",
            "Status": "reserved",
            "ReservedUntil": "2026-02-10T18:19:01.396Z",
            "ReservedBy": "580e17b3-52c5-47d6-8e1e-a7293d72c375",
            "CreatedAt": "2026-02-10T17:43:40.063Z",
            "UpdatedAt": "2026-02-10T18:14:01.396Z"
        }
    ]
}
```

### Error Responses
- **400 Bad Request**: Missing pattern or invalid pattern format.
- **401 Unauthorized**: Missing or invalid authentication token.
- **500 Internal Server Error**: Database or cache connection issues.
