export interface IImageData {
  ID: string;
  Url: string;
  Tags: ITagData[];
}

export interface ITagData {
  Name: string;
  IsGenerated: boolean;
}
