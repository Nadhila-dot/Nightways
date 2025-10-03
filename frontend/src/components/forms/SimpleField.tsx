import React, { KeyboardEvent } from 'react';
import { Input } from "@/components/ui/input";

interface SimpleFormFieldProps {
  label: string;
  name: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
  isLoading?: boolean;
  placeholder?: string;
  onKeyDown?: (e: KeyboardEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
  isTextarea?: boolean;
  rows?: number;
}

export const SimpleFormField: React.FC<SimpleFormFieldProps> = ({
  label,
  name,
  value,
  onChange,
  isLoading = false,
  placeholder,
  onKeyDown,
  isTextarea = true,
  rows = 3
}) => {
  const inputProps = {
    name,
    value,
    onChange,
    onKeyDown,
    placeholder,
    className: "w-full border-2 text-black border-black rounded-lg px-3 py-2 shadow-[2px_2px_0_0_#000] focus:outline-none focus:ring-2 focus:ring-black",
    disabled: isLoading
  };

  return (
    <div>
      <label className="block font-bold mb-1 text-black text-lg">{label}</label>
      <div className="relative">
        {isTextarea ? (
          <textarea
            {...inputProps}
            rows={rows}
            placeholder={placeholder} 
          />
        ) : (
          <Input {...inputProps} />
        )}
        
        {isLoading && (
          <div className="absolute right-3 top-2 text-xs text-gray-500">
            Generating...
          </div>
        )}
      </div>
    </div>
  );
};