export type UserRole = 'viewer' | 'editor' | 'admin';

export interface User {
  id: string;
  email: string;
  role: UserRole;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface Product {
  sku: string;
  name: string;
  volume: number;
  weight: number;
  created_at: string;
  updated_at: string;
}

export interface ShelfItem {
  id: string;
  shelf_id: string;
  sku: string;
  product_name: string;
  quantity: number;
  volume: number;
  created_at: string;
}

export interface Shelf {
  id: string;
  name: string;
  row_index: number;
  col_index: number;
  max_volume: number;
  used_volume: number;
  items: ShelfItem[];
  created_at: string;
  updated_at: string;
}

export interface WarehouseState {
  shelves: Shelf[];
  products: Product[];
  currentShelf: Shelf | null;
  loading: boolean;
  error: string | null;
}
