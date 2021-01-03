import axios from 'axios';
import { IImageData } from './common';
const baseURL = 'http://localhost:8090/';

export let getAllPhotos = (callback: (allImageData: IImageData[]) => void) => {
  axios
    .get(baseURL + 'get-all-photos', {
      headers: {},
    })
    .then(
      (response) => {
        console.log(response);
        callback(response.data);
      },
      (error) => {
        // consider using a toaster popup
        console.log(error);
      },
    );
};

export let uploadPhotos = (files: File[], tags: string[]) => {
  // we upload the photos individually so if one upload fails, the earlier uploads will still go through
  files.forEach((file) => {
    let formData = new FormData();
    formData.append('file', file, file.name);
    formData.append('tags', JSON.stringify(tags));
    axios.post(baseURL + 'upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  });
};
