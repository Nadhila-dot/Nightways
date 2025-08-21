import { Link } from "react-router-dom";

// Incase of errors or etc all the screen blocks are imported here
// Screenblocks will have 403, 404, 500 etc
// This is used to handle errors and display a user-friendly message


export function NotFound({ path }: { path: string }) {
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <h1 className="text-4xl mb-4">{path} was not found.</h1>
      <Link to="/" className="text-blue-500 underline">Go Home</Link>
    </div>
  );
}

