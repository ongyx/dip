export type Metadata = {
  path: {
    asset: string;
    document: string;
  };
};

export const parseMetadata = (): Metadata | null => {
  const raw = document.querySelector("script[id=metadata]")!.textContent;

  if (raw !== null) {
    return JSON.parse(raw);
  }

  return null;
};
