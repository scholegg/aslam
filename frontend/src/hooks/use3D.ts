import { useEffect, useRef } from 'react';
import * as THREE from 'three';
import { Shelf } from '../types';

export const use3D = (containerRef: React.RefObject<HTMLDivElement>, shelves: Shelf[]) => {
  const sceneRef = useRef<THREE.Scene | null>(null);
  const cameraRef = useRef<THREE.Camera | null>(null);
  const rendererRef = useRef<THREE.WebGLRenderer | null>(null);
  const shelvesRef = useRef<THREE.Group[]>([]);

  useEffect(() => {
    if (!containerRef.current) return;

    // Scene setup
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0xf0f0f0);
    sceneRef.current = scene;

    // Camera
    const width = containerRef.current.clientWidth;
    const height = containerRef.current.clientHeight;
    const camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 10000);
    camera.position.set(5, 5, 5);
    camera.lookAt(0, 0, 0);
    cameraRef.current = camera;

    // Renderer
    const renderer = new THREE.WebGLRenderer({ antialias: true });
    renderer.setSize(width, height);
    renderer.shadowMap.enabled = true;
    containerRef.current.appendChild(renderer.domElement);
    rendererRef.current = renderer;

    // Lighting
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.6);
    scene.add(ambientLight);

    const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
    directionalLight.position.set(10, 10, 10);
    directionalLight.castShadow = true;
    directionalLight.shadow.mapSize.width = 2048;
    directionalLight.shadow.mapSize.height = 2048;
    scene.add(directionalLight);

    // Grid ground
    const gridHelper = new THREE.GridHelper(20, 20, 0xcccccc, 0xeeeeee);
    scene.add(gridHelper);

    // Simple orbit controls
    let isRotating = false;
    let previousMousePosition = { x: 0, y: 0 };

    const onMouseDown = (e: MouseEvent) => {
      isRotating = true;
      previousMousePosition = { x: e.clientX, y: e.clientY };
    };

    const onMouseMove = (e: MouseEvent) => {
      if (!isRotating || !cameraRef.current) return;

      const deltaX = e.clientX - previousMousePosition.x;
      const deltaY = e.clientY - previousMousePosition.y;

      const camera = cameraRef.current as THREE.PerspectiveCamera;
      const distance = camera.position.length();

      // Rotate around Y axis
      const theta = Math.atan2(camera.position.z, camera.position.x) + deltaX * 0.005;
      const phi = Math.acos(camera.position.y / distance) + deltaY * 0.005;

      const clampedPhi = Math.max(0.1, Math.min(Math.PI - 0.1, phi));

      camera.position.x = distance * Math.sin(clampedPhi) * Math.cos(theta);
      camera.position.y = distance * Math.cos(clampedPhi);
      camera.position.z = distance * Math.sin(clampedPhi) * Math.sin(theta);
      camera.lookAt(0, 0, 0);

      previousMousePosition = { x: e.clientX, y: e.clientY };
    };

    const onMouseUp = () => {
      isRotating = false;
    };

    // Touch support
    let touchStart = { x: 0, y: 0 };

    const onTouchStart = (e: TouchEvent) => {
      if (e.touches.length === 1) {
        touchStart = { x: e.touches[0].clientX, y: e.touches[0].clientY };
      }
    };

    const onTouchMove = (e: TouchEvent) => {
      if (e.touches.length !== 1 || !cameraRef.current) return;

      const deltaX = e.touches[0].clientX - touchStart.x;
      const deltaY = e.touches[0].clientY - touchStart.y;

      const camera = cameraRef.current as THREE.PerspectiveCamera;
      const distance = camera.position.length();

      const theta = Math.atan2(camera.position.z, camera.position.x) + deltaX * 0.005;
      const phi = Math.acos(camera.position.y / distance) + deltaY * 0.005;

      const clampedPhi = Math.max(0.1, Math.min(Math.PI - 0.1, phi));

      camera.position.x = distance * Math.sin(clampedPhi) * Math.cos(theta);
      camera.position.y = distance * Math.cos(clampedPhi);
      camera.position.z = distance * Math.sin(clampedPhi) * Math.sin(theta);
      camera.lookAt(0, 0, 0);

      touchStart = { x: e.touches[0].clientX, y: e.touches[0].clientY };
    };

    renderer.domElement.addEventListener('mousedown', onMouseDown);
    renderer.domElement.addEventListener('mousemove', onMouseMove);
    renderer.domElement.addEventListener('mouseup', onMouseUp);
    renderer.domElement.addEventListener('touchstart', onTouchStart);
    renderer.domElement.addEventListener('touchmove', onTouchMove);

    // Animation loop
    const animate = () => {
      requestAnimationFrame(animate);
      renderer.render(scene, camera);
    };

    animate();

    // Handle resize
    const handleResize = () => {
      const width = containerRef.current?.clientWidth ?? window.innerWidth;
      const height = containerRef.current?.clientHeight ?? window.innerHeight;

      (cameraRef.current as THREE.PerspectiveCamera).aspect = width / height;
      (cameraRef.current as THREE.PerspectiveCamera).updateProjectionMatrix();
      renderer.setSize(width, height);
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      renderer.domElement.removeEventListener('mousedown', onMouseDown);
      renderer.domElement.removeEventListener('mousemove', onMouseMove);
      renderer.domElement.removeEventListener('mouseup', onMouseUp);
      renderer.domElement.removeEventListener('touchstart', onTouchStart);
      renderer.domElement.removeEventListener('touchmove', onTouchMove);
      containerRef.current?.removeChild(renderer.domElement);
    };
  }, [containerRef]);

  // Update shelves
  useEffect(() => {
    if (!sceneRef.current) return;

    // Remove old shelves
    shelvesRef.current.forEach(shelf => {
      sceneRef.current?.remove(shelf);
    });
    shelvesRef.current = [];

    // Add new shelves
    shelves.forEach((shelfData) => {
      const shelfGroup = new THREE.Group();

      // Create shelf geometry
      const geometry = new THREE.BoxGeometry(2, 0.2, 1);
      const material = new THREE.MeshStandardMaterial({
        color: 0x8b7355,
        metalness: 0.3,
        roughness: 0.7,
      });
      const shelf = new THREE.Mesh(geometry, material);
      shelf.castShadow = true;
      shelf.receiveShadow = true;

      // Position
      shelf.position.set(shelfData.col_index * 2.5, shelfData.row_index * 1.2, 0);

      shelfGroup.add(shelf);

      // Add items representation
      shelfData.items.forEach((item, index) => {
        const boxGeometry = new THREE.BoxGeometry(0.3, 0.3, 0.3);
        const boxMaterial = new THREE.MeshStandardMaterial({
          color: Math.random() * 0xffffff,
          metalness: 0.2,
          roughness: 0.8,
        });
        const box = new THREE.Mesh(boxGeometry, boxMaterial);
        box.castShadow = true;
        box.receiveShadow = true;
        box.position.set(
          -0.6 + (index % 3) * 0.4,
          0.2,
          -0.3 + Math.floor(index / 3) * 0.3
        );
        shelfGroup.add(box);
      });

      sceneRef.current!.add(shelfGroup);
      shelvesRef.current.push(shelfGroup);
    });
  }, [shelves]);

  return {
    scene: sceneRef.current,
    camera: cameraRef.current,
    renderer: rendererRef.current,
  };
};
