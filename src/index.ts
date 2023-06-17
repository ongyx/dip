import "github-markdown-css/github-markdown.css";

import { Handler } from "./handler";
import { parseMetadata } from "./metadata";

const main = () => {
  const metadata = parseMetadata();
  if (metadata === null) {
    console.log("error: failed to read document metadata");
    return;
  }

  const eventUrl = `${metadata.path.asset}/events?stream=${metadata.path.document}`;

  console.log(`connecting to SSE via ${eventUrl}`);

  const handler = new Handler(eventUrl);
  handler.listen();

  console.log(`SSE connected!`);
};

window.onload = main;
