export interface UserID {
  Value: number;
}

export interface SheetID {
  Value: number;
}

export interface ComponentID {
  Value: number;
}

export interface User {
  Id: UserID;
  Name: string;
  Email: string;
  Password?: string;
  RefreshToken?: string;
}

export interface Sheet {
  Id: SheetID;
  UId: UserID;
  Name: string;
  Description: string;
  Template: boolean;
  LastWriteTime?: string;
  Components?: ComponentID[];
}

export enum VarType {
  Container = 0,
  Number = 1,
  Flag = 2,
  Text = 3,
  Media = 4,
}

export interface Component {
  Id: ComponentID;
  UId: UserID;
  SId: SheetID;
  Name: string;
  Description: string;
  Template: boolean;
  VarName: string;
  VarType: VarType;
  VarValue: Uint8Array;
  Style: string;
}

export interface StyleParams {
  x: number;
  y: number;
  width: number;
  height: number;
  collapsed?: boolean;
}

export interface TokenResponse {
  Token: string;
  RefreshToken: string;
}

export interface SheetFilter {
  Template?: boolean;
  Newest?: boolean;
}

export interface ApiError {
  message: string;
  status: number;
}

