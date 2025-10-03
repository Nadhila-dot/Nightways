import { Button } from "@/components/ui/button";
import InvalidSessionModal from "../session/Invalid";
import Header from "@/components/Content/Header";
import PageBlock from "@/components/Content/PageBlock";
import { HelpCircleIcon } from "lucide-react";

export function HelpContainer() {
  return (
    <>
      <PageBlock header="Help" icon={<HelpCircleIcon size={56} />}>
        <div className="max-w-2xl">
          <h2 className="text-3xl font-bold mb-4">How can we help you?</h2>
          <div className="mb-8">
            <p className="mb-2 text-lg">
              Welcome to the AI Education Environment! Here’s how you can create and organize your learning materials:
            </p>
          </div>
          <div className="mb-6">
            <h3 className="text-2xl font-semibold mb-2">Sheets</h3>
            <ul className="list-disc ml-6 mb-4 text-base">
              <li>
                <span className="font-bold">Sheet</span> is the basic unit of an assignment.
              </li>
              <li>
                Each sheet contains a <span className="font-bold">Title</span>, <span className="font-bold">Description</span>, and <span className="font-bold">Name</span>.
              </li>
              <li>
                Every sheet has a hidden base64-encoded <span className="font-bold">ID</span> at the top, storing metadata like source AI, creation date, and misc data.
              </li>
            </ul>
          </div>
          <div className="mb-6">
            <h3 className="text-2xl font-semibold mb-2">Notebooks</h3>
            <ul className="list-disc ml-6 mb-4 text-base">
              <li>
                <span className="font-bold">Notebook</span> groups sheets by subject.
              </li>
              <li>
                Each notebook has a <span className="font-bold">Name</span>, <span className="font-bold">Tags</span> (keywords), <span className="font-bold">Title</span>, and <span className="font-bold">Description</span>.
              </li>
              <li>
                Notebooks include a <span className="font-bold">stack.json</span> file for metadata: source, difficulty, and more.
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-2xl font-semibold mb-2">Tips</h3>
            <ul className="list-disc ml-6 text-base">
              <li>
                Use clear titles and descriptions for easy navigation.
              </li>
              <li>
                Organize sheets into notebooks by subject or topic.
              </li>
              <li>
                Check the <span className="font-bold">stack.json</span> for details about your notebook’s origin and difficulty.
              </li>
            </ul>
          </div>
        </div>
      </PageBlock>
    </>
  );
}