import { create } from 'zustand';
import { WarehouseState, Shelf, Product } from '../types';

interface WarehouseStore extends WarehouseState {
  setShelves: (shelves: Shelf[]) => void;
  setProducts: (products: Product[]) => void;
  setCurrentShelf: (shelf: Shelf | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  updateShelf: (shelf: Shelf) => void;
  addProduct: (product: Product) => void;
}

export const useWarehouseStore = create<WarehouseStore>((set) => ({
  shelves: [],
  products: [],
  currentShelf: null,
  loading: false,
  error: null,

  setShelves: (shelves) => set({ shelves }),
  setProducts: (products) => set({ products }),
  setCurrentShelf: (shelf) => set({ currentShelf: shelf }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),

  updateShelf: (shelf) => set((state) => ({
    shelves: state.shelves.map(s => s.id === shelf.id ? shelf : s),
    currentShelf: state.currentShelf?.id === shelf.id ? shelf : state.currentShelf,
  })),

  addProduct: (product) => set((state) => ({
    products: [...state.products, product],
  })),
}));
