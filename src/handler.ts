import Snackbar from "node-snackbar";

import { newMessage } from "./message";
import { Metadata } from "./metadata";

export class Handler {
  private readonly metadata: Metadata;
  private readonly source: EventSource;

  private readonly content: HTMLElement;

  constructor(metadata: Metadata) {
    this.metadata = metadata;
    this.source = new EventSource(`${metadata.path.asset}/events`);

    this.content = document.querySelector("#content")!;
  }

  public listen = () => {
    this.source.addEventListener("error", this.onerror);
    this.source.addEventListener(this.metadata.path.document, this.onmessage);
  };

  private onerror = () => {
    this.notify("Failed to connect to server.", 0);
  };

  private onmessage = (event: MessageEvent<string>) => {
    const msg = newMessage(event.data);

    const date = new Date(msg.timestamp * 1000);

    this.notify(
      `${this.metadata.path.document} reloaded at ${date.toLocaleTimeString()}`
    );

    this.content.innerHTML = msg.content;
  };

  private notify(text: string, duration?: number) {
    Snackbar.show({ text, duration, actionTextColor: "#FFF8E7" });
  }
}
