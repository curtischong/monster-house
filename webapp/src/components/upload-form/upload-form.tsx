import React from 'react';
import ImageUploader from 'react-images-upload';
import { Button, TagInput, AnchorButton } from '@blueprintjs/core';
import { uploadPhotos } from '../../request';
import { IImageData } from '../../common';
import './upload-form.css';

interface IUploadFormProps {
  setPhotos(allImageData: IImageData[]): void;
  onClose: () => void;
}
interface IUploadFormState {
  pictures: File[];
  tags: React.ReactNode[];
}

export class UploadForm extends React.Component<IUploadFormProps, IUploadFormState> {
  constructor(props: any) {
    super(props);
    this.state = {
      pictures: [],
      tags: [],
    };
    this.onDrop = this.onDrop.bind(this);
    this.onSubmit = this.onSubmit.bind(this);
  }

  onDrop(picture: File[]) {
    let { pictures } = this.state;
    this.setState({
      pictures: pictures.concat(picture),
    });
  }

  handleChangeTags = (tags: React.ReactNode[]) => {
    this.setState({ tags });
  };

  handleClearTags = () => this.handleChangeTags([]);

  clearButton() {
    return (
      <Button icon={this.state.tags.length > 1 ? 'cross' : 'refresh'} minimal={true} onClick={this.handleClearTags} />
    );
  }

  getTags = (): string[] => {
    let tags: string[] = [];
    this.state.tags.forEach((tag) => {
      if (tag !== undefined && tag !== null) {
        tags.push(tag.toString());
      }
    });
    return tags;
  };

  onSubmit() {
    const tags = this.getTags();
    this.handleClearTags();
    uploadPhotos(this.state.pictures, tags, this.props.setPhotos);
    this.props.onClose();
  }

  public render() {
    const { tags } = this.state;
    return (
      <>
        <ImageUploader
          className="image-uploader"
          withIcon={true}
          withPreview={true}
          buttonText="Choose images"
          onChange={this.onDrop}
          imgExtension={['.jpg', '.jpeg', '.gif', '.png']}
          maxFileSize={5242880}
          singleImage={false}
        />
        <TagInput
          className="tag-input"
          onChange={this.handleChangeTags}
          placeholder="Separate tags with commas..."
          rightElement={this.clearButton()}
          values={tags}
        />
        <AnchorButton
          className="pt-button pt-intent-success submit-image-button"
          onClick={this.onSubmit}
          text="Upload"
          type="submit"
        />
      </>
    );
  }
}
