import React, { useRef } from "react";
import { Shelf } from "../types";
import { use3D } from "../hooks/use3D";

interface ShelfViewerProps {
  shelves: Shelf[];
}

export const ShelfViewer: React.FC<ShelfViewerProps> = ({ shelves }) => {
  const containerRef = useRef<HTMLDivElement>(null!);
  use3D(containerRef, shelves);
  return (
    <div ref={containerRef} className="shelf-container">
      {/* conteúdo */}
    </div>
  );

  return (
    <div className="w-full h-screen bg-gray-100">
      <div ref={containerRef} className="w-full h-full" />
      <div className="absolute bottom-4 left-4 bg-white rounded-lg p-4 shadow-lg">
        <p className="text-sm text-gray-600">
          🖱️ Drag to rotate | 📱 Touch to rotate
        </p>
        <p className="text-xs text-gray-500 mt-2">
          Shelves: {shelves.length} | Total Items:{" "}
          {shelves.reduce((acc, s) => acc + s.items.length, 0)}
        </p>
      </div>
    </div>
  );
};
