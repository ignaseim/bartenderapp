// User models
export interface User {
  id: number;
  username: string;
  email: string;
  role: UserRole;
  created_at: string;
  updated_at: string;
}

export type UserRole = 'admin' | 'bartender' | 'guest';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  refresh_token: string;
  user: User;
}

// Ingredient models
export interface Ingredient {
  id: number;
  name: string;
  category: string;
  package_size_ml: number;
  package_cost_cents: number;
  created_at: string;
  updated_at: string;
}

export interface IngredientStock {
  ingredient_id: number;
  quantity_ml: number;
  updated_at: string;
}

// Recipe models
export interface Recipe {
  id: number;
  name: string;
  method: string;
  glass: string;
  garnish: string;
  instructions: string;
  created_by: number;
  created_at: string;
  updated_at: string;
  items?: RecipeItem[];
  cost_cents?: number;
  price_cents?: number;
  can_make?: boolean;
}

export interface RecipeItem {
  recipe_id: number;
  ingredient_id: number;
  amount_ml: number;
  ingredient_name?: string;
}

// Order models
export interface Order {
  id: number;
  customer_id: number | null;
  bartender_id: number | null;
  status: OrderStatus;
  created_at: string;
  updated_at: string;
  items?: OrderItem[];
  total_cents?: number;
}

export type OrderStatus = 'pending' | 'accepted' | 'completed' | 'canceled';

export interface OrderItem {
  order_id: number;
  recipe_id: number;
  quantity: number;
  price_cents: number;
  status: OrderItemStatus;
  recipe_name?: string;
}

export type OrderItemStatus = 'pending' | 'preparing' | 'ready' | 'delivered' | 'canceled';

// Other models
export interface BartenderSkill {
  user_id: number;
  recipe_id: number;
  recipe_name?: string;
}

export interface InventoryTransaction {
  id: number;
  ingredient_id: number;
  quantity_ml: number;
  transaction_type: 'purchase' | 'usage' | 'waste' | 'adjustment';
  reference_id: number | null;
  created_by: number | null;
  created_at: string;
  ingredient_name?: string;
  created_by_name?: string;
}

export interface RevenueReport {
  bartender_id?: number;
  bartender_name?: string;
  order_count: number;
  total_cents: number;
  start_date: string;
  end_date: string;
} 