import React from 'react';
import { Eye, EyeOff } from 'lucide-react';

interface VisibilitySelectorProps {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export const VisibilitySelector: React.FC<VisibilitySelectorProps> = ({
  value,
  onChange
}) => {
  return (
    <div>
      <label className="font-bold mb-1 text-lg flex items-center gap-2">
        Visibility
        <div className="relative w-4 h-4">
          <EyeOff 
            className={`absolute w-4 h-4 transition-all duration-500 ease-out transform ${
              value === "private" 
                ? "opacity-100 scale-100 rotate-0" 
                : "opacity-0 scale-75 rotate-180"
            } text-gray-600`}
          />
          <Eye 
            className={`absolute w-4 h-4 transition-all duration-500 ease-out transform ${
              value === "public" 
                ? "opacity-100 scale-100 rotate-0" 
                : "opacity-0 scale-75 -rotate-180"
            } text-blue-600`}
          />
        </div>
      </label>
      <div className="flex gap-4">
        <label className="flex items-center space-x-2 cursor-pointer">
          <input
            type="radio"
            name="visibility"
            value="private"
            checked={value === "private"}
            onChange={onChange}
            className="w-5 h-5"
          />
          <span>Private</span>
        </label>
        <label className="flex items-center space-x-2 cursor-pointer">
          <input
            type="radio"
            name="visibility"
            value="public"
            checked={value === "public"}
            onChange={onChange}
            className="w-5 h-5"
          />
          <span>Public</span>
        </label>
      </div>
    </div>
  );
};