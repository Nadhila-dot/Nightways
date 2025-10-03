import PageBlock from "@/components/Content/PageBlock";
import LogoutCard from "@/components/Cards/LogoutCard";
import { Button } from "@/components/ui/button";
import { NotebookIcon, Settings2Icon, SettingsIcon, SheetIcon } from "lucide-react";
import SetCard from "@/components/Cards/Setcard";
import CreateSheet from "@/components/Cards/Sheets/SheetCreate";
import NotebookListCard from "@/components/Cards/Notebook/Notebooks";



export function NotebooksContainer() {


  return (
    <PageBlock header="Notebooks" icon={<NotebookIcon size={56} />}>
       <NotebookListCard/>
    </PageBlock>
  );
}