package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StringRequest struct {
	Value string `json:"value" binding:"required"`
}

type StringFilter struct {    
	IsPalindrome        *bool  `form:"is_palindrome" json:"is_palindrome"`
	MinLength           int    `form:"min_length" json:"min_length"`
	MaxLength           int    `form:"max_length" json:"max_length"`
	WordCount           int    `form:"word_count" json:"word_count"`
	ContainsCharacter   string `form:"contains_character" json:"contains_character"`
}

type DataItem struct {
    ID          string             `gorm:"column:sha256_hash;primaryKey" json:"id"` 
    Value       string             `gorm:"column:value" json:"value"`             
    Properties  StringProperties   `gorm:"-" json:"properties"`                   
    CreatedAt   string          `json:"created_at"`
    
    
    Length          int         `gorm:"column:length"`
    IsPalindrome    bool        `gorm:"column:is_palindrome"`
    UniqueCharacters int        `gorm:"column:unique_characters"`
    WordCount       int         `gorm:"column:word_count"`
    CharFreqMapData CharacterFreqMap `gorm:"column:character_frequency_map" json:"-"` 
    Sha256Hash      string      `gorm:"column:sha256_hash"`
}

type FilterAPIResponse struct {
	Data           []FilterAPIResponseItem `json:"data"` 
	Count          int               `json:"count"`
	FiltersApplied StringFilter        `json:"filters_applied"`
}

type FilterAPIResponseItem struct {
    ID          string             `json:"id"`
    Value       string             `json:"value"`
    Properties  StringProperties   `json:"properties"`
    CreatedAt   string             `json:"created_at"`
}


type StringRecord struct {
	ID				string	`json:"id" gorm:"primaryKey;type:varchar(64)"`
	Value			string	`json:"value"`
	Length			int		`json:"length"`
	IsPalindrome	bool	`json:"is_palindrome"`
	UniqueCharaters	int		`json:"unique_characters"`
	WordCount		int		`json:"word_count"`
	Sha256Hash		string	`json:"sha256_hash"`
	CharFreqMap		CharacterFreqMap `json:"character_frequency_map" gorm:"type:jsonb"`
	CreatedAt		string	`json:"created_at"`
}

type NaturalLanguageFilter struct {
	WordCount         int    `json:"word_count,omitempty"`
	IsPalindrome      bool   `json:"is_palindrome,omitempty"`
	MinLength         int    `json:"min_length,omitempty"`
	MaxLength         int    `json:"max_length,omitempty"`
	ContainsCharacter string `json:"contains_character,omitempty"`
}

type ParsedQuery struct {
	Original      string                `json:"original"`
	ParsedFilters NaturalLanguageFilter `json:"parsed_filters"`
}

type NLPResponse struct {
	Data            []FilterAPIResponseItem `json:"data"`
	Count           int               `json:"count"`
	InterpretedQuery ParsedQuery      `json:"interpreted_query"`
}

func createStringHandler(c *gin.Context) {
	var req StringRequest

	if err := c.ShouldBindJSON(&req); err != nil {


		var unmarshallTypeError *json.UnsupportedTypeError

		if errors.As(err, &unmarshallTypeError) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Invalid data type for 'value' (must be string)",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid request body or missing 'value' field",
		})
		return
	}

	

	hash := computeSHA256Hash(req.Value)

	var existing StringRecord
	if err := db.First(&existing, "id = ?", hash).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "String already exists in the system",
		})
		return
	}

	props := analyzeString(req.Value)

	record := StringRecord{
		ID: hash,
		Value: req.Value,
		Length: props.Length,
		IsPalindrome: props.IsPalindrome,
		UniqueCharaters: props.UniqueCharacters,
		WordCount: props.WordCount,
		Sha256Hash: props.Sha256Hash,
		CharFreqMap: CharFreqMap(req.Value),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save record",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": 				record.ID,
		"value": 			record.Value,
		"properties": gin.H{
			"length":			record.Length,
			"is_palindrome":	record.IsPalindrome,
			"unique_characters":	record.UniqueCharaters,
			"word_count":			record.WordCount,
			"sha256_hash":			record.Sha256Hash,
			"character_frequency_map":	record.CharFreqMap,
		},
		"created_at":		record.CreatedAt,
	})
}


func getStringHandler(c *gin.Context) {
	value := c.Param("value")

	var record StringRecord
	if err := db.First(&record, "value = ?", value).Error; err != nil  {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "String does not exist in the system",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": 				record.ID,
		"value": 			record.Value,
		"properties": gin.H{
			"length":			record.Length,
			"is_palindrome":	record.IsPalindrome,
			"unique_characters":	record.UniqueCharaters,
			"word_count":			record.WordCount,
			"sha256_hash":			record.Sha256Hash,
			"character_frequency_map":	record.CharFreqMap,
		},
		"created_at":		record.CreatedAt,
	})
}


func getStringsWithFilter(c *gin.Context) {
	var filter StringFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameter values or types",
		})
		return
	}

	tx := db.Model(&StringRecord{})

	if filter.IsPalindrome != nil {
		tx = tx.Where("is_palindrome = ?", *filter.IsPalindrome)
	}


	if filter.MinLength > 0 {
		tx = tx.Where("length >= ?", filter.MinLength)
	}
	if filter.MaxLength > 0 {
		tx = tx.Where("length <= ?", filter.MaxLength)
	}

	if filter.WordCount > 0 {
		tx = tx.Where("word_count = ?", filter.WordCount)
	}

	if filter.ContainsCharacter != "" {
		tx = tx.Where("value LIKE ?", "%"+filter.ContainsCharacter+"%")
	}


	var totalCount int64
	if err := tx.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not count records",
		})
		return
	}


	var dbResults []StringRecord
	if err := tx.Find(&dbResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch data",
		})
		return
	}


	apiData := make([]FilterAPIResponseItem, len(dbResults))
	for i, item := range dbResults {
		properties := StringProperties{
			Length: item.Length,
			IsPalindrome: item.IsPalindrome,
			UniqueCharacters: item.UniqueCharaters,
			WordCount: item.WordCount,
			Sha256Hash: item.Sha256Hash,
			CharFreqMap: item.CharFreqMap,
		}

		apiData[i] = FilterAPIResponseItem{
			ID: item.ID,
			Value: item.Value,
			Properties: properties,
			CreatedAt: item.CreatedAt,
		}
	}


	response := FilterAPIResponse{
		Data:  apiData,
		Count: int(totalCount),
		FiltersApplied: filter,
	}

	c.JSON(http.StatusOK, response)
}


func extractFilterFromQuery(query string) (NaturalLanguageFilter, error) {
	q := strings.ToLower(query)

	filters := NaturalLanguageFilter{}

	if strings.Contains(q, "palindromic") || strings.Contains(q, "palindrome") {
		filters.IsPalindrome = true
	}
	if strings.Contains(q, "single word") || strings.Contains(q, "one word") {
		filters.WordCount = 1
	}
	if strings.Contains(q, "longer than 10") {
		filters.MinLength = 11
	}
	if strings.Contains(q, "shorter than 5") {
		filters.MaxLength = 4 // Example
	}
	if strings.Contains(q, "letter z") {
		filters.ContainsCharacter = "z"
	}
	if strings.Contains(q, "first vowel") {
		filters.ContainsCharacter = "a"
	}
	
	return filters, nil
}


func getNaturalLanguageFilter(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	naturalQuery := c.Query("query")
	if naturalQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'query' parameter for natural language filter."})
		return
	}

	filters, err := extractFilterFromQuery(naturalQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to interpret query: " + err.Error(),
		})
		return
	}

	tx := db.Model(&StringRecord{}) 
	
	if filters.IsPalindrome {
		tx = tx.Where("is_palindrome = ?", true)
	}
	if filters.MinLength > 0 {
		tx = tx.Where("length >= ?", filters.MinLength)
	}
	if filters.MaxLength > 0 {
		tx = tx.Where("length <= ?", filters.MaxLength)
	}
	if filters.WordCount > 0 {
		tx = tx.Where("word_count = ?", filters.WordCount)
	}
	if filters.ContainsCharacter != "" {
		tx = tx.Where("value LIKE ?", "%"+filters.ContainsCharacter+"%")
	}


	var totalCount int64
    if err := tx.Count(&totalCount).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while counting records: " + err.Error()})
        return
    }

	var dbResults []StringRecord
	if err := tx.Find(&dbResults).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while fetching data: " + err.Error()})
        return
    }


	apiData := make([]FilterAPIResponseItem, len(dbResults))
	for i, record := range dbResults {
		properties := StringProperties{
			Length:          record.Length,
            IsPalindrome:    record.IsPalindrome,
            UniqueCharacters: record.UniqueCharaters,
            WordCount:       record.WordCount,
            Sha256Hash:      record.Sha256Hash,
            CharFreqMap:     record.CharFreqMap,
		}

		apiData[i] = FilterAPIResponseItem{
            ID:          record.ID,
            Value:       record.Value,
            CreatedAt:   record.CreatedAt,
            Properties:  properties,
        }
	}

	response := NLPResponse{
		Data:           apiData,
		Count:          int(totalCount),
		InterpretedQuery: ParsedQuery{
			Original:      naturalQuery,
			ParsedFilters: filters,
		},
	}


	c.JSON(http.StatusOK, response)
}


func deleteStringRecord(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not initialized"})
		return
	}

	stringValue := c.Param("string_value")

	if stringValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing string_value in path."})
		return
	}

	hashID := computeSHA256Hash(stringValue)
	
	result := db.Delete(&StringRecord{}, "id = ?", hashID)

	if result.Error != nil {
		log.Printf("DB Error during deletion of hash %s: %v", hashID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource due to database error."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "String does not exist in the system."})
		return
	}

	c.Status(http.StatusNoContent)
}