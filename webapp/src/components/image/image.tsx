interface IImageProps {
  key: number;
  url: string;
}

export const Image = (props: IImageProps) => {
  return (
    <div key={props.key}>
      <img src={props.url} alt="aimage" width="300" className="image" />
    </div>
  );
};
