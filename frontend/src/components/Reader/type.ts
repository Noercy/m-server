export type ReaderType = "Horizontal" | "Vertical"
export type DisplayMode = "Single" | "Double" | "Auto";
export type ReadDirection = "LeftToRight" | "RightToLeft" | "Vertical";

export interface UserSettings {
    ReaderType:         ReaderType,
    DisplayMode:        DisplayMode,
    ReadDirection:      ReadDirection,
    SeparateFirstPage:  boolean,
}