import React from "react";

export default function JsonResponse() {
  const data = { message: "Hello from component JSON!" };

  {/** A bit of hacky code to send json 👀 */}
  {/** This should jst have the important stuff */}
  return (
    <pre>{JSON.stringify(data, null, 2)}</pre>
  );
}
