import { QueryClient } from '@tanstack/query-core';
import http from '@/http';

const queryClient = new QueryClient();

export interface SheetData {
  subject: string;
  course: string;
  description: string;
}

export interface SheetResponse {
  subject?: string;
  course?: string;
  description?: string;
  tags?: string[];
}

async function fetchTags({ subject, course, description }: SheetData) {
  const { data } = await http.post('/api/v1/sheets/generate-tags', {
    subject,
    course,
    description,
  });
  return data.tags;
}

async function fetchSubject({ course, description, generateTags = false }: 
  Pick<SheetData, 'course' | 'description'> & { generateTags?: boolean }) {
  const { data } = await http.post('/api/v1/sheets/generate-subject', {
    course,
    description,
    generateTags
  });
  return data as SheetResponse;
}

async function fetchCourse({ subject, description, generateTags = false }: 
  Pick<SheetData, 'subject' | 'description'> & { generateTags?: boolean }) {
  const { data } = await http.post('/api/v1/sheets/generate-course', {
    subject,
    description,
    generateTags
  });
  return data as SheetResponse;
}

async function fetchDescription({ subject, course, generateTags = false }: 
  Pick<SheetData, 'subject' | 'course'> & { generateTags?: boolean }) {
  const { data } = await http.post('/api/v1/sheets/generate-description', {
    subject,
    course,
    generateTags
  });
  return data as SheetResponse;
}

export function getTagsQuery(params: SheetData) {
  return queryClient.fetchQuery({
    queryKey: ['tags', params],
    queryFn: () => fetchTags(params),
    staleTime: 60 * 1000, // 1 minute
  });
}

export function getSubjectQuery(params: Pick<SheetData, 'course' | 'description'> & { generateTags?: boolean }) {
  return queryClient.fetchQuery({
    queryKey: ['subject', params],
    queryFn: () => fetchSubject(params),
    staleTime: 60 * 1000,
  });
}

export function getCourseQuery(params: Pick<SheetData, 'subject' | 'description'> & { generateTags?: boolean }) {
  return queryClient.fetchQuery({
    queryKey: ['course', params],
    queryFn: () => fetchCourse(params),
    staleTime: 60 * 1000,
  });
}

export function getDescriptionQuery(params: Pick<SheetData, 'subject' | 'course'> & { generateTags?: boolean }) {
  return queryClient.fetchQuery({
    queryKey: ['description', params],
    queryFn: () => fetchDescription(params),
    staleTime: 60 * 1000,
  });
}