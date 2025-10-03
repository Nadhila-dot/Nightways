import { useTheme } from "next-themes"
import { Toaster as Sonner, ToasterProps } from "sonner"

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = "system" } = useTheme()

  return (
    <Sonner
      theme={theme as ToasterProps["theme"]}
      style={{ 
        overflowWrap: "anywhere",
        fontFamily: "'Space Grotesk', sans-serif",
      }}
      toastOptions={{
        unstyled: true,
        classNames: {
          toast:
            "bg-white !text-black border-black border-4 font-semibold shadow-[8px_8px_0_0_#000] rounded-2xl text-[15px] flex items-center gap-3 p-5 w-[400px] [&:has(button)]:justify-between",
          description: "font-normal",
          actionButton:
            "font-bold border-2 border-black text-[12px] h-7 px-3 bg-black text-white rounded-lg shadow-[3px_3px_0_0_#000] shrink-0 hover:shadow-[1px_1px_0_0_#000] hover:translate-x-[2px] hover:translate-y-[2px] transition-all duration-100 uppercase",
          cancelButton:
            "font-bold border-2 border-black text-[12px] h-7 px-3 bg-white text-black rounded-lg shadow-[3px_3px_0_0_#000] shrink-0 hover:shadow-[1px_1px_0_0_#000] hover:translate-x-[2px] hover:translate-y-[2px] transition-all duration-100 uppercase",
          error: "bg-black text-white border-white shadow-[6px_6px_0_0_#fff]",
          success: "bg-white text-black border-black shadow-[6px_6px_0_0_#000]",
          warning: "bg-black text-white border-white shadow-[6px_6px_0_0_#fff]",
          info: "bg-white text-black border-black shadow-[6px_6px_0_0_#000]",
          loading:
            "bg-white text-black border-black shadow-[6px_6px_0_0_#000] [&[data-sonner-toast]_[data-icon]]:flex [&[data-sonner-toast]_[data-icon]]:size-4 [&[data-sonner-toast]_[data-icon]]:relative [&[data-sonner-toast]_[data-icon]]:justify-start [&[data-sonner-toast]_[data-icon]]:items-center [&[data-sonner-toast]_[data-icon]]:flex-shrink-0",
        },
      }}
      {...props}
    />
  )
}

export { Toaster }