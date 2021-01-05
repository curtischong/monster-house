import { IImageData } from '../../common';
import './image.css';

interface IImageProps {
  imageData: IImageData;
}

export const Image = (props: IImageProps) => {
  let tags: JSX.Element[] = [];
  for (let i = 0; i < props.imageData.Tags.length; i++) {
    const { Name, IsGenerated } = props.imageData.Tags[i];

    let className = 'tag';
    if (IsGenerated) {
      className += ' generated-tag';
    } else {
      className += ' user-tag';
    }

    tags.push(
      <div key={i} className={className}>
        <p className="tag-text">{Name}</p>
      </div>,
    );
  }
  return (
    <div className="image-container">
      <img className="image" src={props.imageData.Url} alt="aimage" width="300" />
      <div className="tags-container">{tags}</div>
    </div>
  );
};
