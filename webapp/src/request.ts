import axios, { AxiosResponse } from 'axios';
import { IImageData } from './common';
const baseURL = 'http://localhost:8090/';

/**
 * Fetches photos from the database that contains at least one tag specified in the query
 */
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

/**
 * Fetches all photos from the database
 */
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

/**
 * Uploads the photos to the server
 * @remarks
 * We upload the photos individually so if one upload fails, the earlier uploads will still go through
 */
export let uploadPhotos = (files: File[], tags: string[], callback: (allImageData: IImageData[]) => void) => {
  let allPromises: Promise<AxiosResponse<any>>[] = [];
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
