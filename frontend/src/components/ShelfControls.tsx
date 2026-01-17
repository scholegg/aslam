import React, { useState } from 'react';
import { useWarehouseStore } from '../stores/warehouseStore';
import { WarehouseService } from '../services/warehouseService';

interface ShelfControlsProps {
  onShelfSelect: (shelfId: string) => void;
}

export const ShelfControls: React.FC<ShelfControlsProps> = ({ onShelfSelect }) => {
  const { shelves, error } = useWarehouseStore();
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [name, setName] = useState('');
  const [maxVolume, setMaxVolume] = useState('100');
  const [loading, setLoading] = useState(false);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const rowIndex = Math.floor(shelves.length / 4);
      const colIndex = shelves.length % 4;
      await WarehouseService.createShelf(name, rowIndex, colIndex, parseFloat(maxVolume));
      setName('');
      setMaxVolume('100');
      setShowCreateForm(false);
    } catch (err) {
      console.error('Failed to create shelf:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-6">
      <h2 className="text-2xl font-bold text-gray-800 mb-4">Shelves</h2>

      {error && (
        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded text-sm">
          {error}
        </div>
      )}

      <div className="space-y-2 max-h-80 overflow-y-auto mb-4">
        {shelves.map((shelf) => (
          <button
            key={shelf.id}
            onClick={() => onShelfSelect(shelf.id)}
            className="w-full text-left p-3 rounded-lg border border-gray-300 hover:bg-blue-50 hover:border-blue-500 transition"
          >
            <p className="font-medium text-gray-800">{shelf.name}</p>
            <p className="text-sm text-gray-600">
              Vol: {shelf.used_volume.toFixed(1)}/{shelf.max_volume} | Items: {shelf.items.length}
            </p>
          </button>
        ))}
      </div>

      {!showCreateForm ? (
        <button
          onClick={() => setShowCreateForm(true)}
          className="w-full bg-blue-600 text-white font-medium py-2 rounded-lg hover:bg-blue-700 transition"
        >
          + New Shelf
        </button>
      ) : (
        <form onSubmit={handleCreate} className="space-y-3">
          <input
            type="text"
            placeholder="Shelf Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
          <input
            type="number"
            placeholder="Max Volume"
            value={maxVolume}
            onChange={(e) => setMaxVolume(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
            step="0.1"
            min="0"
          />
          <div className="flex gap-2">
            <button
              type="submit"
              disabled={loading}
              className="flex-1 bg-green-600 text-white font-medium py-2 rounded-lg hover:bg-green-700 disabled:bg-green-400 transition"
            >
              Create
            </button>
            <button
              type="button"
              onClick={() => setShowCreateForm(false)}
              className="flex-1 bg-gray-400 text-white font-medium py-2 rounded-lg hover:bg-gray-500 transition"
            >
              Cancel
            </button>
          </div>
        </form>
      )}
    </div>
  );
};
