import { QueryClient } from '@tanstack/query-core';
import http from '@/http';

const queryClient = new QueryClient();

async function fetchSystemInfo() {
  const { data } = await http.get('/api/v1/system');
  return data;
}

export function getSystemInfo() {
  return queryClient.fetchQuery({
    queryKey: ['systemInfo'],
    queryFn: fetchSystemInfo,
    staleTime: 500 * 60, // 1/2 minute
  });
}

