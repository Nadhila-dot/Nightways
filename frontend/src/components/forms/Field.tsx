import React, { KeyboardEvent } from 'react';
import { Input } from "@/components/ui/input";
import { Pill } from "@/components/ui/pill";

interface FormFieldProps {
  label: string;
  name: string;
  value: string;
  placeholder: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
  onKeyDown: (e: KeyboardEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
  isLoading: boolean;
  showTip: boolean;
  tipText: string;
  tipAction: string;
  onTipClick: () => void;
  isTextarea?: boolean;
  rows?: number;
}

export const FormField: React.FC<FormFieldProps> = ({
  label,
  name,
  value,
  placeholder,
  onChange,
  onKeyDown,
  isLoading,
  showTip,
  tipText,
  tipAction,
  onTipClick,
  isTextarea = false,
  rows = 3
}) => {
  const inputProps = {
    name,
    value,
    onChange,
    onKeyDown,
    placeholder,
    className: "w-full border-2 border-black rounded-lg px-3 py-2 shadow-[2px_2px_0_0_#000] focus:outline-none focus:ring-2 focus:ring-black",
    disabled: isLoading
  };

  return (
    <div>
      <label className="block font-bold mb-1 text-lg">{label}</label>
      <div className="relative">
        {isTextarea ? (
          <textarea {...inputProps} rows={rows} />
        ) : (
          <Input {...inputProps} />
        )}
        
        {isLoading && (
          <div className="absolute right-3 top-2 text-xs text-gray-500">
            Generating...
          </div>
        )}
        
        {showTip && !isLoading && (
          <div className="absolute right-3 top-2">
            <Pill 
              text={tipText}
              action={tipAction}
              onClick={onTipClick}
            />
          </div>
        )}
      </div>
    </div>
  );
};