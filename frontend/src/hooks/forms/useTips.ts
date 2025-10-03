import { useState, useEffect } from 'react';
import { FormData } from './useFormGeneration';

export interface TipState {
  subject: boolean;
  course: boolean;
  description: boolean;
  tags: boolean;
}

export const useTips = (formData: FormData) => {
  const [showTips, setShowTips] = useState<TipState>({
    subject: false,
    course: false,
    description: false,
    tags: false
  });

  useEffect(() => {
    // Subject tip
    if (formData.subject.trim().length === 0 &&
        formData.course.trim().length >= 5) {
      setShowTips(prev => ({ ...prev, subject: true }));
    } else {
      setShowTips(prev => ({ ...prev, subject: false }));
    }
    
    // Course tip
    if (formData.course.trim().length === 0 &&
        formData.subject.trim().length >= 5) {
      setShowTips(prev => ({ ...prev, course: true }));
    } else {
      setShowTips(prev => ({ ...prev, course: false }));
    }
    
    // Description tip
    if (formData.description.trim().length === 0 &&
        formData.subject.trim().length >= 5 &&
        formData.course.trim().length >= 5) {
      setShowTips(prev => ({ ...prev, description: true }));
    } else {
      setShowTips(prev => ({ ...prev, description: false }));
    }
    
    // Tags tip
    if (formData.tags.trim().length === 0 &&
        formData.subject.trim().length >= 5 &&
        formData.course.trim().length >= 5) {
      setShowTips(prev => ({ ...prev, tags: true }));
    } else {
      setShowTips(prev => ({ ...prev, tags: false }));
    }
  }, [formData.subject, formData.course, formData.description, formData.tags]);

  return { showTips, setShowTips };
};