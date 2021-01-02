import React from 'react';
import { UploadForm } from './components/upload-form';
import { Button, Dialog } from '@blueprintjs/core';
import './App.css';

interface IAppProps {}
interface IAppState {
  isUploadOverlayOpen: boolean;
}

class App extends React.Component<IAppProps, IAppState> {
  constructor(props: any) {
    super(props);
    this.state = {
      isUploadOverlayOpen: false,
    };
    this.toggleOverlay = this.toggleOverlay.bind(this);
  }

  public render() {
    const { isUploadOverlayOpen } = this.state;
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
      </div>
    );
  }

  toggleOverlay() {
    this.setState({
      isUploadOverlayOpen: !this.state.isUploadOverlayOpen,
    });
  }
}

export default App;
