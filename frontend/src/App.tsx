import { BrowserRouter as Router, Routes, Route, Link, useLocation } from "react-router-dom";
import CountBtn from "@/components/count-btn";
import ReactSVG from "@/assets/react.svg";
import { Badge } from "@/components/ui/badge";
import { NotFound } from "./components/ScreenBlock";

function Home() {
  return (
    <main className="flex flex-col items-center justify-center h-screen">
      <div className="flex flex-col items-center gap-y-4">
        <div className="inline-flex items-center gap-x-4">
          <img src={ReactSVG} alt="React Logo" className="w-32" />
          <span className="text-6xl">+</span>
          <img src={"/vite.svg"} alt="Vite Logo" className="w-32" />
        </div>
        <a
          href="https://ui.shadcn.com"
          rel="noopener noreferrer nofollow"
          target="_blank"
        >
          <Badge variant="outline">shadcn/ui</Badge>
        </a>
        <CountBtn />
      </div>
    </main>
  );
}

function About() {
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-4xl mb-4">About Page</h1>
      <Link to="/" className="text-blue-500 underline">Go Home</Link>
    </div>
  );
}


function App() {


  return (
    <Router>
      <nav className="absolute top-4 left-4">
        <Link to="/" className="mr-4">Home</Link>
        <Link to="/about">About</Link>
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
        {/*//@ts-ignore */}
        <Route path="*" element={<NotFound path={window.location.pathname} />} />
      </Routes>
    </Router>
  );
}

export default App;