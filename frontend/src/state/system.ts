import { create } from 'zustand';
import { getSystemInfo } from '@/scripts/getSystem';

// Define the shape of the system info
interface SystemInfo {
  build: string;
  date: string;
  author: string;
  // Add other fields as needed
}

// Define the store state
interface SystemState {
  data: SystemInfo | null;
  loading: boolean;
  error: string | null;
  fetchSystemInfo: () => Promise<void>;
}

// Create the Zustand store
export const useSystemStore = create<SystemState>((set) => ({
  data: null,
  loading: false,
  error: null,
  fetchSystemInfo: async () => {
    set({ loading: true, error: null });
    try {
      const data = await getSystemInfo();
      set({ data, loading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to fetch system info', loading: false });
    }
  },
}));