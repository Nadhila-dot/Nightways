import { Button } from "@/components/ui/button";
import InvalidSessionModal from "../session/Invalid";
import Header from "@/components/Content/Header";
import PageBlock from "@/components/Content/PageBlock";
import { HomeIcon } from "lucide-react";
import SheetListCard from "@/components/Cards/Sheets/SheetList";


export function HomeContainer() {


  return (
    <>
    <PageBlock header="Home" icon={<HomeIcon size={56} />}>
        <Header/>
        <div>
          <SheetListCard/>
        </div>
    </PageBlock>
    </>
  );
}