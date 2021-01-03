import { IImageData } from '../../common';
import './image.css';

interface IImageProps {
  imageData: IImageData;
}

export const Image = (props: IImageProps) => {
  let tags: JSX.Element[] = [];
  for (let i = 0; i < props.imageData.Tags.length; i++) {
    const tagName = props.imageData.Tags[i];
    tags.push(
      <div key={i} className="tag">
        {tagName}
      </div>,
    );
  }
  return (
    <div className="image-container">
      <img src={props.imageData.Url} alt="aimage" width="300" />
      <div className="tags-container">{tags}</div>
    </div>
  );
};
