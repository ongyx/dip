import "./asset";

import { Handler } from "./handler";
import { newMetadata } from "./metadata";

const main = () => {
  const metadata = newMetadata();
  if (metadata === null) {
    console.log("error: failed to read document metadata");
    return;
  }

  console.log(`connecting to SSE...`);

  const handler = new Handler(metadata);

  console.log(`SSE connected!`);

  handler.listen();
};

window.onload = main;
