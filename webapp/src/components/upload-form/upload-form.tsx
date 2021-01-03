import React from 'react';
import ImageUploader from 'react-images-upload';
import { Button, TagInput, AnchorButton } from '@blueprintjs/core';
import { uploadPhotos } from '../../request';

interface IUploadFormProps {
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
    uploadPhotos(this.state.pictures, tags);
    this.props.onClose();
  }

  public render() {
    const { tags } = this.state;
    return (
      <>
        <ImageUploader
          withIcon={true}
          withPreview={true}
          buttonText="Choose images"
          onChange={this.onDrop}
          imgExtension={['.jpg', '.gif', '.png', '.gif']}
          maxFileSize={5242880}
          singleImage={false}
        />
        <TagInput
          onChange={this.handleChangeTags}
          placeholder="Separate values with commas..."
          rightElement={this.clearButton()}
          values={tags}
        />
        <AnchorButton className="pt-button pt-intent-success" onClick={this.onSubmit} text="Upload" type="submit" />
      </>
    );
  }
}
