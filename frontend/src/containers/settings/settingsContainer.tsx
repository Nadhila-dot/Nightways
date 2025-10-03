import PageBlock from "@/components/Content/PageBlock";
import LogoutCard from "@/components/Cards/LogoutCard";
import { Button } from "@/components/ui/button";
import { Settings2Icon, SettingsIcon } from "lucide-react";
import SetCard from "@/components/Cards/Setcard";
import RefreshCache from "@/components/Cards/Education/RefershCache";


export function SettingsContainer() {


  return (
    <PageBlock header="Settings" icon={<SettingsIcon size={56} />}>
      
        <div className="max-w-3xl ">
            <SetCard/>
        </div>
        <div className="max-w-3xl ">
            <RefreshCache/>
        </div>
        <div className="max-w-3xl ">
            <LogoutCard/>
        </div>
        
     
    </PageBlock>
  );
}