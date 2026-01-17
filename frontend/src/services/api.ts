import axios, { AxiosInstance, AxiosError } from "axios";
import { AuthResponse, User, Product, Shelf, ShelfItem } from "../types";

class APIClient {
  private client: AxiosInstance;
  private token: string | null = null;

  constructor() {
    const baseURL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

    this.client = axios.create({
      baseURL,
      headers: {
        "Content-Type": "application/json",
      },
    });

    this.loadToken();

    // Interceptor para adicionar token
    this.client.interceptors.request.use(
      (config) => {
        if (this.token) {
          config.headers.Authorization = `Bearer ${this.token}`;
        }
        return config;
      },
      (error) => Promise.reject(error),
    );

    // Interceptor para erros
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          this.clearToken();
          window.location.href = "/login";
        }
        return Promise.reject(error);
      },
    );
  }

  private loadToken(): void {
    this.token = localStorage.getItem("auth_token");
  }

  setToken(token: string): void {
    this.token = token;
    localStorage.setItem("auth_token", token);
  }

  clearToken(): void {
    this.token = null;
    localStorage.removeItem("auth_token");
  }

  isAuthenticated(): boolean {
    return !!this.token;
  }

  // Auth endpoints
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await this.client.post<AuthResponse>("/auth/login", {
      email,
      password,
    });
    return response.data;
  }

  async getProfile(): Promise<User> {
    const response = await this.client.get<User>("/auth/profile");
    return response.data;
  }

  async updateProfile(email: string): Promise<User> {
    const response = await this.client.put<User>("/auth/profile", { email });
    return response.data;
  }

  // User management endpoints
  async createUser(
    email: string,
    password: string,
    role: string,
  ): Promise<User> {
    const response = await this.client.post<User>("/users", {
      email,
      password,
      role,
    });
    return response.data;
  }

  async listUsers(): Promise<User[]> {
    const response = await this.client.get<{ users: User[] }>("/users");
    return response.data.users;
  }

  async deleteUser(id: string): Promise<void> {
    await this.client.delete(`/users/${id}`);
  }

  // Product endpoints
  async createProduct(
    sku: string,
    name: string,
    volume: number,
    weight: number,
  ): Promise<Product> {
    const response = await this.client.post<Product>("/products", {
      sku,
      name,
      volume,
      weight,
    });
    return response.data;
  }

  async getProduct(sku: string): Promise<Product> {
    const response = await this.client.get<Product>(`/products/${sku}`);
    return response.data;
  }

  async listProducts(): Promise<Product[]> {
    const response = await this.client.get<{ products: Product[] }>(
      "/products",
    );
    return response.data.products;
  }

  async updateProduct(
    sku: string,
    name?: string,
    volume?: number,
    weight?: number,
  ): Promise<Product> {
    const response = await this.client.put<Product>(`/products/${sku}`, {
      ...(name && { name }),
      ...(volume && { volume }),
      ...(weight && { weight }),
    });
    return response.data;
  }

  async deleteProduct(sku: string): Promise<void> {
    await this.client.delete(`/products/${sku}`);
  }

  // Shelf endpoints
  async createShelf(
    name: string,
    rowIndex: number,
    colIndex: number,
    maxVolume: number,
  ): Promise<Shelf> {
    const response = await this.client.post<Shelf>("/shelves", {
      name,
      row_index: rowIndex,
      col_index: colIndex,
      max_volume: maxVolume,
    });
    return response.data;
  }

  async getShelf(id: string): Promise<Shelf> {
    const response = await this.client.get<Shelf>(`/shelves/${id}`);
    return response.data;
  }

  async listShelves(): Promise<Shelf[]> {
    const response = await this.client.get<{ shelves: Shelf[] }>("/shelves");
    return response.data.shelves;
  }

  async updateShelf(
    id: string,
    name?: string,
    maxVolume?: number,
  ): Promise<Shelf> {
    const response = await this.client.put<Shelf>(`/shelves/${id}`, {
      ...(name && { name }),
      ...(maxVolume && { max_volume: maxVolume }),
    });
    return response.data;
  }

  async deleteShelf(id: string): Promise<void> {
    await this.client.delete(`/shelves/${id}`);
  }

  async addItemToShelf(
    shelfId: string,
    sku: string,
    quantity: number,
  ): Promise<ShelfItem> {
    const response = await this.client.post<ShelfItem>(
      `/shelves/${shelfId}/items`,
      {
        sku,
        quantity,
      },
    );
    return response.data;
  }

  async removeItemFromShelf(shelfId: string, itemId: string): Promise<void> {
    await this.client.delete(`/shelves/${shelfId}/items/${itemId}`);
  }

  async updateItemQuantity(
    shelfId: string,
    itemId: string,
    quantity: number,
  ): Promise<void> {
    await this.client.put(`/shelves/${shelfId}/items/${itemId}`, {
      quantity,
    });
  }
}

export const apiClient = new APIClient();
