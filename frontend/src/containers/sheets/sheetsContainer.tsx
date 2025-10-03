import PageBlock from "@/components/Content/PageBlock";
import LogoutCard from "@/components/Cards/LogoutCard";
import { Button } from "@/components/ui/button";
import { Settings2Icon, SettingsIcon, SheetIcon } from "lucide-react";
import SetCard from "@/components/Cards/Setcard";
import CreateSheet from "@/components/Cards/Sheets/SheetCreate";


export function SheetsContainer() {


  return (
    <PageBlock header="Sheets" icon={<SheetIcon size={56} />}>
       <CreateSheet/>
    </PageBlock>
  );
}