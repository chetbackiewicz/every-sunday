# REST API Endpoint Design Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: RESTful API with Gin Framework

---

## Recommendation Summary

**Architecture: RESTful API with JSON**  
**Framework: Gin (from research doc 08)**  
**Base URL: `/api/v1`**  
**Auth: JWT Bearer tokens (future phase)**

### Core Principles

1. ✅ **Resource-based URLs** - Nouns, not verbs (`/budgets`, not `/getBudget`)
2. ✅ **HTTP methods** - GET (read), POST (create), PUT (update), DELETE (delete)
3. ✅ **Status codes** - Semantic HTTP codes (200, 201, 204, 400, 404, 409, 500)
4. ✅ **JSON everywhere** - Requests and responses in JSON format
5. ✅ **Consistent errors** - Standard error response structure
6. ✅ **Versioning** - `/api/v1` prefix for future compatibility
7. ✅ **Stateless** - Each request contains all necessary information

---

## Core Requirements

From spike document and copilot-instructions.md:
- CRUD operations for budgets, categories, transactions
- CSV file upload (up to 10 files per month)
- Cell reference tracking (linking transactions to category actuals)
- Multi-month budget support
- Hierarchical category tree (3 levels deep)
- User isolation (multi-user support)

---

## URL Structure

### Base URL

```
https://api.budgetapp.com/api/v1
```

**Rationale**:
- `/api` - Distinguishes API from frontend routes
- `/v1` - API versioning for future backward compatibility

---

## Authentication (Future Phase)

### Endpoints

```
POST   /api/v1/auth/register    # Create new user account
POST   /api/v1/auth/login       # Login, returns JWT token
POST   /api/v1/auth/logout      # Logout (invalidate token)
POST   /api/v1/auth/refresh     # Refresh expired JWT token
GET    /api/v1/auth/me          # Get current user info
```

### Request Examples

**Register**:
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response** (201 Created):
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "createdAt": "2025-10-19T14:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Login**:
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response** (200 OK):
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expiresAt": "2025-10-20T14:30:00Z"
}
```

### Authentication Header

All authenticated endpoints require:
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## Budget Endpoints

### Get Budget for Month

```
GET /api/v1/budgets/:month
```

**Path Parameters**:
- `month` - Format: `YYYY-MM` (e.g., `2025-10`)

**Response** (200 OK):
```json
{
  "id": 123,
  "userId": 1,
  "month": "2025-10",
  "name": "October 2025 Budget",
  "createdAt": "2025-10-01T10:00:00Z",
  "updatedAt": "2025-10-15T14:30:00Z"
}
```

**Error** (404 Not Found):
```json
{
  "error": {
    "code": "BUDGET_NOT_FOUND",
    "message": "Budget for month 2025-10 not found",
    "field": "month"
  }
}
```

---

### Create Budget for Month

```
POST /api/v1/budgets
```

**Request**:
```json
{
  "month": "2025-10",
  "name": "October 2025 Budget"
}
```

**Response** (201 Created):
```json
{
  "id": 123,
  "userId": 1,
  "month": "2025-10",
  "name": "October 2025 Budget",
  "createdAt": "2025-10-19T14:30:00Z",
  "updatedAt": "2025-10-19T14:30:00Z"
}
```

**Error** (409 Conflict - budget already exists):
```json
{
  "error": {
    "code": "BUDGET_ALREADY_EXISTS",
    "message": "Budget for month 2025-10 already exists",
    "field": "month"
  }
}
```

---

### Update Budget

```
PUT /api/v1/budgets/:month
```

**Request**:
```json
{
  "name": "October Personal Budget"
}
```

**Response** (200 OK):
```json
{
  "id": 123,
  "userId": 1,
  "month": "2025-10",
  "name": "October Personal Budget",
  "createdAt": "2025-10-01T10:00:00Z",
  "updatedAt": "2025-10-19T14:35:00Z"
}
```

---

### Delete Budget

```
DELETE /api/v1/budgets/:month
```

**Response** (204 No Content)

**Error** (404 Not Found):
```json
{
  "error": {
    "code": "BUDGET_NOT_FOUND",
    "message": "Budget for month 2025-10 not found"
  }
}
```

---

### List User's Budgets

```
GET /api/v1/budgets
```

**Query Parameters**:
- `limit` (optional) - Max number of results (default: 12)
- `offset` (optional) - Pagination offset (default: 0)
- `sort` (optional) - Sort order: `asc` or `desc` (default: `desc`)

**Response** (200 OK):
```json
{
  "budgets": [
    {
      "id": 125,
      "month": "2025-11",
      "name": "November 2025 Budget",
      "createdAt": "2025-11-01T10:00:00Z",
      "updatedAt": "2025-11-01T10:00:00Z"
    },
    {
      "id": 123,
      "month": "2025-10",
      "name": "October 2025 Budget",
      "createdAt": "2025-10-01T10:00:00Z",
      "updatedAt": "2025-10-19T14:35:00Z"
    }
  ],
  "total": 2,
  "limit": 12,
  "offset": 0
}
```

---

## Category Endpoints

### List Categories for Budget

```
GET /api/v1/budgets/:month/categories
```

**Response** (200 OK):
```json
{
  "categories": [
    {
      "id": 1,
      "budgetId": 123,
      "name": "Income",
      "projected": "5000.00",
      "manualActual": null,
      "calculatedActual": "4840.48",
      "children": [
        {
          "id": 2,
          "name": "Salary",
          "projected": "4500.00",
          "manualActual": null,
          "calculatedActual": "4840.48",
          "children": []
        }
      ]
    },
    {
      "id": 10,
      "budgetId": 123,
      "name": "Expenses",
      "projected": "3500.00",
      "manualActual": null,
      "calculatedActual": "2825.68",
      "children": [
        {
          "id": 11,
          "name": "Housing",
          "projected": "1500.00",
          "manualActual": null,
          "calculatedActual": "1500.00",
          "children": [
            {
              "id": 12,
              "name": "Mortgage",
              "projected": "1200.00",
              "manualActual": null,
              "calculatedActual": "1200.00",
              "children": []
            },
            {
              "id": 13,
              "name": "Insurance",
              "projected": "300.00",
              "manualActual": null,
              "calculatedActual": "300.00",
              "children": []
            }
          ]
        }
      ]
    }
  ],
  "total": 2
}
```

**Notes**:
- Returns hierarchical tree structure (up to 3 levels)
- `calculatedActual` is sum of linked transactions
- `manualActual` overrides `calculatedActual` if set

---

### Create Category

```
POST /api/v1/budgets/:month/categories
```

**Request**:
```json
{
  "name": "Salary",
  "projected": "4500.00",
  "parentId": 1
}
```

**Validation Rules**:
- `name` - Required, max 255 characters
- `projected` - Required, must be >= 0, max 12 digits with 2 decimals
- `parentId` - Optional, if provided must reference existing category
- Cannot nest more than 3 levels deep

**Response** (201 Created):
```json
{
  "id": 2,
  "budgetId": 123,
  "name": "Salary",
  "projected": "4500.00",
  "manualActual": null,
  "calculatedActual": "0.00",
  "parentId": 1,
  "depth": 1,
  "createdAt": "2025-10-19T14:40:00Z",
  "updatedAt": "2025-10-19T14:40:00Z"
}
```

**Error** (400 Bad Request - max depth exceeded):
```json
{
  "error": {
    "code": "MAX_DEPTH_EXCEEDED",
    "message": "Cannot nest categories more than 3 levels deep",
    "field": "parentId"
  }
}
```

---

### Update Category

```
PUT /api/v1/categories/:id
```

**Request**:
```json
{
  "name": "Base Salary",
  "projected": "4800.00",
  "manualActual": "4840.48"
}
```

**Response** (200 OK):
```json
{
  "id": 2,
  "budgetId": 123,
  "name": "Base Salary",
  "projected": "4800.00",
  "manualActual": "4840.48",
  "calculatedActual": "4840.48",
  "parentId": 1,
  "depth": 1,
  "updatedAt": "2025-10-19T14:45:00Z"
}
```

---

### Delete Category

```
DELETE /api/v1/categories/:id
```

**Response** (204 No Content)

**Notes**:
- Deletes category and all descendants (CASCADE)
- Unlinks all cell references to this category

**Error** (404 Not Found):
```json
{
  "error": {
    "code": "CATEGORY_NOT_FOUND",
    "message": "Category with ID 999 not found"
  }
}
```

---

## CSV File Endpoints

### Upload CSV File

```
POST /api/v1/budgets/:month/files
Content-Type: multipart/form-data
```

**Request**:
```
multipart/form-data boundary
--boundary
Content-Disposition: form-data; name="file"; filename="Chase7547_Activity_20251019.CSV"
Content-Type: text/csv

Transaction Date,Post Date,Description,Category,Type,Amount,Memo
10/15/2025,10/16/2025,GITHUB INC PAYROLL,,ACH_CREDIT,2840.48,
10/14/2025,10/14/2025,Payment to Chase,,LOAN_PMT,-325.68,
--boundary--
```

**Validation Rules**:
- File extension must be `.csv` or `.CSV`
- File size <= 10MB
- Max 10 files per budget
- CSV must have required columns: `Transaction Date`, `Post Date`, `Description`, `Amount`

**Response** (201 Created):
```json
{
  "file": {
    "id": 45,
    "budgetId": 123,
    "filename": "Chase7547_Activity_20251019.CSV",
    "fileSize": 2048,
    "rowCount": 42,
    "uploadDate": "2025-10-19T14:50:00Z"
  },
  "transactions": [
    {
      "id": 501,
      "fileId": 45,
      "rowNumber": 1,
      "transactionDate": "2025-10-15",
      "postDate": "2025-10-16",
      "description": "GITHUB INC PAYROLL",
      "category": "",
      "type": "ACH_CREDIT",
      "amount": "2840.48",
      "note": ""
    },
    {
      "id": 502,
      "fileId": 45,
      "rowNumber": 2,
      "transactionDate": "2025-10-14",
      "postDate": "2025-10-14",
      "description": "Payment to Chase",
      "category": "",
      "type": "LOAN_PMT",
      "amount": "-325.68",
      "note": ""
    }
  ]
}
```

**Error** (400 Bad Request - invalid file type):
```json
{
  "error": {
    "code": "INVALID_FILE_TYPE",
    "message": "Only CSV files are allowed",
    "field": "file"
  }
}
```

**Error** (409 Conflict - max files reached):
```json
{
  "error": {
    "code": "MAX_FILES_EXCEEDED",
    "message": "Cannot upload more than 10 CSV files per budget",
    "field": "file"
  }
}
```

---

### List CSV Files for Budget

```
GET /api/v1/budgets/:month/files
```

**Response** (200 OK):
```json
{
  "files": [
    {
      "id": 45,
      "budgetId": 123,
      "filename": "Chase7547_Activity_20251019.CSV",
      "fileSize": 2048,
      "rowCount": 42,
      "uploadDate": "2025-10-19T14:50:00Z"
    },
    {
      "id": 44,
      "budgetId": 123,
      "filename": "Chase7547_Activity_20251001.CSV",
      "fileSize": 1536,
      "rowCount": 28,
      "uploadDate": "2025-10-05T10:15:00Z"
    }
  ],
  "total": 2
}
```

---

### Delete CSV File

```
DELETE /api/v1/files/:id
```

**Response** (204 No Content)

**Notes**:
- Deletes file record and all associated transactions (CASCADE)
- Removes all cell references to transactions from this file

---

## Transaction Endpoints

### List Transactions for Budget

```
GET /api/v1/budgets/:month/transactions
```

**Query Parameters**:
- `fileId` (optional) - Filter by CSV file ID
- `limit` (optional) - Max results (default: 100)
- `offset` (optional) - Pagination offset (default: 0)

**Response** (200 OK):
```json
{
  "transactions": [
    {
      "id": 501,
      "fileId": 45,
      "fileName": "Chase7547_Activity_20251019.CSV",
      "rowNumber": 1,
      "transactionDate": "2025-10-15",
      "postDate": "2025-10-16",
      "description": "GITHUB INC PAYROLL",
      "category": "",
      "type": "ACH_CREDIT",
      "amount": "2840.48",
      "note": "",
      "linkedCategories": [],
      "createdAt": "2025-10-19T14:50:00Z",
      "updatedAt": "2025-10-19T14:50:00Z"
    },
    {
      "id": 502,
      "fileId": 45,
      "fileName": "Chase7547_Activity_20251019.CSV",
      "rowNumber": 2,
      "transactionDate": "2025-10-14",
      "postDate": "2025-10-14",
      "description": "Payment to Chase",
      "category": "",
      "type": "LOAN_PMT",
      "amount": "-325.68",
      "note": "",
      "linkedCategories": [
        {
          "categoryId": 12,
          "categoryName": "Mortgage",
          "amountSnapshot": "325.68"
        }
      ],
      "createdAt": "2025-10-19T14:50:00Z",
      "updatedAt": "2025-10-19T14:55:00Z"
    }
  ],
  "total": 2,
  "limit": 100,
  "offset": 0
}
```

---

### Update Transaction

```
PUT /api/v1/transactions/:id
```

**Request**:
```json
{
  "description": "GitHub October Payroll",
  "category": "Salary",
  "note": "Bi-weekly paycheck"
}
```

**Notes**:
- All editable fields are optional
- Cannot update `amount`, `transactionDate`, `postDate` (CSV source data)
- Can edit `description`, `category`, `type`, `note`

**Response** (200 OK):
```json
{
  "id": 501,
  "fileId": 45,
  "rowNumber": 1,
  "transactionDate": "2025-10-15",
  "postDate": "2025-10-16",
  "description": "GitHub October Payroll",
  "category": "Salary",
  "type": "ACH_CREDIT",
  "amount": "2840.48",
  "note": "Bi-weekly paycheck",
  "updatedAt": "2025-10-19T15:00:00Z"
}
```

---

### Delete Transaction

```
DELETE /api/v1/transactions/:id
```

**Response** (204 No Content)

**Notes**:
- Removes transaction from database
- Removes all cell references linking this transaction to categories
- Updates category `calculatedActual` values

---

## Cell Reference Endpoints

### Link Transaction to Category

```
POST /api/v1/cell-references
```

**Request**:
```json
{
  "categoryId": 12,
  "transactionId": 502
}
```

**Validation Rules**:
- `categoryId` must reference existing category
- `transactionId` must reference existing transaction
- Cannot link same transaction to same category twice (unique constraint)

**Response** (201 Created):
```json
{
  "id": 78,
  "categoryId": 12,
  "transactionId": 502,
  "amountSnapshot": "325.68",
  "createdAt": "2025-10-19T14:55:00Z"
}
```

**Error** (409 Conflict - already linked):
```json
{
  "error": {
    "code": "REFERENCE_ALREADY_EXISTS",
    "message": "Transaction 502 is already linked to category 12",
    "field": "transactionId"
  }
}
```

---

### Unlink Transaction from Category

```
DELETE /api/v1/cell-references/:id
```

**Alternative (by IDs)**:
```
DELETE /api/v1/cell-references?categoryId=12&transactionId=502
```

**Response** (204 No Content)

---

### List Cell References for Category

```
GET /api/v1/categories/:id/cell-references
```

**Response** (200 OK):
```json
{
  "references": [
    {
      "id": 78,
      "categoryId": 12,
      "transactionId": 502,
      "amountSnapshot": "325.68",
      "transaction": {
        "description": "Payment to Chase",
        "transactionDate": "2025-10-14",
        "amount": "325.68"
      },
      "createdAt": "2025-10-19T14:55:00Z"
    }
  ],
  "total": 1,
  "sumAmount": "325.68"
}
```

---

## Budget Summary Endpoint

### Get Budget Summary with Calculations

```
GET /api/v1/budgets/:month/summary
```

**Response** (200 OK):
```json
{
  "budget": {
    "id": 123,
    "month": "2025-10",
    "name": "October 2025 Budget"
  },
  "summary": {
    "income": {
      "projected": "5000.00",
      "actual": "4840.48",
      "difference": "-159.52"
    },
    "savingsPreTax": {
      "projected": "500.00",
      "actual": "0.00",
      "difference": "-500.00"
    },
    "savingsPostTax": {
      "projected": "300.00",
      "actual": "0.00",
      "difference": "-300.00"
    },
    "expenses": {
      "projected": "3500.00",
      "actual": "2825.68",
      "difference": "-674.32"
    },
    "remaining": {
      "projected": "700.00",
      "actual": "2014.80",
      "difference": "1314.80"
    }
  },
  "categories": [
    {
      "id": 1,
      "name": "Income",
      "projected": "5000.00",
      "actual": "4840.48",
      "difference": "-159.52",
      "children": [...]
    }
  ]
}
```

**Calculation Rules**:
- `actual` = `manualActual` ?? `calculatedActual` (manual override takes precedence)
- `calculatedActual` = SUM of linked transaction amounts + child category actuals
- `difference` = `actual` - `projected`
- `remaining` = `income` - `savingsPreTax` - `savingsPostTax` - `expenses`

---

## Error Response Format

### Standard Error Schema

All errors follow this structure:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "field": "fieldName",
    "details": {}
  }
}
```

**Fields**:
- `code` - Machine-readable error identifier (UPPER_SNAKE_CASE)
- `message` - User-friendly description
- `field` (optional) - Which request field caused the error
- `details` (optional) - Additional error context

---

### HTTP Status Codes

#### Success Codes

| Code | Name       | Usage                                |
| ---- | ---------- | ------------------------------------ |
| 200  | OK         | Successful GET, PUT requests         |
| 201  | Created    | Successful POST (resource created)   |
| 204  | No Content | Successful DELETE (no response body) |

#### Client Error Codes

| Code | Name                 | Usage                                        |
| ---- | -------------------- | -------------------------------------------- |
| 400  | Bad Request          | Invalid request format, validation failed    |
| 401  | Unauthorized         | Missing or invalid authentication token      |
| 403  | Forbidden            | Authenticated but not authorized             |
| 404  | Not Found            | Resource doesn't exist                       |
| 409  | Conflict             | Resource already exists, constraint violated |
| 422  | Unprocessable Entity | Validation error (semantic)                  |

#### Server Error Codes

| Code | Name                  | Usage                           |
| ---- | --------------------- | ------------------------------- |
| 500  | Internal Server Error | Unexpected server error         |
| 503  | Service Unavailable   | Database down, maintenance mode |

---

### Common Error Codes

**Authentication**:
- `UNAUTHORIZED` - Missing authentication token
- `INVALID_TOKEN` - JWT token expired or invalid
- `FORBIDDEN` - User doesn't have permission

**Validation**:
- `VALIDATION_ERROR` - Request validation failed
- `INVALID_FORMAT` - Date/number format incorrect
- `FIELD_REQUIRED` - Required field missing
- `FIELD_TOO_LONG` - String exceeds max length

**Resource Errors**:
- `RESOURCE_NOT_FOUND` - Budget, category, transaction not found
- `RESOURCE_ALREADY_EXISTS` - Unique constraint violated
- `MAX_DEPTH_EXCEEDED` - Category nesting too deep
- `MAX_FILES_EXCEEDED` - Too many CSV files uploaded

**Database Errors**:
- `DATABASE_ERROR` - Generic database error
- `CONSTRAINT_VIOLATION` - Foreign key or unique constraint

---

### Error Examples

**400 Bad Request** - Validation Error:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "name": ["Field is required", "Must be at least 3 characters"],
      "projected": ["Must be a number greater than or equal to 0"]
    }
  }
}
```

**401 Unauthorized**:
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required. Please provide a valid Bearer token."
  }
}
```

**404 Not Found**:
```json
{
  "error": {
    "code": "BUDGET_NOT_FOUND",
    "message": "Budget for month 2025-10 not found"
  }
}
```

**409 Conflict**:
```json
{
  "error": {
    "code": "BUDGET_ALREADY_EXISTS",
    "message": "Budget for month 2025-10 already exists. Use PUT to update.",
    "field": "month"
  }
}
```

**500 Internal Server Error**:
```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An unexpected error occurred. Please try again later.",
    "details": {
      "requestId": "abc123-def456"
    }
  }
}
```

---

## Gin Implementation Examples

### Handler Structure

```go
// controllers/budget.go
package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type BudgetController struct {
    service *services.BudgetService
}

func NewBudgetController(service *services.BudgetService) *BudgetController {
    return &BudgetController{service: service}
}

// GET /api/v1/budgets/:month
func (bc *BudgetController) GetBudget(c *gin.Context) {
    month := c.Param("month")
    userID := getUserID(c) // From auth middleware
    
    budget, err := bc.service.GetBudgetByMonth(userID, month)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            c.JSON(http.StatusNotFound, gin.H{
                "error": gin.H{
                    "code": "BUDGET_NOT_FOUND",
                    "message": fmt.Sprintf("Budget for month %s not found", month),
                },
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code": "INTERNAL_ERROR",
                "message": "Failed to retrieve budget",
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, budget)
}

// POST /api/v1/budgets
func (bc *BudgetController) CreateBudget(c *gin.Context) {
    var req CreateBudgetRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code": "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }
    
    userID := getUserID(c)
    budget, err := bc.service.CreateBudget(userID, req.Month, req.Name)
    
    if err != nil {
        if errors.Is(err, ErrBudgetExists) {
            c.JSON(http.StatusConflict, gin.H{
                "error": gin.H{
                    "code": "BUDGET_ALREADY_EXISTS",
                    "message": fmt.Sprintf("Budget for month %s already exists", req.Month),
                    "field": "month",
                },
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code": "INTERNAL_ERROR",
                "message": "Failed to create budget",
            },
        })
        return
    }
    
    c.JSON(http.StatusCreated, budget)
}
```

---

### Request Validation

```go
// models/requests.go
package models

type CreateBudgetRequest struct {
    Month string `json:"month" binding:"required,len=7"`
    Name  string `json:"name" binding:"max=255"`
}

type CreateCategoryRequest struct {
    Name      string   `json:"name" binding:"required,max=255"`
    Projected float64  `json:"projected" binding:"required,gte=0"`
    ParentID  *int     `json:"parentId" binding:"omitempty"`
}

type UpdateTransactionRequest struct {
    Description string `json:"description" binding:"omitempty,max=500"`
    Category    string `json:"category" binding:"omitempty,max=255"`
    Type        string `json:"type" binding:"omitempty,max=50"`
    Note        string `json:"note" binding:"omitempty,max=1000"`
}
```

**Validation Tags**:
- `binding:"required"` - Field must be present
- `binding:"gte=0"` - Greater than or equal to 0
- `binding:"max=255"` - Max string length
- `binding:"len=7"` - Exact length (for YYYY-MM format)
- `binding:"omitempty"` - Optional field

---

### File Upload Handler

```go
// controllers/upload.go
package controllers

import (
    "net/http"
    "path/filepath"
    "github.com/gin-gonic/gin"
)

const MaxFileSize = 10 * 1024 * 1024 // 10MB
const MaxFilesPerBudget = 10

type UploadController struct {
    service *services.UploadService
}

func (uc *UploadController) UploadCSV(c *gin.Context) {
    month := c.Param("month")
    userID := getUserID(c)
    
    // Get file from form
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code": "NO_FILE_UPLOADED",
                "message": "No file uploaded",
                "field": "file",
            },
        })
        return
    }
    
    // Validate file extension
    ext := filepath.Ext(file.Filename)
    if ext != ".csv" && ext != ".CSV" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code": "INVALID_FILE_TYPE",
                "message": "Only CSV files are allowed",
                "field": "file",
            },
        })
        return
    }
    
    // Validate file size
    if file.Size > MaxFileSize {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code": "FILE_TOO_LARGE",
                "message": fmt.Sprintf("File size must be less than %dMB", MaxFileSize/(1024*1024)),
                "field": "file",
            },
        })
        return
    }
    
    // Check max files
    fileCount, err := uc.service.GetFileCount(userID, month)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code": "INTERNAL_ERROR",
                "message": "Failed to check file count",
            },
        })
        return
    }
    
    if fileCount >= MaxFilesPerBudget {
        c.JSON(http.StatusConflict, gin.H{
            "error": gin.H{
                "code": "MAX_FILES_EXCEEDED",
                "message": fmt.Sprintf("Cannot upload more than %d CSV files per budget", MaxFilesPerBudget),
                "field": "file",
            },
        })
        return
    }
    
    // Process CSV file
    result, err := uc.service.ProcessCSV(userID, month, file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code": "PROCESSING_ERROR",
                "message": "Failed to process CSV file",
            },
        })
        return
    }
    
    c.JSON(http.StatusCreated, result)
}
```

---

### Error Handling Middleware

```go
// middleware/errors.go
package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // Process request
        
        // Check if there were any errors
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            // Log error for debugging
            log.Printf("Error: %v", err.Err)
            
            // Return generic error to client (don't expose internals)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": gin.H{
                    "code": "INTERNAL_ERROR",
                    "message": "An unexpected error occurred",
                },
            })
        }
    }
}

// Usage in main.go
func main() {
    r := gin.Default()
    r.Use(middleware.ErrorHandler())
    // ... routes
}
```

---

### CORS Middleware

```go
// middleware/cors.go
package middleware

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // React dev server
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

---

### Route Setup

```go
// routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "myapp/controllers"
    "myapp/middleware"
)

func SetupRoutes(r *gin.Engine, controllers *controllers.Controllers) {
    // Middleware
    r.Use(middleware.CORS())
    r.Use(middleware.ErrorHandler())
    
    // API v1 routes
    v1 := r.Group("/api/v1")
    {
        // Public routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", controllers.Auth.Register)
            auth.POST("/login", controllers.Auth.Login)
        }
        
        // Protected routes (require authentication)
        protected := v1.Group("")
        protected.Use(middleware.AuthRequired())
        {
            // Budgets
            budgets := protected.Group("/budgets")
            {
                budgets.GET("", controllers.Budget.List)
                budgets.GET("/:month", controllers.Budget.Get)
                budgets.POST("", controllers.Budget.Create)
                budgets.PUT("/:month", controllers.Budget.Update)
                budgets.DELETE("/:month", controllers.Budget.Delete)
                budgets.GET("/:month/summary", controllers.Budget.GetSummary)
                
                // Categories
                budgets.GET("/:month/categories", controllers.Category.List)
                budgets.POST("/:month/categories", controllers.Category.Create)
                
                // Files
                budgets.GET("/:month/files", controllers.Upload.ListFiles)
                budgets.POST("/:month/files", controllers.Upload.UploadCSV)
                
                // Transactions
                budgets.GET("/:month/transactions", controllers.Transaction.List)
            }
            
            // Categories
            categories := protected.Group("/categories")
            {
                categories.PUT("/:id", controllers.Category.Update)
                categories.DELETE("/:id", controllers.Category.Delete)
                categories.GET("/:id/cell-references", controllers.CellReference.ListForCategory)
            }
            
            // Transactions
            transactions := protected.Group("/transactions")
            {
                transactions.PUT("/:id", controllers.Transaction.Update)
                transactions.DELETE("/:id", controllers.Transaction.Delete)
            }
            
            // Files
            files := protected.Group("/files")
            {
                files.DELETE("/:id", controllers.Upload.DeleteFile)
            }
            
            // Cell References
            cellRefs := protected.Group("/cell-references")
            {
                cellRefs.POST("", controllers.CellReference.Create)
                cellRefs.DELETE("/:id", controllers.CellReference.Delete)
            }
        }
    }
}
```

---

### Main Application

```go
// main.go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "myapp/config"
    "myapp/controllers"
    "myapp/db"
    "myapp/routes"
    "myapp/services"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize database
    database, err := db.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()
    
    // Run migrations
    if err := db.RunMigrations(database); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Initialize services
    budgetService := services.NewBudgetService(database)
    categoryService := services.NewCategoryService(database)
    transactionService := services.NewTransactionService(database)
    uploadService := services.NewUploadService(database)
    cellRefService := services.NewCellReferenceService(database)
    
    // Initialize controllers
    ctrl := &controllers.Controllers{
        Budget:        controllers.NewBudgetController(budgetService),
        Category:      controllers.NewCategoryController(categoryService),
        Transaction:   controllers.NewTransactionController(transactionService),
        Upload:        controllers.NewUploadController(uploadService),
        CellReference: controllers.NewCellReferenceController(cellRefService),
    }
    
    // Setup Gin
    r := gin.Default()
    routes.SetupRoutes(r, ctrl)
    
    // Start server
    log.Printf("Server starting on %s", cfg.ServerAddr)
    if err := r.Run(cfg.ServerAddr); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

---

## API Design Best Practices

### 1. Resource-Based URLs

✅ **Do**:
```
GET    /api/v1/budgets/2025-10
POST   /api/v1/budgets
PUT    /api/v1/budgets/2025-10
DELETE /api/v1/budgets/2025-10
```

❌ **Don't**:
```
GET    /api/v1/getBudget?month=2025-10
POST   /api/v1/createBudget
POST   /api/v1/updateBudget
POST   /api/v1/deleteBudget
```

---

### 2. Use HTTP Methods Semantically

- **GET** - Retrieve resource (idempotent, safe, cacheable)
- **POST** - Create resource (non-idempotent)
- **PUT** - Update entire resource (idempotent)
- **PATCH** - Partial update (idempotent)
- **DELETE** - Remove resource (idempotent)

---

### 3. Consistent Naming

- Use **plural nouns** for collections (`/budgets`, not `/budget`)
- Use **kebab-case** for multi-word resources (`/cell-references`)
- Use **camelCase** for JSON fields (`userId`, not `user_id`)

---

### 4. Nested Resources

✅ **Do** (show relationship):
```
GET /api/v1/budgets/2025-10/categories
POST /api/v1/budgets/2025-10/files
```

⚠️ **Limit nesting** to 2 levels:
```
❌ /api/v1/budgets/2025-10/categories/12/cell-references/78
✅ /api/v1/cell-references/78
```

---

### 5. Filtering, Sorting, Pagination

```
GET /api/v1/transactions?fileId=45&limit=50&offset=0
GET /api/v1/budgets?sort=desc&limit=12
```

---

### 6. Versioning

Always version APIs:
```
/api/v1/budgets    # Current
/api/v2/budgets    # Future breaking changes
```

---

### 7. Error Consistency

Always return same error structure:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "field": "fieldName"
  }
}
```

---

## Testing Strategy

### Unit Tests

```go
// controllers/budget_test.go
package controllers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestCreateBudget(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    // Setup
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    reqBody := map[string]string{
        "month": "2025-10",
        "name": "Test Budget",
    }
    jsonBytes, _ := json.Marshal(reqBody)
    
    c.Request = httptest.NewRequest("POST", "/api/v1/budgets", bytes.NewBuffer(jsonBytes))
    c.Request.Header.Set("Content-Type", "application/json")
    
    // Execute
    controller.CreateBudget(c)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "2025-10", response["month"])
}
```

---

### Integration Tests

```go
// tests/integration/budget_test.go
package integration

import (
    "testing"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
)

func TestBudgetWorkflow(t *testing.T) {
    // Setup test server with test database
    server := setupTestServer(t)
    defer cleanupTestDB(t)
    
    // Test: Create budget
    resp := server.POST("/api/v1/budgets", map[string]string{
        "month": "2025-10",
        "name": "Test Budget",
    })
    assert.Equal(t, 201, resp.Code)
    
    budgetID := extractID(resp.Body)
    
    // Test: Get budget
    resp = server.GET("/api/v1/budgets/2025-10")
    assert.Equal(t, 200, resp.Code)
    
    // Test: Create category
    resp = server.POST("/api/v1/budgets/2025-10/categories", map[string]interface{}{
        "name": "Income",
        "projected": 5000.00,
    })
    assert.Equal(t, 201, resp.Code)
    
    // Test: Delete budget (cascade deletes categories)
    resp = server.DELETE("/api/v1/budgets/2025-10")
    assert.Equal(t, 204, resp.Code)
}
```

---

## Performance Considerations

### 1. Database Query Optimization

- Use indexes on foreign keys (already in schema from doc 10)
- Composite indexes for common queries (`user_id`, `month`)
- Limit N+1 queries with eager loading

### 2. Pagination

Always paginate large result sets:
```go
func ListTransactions(c *gin.Context) {
    limit := c.DefaultQuery("limit", "100")
    offset := c.DefaultQuery("offset", "0")
    
    transactions, err := service.GetTransactions(limit, offset)
    // ...
}
```

### 3. Caching

Use Redis for frequently accessed data:
- Budget summaries
- Category trees
- User sessions

### 4. Rate Limiting

```go
import "github.com/gin-contrib/ratelimit"

func main() {
    r := gin.Default()
    r.Use(ratelimit.Rate(time.Minute, 60)) // 60 requests per minute
}
```

---

## Security Considerations

### 1. Authentication

- Use JWT tokens with short expiry (15 minutes)
- Refresh tokens with longer expiry (7 days)
- Store refresh tokens in httpOnly cookies

### 2. Authorization

- Always check user owns resource before modification
- Use middleware to inject `userId` from token

```go
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        claims, err := validateJWT(token)
        if err != nil {
            c.AbortWithStatus(401)
            return
        }
        c.Set("userId", claims.UserID)
        c.Next()
    }
}
```

### 3. Input Validation

- Validate all inputs with struct tags
- Sanitize user input (SQL injection protection via parameterized queries)
- Validate file uploads (size, type, content)

### 4. CORS

- Whitelist frontend origin only
- Don't use `AllowOrigins: ["*"]` in production

---

## Deployment Checklist

### Production Configuration

- [ ] Use environment variables for secrets (DATABASE_URL, JWT_SECRET)
- [ ] Enable HTTPS (TLS certificates)
- [ ] Configure CORS for production frontend URL
- [ ] Set Gin to release mode: `gin.SetMode(gin.ReleaseMode)`
- [ ] Use structured logging (not `log.Printf`)
- [ ] Set up health check endpoint (`GET /health`)
- [ ] Configure graceful shutdown
- [ ] Set up monitoring (Prometheus metrics)
- [ ] Rate limiting per IP address
- [ ] Request timeout middleware (30 seconds)

---

## Pros of This Design

- ✅ **RESTful** - Follows REST principles (resources, HTTP methods, status codes)
- ✅ **Consistent** - Same patterns for all endpoints
- ✅ **Extensible** - Easy to add new resources
- ✅ **Versioned** - `/api/v1` allows future changes
- ✅ **Type-safe** - Gin binding with struct validation
- ✅ **Testable** - Clear separation of concerns (controllers, services)
- ✅ **Documented** - Clear request/response examples
- ✅ **Error handling** - Consistent error format
- ✅ **Secure** - Authentication, authorization, input validation

---

## Cons/Challenges

- ⚠️ **Over-fetching** - REST can return more data than needed (GraphQL alternative)
- ⚠️ **Multiple requests** - Need separate calls for related data (vs GraphQL single query)
- ⚠️ **URL versioning** - Requires maintaining multiple API versions
- ⚠️ **File upload** - Multipart handling more complex than JSON
- ⚠️ **Real-time updates** - REST doesn't support push notifications (need WebSockets)

---

## Alternative Approaches (Not Selected)

### GraphQL

**Pros**:
- Single endpoint for all queries
- Client specifies exact data needed
- No over-fetching or under-fetching

**Cons**:
- More complex to implement
- Harder to cache
- Overkill for simple CRUD API

### gRPC

**Pros**:
- Binary protocol (faster than JSON)
- Type-safe with Protocol Buffers
- Bi-directional streaming

**Cons**:
- Not browser-friendly (requires gRPC-Web)
- Harder to debug (binary format)
- Unnecessary for web app

**Recommendation**: Stick with REST for simplicity and browser compatibility

---

## Next Steps

1. **Create API handlers** - Implement controllers from examples above
2. **Write tests** - Unit and integration tests for all endpoints
3. **Add authentication** - JWT middleware and auth endpoints
4. **Document with Swagger** - Auto-generate API documentation
5. **Set up CI/CD** - Automated testing and deployment
6. **Load testing** - Benchmark with `wrk` or `ab`

---

## References

- [REST API Tutorial](https://restfulapi.net/)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [HTTP Status Codes](https://httpstatuses.com/)
- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [Richardson Maturity Model](https://martinfowler.com/articles/richardsonMaturityModel.html)
- [REST API Best Practices](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
- [Gin Examples](https://github.com/gin-gonic/examples)

---

**Next Research**: Database integration (ORM evaluation - GORM vs sqlx vs sqlc)
