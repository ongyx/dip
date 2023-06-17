export class Handler {
  private readonly source: EventSource;

  constructor(eventUrl: string) {
    this.source = new EventSource(eventUrl);
  }

  public listen() {
    this.source.onopen = this.onopen;
    this.source.onmessage = this.onmessage;
    this.source.onerror = this.onerror;
  }

  private onopen() {
    document.querySelector("body")!.prepend("connection established!");
  }

  private onmessage(event: MessageEvent<string>) {
    document.querySelector("article")!.innerHTML = event.data;
  }

  private onerror(err: Event) {
    document
      .querySelector("body")!
      .prepend(`failed to connect to server: ${err}`);
  }
}
