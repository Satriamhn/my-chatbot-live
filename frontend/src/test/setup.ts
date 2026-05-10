import '@testing-library/jest-dom';
import { vi } from 'vitest';

class ResizeObserverMock {
  observe() {}
  unobserve() {}
  disconnect() {}
}

window.ResizeObserver = ResizeObserverMock;

window.scrollTo = vi.fn();
window.HTMLElement.prototype.scrollIntoView = vi.fn();

// Mock jvectormap to avoid VML/canvas issues in JSDOM
vi.mock('@react-jvectormap/core', () => ({
  VectorMap: () => null,
}));
vi.mock('@react-jvectormap/world', () => ({
  worldMill: {},
}));

// Mock ApexCharts
vi.mock('react-apexcharts', () => ({
  default: () => null,
}));
