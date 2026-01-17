import { apiClient } from './api';
import { useWarehouseStore } from '../stores/warehouseStore';

export class WarehouseService {
  static async loadShelves(): Promise<void> {
    try {
      const store = useWarehouseStore.getState();
      store.setLoading(true);
      store.setError(null);

      const shelves = await apiClient.listShelves();
      store.setShelves(shelves);
    } catch (error) {
      const store = useWarehouseStore.getState();
      store.setError('Failed to load shelves');
      console.error('Error loading shelves:', error);
    } finally {
      const store = useWarehouseStore.getState();
      store.setLoading(false);
    }
  }

  static async loadProducts(): Promise<void> {
    try {
      const store = useWarehouseStore.getState();
      store.setLoading(true);
      store.setError(null);

      const products = await apiClient.listProducts();
      store.setProducts(products);
    } catch (error) {
      const store = useWarehouseStore.getState();
      store.setError('Failed to load products');
      console.error('Error loading products:', error);
    } finally {
      const store = useWarehouseStore.getState();
      store.setLoading(false);
    }
  }

  static async createProduct(sku: string, name: string, volume: number, weight: number): Promise<void> {
    try {
      const store = useWarehouseStore.getState();
      const product = await apiClient.createProduct(sku, name, volume, weight);
      store.addProduct(product);
    } catch (error) {
      console.error('Error creating product:', error);
      throw error;
    }
  }

  static async createShelf(name: string, rowIndex: number, colIndex: number, maxVolume: number): Promise<void> {
    try {
      const shelf = await apiClient.createShelf(name, rowIndex, colIndex, maxVolume);
      const store = useWarehouseStore.getState();
      store.setShelves([...store.shelves, shelf]);
    } catch (error) {
      console.error('Error creating shelf:', error);
      throw error;
    }
  }

  static async addItemToShelf(shelfId: string, sku: string, quantity: number): Promise<void> {
    try {
      await apiClient.addItemToShelf(shelfId, sku, quantity);
      const shelf = await apiClient.getShelf(shelfId);
      const store = useWarehouseStore.getState();
      store.updateShelf(shelf);
    } catch (error) {
      console.error('Error adding item to shelf:', error);
      throw error;
    }
  }

  static async removeItemFromShelf(shelfId: string, itemId: string): Promise<void> {
    try {
      await apiClient.removeItemFromShelf(shelfId, itemId);
      const shelf = await apiClient.getShelf(shelfId);
      const store = useWarehouseStore.getState();
      store.updateShelf(shelf);
    } catch (error) {
      console.error('Error removing item from shelf:', error);
      throw error;
    }
  }

  static async updateItemQuantity(shelfId: string, itemId: string, quantity: number): Promise<void> {
    try {
      await apiClient.updateItemQuantity(shelfId, itemId, quantity);
      const shelf = await apiClient.getShelf(shelfId);
      const store = useWarehouseStore.getState();
      store.updateShelf(shelf);
    } catch (error) {
      console.error('Error updating item quantity:', error);
      throw error;
    }
  }
}
