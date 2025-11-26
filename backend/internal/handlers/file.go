package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chetbackiewicz/every-sunday-backend/internal/models"
	"github.com/chetbackiewicz/every-sunday-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

const (
	MaxFileSize       = 10 * 1024 * 1024 // 10MB
	MaxFilesPerBudget = 10
)

type FileHandler struct {
	repos *repositories.Repositories
}

func NewFileHandler(repos *repositories.Repositories) *FileHandler {
	return &FileHandler{repos: repos}
}

func (h *FileHandler) Upload(c *gin.Context) {
	month := c.Param("month")
	// In a real app with auth, we would get user ID from context
	userIDVal, exists := c.Get("userID")
	if !exists {
		// Fallback for dev if auth middleware isn't strictly enforcing or for testing
		userIDVal = 1 
	}
	userID := userIDVal.(int)

	// 1. Get Monthly Budget ID
	budget, err := h.repos.Budget.GetByMonth(c, userID, month)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	// 2. Check file count
	count, err := h.repos.File.CountByBudget(c, budget.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check file count"})
		return
	}
	if count >= MaxFilesPerBudget {
		c.JSON(http.StatusConflict, gin.H{"error": "Max files exceeded"})
		return
	}

	// 3. Get file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	if fileHeader.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// 4. Read and Parse CSV
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Assuming header row exists
	header, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read CSV header"})
		return
	}

	// Validate headers roughly
	headerMap := make(map[string]int)
	for i, col := range header {
		headerMap[strings.TrimSpace(strings.ToLower(col))] = i
	}

	// Relaxed requirements based on typical CSVs, but ensuring we have basics
	requiredCols := []string{"date", "description", "amount"}
	for _, col := range requiredCols {
		found := false
		for h := range headerMap {
			if strings.Contains(h, col) {
				found = true
				break
			}
		}
		if !found {
			// Check for alternatives (e.g. "Transaction Date" vs "Date")
			// Simple check for now
			if col == "date" {
				if _, ok := headerMap["transaction date"]; ok { continue }
				if _, ok := headerMap["post date"]; ok { continue }
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Missing required column containing: %s", col)})
			return
		}
	}
	
	// Remap strict keys if found
	getDateIdx := func() int {
		if idx, ok := headerMap["transaction date"]; ok { return idx }
		if idx, ok := headerMap["date"]; ok { return idx }
		return -1
	}
	dateIdx := getDateIdx()

	getPostDateIdx := func() int {
		if idx, ok := headerMap["post date"]; ok { return idx }
		return -1
	}
	postDateIdx := getPostDateIdx()

	getDescIdx := func() int {
		if idx, ok := headerMap["description"]; ok { return idx }
		return -1
	}
	descIdx := getDescIdx()

	getAmountIdx := func() int {
		if idx, ok := headerMap["amount"]; ok { return idx }
		return -1
	}
	amountIdx := getAmountIdx()

	getCatIdx := func() int {
		if idx, ok := headerMap["category"]; ok { return idx }
		return -1
	}
	catIdx := getCatIdx()

	getTypeIdx := func() int {
		if idx, ok := headerMap["type"]; ok { return idx }
		return -1
	}
	typeIdx := getTypeIdx()

	var transactions []*models.Transaction
	rowNum := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to read row %d", rowNum+1)})
			return
		}
		rowNum++

		tx := &models.Transaction{
			RowNumber: rowNum,
		}

		// Parse Date
		if dateIdx != -1 && dateIdx < len(record) {
			tx.TransactionDate, err = parseDate(record[dateIdx])
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid date format in row %d: %s", rowNum, record[dateIdx])})
				return
			}
		}
		if postDateIdx != -1 && postDateIdx < len(record) {
			tx.PostDate, err = parseDate(record[postDateIdx])
			if err != nil {
				tx.PostDate = tx.TransactionDate 
			}
		} else {
			tx.PostDate = tx.TransactionDate
		}

		// Parse Description
		if descIdx != -1 && descIdx < len(record) {
			tx.Description = record[descIdx]
		}

		// Parse Amount
		if amountIdx != -1 && amountIdx < len(record) {
			// Handle currency symbols if present
			amtStr := strings.ReplaceAll(record[amountIdx], "$", "")
			amtStr = strings.ReplaceAll(amtStr, ",", "")
			tx.Amount, err = strconv.ParseFloat(strings.TrimSpace(amtStr), 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid amount in row %d", rowNum)})
				return
			}
		}

		// Parse Category
		if catIdx != -1 && catIdx < len(record) {
			tx.Category = record[catIdx]
		}

		// Parse Type
		if typeIdx != -1 && typeIdx < len(record) {
			tx.Type = record[typeIdx]
		}

		transactions = append(transactions, tx)
	}

	// 5. Create File Record
	csvFile := &models.CSVFile{
		MonthlyBudgetID: budget.ID,
		Filename:        fileHeader.Filename,
		FileSize:        int(fileHeader.Size),
		RowCount:        rowNum,
	}

	csvFile, err = h.repos.File.Create(c, csvFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file record"})
		return
	}

	// 6. Save Transactions
	for _, tx := range transactions {
		tx.CSVFileID = csvFile.ID
	}
	
	err = h.repos.Transaction.CreateBatch(c, transactions)
	if err != nil {
		// Cleanup file on failure
		_ = h.repos.File.Delete(c, csvFile.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save transactions"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file": csvFile,
		"transaction_count": len(transactions),
	})
}

// ListFiles lists all files for a budget
func (h *FileHandler) ListFiles(c *gin.Context) {
	month := c.Param("month")
	userIDVal, exists := c.Get("userID")
	if !exists {
		userIDVal = 1
	}
	userID := userIDVal.(int)

	budget, err := h.repos.Budget.GetByMonth(c, userID, month)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	files, err := h.repos.File.List(c, budget.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// DeleteFile deletes a file
func (h *FileHandler) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Validation of ownership should happen here in a real app
	
	err = h.repos.File.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"01/02/2006",
		"2006-01-02",
		"1/2/2006",
		"01/02/06",
	}
	
	dateStr = strings.TrimSpace(dateStr)
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format")
}
