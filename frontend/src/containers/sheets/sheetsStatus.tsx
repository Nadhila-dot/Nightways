import React from "react";
import { useSearchParams } from "react-router-dom";
import PageBlock from "@/components/Content/PageBlock";
import { ClockIcon, SheetIcon } from "lucide-react";
import JobUpdates from "@/components/Cards/Sheets/Status";


export function SheetsStatusContainer() {
  const [searchParams] = useSearchParams();
  const jobId = searchParams.get("jobId"); // Get jobId from query string
  

  return (
    <PageBlock header={"Sheet Updates"} icon={<ClockIcon size={56} />}>
      {jobId ? (
        <JobUpdates jobId={jobId} />
      ) : (
        <div className="text-center text-lg font-medium">
          No job ID provided. Please specify a job ID in the query string (?jobId=).
        </div>
      )}
    </PageBlock>
  );
}