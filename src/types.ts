export type DateTimeString = string;

export interface Tracker {
  id: string;
  description: string;
  start: DateTimeString;
  stop?: DateTimeString;
  elapsed?: any;
}
