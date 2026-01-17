import React, { useState } from 'react';
import { Shelf, Product } from '../types';
import { WarehouseService } from '../services/warehouseService';

interface ItemPanelProps {
  shelf: Shelf | null;
  products: Product[];
  userRole: string;
}

export const ItemPanel: React.FC<ItemPanelProps> = ({ shelf, products, userRole }) => {
  const [selectedSKU, setSelectedSKU] = useState('');
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(false);

  const canEdit = userRole === 'editor' || userRole === 'admin';

  const handleAddItem = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!shelf) return;

    setLoading(true);
    try {
      await WarehouseService.addItemToShelf(shelf.id, selectedSKU, quantity);
      setSelectedSKU('');
      setQuantity(1);
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to add item');
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveItem = async (itemId: string) => {
    if (!shelf) return;

    if (!confirm('Remove this item?')) return;

    setLoading(true);
    try {
      await WarehouseService.removeItemFromShelf(shelf.id, itemId);
    } catch (err) {
      alert('Failed to remove item');
    } finally {
      setLoading(false);
    }
  };

  if (!shelf) {
    return (
      <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-6">
        <p className="text-gray-500 text-center">Select a shelf to view items</p>
      </div>
    );
  }

  return (
    <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-6">
      <h2 className="text-2xl font-bold text-gray-800 mb-4">{shelf.name}</h2>

      <div className="mb-4 p-3 bg-blue-50 rounded-lg">
        <p className="text-sm text-gray-700">
          <strong>Volume:</strong> {shelf.used_volume.toFixed(1)}/{shelf.max_volume}
        </p>
        <p className="text-sm text-gray-700">
          <strong>Items:</strong> {shelf.items.length}
        </p>
      </div>

      {canEdit && (
        <form onSubmit={handleAddItem} className="mb-6 space-y-3 border-b pb-6">
          <select
            value={selectedSKU}
            onChange={(e) => setSelectedSKU(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          >
            <option value="">Select Product</option>
            {products.map((p) => (
              <option key={p.sku} value={p.sku}>
                {p.name} ({p.sku})
              </option>
            ))}
          </select>

          <input
            type="number"
            placeholder="Quantity"
            value={quantity}
            onChange={(e) => setQuantity(parseInt(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
            min="1"
          />

          <button
            type="submit"
            disabled={loading || !selectedSKU}
            className="w-full bg-green-600 text-white font-medium py-2 rounded-lg hover:bg-green-700 disabled:bg-green-400 transition"
          >
            Add Item
          </button>
        </form>
      )}

      <div className="space-y-2 max-h-64 overflow-y-auto">
        {shelf.items.length === 0 ? (
          <p className="text-gray-500 text-center py-4">No items in this shelf</p>
        ) : (
          shelf.items.map((item) => (
            <div
              key={item.id}
              className="p-3 border border-gray-200 rounded-lg flex justify-between items-start"
            >
              <div className="flex-1">
                <p className="font-medium text-gray-800">{item.product_name}</p>
                <p className="text-sm text-gray-600">{item.sku}</p>
                <p className="text-sm text-gray-600">Qty: {item.quantity}</p>
              </div>
              {canEdit && (
                <button
                  onClick={() => handleRemoveItem(item.id)}
                  disabled={loading}
                  className="ml-2 px-3 py-1 bg-red-500 text-white rounded text-sm hover:bg-red-600 disabled:bg-red-300 transition"
                >
                  Remove
                </button>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};
