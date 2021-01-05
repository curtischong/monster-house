import React from 'react';
import { IImageData } from '../../common';
import { getAllPhotos, getPhotos } from '../../request';
import './search-field.css';

interface ISearchFieldProps {
  setPhotos(allImageData: IImageData[]): void;
}
interface ISearchFieldState {
  value: string;
}

export class SearchField extends React.Component<ISearchFieldProps, ISearchFieldState> {
  constructor(props: any) {
    super(props);
    this.state = {
      value: '',
    };
    this.handleChange = this.handleChange.bind(this);
    this.keyPress = this.keyPress.bind(this);
  }

  handleChange(event: React.ChangeEvent<HTMLInputElement>) {
    this.setState({ value: event.target.value });
  }
  keyPress(event: React.KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter') {
      if (this.state.value === '') {
        getAllPhotos(this.props.setPhotos);
      } else {
        getPhotos(this.state.value, this.props.setPhotos);
      }
    }
  }

  render() {
    return (
      <div className="bp3-input-group search-field">
        <span className="bp3-icon bp3-icon-search"></span>
        <input
          className="bp3-input"
          type="search"
          placeholder="Search for specific tags separated with ','"
          dir="auto"
          onChange={this.handleChange}
          onKeyDown={this.keyPress}
        />
      </div>
    );
  }
}
