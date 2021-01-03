import React from 'react';
import { Button, Dialog } from '@blueprintjs/core';
import './App.css';
import { getAllPhotos } from './request';
import { Image } from './components/image/image';
import { UploadForm } from './components/upload-form/upload-form';
import { SearchField } from './components/search-field/search-field';
import { IImageData } from './common';

interface IAppProps {}
interface IAppState {
  isUploadOverlayOpen: boolean;
  allImageData: IImageData[];
}

class App extends React.Component<IAppProps, IAppState> {
  constructor(props: any) {
    super(props);
    this.state = {
      isUploadOverlayOpen: false,
      allImageData: [],
    };
    this.toggleOverlay = this.toggleOverlay.bind(this);
    this.setPhotos = this.setPhotos.bind(this);

    getAllPhotos(this.setPhotos);
  }

  public render() {
    const { isUploadOverlayOpen, allImageData } = this.state;
    const photos = this.generatePhotos(allImageData);
    return (
      <div className="App">
        <header className="App-header">
          <h1>Monster House</h1>
        </header>
        <div>
          <Button text="Upload photos" onClick={this.toggleOverlay} />
          <SearchField setPhotos={this.setPhotos} />
          <Dialog isOpen={isUploadOverlayOpen} onClose={this.toggleOverlay}>
            <UploadForm onClose={this.toggleOverlay} />
          </Dialog>
        </div>
        <div className="images-container">{photos}</div>
      </div>
    );
  }

  generatePhotos(allImageData: IImageData[]): JSX.Element[] {
    let images: JSX.Element[] = [];
    for (let i = 0; i < allImageData.length; i++) {
      images.push(<Image key={i} imageData={allImageData[i]}></Image>);
    }
    return images;
  }

  setPhotos(allImageData: IImageData[]) {
    this.setState({
      allImageData: allImageData,
    });
  }

  toggleOverlay() {
    this.setState({
      isUploadOverlayOpen: !this.state.isUploadOverlayOpen,
    });
  }
}

export default App;
