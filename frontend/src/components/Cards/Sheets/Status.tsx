import React, { useState, useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { connectToJobWebSocket } from "@/scripts/sheets/statuswebsocket";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Loader } from "lucide-react";
import MarkdownCard from "@/components/Cards/Markdown/markdown";
import { NeoButton } from "@/components/ui/neo-button";

type JobUpdate = {
  message: string;
  type: string;
  timestamp: string;
  step?: string;
  stage?: string;
  progress?: { current: number; total: number };
  meta?: string;
  result?: any;
  data?: any;
  error?: string;
  modal?: {
    heading: string;
    content: string;
    optional: boolean;
  };
};

export default function JobUpdates({ jobId }: { jobId: string }) {
  const [jobStatus, setJobStatus] = useState("processing");
  const [updates, setUpdates] = useState<JobUpdate[]>([]);
  const [jobTime, setJobTime] = useState("0s");
  const [startTime] = useState(Date.now());
  const [isCompleted, setIsCompleted] = useState(false);
  const [resultData, setResultData] = useState<any>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalContent, setModalContent] = useState<{heading: string, content: string, optional: boolean}>({
    heading: "",
    content: "",
    optional: true
  });
  const [isMarkdownLoading, setIsMarkdownLoading] = useState(true);
  // New state to track queued modals and non-optional modal status
  const [modalQueue, setModalQueue] = useState<Array<{heading: string, content: string, optional: boolean}>>([]);
  const [hasNonOptionalModal, setHasNonOptionalModal] = useState(false);

  useEffect(() => {
    if (isCompleted) return;
    const timer = setInterval(() => {
      const elapsed = Math.floor((Date.now() - startTime) / 1000);
      setJobTime(`${elapsed}s`);
    }, 1000);
    return () => clearInterval(timer);
  }, [startTime, isCompleted]);

  useEffect(() => {
    if (modalContent.content) {
      const timer = setTimeout(() => setIsMarkdownLoading(false), 600);
      return () => clearTimeout(timer);
    }
  }, [modalContent.content]);

  // New effect to handle modal queue
  useEffect(() => {
    if (!isModalOpen && modalQueue.length > 0) {
      const nextModal = modalQueue[0];
      setModalContent(nextModal);
      setIsModalOpen(true);
      setIsMarkdownLoading(true);
      setHasNonOptionalModal(!nextModal.optional);
      setModalQueue(prev => prev.slice(1));
    }
  }, [isModalOpen, modalQueue]);

  // Function to handle modal close with non-optional logic
  const handleModalClose = () => {
    setIsModalOpen(false);
    setHasNonOptionalModal(false);
  };

  // Function to handle showing a modal with queue logic
  const showModal = (modalData: {heading: string, content: string, optional: boolean}) => {
    if (isModalOpen) {
      // If a non-optional modal is open, queue the new modal
      setModalQueue(prev => [...prev, modalData]);
    } else {
      // Otherwise show it immediately
      setModalContent(modalData);
      setIsModalOpen(true);
      setIsMarkdownLoading(true);
      setHasNonOptionalModal(!modalData.optional);
      setTimeout(() => setIsMarkdownLoading(false), 600);
    }
  };

  useEffect(() => {
    const sessionId = localStorage.getItem("session");
    const ws = connectToJobWebSocket(jobId, sessionId, (data) => {
      // Extract the main data payload
      const payload = data.data || {};
      
      // Parse different update types
      const type = payload.type || data.type || "processing";
      let message = payload.message || data.message || "No message";
      const stage = payload.stage;
      const step = payload.step;
      const progress = payload.progress;
      const meta = payload.meta;
      const result = data.result;
      const error = payload.error || data.error;
      const modal = payload.modal;
      
      // Enhance message for stage updates
      if (type === "stage" && stage && step) {
        message = `${stage}: ${step}`;
      }

      const update: JobUpdate = {
        message,
        type,
        timestamp: new Date().toLocaleTimeString(),
        step,
        stage,
        progress,
        meta,
        result,
        data: payload,
        error,
        modal
      };

      setUpdates((prev) => [...prev, update]);

      // Handle modal content for review-out type
      if (type === "review-out" && modal && modal.content) {
        showModal({
          heading: modal.heading || "Review Content",
          content: modal.content,
          optional: modal.optional !== undefined ? modal.optional : true
        });
      }

      // Update job status based on type
      if (type === "completed") {
        setJobStatus("completed");
        setIsCompleted(true);
        if (result) {
          setResultData(result);
        }
      } else if (type === "processing" || type === "progress") {
        setJobStatus("processing");
      } else if (type === "error") {
        setJobStatus("error");
        setIsCompleted(true);
      }
    });

    return () => ws.close();
  }, [jobId]);

  const getStatusIcon = () => {
    switch (jobStatus) {
      case "completed":
        return "✓";
      case "error":
        return "✗";
      default:
        return "⏱";
    }
  };

  const getStatusText = () => {
    switch (jobStatus) {
      case "completed":
        return "Completed";
      case "error":
        return "Failed";
      default:
        return "Processing";
    }
  };

  const openPdfInNewTab = (url: string) => {
    if (url) {
      window.open(url, '_blank');
    }
  };

  const getUpdateBadgeColor = (type: string) => {
    switch (type) {
      case "completed":
        return "bg-green-400";
      case "error":
        return "bg-red-400";
      case "stage":
        return "bg-blue-400";
      case "review-out":
        return "bg-purple-400";
      case "start":
        return "bg-teal-400";
      default:
        return "bg-yellow-400";
    }
  };

  return (
    <>
      <Card className="max-w-full border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] bg-white text-black p-0">
        <CardContent className="mt-1 px-2 py-2">
          <h1 className="text-5xl font-extrabold tracking-tight mb-4" style={{ fontFamily: "'Space Grotesk', sans-serif" }}>
            Updates on {jobId}
          </h1>
          <div className="mb-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-base font-medium">
                Status: <span className={`font-bold ${jobStatus === 'completed' ? 'text-green-600' : jobStatus === 'error' ? 'text-red-600' : 'text-yellow-600'}`}>
                  {getStatusText()}
                </span>
              </span>
              <div className="flex items-center gap-2">
                <span className="text-sm font-bold">{jobTime}</span>
                <span className="text-lg">{getStatusIcon()}</span>
              </div>
            </div>
          </div>

          {/* Show result actions when job is completed */}
          {jobStatus === "completed" && resultData && resultData.pdf_url && (
            <div className="mb-4 p-4 bg-green-50 border-2 border-green-500 rounded-lg">
              <h3 className="text-lg font-bold mb-2">Job Completed Successfully!</h3>
              <button 
                onClick={() => openPdfInNewTab(resultData.pdf_url)}
                className="px-4 py-2 bg-blue-600 text-white font-medium rounded-md hover:bg-blue-700 transition-colors"
              >
                Open PDF Document
              </button>
              {resultData.metadata && (
                <div className="mt-3 p-3 bg-gray-50 border border-gray-200 rounded-md">
                  <h4 className="text-sm font-bold mb-1">Document Metadata:</h4>
                  <pre className="text-xs overflow-auto max-h-40">{JSON.stringify(resultData.metadata, null, 2)}</pre>
                </div>
              )}
            </div>
          )}
          
          <div className="space-y-2">
            {updates.length === 0 ? (
              <div className="p-3 bg-gray-100 border-2 border-black rounded text-gray-600 font-medium">
                Waiting for updates...
              </div>
            ) : (
              updates.map((update, index) => (
                <div key={index} className="flex flex-col gap-1 p-3 bg-gray-50 border-2 border-black rounded shadow-[2px_2px_0_0_#000]">
                  <div className="flex items-center gap-3">
                    <div className={`px-2 py-1 ${getUpdateBadgeColor(update.type)} border border-black rounded text-xs font-bold`}>
                      {update.type.toUpperCase()}
                    </div>
                    <span className="font-medium text-sm">{update.message}</span>
                    {update.stage && update.type !== "stage" && (
                      <span className="ml-2 text-xs bg-blue-100 border border-blue-400 rounded px-2 py-1 font-semibold">
                        {update.stage}
                      </span>
                    )}
                    {update.step && update.type !== "stage" && (
                      <span className="ml-2 text-xs bg-blue-100 border border-blue-400 rounded px-2 py-1 font-semibold">
                        {update.step}
                      </span>
                    )}
                    {update.progress && (
                      <span className="ml-2 text-xs bg-gray-200 border border-gray-400 rounded px-2 py-1 font-semibold">
                        Progress: {update.progress.current}/{update.progress.total}
                      </span>
                    )}
                    {update.meta && (
                      <span className="ml-2 text-xs bg-purple-100 border border-purple-400 rounded px-2 py-1 font-semibold">
                        {update.meta}
                      </span>
                    )}
                  </div>
                  
                  {/* Display modal content button for review-out updates */}
                  {update.type === "review-out" && update.modal && (
                    <div className="mt-2">
                      <button 
                        onClick={() => {
                          // Only allow showing content if no non-optional modal is active
                          if (!hasNonOptionalModal) {
                            showModal({
                              heading: update.modal!.heading || "Review Content",
                              content: update.modal!.content,
                              optional: update.modal!.optional !== undefined ? update.modal!.optional : true
                            });
                          }
                        }}
                        className={`px-3 py-1 ${hasNonOptionalModal ? 'bg-gray-400' : 'bg-purple-600 hover:bg-purple-700'} text-white font-medium text-sm rounded-md transition-colors`}
                        disabled={hasNonOptionalModal}
                      >
                        Show Content {hasNonOptionalModal ? "(Complete required action first)" : ""}
                      </button>
                    </div>
                  )}
                  
                  {/* Display error details */}
                  {update.type === "error" && update.error && (
                    <div className="mt-2 p-2 bg-red-50 border border-red-300 rounded text-xs font-mono">
                      <div className="font-semibold">Error details:</div>
                      <div className="overflow-x-auto">{update.error}</div>
                    </div>
                  )}
                  
                  {/* Display retry information for errors */}
                  {update.type === "error" && update.data && (
                    <div className="mt-2 p-2 bg-yellow-50 border border-yellow-300 rounded text-xs">
                      {update.data.retries !== undefined && update.data.maxRetry !== undefined && (
                        <div className="font-medium">
                          Retry attempt {update.data.retries} of {update.data.maxRetry} 
                          {update.data.willRetry ? " - Will retry" : " - Will not retry"}
                        </div>
                      )}
                    </div>
                  )}
                  
                  <div className="text-xs font-medium text-gray-500">
                    {update.timestamp}
                  </div>
                </div>
              ))
            )}
          </div>
        </CardContent>
      </Card>

      {/* Dialog for displaying markdown content */}
      <Dialog open={isModalOpen} onOpenChange={(open) => {
        // Only allow closing if modal is optional
        if (modalContent.optional || !open) {
          setIsModalOpen(open);
        }
      }}>
        <DialogContent className="max-w-3xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <div className="flex items-center justify-between mb-2">
              <DialogTitle>{modalContent.heading}</DialogTitle>
              <NeoButton
                color={modalContent.optional ? "blue" : "green"}
                title={modalContent.optional ? "Close" : "Understood"}
                onClick={handleModalClose}
                className="px-4 py-1 text-sm"
              />
            </div>
          </DialogHeader>
          <div className="py-4">
            {isMarkdownLoading ? (
              <div className="flex items-center justify-center h-32">
                <Loader className="w-8 h-8 animate-spin text-gray-600" />
                <span className="ml-3 font-medium text-lg">Parsing Markdown...</span>
              </div>
            ) : (
              <MarkdownCard content={modalContent.content} />
            )}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}