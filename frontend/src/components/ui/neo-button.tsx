import { ButtonHTMLAttributes } from 'react'

interface NeoButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  color?: 'blue' | 'red' | 'green' | 'yellow' | 'purple'
  title: string
}

export const NeoButton = ({ 
  color = 'blue',
  title,
  className = '',
  ...props 
}: NeoButtonProps) => {
  const colorMap = {
    blue: 'bg-blue-600 hover:bg-blue-700',
    red: 'bg-red-600 hover:bg-red-700',
    green: 'bg-green-600 hover:bg-green-700',
    yellow: 'bg-yellow-600 hover:bg-yellow-700',
    purple: 'bg-purple-600 hover:bg-purple-700'
  }

  return (
    <button
      className={`
        ${colorMap[color]}
        text-white 
        border-2 
        border-black 
        rounded-lg 
        px-6 
        py-2 
        font-bold 
        shadow-[4px_4px_0_0_#000] 
        transition-all
        active:shadow-[0px_0px_0_0_#000]
        active:translate-x-[4px]
        active:translate-y-[4px]
        ${className}
      `}
      {...props}
    >
      {title}
    </button>
  )
}