import axios from 'axios';
import { IImageData } from './common';
const baseURL = 'http://localhost:8090/';

export let getPhotos = (query: string, callback: (allImageData: IImageData[]) => void) => {
  axios
    .get(baseURL + 'get-photos', {
      params: {
        query: query,
      },
      headers: {},
    })
    .then(
      (response) => {
        callback(response.data);
      },
      (error) => {
        // consider using a toaster popup
        console.log(error);
      },
    );
};

export let getAllPhotos = (callback: (allImageData: IImageData[]) => void) => {
  axios
    .get(baseURL + 'get-all-photos', {
      headers: {},
    })
    .then(
      (response) => {
        callback(response.data);
      },
      (error) => {
        // consider using a toaster popup
        console.log(error);
      },
    );
};

export let uploadPhotos = (files: File[], tags: string[], callback: (allImageData: IImageData[]) => void) => {
  // we upload the photos individually so if one upload fails, the earlier uploads will still go through

  let allPromises: Promise<any>[] = [];
  files.forEach((file) => {
    let formData = new FormData();
    formData.append('file', file, file.name);
    formData.append('tags', JSON.stringify(tags));
    allPromises.push(
      axios.post(baseURL + 'upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }),
    );
  });

  Promise.all(allPromises).then(() => {
    getAllPhotos(callback);
  });
};
