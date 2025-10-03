import { QueryClient } from '@tanstack/query-core';
import http from '@/http';

const queryClient = new QueryClient();

async function fetchTags({ subject, course, description }: { subject: string; course: string; description: string }) {
  const { data } = await http.post('/api/v1/sheets/generate-tags', {
    subject,
    course,
    description,
  });
  return data.tags;
}

export function getTagsQuery(params: { subject: string; course: string; description: string }) {
  return queryClient.fetchQuery({
    queryKey: ['tags', params],
    queryFn: () => fetchTags(params),
    staleTime: 60 * 1000, // 1 minute
  });
}