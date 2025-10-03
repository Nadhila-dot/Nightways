import React, { useState, KeyboardEvent } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { NeoButton } from "@/components/ui/neo-button";
import { FormField } from "@/components/forms/Field";
import { VisibilitySelector } from "@/components/forms/VisibilitySelector";
import { useFormGeneration, FormData } from "@/hooks/forms/useFormGeneration";
import { useTips } from "@/hooks/forms/useTips";
import { SimpleFormField } from "@/components/forms/SimpleField";
import { toast } from "sonner";
import { createSheet } from "@/scripts/sheets";

export default function CreateSheet() {
  const [formData, setFormData] = useState<FormData>({
    subject: "",
    course: "",
    description: "",
    tags: "",
    curriculum: "",
    specialInstructions: "",
    visibility: "private"
  });

  const {
    loadingState,
    generateSubject,
    generateCourse,
    generateDescription,
    generateTags,
    generateMissingFields,
  } = useFormGeneration(formData, setFormData);

  const { showTips } = useTips(formData);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement | HTMLTextAreaElement>, field: string) => {
    if (e.key === 'Tab' && !e.shiftKey) {
      const fieldActions = {
        subject: () => showTips.subject && generateSubject(),
        course: () => showTips.course && generateCourse(),
        description: () => showTips.description && generateDescription(),
        tags: () => showTips.tags && generateTags(),
      };

      const action = fieldActions[field as keyof typeof fieldActions];
      //@ts-ignore
      if (action && action()) {
        e.preventDefault();
      }
    }
  };

  const handleCreateSheet = async () => {
  console.log("Create sheet with data:", formData);

  toast.dismiss();
  try {
    const result = await createSheet(formData); // Call the API to create the sheet
    console.log(result);

    if (result.jobId) {
      toast.success(`Sheet created successfully! Redirecting to status page...`);
      // Redirect to the status page with the job ID
      setTimeout(() => {
        window.location.href = `/sheets/status?jobId=${result.jobId}`;
      }, 2000); // Add a slight delay for the toast to display
    } else {
      toast.error("Failed to retrieve job ID.");
    }
  } catch (error) {
    toast.error("Failed to create sheet.");
    console.error(error);
  }
};

  return (
    <Card className="max-w-full border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] bg-white text-black p-0">
      <CardContent className="mt-1 px-4 py-4">
        <h1 className="text-5xl font-extrabold tracking-tight mb-2" style={{ fontFamily: "'Space Grotesk', sans-serif" }}>
          Create a new sheet
        </h1>
        <div className="mb-6 text-base font-medium">
          A sheet is the simplest form of a data structure, It's a blank thought, waiting to be filled with your new paper / idea.
        </div>

        <div className="space-y-5">
          <FormField
            label="Subject"
            name="subject"
            value={formData.subject}
            placeholder="e.g. Mathematics, Computer Science, Physics"
            onChange={handleChange}
            onKeyDown={(e) => handleKeyDown(e, 'subject')}
            isLoading={loadingState.subject}
            showTip={showTips.subject}
            tipText="Press Tab"
            tipAction="generate subject"
            onTipClick={generateSubject}
            isTextarea={true}
            rows={1}
          />

          <FormField
            label="Course"
            name="course"
            value={formData.course}
            placeholder="e.g. Calculus 101, Introduction to AI"
            onChange={handleChange}
            onKeyDown={(e) => handleKeyDown(e, 'course')}
            isLoading={loadingState.course}
            showTip={showTips.course}
            tipText="Press Tab"
            tipAction="generate course"
            onTipClick={generateCourse}
            isTextarea={true}
            rows={1}
          />

          <FormField
            label="Description"
            name="description"
            value={formData.description}
            placeholder="What's this sheet about? Add some context here..."
            onChange={handleChange}
            onKeyDown={(e) => handleKeyDown(e, 'description')}
            isLoading={loadingState.description}
            showTip={showTips.description}
            tipText="Press Tab"
            tipAction="generate description"
            onTipClick={generateDescription}
            isTextarea={true}
            rows={3}
          />

          <FormField
            label="Tags"
            name="tags"
            value={formData.tags}
            placeholder="Separate tags with commas"
            onChange={handleChange}
            onKeyDown={(e) => handleKeyDown(e, 'tags')}
            isLoading={loadingState.tags}
            showTip={showTips.tags}
            tipText="Press Tab"
            tipAction="generate tags"
            onTipClick={generateTags}
            rows={1}
            isTextarea={true}
          />

          <h1 className="block font-bold mb-1 text-3xl">
            Optional
          </h1>

          <SimpleFormField
            label="Guidence / Curriculum"
            name="curriculum"
            value={formData.curriculum}
            placeholder="Tell vela more about the curicullum your following and what it needs to base it's information off of."
            onChange={handleChange} // change handler
            rows={2}
          />

          <SimpleFormField
            label="Special Instructions"
            name="specialInstructions"
            value={formData.specialInstructions}
            placeholder="Any special instructions for vela to follow when creating this sheet?"
            onChange={handleChange} // change handler
            rows={2}
          />

          <VisibilitySelector
            value={formData.visibility}
            onChange={handleChange}
          />
        </div>

        <div className="mt-8 flex gap-4">
          <NeoButton
            color="green"
            title="Fill Missing Fields"
            onClick={generateMissingFields}
            className="mr-4"
          />
          <NeoButton
            color="red"
            title="Create Sheet"
            onClick={handleCreateSheet}
          />
        </div>
      </CardContent>
    </Card>
  );
}