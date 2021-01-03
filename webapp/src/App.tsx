import React from 'react';
import { Button, Dialog } from '@blueprintjs/core';
import './App.css';
import { getAllPhotos } from './request';
import { Image } from './components/image/image';
import { UploadForm } from './components/upload-form/upload-form';

interface IAppProps {}
interface IAppState {
  isUploadOverlayOpen: boolean;
  imageUrls: string[];
}

class App extends React.Component<IAppProps, IAppState> {
  constructor(props: any) {
    super(props);
    this.state = {
      isUploadOverlayOpen: false,
      imageUrls: [],
    };
    this.toggleOverlay = this.toggleOverlay.bind(this);
    this.setPhotos = this.setPhotos.bind(this);

    getAllPhotos(this.setPhotos);
  }

  public render() {
    const { isUploadOverlayOpen, imageUrls } = this.state;
    const photos = this.generatePhotos(imageUrls);
    return (
      <div className="App">
        <header className="App-header">
          <h1>Monster House</h1>
        </header>
        <div>
          <Button text="Upload photos" onClick={this.toggleOverlay} />
          <Dialog isOpen={isUploadOverlayOpen} onClose={this.toggleOverlay}>
            <UploadForm onClose={this.toggleOverlay} />
          </Dialog>
        </div>
        <div className="images-container">{photos}</div>
      </div>
    );
  }

  generatePhotos(imageUrls: string[]): JSX.Element[] {
    let images: JSX.Element[] = [];
    for (let i = 0; i < imageUrls.length; i++) {
      images.push(<Image key={i} url={imageUrls[i]}></Image>);
    }
    return images;
  }

  setPhotos(imageUrls: string[]) {
    this.setState({
      imageUrls: imageUrls,
    });
  }

  toggleOverlay() {
    this.setState({
      isUploadOverlayOpen: !this.state.isUploadOverlayOpen,
    });
  }
}

export default App;
