package models

import (
	"time"
)

// User represents a system user
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Ingredient represents a cocktail ingredient
type Ingredient struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	PackageSizeML   float64   `json:"package_size_ml"`
	PackageCostCents int      `json:"package_cost_cents"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// IngredientStock represents the current stock of an ingredient
type IngredientStock struct {
	IngredientID int       `json:"ingredient_id"`
	QuantityML   float64   `json:"quantity_ml"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Recipe represents a cocktail recipe
type Recipe struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Method       string    `json:"method"`
	Glass        string    `json:"glass"`
	Garnish      string    `json:"garnish"`
	Instructions string    `json:"instructions"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Joined fields
	Items        []RecipeItem `json:"items,omitempty"`
	CostCents    int          `json:"cost_cents,omitempty"`
	PriceCents   int          `json:"price_cents,omitempty"`
	CanMake      bool         `json:"can_make,omitempty"`
}

// RecipeItem represents an ingredient in a recipe
type RecipeItem struct {
	RecipeID     int     `json:"recipe_id"`
	IngredientID int     `json:"ingredient_id"`
	AmountML     float64 `json:"amount_ml"`
	
	// Joined fields
	IngredientName string `json:"ingredient_name,omitempty"`
}

// Order represents a customer order
type Order struct {
	ID          int       `json:"id"`
	CustomerID  *int      `json:"customer_id"`
	BartenderID *int      `json:"bartender_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Joined fields
	Items []OrderItem `json:"items,omitempty"`
	Total int         `json:"total_cents,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	OrderID     int    `json:"order_id"`
	RecipeID    int    `json:"recipe_id"`
	Quantity    int    `json:"quantity"`
	PriceCents  int    `json:"price_cents"`
	Status      string `json:"status"`
	
	// Joined fields
	RecipeName  string `json:"recipe_name,omitempty"`
}

// InventoryTransaction represents a change in ingredient stock
type InventoryTransaction struct {
	ID              int       `json:"id"`
	IngredientID    int       `json:"ingredient_id"`
	QuantityML      float64   `json:"quantity_ml"`
	TransactionType string    `json:"transaction_type"`
	ReferenceID     *int      `json:"reference_id"`
	CreatedBy       *int      `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	
	// Joined fields
	IngredientName string `json:"ingredient_name,omitempty"`
	CreatedByName  string `json:"created_by_name,omitempty"`
}

// BartenderSkill represents a cocktail a bartender can make
type BartenderSkill struct {
	UserID   int `json:"user_id"`
	RecipeID int `json:"recipe_id"`
	
	// Joined fields
	RecipeName string `json:"recipe_name,omitempty"`
}

// Auth related models

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a successful login
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// RecipeCost represents a calculated recipe cost
type RecipeCost struct {
	RecipeID         int    `json:"recipe_id"`
	Name             string `json:"name"`
	IngredientCostCents int  `json:"ingredient_cost_cents"`
	SuggestedPriceCents int `json:"suggested_price_cents"`
}

// RevenueReport represents a revenue report for a time period
type RevenueReport struct {
	BartenderID   *int   `json:"bartender_id,omitempty"`
	BartenderName string `json:"bartender_name,omitempty"`
	OrderCount    int    `json:"order_count"`
	TotalCents    int    `json:"total_cents"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
} 