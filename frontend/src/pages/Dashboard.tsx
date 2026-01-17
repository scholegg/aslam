import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ShelfViewer } from '../components/ShelfViewer';
import { ShelfControls } from '../components/ShelfControls';
import { ItemPanel } from '../components/ItemPanel';
import { useWarehouseStore } from '../stores/warehouseStore';
import { WarehouseService } from '../services/warehouseService';
import { apiClient } from '../services/api';
import { User } from '../types';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { shelves, products, currentShelf, setCurrentShelf } = useWarehouseStore();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const initializeDashboard = async () => {
      try {
        if (!apiClient.isAuthenticated()) {
          navigate('/login');
          return;
        }

        const userData = await apiClient.getProfile();
        setUser(userData);

        await Promise.all([
          WarehouseService.loadShelves(),
          WarehouseService.loadProducts(),
        ]);
      } catch (error) {
        console.error('Failed to initialize dashboard:', error);
        navigate('/login');
      } finally {
        setLoading(false);
      }
    };

    initializeDashboard();
  }, [navigate]);

  const handleLogout = () => {
    apiClient.clearToken();
    navigate('/login');
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-100">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-full mx-auto px-6 py-4 flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-800">Aslam</h1>
            <p className="text-sm text-gray-600">{user?.role.toUpperCase()}</p>
          </div>
          <div className="flex items-center gap-4">
            <span className="text-gray-700">{user?.email}</span>
            <button
              onClick={handleLogout}
              className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex flex-1 gap-4 p-4 overflow-hidden">
        {/* Left Sidebar */}
        <div className="w-64 space-y-4 overflow-y-auto">
          <ShelfControls onShelfSelect={(shelfId) => {
            const shelf = shelves.find(s => s.id === shelfId);
            setCurrentShelf(shelf || null);
          }} />
        </div>

        {/* 3D Viewer */}
        <div className="flex-1 rounded-lg overflow-hidden bg-white shadow">
          <ShelfViewer shelves={shelves} />
        </div>

        {/* Right Sidebar */}
        <div className="w-64 space-y-4 overflow-y-auto">
          <ItemPanel shelf={currentShelf} products={products} userRole={user?.role || 'viewer'} />
        </div>
      </div>
    </div>
  );
};
