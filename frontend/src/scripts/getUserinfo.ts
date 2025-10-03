import { QueryClient } from '@tanstack/query-core';
import http from '@/http';

const queryClient = new QueryClient();

async function fetchUserInfo(session: string) {
  const { data } = await http.get(`/api/v1/session/${session}`);
  return data;
}

export function getUserInfo(session: string) {
  return queryClient.fetchQuery({
    queryKey: ['userInfo', session],
    queryFn: () => fetchUserInfo(session),
    staleTime: 1000 * 60, // 1 minute
  });
}