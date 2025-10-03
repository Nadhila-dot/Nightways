import React from "react";

type PageBlockProps = {
  header: string;
  icon?: React.ReactNode;
  children: React.ReactNode;
};

const PageBlock: React.FC<PageBlockProps> = ({ header, icon, children }) => {
  return (
    <div className="relative px-0 pt-16 pb-8 min-h-[60vh] rounded-2xl overflow-hidden">
      <h1
        className="absolute top-6 left-0 m-0 flex items-center gap-4 text-7xl font-extrabold text-white font-spaceGrotesk tracking-wider z-10 select-none"
        style={{
          textShadow: "2px 2px 0 white"
        }}
      >
        {icon && <span className="text-white">{icon}</span>}
        <span>{header}</span>
      </h1>
      <div className="mt-20 ml-0">
        {children}
      </div>
    </div>
  );
};

export default PageBlock;