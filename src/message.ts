export type Message = {
  content: string;
  timestamp: number;
};

export const newMessage = (raw: string): Message => {
  return JSON.parse(raw);
};
